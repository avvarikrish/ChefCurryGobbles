package ccgobblesserver

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/avvarikrish/chefcurrygobbles/ccgobbles_server/config"
	bc "github.com/avvarikrish/chefcurrygobbles/pkg/bsonconversion"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/ccgobbles_server"
)

// CcgobblesServer represents a new instance of the server
type CcgobblesServer struct {
	cfg         config.CcgobblesServerConfig
	initialized bool

	mongoClient     *mongo.Client
	db              *mongo.Database
	userCollection  *mongo.Collection
	restCollection  *mongo.Collection
	orderCollection *mongo.Collection
}

type user struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
}

type restaurant struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	RestID string             `bson:"rest_id"`
	Name   string             `bson:"name"`
}

// New returns a new initialized instance of CCGobblesServer.
func New(file string) *CcgobblesServer {
	return &CcgobblesServer{
		cfg: config.NewConfig(file),
	}
}

// Start starts the CCGoblesServer service.
func (c *CcgobblesServer) Start() error {
	if err := c.initialize(); err != nil {
		return fmt.Errorf("failed to initialize app: %v", err)
	}

	return c.startGRPC()
}

func (c *CcgobblesServer) initialize() error {
	if err := c.setupMongo(); err != nil {
		return fmt.Errorf("error while connecting to mongo: %v", err)
	}

	c.initialized = true
	return nil
}

func (c *CcgobblesServer) setupMongo() error {
	log.Info("creating mongo client, db, collections")

	client, err := mongo.NewClient(options.Client().ApplyURI(c.cfg.Mongo.MongoServer))
	if err != nil {
		return fmt.Errorf("error while creating mongo client: %v", err)
	}

	c.mongoClient = client
	c.db = c.mongoClient.Database(c.cfg.Mongo.Database)
	c.userCollection = c.db.Collection(c.cfg.Mongo.Collections.Users)
	c.restCollection = c.db.Collection(c.cfg.Mongo.Collections.Restaurants)
	c.orderCollection = c.db.Collection(c.cfg.Mongo.Collections.Orders)

	return nil
}

// RegisterUser creates a new user
func (c *CcgobblesServer) RegisterUser(_ context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	log.Info("Registering user")

	userReq := req.GetUser()
	data := &user{}
	filter := bson.M{"email": userReq.GetEmail()}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	checkRes := c.userCollection.FindOne(ctx, filter)
	if err := checkRes.Decode(data); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("email already exists: %v\n", userReq.GetEmail()))
	}

	insCtx, insCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer insCancel()
	_, err := c.userCollection.InsertOne(insCtx, bc.CreateUserBson(userReq))
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while inserting into mongo: %v\n", err))
	}

	return &pb.RegisterUserResponse{
		Response: proto.String("Successfully Added User"),
	}, nil
}

// LoginUser logs in a user if they enter the right password and email
func (c *CcgobblesServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	log.Info("Sign in request")

	loginPassword := req.GetPassword()
	loginEmail := req.GetEmail()
	data := &user{}
	filter := bson.M{"email": loginEmail}
	res := c.userCollection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Email not found: %v\n", loginEmail))
	}

	if loginPassword == data.Password {
		return &pb.LoginUserResponse{
			Response: proto.Bool(true),
		}, nil
	}

	return &pb.LoginUserResponse{
		Response: proto.Bool(false),
	}, nil
}

// UpdateUser updates user info in the db
func (c *CcgobblesServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	fmt.Println("Attempting user update")

	userToUpdate := req.GetUser()
	email := userToUpdate.GetEmail()
	filter := bson.M{"email": email}
	res := c.userCollection.FindOne(ctx, filter)
	if err := res.Decode(&user{}); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Email not found: %v\n", email))
	}
	_, err := c.userCollection.ReplaceOne(ctx, filter, bc.CreateUserBson(userToUpdate))
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot update: %v\n", err))
	}
	return &pb.UpdateUserResponse{
		Response: proto.String("Successfully Updated User"),
	}, nil
}

// DeleteUser deletes the specified user from the db
func (c *CcgobblesServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	fmt.Println("Deleting User")

	emailToDelete := req.GetEmail()
	filter := bson.M{"email": emailToDelete}
	delRes, delErr := c.userCollection.DeleteOne(ctx, filter)
	if delErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v\n", delErr))
	}
	if delRes.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Email not found"))
	}
	return &pb.DeleteUserResponse{
		Response: proto.String("Successfully Deleted User"),
	}, nil
}

// AddRestaurant adds a restaurant to the db
func (c *CcgobblesServer) AddRestaurant(ctx context.Context, req *pb.AddRestaurantRequest) (*pb.AddRestaurantResponse, error) {
	fmt.Println("Adding restaurant")

	resReq := req.GetRestaurant()
	filter := bson.M{"rest_id": resReq.GetRestId()}
	checkExist := c.restCollection.FindOne(ctx, filter)
	if err := checkExist.Decode(&restaurant{}); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("Restaurant already exists: %v\n", resReq.GetRestId()))
	}
	_, err := c.restCollection.InsertOne(ctx, bc.CreateRestBson(resReq))
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v\n", err))
	}
	return &pb.AddRestaurantResponse{
		Response: proto.String("Successfully Added Restaurant"),
	}, nil
}

// CreateOrder creates an order given a user, restaurant, and menu items
func (c *CcgobblesServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// need to add a lot of checks

	fmt.Println("Creating Order")

	dt := time.Now()

	// check if restaurant exists
	resFilter := bson.M{"rest_id": req.GetRestId()}
	userFilter := bson.M{"email": req.GetEmail()}
	checkExist := c.restCollection.FindOne(ctx, resFilter)
	checkUserExist := c.userCollection.FindOne(ctx, userFilter)
	if err := checkExist.Decode(&restaurant{}); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Restaurant does not exist: %v\n", req.GetRestId()))
	}
	if err := checkUserExist.Decode(&user{}); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("User does not exist: %v\n", req.GetEmail()))
	}

	if len(req.GetOrderItem()) <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Need more than 1 item"))
	}

	_, err := c.orderCollection.InsertOne(ctx, bc.CreateOrderBson(req, dt))

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v\n", err))
	}

	return &pb.CreateOrderResponse{
		Response: proto.String("Successfully created order"),
	}, nil
}
