package main

import (
	"log"

	ccgobblesserver "github.com/avvarikrish/chefcurrygobbles/ccgobbles_server"
)

func main() {
	config := "config/ccgobbles_server.yml"
	log.Fatalf("Failed to start: %v", ccgobblesserver.New(config).Start())
}
