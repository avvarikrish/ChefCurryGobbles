package main

import (
	"log"

	metricsserver "github.com/avvarikrish/chefcurrygobbles/metrics_server"
)

func main() {
	log.Fatalf("Failed to start: %v", metricsserver.New().Start())
}
