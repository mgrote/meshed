package dbclient

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"meshed/meshnode/configurations"
	"time"
)

func InitDbConnection(dbConfig configurations.DbConfig) *mongo.Client {
	log.Println("mesh connecting to database", dbConfig.Dbname)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	opts := options.Client()
	opts.ApplyURI(dbConfig.Dburl)
	opts.SetMaxPoolSize(100)
	MeshMongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(dbConfig.Dburl))
	// dont ping if something went wrong
	if MeshMongoClient != nil && err == nil{
		err = MeshMongoClient.Ping(ctx, readpref.Primary())
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	MeshMongoClient.Database(dbConfig.Dbname)
	return MeshMongoClient
}
