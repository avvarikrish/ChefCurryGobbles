package metrics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/avvarikrish/chefcurrygobbles/pkg/algo"
	pb "github.com/avvarikrish/chefcurrygobbles/proto/metrics_server"
)

// MetricServer represents an instance of the metrics server
type MetricServer struct{}

var orderCollection *mongo.Collection

type metric interface {
	calculate(t time.Duration) float64
}

type average struct{}
type percentile struct {
	value float64
}

type metricRun struct {
	CalcType metric
	Time     time.Duration
	Name     string
	Help     string
	Gauge    *prometheus.Gauge
}

type order struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Items   []interface{}      `bson:"items" json:"items"`
	Time    time.Time          `bson:"time" json:"time"`
	OrderID string             `bson:"order_id" json:"order_id"`
	Email   string             `bson:"email" json:"email"`
	RestID  string             `bson:"rest_id" json:"rest_id"`
}

// New returns a new initialized instance of MetricServer.
func New() *MetricServer {
	return &MetricServer{}
}

// Start enables the Metrics Server service.
func (m *MetricServer) Start() error {
	// start grpc
	return m.startGRPC()

	// start actual program
	// return m.startProgram()
}

// RunMetricScrape starts a new metric scrape based on the duration and calculation type
func (m *MetricServer) RunMetricScrape(ctx context.Context, req *pb.RunMetricScrapeRequest) (*pb.RunMetricScrapeResponse, error) {
	fmt.Println("Creating a new metric scrape process")

	metricRun, err := getMetricRun(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid argument")
	}
	go startScrape(metricRun)

	// select {}
	return &pb.RunMetricScrapeResponse{
		Response: proto.String("Successfully started scrape"),
	}, nil

}

func startScrape(mr *metricRun) error {
	for {
		fmt.Println("Setting value")
		(*mr.Gauge).Set(mr.CalcType.calculate(mr.Time))
		time.Sleep(mr.Time)
	}
}

func getMetricRun(req *pb.RunMetricScrapeRequest) (*metricRun, error) {
	metric := getMetric(req.GetCalcType())
	name := req.GetName()
	help := req.GetHelp()
	if metric == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid argument")
	}
	gauge := registerMetric(name, help)
	time := getTime(req.GetTime())

	return &metricRun{
		CalcType: metric,
		Time:     time,
		Name:     name,
		Help:     help,
		Gauge:    gauge,
	}, nil
}

func registerMetric(name string, help string) *prometheus.Gauge {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
	prometheus.MustRegister(metric)
	return &metric
}

func getMetric(m pb.CalcType) metric {
	switch m {
	case pb.CalcType_CALC_AVERAGE:
		return average{}
	case pb.CalcType_CALC_PERCENTILE:
		return percentile{value: 0.95}
	default:
		return nil
	}
}

func getTime(t string) time.Duration {
	switch t {
	case "s":
		return time.Second * 5
	case "m":
		return time.Minute
	case "h":
		return time.Hour
	default:
		return time.Second
	}
}

func (average) calculate(t time.Duration) float64 {
	now := time.Now()
	filter := bson.M{"time": bson.M{"$lt": now, "$gte": now.Add(-t)}}

	cur, err := orderCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error while finding collection: %v\n", err)
	}
	defer cur.Close(context.Background())
	sum := 0
	size := 0
	for cur.Next(context.Background()) {
		result := order{}
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		sum += len(result.Items)
		size++
	}
	if err := cur.Err(); err != nil {
		log.Fatalf("Error with cursor: %v\n", err)
	}

	if size > 0 {
		return float64(sum) / float64(size)
	}
	return float64(0)
}

func (p percentile) calculate(t time.Duration) float64 {
	now := time.Now()
	filter := bson.M{"time": bson.M{"$lt": now, "$gt": now.Add(-t)}}

	cur, err := orderCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error while finding collection: %v\n", err)
	}
	defer cur.Close(context.Background())
	var itemCounts []int
	for cur.Next(context.Background()) {
		result := order{}
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		itemCounts = append(itemCounts, len(result.Items))
	}
	if err := cur.Err(); err != nil {
		log.Fatalf("Error with cursor: %v\n", err)
	}
	size := len(itemCounts)
	if size == 0 {
		return float64(0)
	}
	return float64(algo.Ksmallest(itemCounts, int(p.value*float64(size)), 0, size-1))
}
