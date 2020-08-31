package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/avvarikrish/chefcurrygobbles/pkg/input"

	mpb "github.com/avvarikrish/chefcurrygobbles/proto/metrics_server"
)

type Metric struct {
	Name string
	Help string
	Time string
}

func main() {
	fmt.Println("Hello I'm a metrics client")

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	mm, err := grpc.Dial("localhost:50052", opts...)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer mm.Close()

	m := mpb.NewMetricsClient(mm)
	sendToMetrics(m, os.Args[1])
}

func sendToMetrics(m mpb.MetricsClient, input_metric string) {
	t := reflect.TypeOf(Metric{})
	v := reflect.New(t)
	input.ReadInput(t, v.Elem())
	metric := v.Interface().(*Metric)

	req := &mpb.RunMetricScrapeRequest{
		Time:     proto.String(metric.Time),
		CalcType: getCalcType(input_metric),
		Name:     proto.String(metric.Name),
		Help:     proto.String(metric.Help),
	}

	res, err := m.RunMetricScrape(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while running metric scrape: %v\n", err)
	}
	fmt.Printf("Success: %v\n", res)
}

func getCalcType(c string) *mpb.CalcType {
	switch c {
	case "average":
		return enumToCalcPointer(mpb.CalcType_CALC_AVERAGE)
	case "percentile":
		return enumToCalcPointer(mpb.CalcType_CALC_PERCENTILE)
	default:
		return enumToCalcPointer(mpb.CalcType_CALC_UNKNOWN)
	}
}

func enumToCalcPointer(c mpb.CalcType) *mpb.CalcType {
	return &c
}
