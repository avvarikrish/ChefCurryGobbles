package restaurant_server

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
	"github.com/avvarikrish/chefcurrygobbles/restaurant_server/config"

	pb "github.com/avvarikrish/chefcurrygobbles/proto/restaurant_server"
)

// RestaurantServer represents a new instance of the server
type RestaurantServer struct {
	cfg         config.RestaurantServerConfig
	initialized bool

	mongoClient     *mongo.Client
	db              *mongo.Database
	userCollection  *mongo.Collection
	restCollection  *mongo.Collection
	orderCollection *mongo.Collection
}

// New returns a new initialized instance of RestaurantServer.
func New(file string) *RestaurantServer {
	return &RestaurantServer{
		cfg: config.NewConfig(file),
	}
}

// Start starts the CCGoblesServer service.
func (r *RestaurantServer) Start() error {
	if err := r.initialize(); err != nil {
		return fmt.Errorf("failed to initialize app: %v", err)
	}

	return r.startGRPC()
}

func (r *RestaurantServer) initialize() error {
	if err := r.setupMongo(); err != nil {
		return fmt.Errorf("error while connecting to mongo: %v", err)
	}

	r.initialized = true
	return nil
}

func (r *RestaurantServer) setupMongo() error {
	log.Info("creating mongo client, db, collections")

	client, err := mongo.NewClient(options.Client().ApplyURI(r.cfg.Mongo.MongoServer))
	if err != nil {
		return fmt.Errorf("error while creating mongo client: %v", err)
	}

	r.mongoClient = client
	r.db = r.mongoClient.Database(r.cfg.Mongo.Database)
	r.userCollection = r.db.Collection(r.cfg.Mongo.Collections.Users)
	r.restCollection = r.db.Collection(r.cfg.Mongo.Collections.Restaurants)
	r.orderCollection = r.db.Collection(r.cfg.Mongo.Collections.Orders)

	return nil
}

// AddRestaurant adds a restaurant to the db
func (r *RestaurantServer) AddRestaurant(_ context.Context, req *pb.AddRestaurantRequest) (*pb.AddRestaurantResponse, error) {
	log.Info("Adding restaurant")

	resReq := req.GetRestaurant()
	emailFilterRes := r.restExists(bson.M{"email": resReq.GetEmail()})
	phoneFilterRes := r.restExists(bson.M{"phone": resReq.GetPhone()})
	if emailFilterRes != nil || phoneFilterRes != nil {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("restaurant already exists"))
	}

	insCtx, insCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer insCancel()
	_, err := r.restCollection.InsertOne(insCtx, bc.CreateRestBson(resReq))
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while inserting restaurant: %v", err))
	}

	return &pb.AddRestaurantResponse{
		Response: "Successfully added restaurant",
	}, nil
}

func (r *RestaurantServer) restExists(filter bson.M) *bc.Restaurant {
	resObj := &bc.Restaurant{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	checkExist := r.restCollection.FindOne(ctx, filter)
	if err := checkExist.Decode(resObj); err == nil {
		return resObj
	}

	return nil
}
