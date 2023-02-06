package mongodb

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client

func Init(ctx context.Context, uri string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	return client
}

func Ready(ctx context.Context) error {

	if MongoClient == nil {
		return errors.New("MongoDB not initialized")
	}

	if err := MongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	return nil
}

func Close(ctx context.Context) {
	if MongoClient == nil {
		log.Info("MongoDB already closed!")
		return
	}

	if err := MongoClient.Disconnect(ctx); err != nil {
		log.Error(err)
		return
	}
	log.Info("MongoDB disconnected.")
}
