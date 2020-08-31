package ccgobblesserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

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
		err := connectToMongo()
		if err != nil {
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
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("could not ping to mongo db service: %v", err)
	}
	fmt.Println("Successfully connected to mongo")

	userCollection = client.Database("chefcurrygobbles").Collection("Users")
	restCollection = client.Database("chefcurrygobbles").Collection("Restaurants")
	orderCollection = client.Database("chefcurrygobbles").Collection("Orders")
	return nil
}
