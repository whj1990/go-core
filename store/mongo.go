package store

import (
	"fmt"

	"github.com/whj1990/go-core/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
)

func NewMongoDB() (*mongo.Client, error) {
	return openMongoDB(
		config.GetNacosConfigData().Mongo.Address,
		config.GetNacosConfigData().Mongo.Username,
		config.GetNacosConfigData().Mongo.Password,
	)
}

func openMongoDB(address, username, password string) (*mongo.Client, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", username, password, address)))
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}
