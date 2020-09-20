package orders_server

import (
	"context"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/orders_server"
)

func (o *OrdersServer) startGRPC() error {
	if !o.initialized {
		return fmt.Errorf("server not initialized")
	}

	lis, err := net.Listen(o.cfg.Server.Network, o.cfg.Server.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterOrdersServer(server, o)

	go func() {
		if err := o.connectToMongo(); err != nil {
			log.Fatalf("Error while connecting to Mongo: %v", err)
		}
	}()
	go func() {
		log.Info("Starting Server")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()

	// properly close everything
	select {}
}

func (o *OrdersServer) connectToMongo() error {
	log.Info("Connecting to Mongo")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := o.mongoClient.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = o.mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("could not ping to mongo db service: %v", err)
	}

	log.Info("Successfully connected to mongo")

	return nil
}
