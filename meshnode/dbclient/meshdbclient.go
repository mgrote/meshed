package dbclient

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"meshed/configuration/configurations"
	"meshed/meshnode/mesh"
	"meshed/meshnode/model"
	"time"
)

var MeshMongoClient *mongo.Client
var meshDbConfig configurations.DbConfig
const meshDbConfigFile = "/Users/michaelgrote/etc/go/mesh.db.properties.ini"

const ErrorDocumentNotFound = "docnotfound"

func init() {
	initMeshDatabase(meshDbConfigFile)
}

func ReinitMeshDbClientWithConfig(pathToConfigFile string) {
	log.Println("Reinit database with config", pathToConfigFile)
	initMeshDatabase(pathToConfigFile)
}

func initMeshDatabase(pathToConfigFile string) {
	meshDbConfig = configurations.ReadConfig(pathToConfigFile)
	MeshMongoClient = InitDbConnection(meshDbConfig)
}

// Insert a new database object
func Insert(doc mesh.MeshNode) {
	log.Println("inserting", doc.GetID(), doc.GetClass(), doc.GetContent())
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := MeshMongoClient.Database(meshDbConfig.Dbname).Collection(doc.GetClass())
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
		collection := MeshMongoClient.Database(meshDbConfig.Dbname).Collection(doc.GetClass())
		filter := bson.M{"_id": doc.GetID()}
		_ , err := collection.ReplaceOne(ctx, filter, doc)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func FindFromIDList(className string, nodeIdList []primitive.ObjectID) []mesh.MeshNode {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := MeshMongoClient.Database(meshDbConfig.Dbname).Collection(className)
	findIn := bson.M{"$in" : nodeIdList}
	filter := bson.M{"_id": findIn}
	findOptions := options.Find()
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	return mapNodes(cursor, ctx, className, int64(len(nodeIdList)))
}

// Delivers a document in a collection (className) by id
func FindById(className string, id primitive.ObjectID) (mesh.MeshNode, error) {
	log.Println("find", className, "by id", id.Hex())
	return findOne(className, bson.M{"_id": id})
}

func FindOneByProperty(className string, property string, value string) (mesh.MeshNode, error) {
	log.Println("find", className, "by", property, ":", value)
	return findOne(className, bson.M{property: value})
}

func findOne(className string, filter bson.M) (mesh.MeshNode, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := MeshMongoClient.Database(meshDbConfig.Dbname).Collection(className)

	creator := model.GetTypeConverter(className)
	node := creator()

	err := collection.FindOne(ctx, filter).Decode(*node)
	if err != nil {
		log.Println(filter, "not found in database")
		return nil, errors.New(ErrorDocumentNotFound)
	}
	var n mesh.MeshNode
	n = *node
	contentConverter := model.GetContentConverter(className)
	content := contentConverter(n.GetContent().(primitive.D).Map())
	n.SetContent(content)
	return n, nil
}

// Delivers all documents in a collection.
// If collection not exists, false is returned.
func FindAllByClassName(className string) ([]mesh.MeshNode, bool) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := MeshMongoClient.Database(meshDbConfig.Dbname).Collection(className)
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

func mapNodes(cursor *mongo.Cursor, ctx context.Context, className string, initialLength int64) []mesh.MeshNode {
	resultList := make([]mesh.MeshNode, initialLength)
	for cursor.Next(ctx) {
		creator := model.GetTypeConverter(className)
		node := creator()
		err := cursor.Decode(*node)
		contentNode := *node
		contentConverter := model.GetContentConverter(className)
		content := contentConverter(contentNode.GetContent().(primitive.D).Map())
		contentNode.SetContent(content)
		if err != nil {
			log.Fatal(err)
		}
		resultList = append(resultList, contentNode)
	}
	return resultList
}


