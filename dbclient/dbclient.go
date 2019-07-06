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
	"meshnode/mesh"
	"meshnode/model"
	"time"
)

var MongoClient *mongo.Client

const dbname = "meshdb"

func init() {
	log.Println("mesh connecting to database", dbname)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://meshdbuser:misAnberf$QuinAfkoarc6@localhost:27017"))
	err = MongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

// Insert a new database object
func Insert(doc mesh.MeshNode) {
	log.Println("inserting", doc.GetID(), doc.GetClass(), doc.GetContent())
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := MongoClient.Database(dbname).Collection(doc.GetClass())
	// increase db document version
	version := doc.GetVersion() + 1
	doc.SetVersion(version)
	result, err := collection.InsertOne(ctx, doc)
	if result != nil {
		// set db id to reference if not exists
		doc.SetID(result.InsertedID.(primitive.ObjectID))
	} else {
		log.Fatal("could not write to database", err)
	}
	log.Println("saved", doc.GetID(), doc.GetVersion(), doc.GetClass())
}

// Save saves the dbo as it is, there is no merge with any existing document,
// existing documents will be overwritten with this doc.
// If the doc was never written to the database, it will be created with a new id.
func Save(doc mesh.MeshNode) {
	if doc.GetID() == primitive.NilObjectID {
		Insert(doc)
	} else {
		log.Println("updating", doc.GetID(), doc.GetClass(), doc.GetContent())
		// increase db document version
		version := doc.GetVersion() + 1
		doc.SetVersion(version)
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		collection := MongoClient.Database(dbname).Collection(doc.GetClass())
		filter := bson.M{"_id": doc.GetID()}
		_ , err := collection.ReplaceOne(ctx, filter, doc)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func FindFromIDList(className string, nodeIdList []primitive.ObjectID) []mesh.MeshNode {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := MongoClient.Database(dbname).Collection(className)
	findIn := bson.M{"$in" : nodeIdList}
	filter := bson.M{"_id": findIn}
	findOptions := options.Find()
	creator := model.GetType(className)
	resultList := make([]mesh.MeshNode, len(nodeIdList))
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(ctx) {
		node := creator()
		err := cursor.Decode(&node)
		if err != nil {
			log.Fatal(err)
		}
		resultList = append(resultList, node)
	}
	return resultList
}

func FindById(className string, id primitive.ObjectID) mesh.MeshNode {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := MongoClient.Database(dbname).Collection(className)
	filter := bson.M{"_id": id}
	creator := model.GetType(className)
	node := creator()
	err := collection.FindOne(ctx, filter).Decode(&node)
	if err != nil {
		log.Fatal(err)
	}
	return node
}
