package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Ctx  *context.Context
	Host string
	Port int
}

func NewMongo(cfg Config) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return client, err
	}

	err = client.Connect(*cfg.Ctx)

	return client, err
}
