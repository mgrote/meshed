package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/mgrote/meshed/configurations"
	"github.com/mgrote/meshed/mesh"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var MeshMongoClient *mongo.Client
var meshDbConfig *configurations.DbConfig

const ErrorDocumentNotFound = "docnotfound"

func InitDbConnection(dbConfig *configurations.DbConfig) (*mongo.Client, error) {
	log.Println("mesh connecting to database", dbConfig.DbName)

	var err error
	opts := options.Client()
	opts.ApplyURI(dbConfig.DbURL)
	opts.SetMaxPoolSize(100)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	MeshMongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(dbConfig.DbURL))
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	err = MeshMongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	log.Println("connected to MongoDB!")
	MeshMongoClient.Database(dbConfig.DbName)
	return MeshMongoClient, nil
}

func InitMeshDbClientWithConfig(configFileName string) (*mongo.Client, error) {
	dbConfig, err := configurations.ReadDbConfig(configFileName)
	if err != nil {

	}
	return InitDbConnection(dbConfig)
}

// TODO use references !!!!!!!!!

// Insert a new database object
func Insert(doc mesh.Node) error {
	log.Println("inserting", doc.GetID(), doc.GetClass(), doc.GetContent())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := MeshMongoClient.Database(meshDbConfig.DbName).Collection(doc.GetClass())
	// increase db document version
	version := doc.GetVersion() + 1
	doc.SetVersion(version)
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("could not insert document %v into database %w", doc.GetClass(), err)
	}
	// set db id to reference if not exists
	doc.SetID(result.InsertedID.(primitive.ObjectID))
	log.Println("saved", doc.GetID(), doc.GetVersion(), doc.GetClass())
	return nil
}

// Save saves the dbo as it is, there is no merge with any existing document,
// existing documents will be overwritten with this doc.
// If the doc was never written to the database, it will be created with a new id.
func Save(doc mesh.Node) error {
	if doc.GetID() == primitive.NilObjectID {
		return Insert(doc)
	}
	log.Println("updating", doc.GetID(), doc.GetClass(), doc.GetContent())
	// increase db document version
	version := doc.GetVersion() + 1
	doc.SetVersion(version)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := MeshMongoClient.Database(meshDbConfig.DbName).Collection(doc.GetClass())
	filter := bson.M{"_id": doc.GetID()}
	_, err := collection.ReplaceOne(ctx, filter, doc)
	return err
}

func FindFromIDList(className string, nodeIdList []primitive.ObjectID) []mesh.Node {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := MeshMongoClient.Database(meshDbConfig.DbName).Collection(className)
	findIn := bson.M{"$in": nodeIdList}
	filter := bson.M{"_id": findIn}
	findOptions := options.Find()
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	return mapNodes(cursor, ctx, className, int64(len(nodeIdList)))
}

// FindById delivers a document in a collection (className) by id
func FindById(className string, ID primitive.ObjectID) (mesh.Node, error) {
	log.Println("find", className, "by ID", ID.Hex())
	return findOne(className, bson.M{"_id": ID})
}

func FindOneByProperty(className string, property string, value string) (mesh.Node, error) {
	log.Println("find", className, "by", property, ":", value)
	return findOne(className, bson.M{property: value})
}

// FindAllByClassName delivers all documents in a collection.
// If collection not exists, false is returned.
func FindAllByClassName(className string) ([]mesh.Node, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := MeshMongoClient.Database(meshDbConfig.DbName).Collection(className)
	numDocs, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if numDocs == 0 {
		return nil, false
	}
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	return mapNodes(cursor, ctx, className, numDocs), true
}

func findOne(className string, filter bson.M) (mesh.Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := MeshMongoClient.Database(meshDbConfig.DbName).Collection(className)

	creator := mesh.GetTypeConverter(className)
	node := creator()

	err := collection.FindOne(ctx, filter).Decode(*node)
	if err != nil {
		log.Println(filter, "not found in database")
		return nil, errors.New(ErrorDocumentNotFound)
	}
	var n mesh.Node
	n = *node
	contentConverter := mesh.GetContentConverter(className)
	content := contentConverter(n.GetContent().(primitive.D).Map())
	n.SetContent(content)
	return n, nil
}

func mapNodes(cursor *mongo.Cursor, ctx context.Context, className string, initialLength int64) []mesh.Node {
	resultList := make([]mesh.Node, initialLength)
	for cursor.Next(ctx) {
		creator := mesh.GetTypeConverter(className)
		node := creator()
		err := cursor.Decode(*node)
		contentNode := *node
		contentConverter := mesh.GetContentConverter(className)
		content := contentConverter(contentNode.GetContent().(primitive.D).Map())
		contentNode.SetContent(content)
		if err != nil {
			log.Fatal(err)
		}
		resultList = append(resultList, contentNode)
	}
	return resultList
}
