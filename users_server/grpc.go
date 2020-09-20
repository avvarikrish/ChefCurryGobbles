package users_server

import (
	"context"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/users_server"
)

func (u *UsersServer) startGRPC() error {
	if !u.initialized {
		return fmt.Errorf("server not initialized")
	}

	lis, err := net.Listen(u.cfg.Server.Network, u.cfg.Server.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterUsersServer(server, u)

	go func() {
		if err := u.connectToMongo(); err != nil {
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

func (u *UsersServer) connectToMongo() error {
	log.Info("Connecting to Mongo")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := u.mongoClient.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = u.mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("could not ping to mongo db service: %v", err)
	}

	log.Info("Successfully connected to mongo")

	return nil
}
