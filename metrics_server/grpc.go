package metrics

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/metrics_server"
)

func (m *MetricServer) startGRPC() error {
	err := connectToMongo()
	if err != nil {
		log.Fatalf("Error while connecting to Mongo: %v", err)
	}
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterMetricsServer(server, m)
	go func() {
		fmt.Println("Connecting to 8080")
		err := connectToPrometheus()
		if err != nil {
			log.Fatalf("Error while connecting to Prometheus: %v", err)
		}
	}()
	go func() {
		fmt.Println("Starting Server")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()
	select {}
}

func connectToPrometheus() error {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return fmt.Errorf("Prometheus http error: %v", err)
	}
	return nil
}

func connectToMongo() error {
	fmt.Println("Connecting to Mongo")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		return fmt.Errorf("error while connecting to mongodb: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %v", err)
	}
	// ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("could not ping to mongo db service: %v", err)
	}
	fmt.Println("Successfully connected to mongo")

	orderCollection = client.Database("chefcurrygobbles").Collection("Orders")
	return nil
}
