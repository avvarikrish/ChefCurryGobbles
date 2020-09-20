package main

import (
	"log"

	users_server "github.com/avvarikrish/chefcurrygobbles/users_server"
)

func main() {
	config := "config/users_server.yml"
	log.Fatalf("failed to start: %v", users_server.New(config).Start())
}
