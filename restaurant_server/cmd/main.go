package main

import (
	"log"

	"github.com/avvarikrish/chefcurrygobbles/restaurant_server"
)

func main() {
	config := "config/restaurant_server.yml"
	log.Fatalf("failed to start: %v", restaurant_server.New(config).Start())
}
