package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectDB() error {
	mongoURI := os.Getenv("DB_URL")

	if mongoURI == "" {
		return fmt.Errorf("error trying to get DB_URL from .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not ping to MongoDB: %v", err)
	}
	MongoClient = client
	MongoDB = client.Database("pureweb2")

	log.Println("Connected to MongoDB")
	return nil
}
