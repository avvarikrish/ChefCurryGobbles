package ccgobblesserver

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
	data := &bc.User{}
	filter := bson.M{"email": userReq.GetEmail()}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	checkRes := c.userCollection.FindOne(ctx, filter)
	if err := checkRes.Decode(data); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("email already exists: %v", userReq.GetEmail()))
	}

	insCtx, insCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer insCancel()
	_, err := c.userCollection.InsertOne(insCtx, bc.CreateUserBson(userReq))
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while inserting into mongo: %v", err))
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
	data := &bc.User{}
	filter := bson.M{"email": loginEmail}
	res := c.userCollection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("email not found: %v", loginEmail))
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
func (c *CcgobblesServer) UpdateUser(_ context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Info("Update user")

	email := req.GetOldEmail()
	filter := bson.M{"email": email}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := c.userCollection.FindOne(ctx, filter)
	if err := res.Decode(&bc.User{}); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("email not found: %v", email))
	}

	replaceCtx, replaceCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer replaceCancel()
	_, err := c.userCollection.ReplaceOne(replaceCtx, filter, bc.CreateUserBson(req.GetUser()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while replacing in mongo: %v", err))
	}

	return &pb.UpdateUserResponse{
		Response: proto.String("Successfully updated user"),
	}, nil
}

// DeleteUser deletes the specified user from the db
func (c *CcgobblesServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	log.Info("Deleting user")

	emailToDelete := req.GetEmail()
	filter := bson.M{"email": emailToDelete}
	delRes, delErr := c.userCollection.DeleteOne(ctx, filter)
	if delErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while deleting in mongo: %v", delErr))
	}
	if delRes.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("email not found"))
	}

	return &pb.DeleteUserResponse{
		Response: proto.String("Successfully deleted user"),
	}, nil
}

// AddRestaurant adds a restaurant to the db
func (c *CcgobblesServer) AddRestaurant(_ context.Context, req *pb.AddRestaurantRequest) (*pb.AddRestaurantResponse, error) {
	log.Info("Adding restaurant")

	resReq := req.GetRestaurant()
	emailFilterRes := c.restExists(bson.M{"email": resReq.GetEmail()})
	phoneFilterRes := c.restExists(bson.M{"phone": resReq.GetPhone()})
	if emailFilterRes != nil || phoneFilterRes != nil {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("restaurant already exists"))
	}

	insCtx, insCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer insCancel()
	_, err := c.restCollection.InsertOne(insCtx, bc.CreateRestBson(resReq))
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while inserting restaurant: %v", err))
	}

	return &pb.AddRestaurantResponse{
		Response: proto.String("Successfully added restaurant"),
	}, nil
}

// CreateOrder creates an order given a user, restaurant, and menu items
func (c *CcgobblesServer) CreateOrder(_ context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log.Info("Creating Order")

	orderReq := req.GetOrder()
	resPhoneFilterRes := c.restExists(bson.M{"phone": orderReq.GetRestPhone()})
	resEmailFilterRes := c.restExists(bson.M{"email": orderReq.GetRestEmail()})
	userFilter := bson.M{"email": orderReq.GetEmail()}

	ctxUser, cancelUser := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelUser()
	checkUserExist := c.userCollection.FindOne(ctxUser, userFilter)

	resObj := &bc.Restaurant{}
	userObj := &bc.User{}
	if resPhoneFilterRes != nil {
		resObj = resPhoneFilterRes
	} else if resEmailFilterRes != nil {
		resObj = resEmailFilterRes
	} else {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("restaurant does not exist"))
	}

	if err := checkUserExist.Decode(userObj); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user does not exist: %v", orderReq.GetEmail()))
	}

	if len(orderReq.GetOrderItem()) <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("need more than 1 item"))
	}

	ctxIns, cancelIns := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelIns()
	_, err := c.orderCollection.InsertOne(ctxIns, bc.CreateOrderBson(req, resObj.ID, userObj.ID))

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	return &pb.CreateOrderResponse{
		Response: proto.String("Successfully created order"),
	}, nil
}

func (c *CcgobblesServer) restExists(filter bson.M) *bc.Restaurant {
	resObj := &bc.Restaurant{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	checkExist := c.restCollection.FindOne(ctx, filter)
	if err := checkExist.Decode(resObj); err == nil {
		return resObj
	}

	return nil
}
