package ccgobblesserver

import (
	"context"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/ccgobbles_server"
)

func (c *CcgobblesServer) startGRPC() error {
	if !c.initialized {
		return fmt.Errorf("server not initialized")
	}

	lis, err := net.Listen(c.cfg.Server.Network, c.cfg.Server.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterCCGobblesServer(server, c)

	go func() {
		if err := c.connectToMongo(); err != nil {
			log.Fatalf("Error while connecting to Mongo: %v", err)
		}
	}()
	go func() {
		fmt.Println("Starting Server")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()

	// properly close everything
	select {}
}

func (c *CcgobblesServer) connectToMongo() error {
	log.Info("Connecting to Mongo")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := c.mongoClient.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = c.mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("could not ping to mongo db service: %v", err)
	}

	log.Info("Successfully connected to mongo")

	return nil
}
