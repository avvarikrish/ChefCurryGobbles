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

	mongoClient    *mongo.Client
	db             *mongo.Database
	restCollection *mongo.Collection
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
	r.restCollection = r.db.Collection(r.cfg.Mongo.Collections.Restaurants)

	return nil
}

// AddRestaurant adds a restaurant to the db
func (r *RestaurantServer) AddRestaurant(_ context.Context, req *pb.AddRestaurantRequest) (*pb.AddRestaurantResponse, error) {
	log.Info("Adding restaurant")

	resReq := req.GetRestaurant()
	if r.restExists(resReq.GetEmail(), resReq.GetPhone()) != nil {
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

func (r *RestaurantServer) UpdateRestaurant(ctx context.Context, req *pb.UpdateRestaurantRequest) (*pb.UpdateRestaurantResponse, error) {
	log.Info("Update restaurant")

	filter := r.restExists(req.GetOldEmail(), req.GetOldPhone())
	if filter == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("restaurant not found"))
	}

	replaceCtx, replaceCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer replaceCancel()

	_, err := r.restCollection.ReplaceOne(replaceCtx, filter, bc.CreateRestBson(req.GetRestaurant()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while replacing in mongo: %v", err))
	}

	return &pb.UpdateRestaurantResponse{
		Response: "Successfully updated restaurant",
	}, nil
}

func (r *RestaurantServer) DeleteRestaurant(ctx context.Context, req *pb.DeleteRestaurantRequest) (*pb.DeleteRestaurantResponse, error) {
	log.Info("Deleting restaurant")

	filter := r.restExists(req.GetEmail(), req.GetPhone())
	if filter == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("restaurant not found"))
	}

	_, delErr := r.restCollection.DeleteOne(ctx, filter)
	if delErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error while deleting in mongo: %v", delErr))
	}

	return &pb.DeleteRestaurantResponse{
		Response: "Successfully deleted restaurant",
	}, nil
}

func (r *RestaurantServer) restExists(email string, phone string) bson.M {
	emailFilterRes := bson.M{"email": email}
	phoneFilterRes := bson.M{"phone": phone}
	if r.checkRes(emailFilterRes) != nil {
		return emailFilterRes
	} else if r.checkRes(phoneFilterRes) != nil {
		return phoneFilterRes
	}

	return nil
}

func (r *RestaurantServer) checkRes(filter bson.M) *bc.Restaurant {
	resObj := &bc.Restaurant{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	checkExist := r.restCollection.FindOne(ctx, filter)
	if err := checkExist.Decode(resObj); err == nil {
		return resObj
	}

	return nil
}
