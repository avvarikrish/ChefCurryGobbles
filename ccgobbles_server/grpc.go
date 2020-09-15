package ccgobblesserver

import (
	"context"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/ccgobbles_server"
)

func (c *CcgobblesServer) startGRPC() error {
	lis, err := net.Listen("tcp", ":50051")
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

	client, err := mongo.NewClient(options.Client().ApplyURI(c.cfg.Mongo.MongoServer))
	if err != nil {
		return fmt.Errorf("error while connecting to mongodb: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("could not ping to mongo db service: %v", err)
	}

	log.Info("Successfully connected to mongo")

	userCollection = client.Database(c.cfg.Mongo.Database).Collection(c.cfg.Mongo.Collections.Users)
	restCollection = client.Database(c.cfg.Mongo.Database).Collection(c.cfg.Mongo.Collections.Restaurants)
	orderCollection = client.Database(c.cfg.Mongo.Database).Collection(c.cfg.Mongo.Collections.Orders)

	return nil
}
