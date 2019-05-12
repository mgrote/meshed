package dbclient

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"mongodbtest/dbo"
	"mongodbtest/dbtodo"
	"time"
)

var MongoClient *mongo.Client

const dbname = "meshdb"

func init() {
	log.Println("mesh connecting to database", dbname)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://meshdbuser:isAnberf$QuinAfkoarc6@localhost:27017"))
	err = MongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}
