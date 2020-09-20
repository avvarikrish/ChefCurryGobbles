package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type OrdersServerConfig struct {
	Mongo  mongo
	Server server
}

type server struct {
	Network string
	Address string
}

type mongo struct {
	MongoServer string
	Database    string
	Collections mongoCollections
}

type mongoCollections struct {
	Users       string
	Restaurants string
	Orders      string
}

// NewConfig initializes a new config object
func NewConfig(file string) OrdersServerConfig {
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		log.Errorf("error reading config: %v", err)
	}
	return OrdersServerConfig{
		Mongo:  newMongo(v.Sub("mongo")),
		Server: newServer(v.Sub("server")),
	}
}

func newServer(v *viper.Viper) server {
	return server{
		Network: v.GetString("network"),
		Address: v.GetString("address"),
	}
}

func newMongo(v *viper.Viper) mongo {
	fmt.Println(v.GetString("database"))
	m := v.GetStringMapString("collections")
	return mongo{
		MongoServer: v.GetString("mongo_server"),
		Database:    v.GetString("database"),
		Collections: mongoCollections{
			Users:       m["users"],
			Restaurants: m["restaurants"],
			Orders:      m["orders"],
		},
	}
}
