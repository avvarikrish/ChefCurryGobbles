package main

import (
	"log"

	"github.com/avvarikrish/chefcurrygobbles/orders_server"
)

func main() {
	config := "config/orders_server.yml"
	log.Fatalf("Failed to start: %v", orders_server.New(config).Start())
}
