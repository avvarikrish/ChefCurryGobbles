package users_server

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bc "github.com/avvarikrish/chefcurrygobbles/pkg/bsonconversion"
	"github.com/avvarikrish/chefcurrygobbles/users_server/config"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/users_server"
)

// UsersServer represents a new instance of the server
type UsersServer struct {
	cfg         config.UsersServerConfig
	initialized bool

	mongoClient     *mongo.Client
	db              *mongo.Database
	userCollection  *mongo.Collection
	restCollection  *mongo.Collection
	orderCollection *mongo.Collection
}

// New returns a new initialized instance of UsersServer.
func New(file string) *UsersServer {
	return &UsersServer{
		cfg: config.NewConfig(file),
	}
}

// Start starts the CCGoblesServer service.
func (u *UsersServer) Start() error {
	if err := u.initialize(); err != nil {
		return fmt.Errorf("failed to initialize app: %v", err)
	}

	return u.startGRPC()
}

func (u *UsersServer) initialize() error {
	if err := u.setupMongo(); err != nil {
		return fmt.Errorf("error while connecting to mongo: %v", err)
	}

	u.initialized = true
	return nil
}

func (u *UsersServer) setupMongo() error {
	log.Info("creating mongo client, db, collections")

	client, err := mongo.NewClient(options.Client().ApplyURI(u.cfg.Mongo.MongoServer))
	if err != nil {
		return fmt.Errorf("error while creating mongo client: %v", err)
	}

	u.mongoClient = client
	u.db = u.mongoClient.Database(u.cfg.Mongo.Database)
	u.userCollection = u.db.Collection(u.cfg.Mongo.Collections.Users)
	u.restCollection = u.db.Collection(u.cfg.Mongo.Collections.Restaurants)
	u.orderCollection = u.db.Collection(u.cfg.Mongo.Collections.Orders)

	return nil
}

// RegisterUser creates a new user
func (u *UsersServer) RegisterUser(_ context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	log.Info("Registering user")

	userReq := req.GetUser()
	data := &bc.User{}
	filter := bson.M{"email": userReq.GetEmail()}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	checkRes := u.userCollection.FindOne(ctx, filter)
	if err := checkRes.Decode(data); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("email already exists: %v", userReq.GetEmail()))
	}

	insCtx, insCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer insCancel()
	_, err := u.userCollection.InsertOne(insCtx, bc.CreateUserBson(userReq))
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while inserting into mongo: %v", err))
	}

	return &pb.RegisterUserResponse{
		Response: "Successfully Added User",
	}, nil
}

// LoginUser logs in a user if they enter the right password and email
func (u *UsersServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	log.Info("Sign in request")

	loginPassword := req.GetPassword()
	loginEmail := req.GetEmail()
	data := &bc.User{}
	filter := bson.M{"email": loginEmail}
	res := u.userCollection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("email not found: %v", loginEmail))
	}

	if loginPassword == data.Password {
		return &pb.LoginUserResponse{
			Response: true,
		}, nil
	}

	return &pb.LoginUserResponse{
		Response: false,
	}, nil
}

// UpdateUser updates user info in the db
func (u *UsersServer) UpdateUser(_ context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Info("Update user")

	email := req.GetOldEmail()
	filter := bson.M{"email": email}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := u.userCollection.FindOne(ctx, filter)
	if err := res.Decode(&bc.User{}); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("email not found: %v", email))
	}

	replaceCtx, replaceCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer replaceCancel()
	_, err := u.userCollection.ReplaceOne(replaceCtx, filter, bc.CreateUserBson(req.GetUser()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while replacing in mongo: %v", err))
	}

	return &pb.UpdateUserResponse{
		Response: "Successfully updated user",
	}, nil
}

// DeleteUser deletes the specified user from the db
func (u *UsersServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	log.Info("Deleting user")

	emailToDelete := req.GetEmail()
	filter := bson.M{"email": emailToDelete}
	delRes, delErr := u.userCollection.DeleteOne(ctx, filter)
	if delErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while deleting in mongo: %v", delErr))
	}
	if delRes.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("email not found"))
	}

	return &pb.DeleteUserResponse{
		Response: "Successfully deleted user",
	}, nil
}
