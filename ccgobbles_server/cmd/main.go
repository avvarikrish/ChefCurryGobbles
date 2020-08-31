package main

import (
	"log"

	ccgobblesserver "github.com/avvarikrish/chefcurrygobbles/ccgobbles_server"
)

func main() {
	log.Fatalf("Failed to start: %v", ccgobblesserver.New().Start())
}
