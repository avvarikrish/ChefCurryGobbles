package orders_server

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

	"github.com/avvarikrish/chefcurrygobbles/orders_server/config"
	bc "github.com/avvarikrish/chefcurrygobbles/pkg/bsonconversion"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/orders_server"
)

// OrdersServer represents a new instance of the server
type OrdersServer struct {
	cfg         config.OrdersServerConfig
	initialized bool

	mongoClient     *mongo.Client
	db              *mongo.Database
	userCollection  *mongo.Collection
	restCollection  *mongo.Collection
	orderCollection *mongo.Collection
}

// New returns a new initialized instance of OrdersServer.
func New(file string) *OrdersServer {
	return &OrdersServer{
		cfg: config.NewConfig(file),
	}
}

// Start starts the CCGoblesServer service.
func (o *OrdersServer) Start() error {
	if err := o.initialize(); err != nil {
		return fmt.Errorf("failed to initialize app: %v", err)
	}

	return o.startGRPC()
}

func (o *OrdersServer) initialize() error {
	if err := o.setupMongo(); err != nil {
		return fmt.Errorf("error while connecting to mongo: %v", err)
	}

	o.initialized = true
	return nil
}

func (o *OrdersServer) setupMongo() error {
	log.Info("creating mongo client, db, collections")

	client, err := mongo.NewClient(options.Client().ApplyURI(o.cfg.Mongo.MongoServer))
	if err != nil {
		return fmt.Errorf("error while creating mongo client: %v", err)
	}

	o.mongoClient = client
	o.db = o.mongoClient.Database(o.cfg.Mongo.Database)
	o.userCollection = o.db.Collection(o.cfg.Mongo.Collections.Users)
	o.restCollection = o.db.Collection(o.cfg.Mongo.Collections.Restaurants)
	o.orderCollection = o.db.Collection(o.cfg.Mongo.Collections.Orders)

	return nil
}

// CreateOrder creates an order given a user, restaurant, and menu items
func (o *OrdersServer) CreateOrder(_ context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log.Info("Creating Order")

	orderReq := req.GetOrder()
	resPhoneFilterRes := o.restExists(bson.M{"phone": orderReq.GetRestPhone()})
	resEmailFilterRes := o.restExists(bson.M{"email": orderReq.GetRestEmail()})
	userFilter := bson.M{"email": orderReq.GetEmail()}

	ctxUser, cancelUser := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelUser()
	checkUserExist := o.userCollection.FindOne(ctxUser, userFilter)

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
	_, err := o.orderCollection.InsertOne(ctxIns, bc.CreateOrderBson(req, resObj.ID, userObj.ID))

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error: %v", err))
	}

	return &pb.CreateOrderResponse{
		Response: "Successfully created order",
	}, nil
}

func (o *OrdersServer) restExists(filter bson.M) *bc.Restaurant {
	resObj := &bc.Restaurant{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	checkExist := o.restCollection.FindOne(ctx, filter)
	if err := checkExist.Decode(resObj); err == nil {
		return resObj
	}

	return nil
}
