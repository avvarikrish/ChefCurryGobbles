package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCcgobblesServerConfig(t *testing.T) {
	got := NewConfig("./basic.yml")
	want := CcgobblesServerConfig{
		Mongo: mongo{
			MongoServer: "mongodb://mongodb:27017",
			Database:    "chefcurrygobbles",
			Collections: mongoCollections{
				Users:       "Users",
				Restaurants: "Restaurants",
				Orders:      "Orders",
			},
		},
		Server: server{
			Network: "tcp",
			Address: ":50051",
		},
	}

	assert.New(t).Equal(want, got)
}
