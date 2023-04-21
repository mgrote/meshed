package mesh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mgrote/meshed/configurations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

const (
	DefaultConfigPath     = "../config/mesh.db.properties.ini"
	ErrorDocumentNotFound = "documentNotFound"
)

type NodeService interface {
	Insert(doc Node) error
	Save(doc Node) error
	FindNodeById(className string, id primitive.ObjectID) (Node, error)
	FindNodesFromIDList(className string, nodeIdList []primitive.ObjectID) []Node
	FindNodesByClassName(className string) ([]Node, bool)
	FindNodeByProperty(className string, property string, value string) (Node, error)
	StoreBlob(file, filename string) (ID primitive.ObjectID, fileSize int64, retErr error)
	RetrieveBlobByName(fileNameInDb string, downloadPath string) error
}

var Service NodeService

func InitApiWithConfig(configFileName string) error {
	config, err := configurations.ReadDbConfig(configFileName)
	if err != nil {
		return fmt.Errorf("node service read config: %w", err)
	}
	Service, err = NewNodeServiceMongoDB(config)
	if err != nil {
		return fmt.Errorf("could not init mesh api: %w", err)
	}
	return nil
}

func InitApi() error {
	return InitApiWithConfig(DefaultConfigPath)
}

func NewNodeServiceMongoDB(config *configurations.DbConfig) (NodeService, error) {
	service, err := initMongoDbConnection(config)
	if err != nil {
		return nil, fmt.Errorf("node service connect to database: %w", err)
	}
	return service, nil
}

type NodeServiceMongoDB struct {
	meshDbClient   *mongo.Client
	meshDbName     string
	blobDbName     string
	blobBucketOpts *options.BucketOptions
}

func initMongoDbConnection(dbConfig *configurations.DbConfig) (*NodeServiceMongoDB, error) {
	mongoServerAPI := options.ServerAPI(options.ServerAPIVersion1).
		SetStrict(true).
		SetDeprecationErrors(true)
	opts := options.Client().ApplyURI(dbConfig.DbURL).
		SetAppName("meshdb").
		SetMaxPoolSize(100).
		SetServerAPIOptions(mongoServerAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	meshDbClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	err = meshDbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	log.Println("connected to MongoDB!")

	// TODO Creates this a database, if yes, why?
	//meshDbClient.Database(dbConfig.MeshDbName)

	return &NodeServiceMongoDB{
		meshDbClient:   meshDbClient,
		meshDbName:     dbConfig.MeshDbName,
		blobDbName:     dbConfig.BlobDbName,
		blobBucketOpts: options.GridFSBucket().SetName(dbConfig.BlobBucketName),
	}, nil
}

func (n *NodeServiceMongoDB) Insert(doc Node) error {
	log.Println("inserting", doc.GetID(), doc.GetClass(), doc.GetContent())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := n.meshDbClient.Database(n.meshDbName).Collection(doc.GetClass())
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

func (n *NodeServiceMongoDB) Save(doc Node) error {
	if doc.GetID() == primitive.NilObjectID {
		return n.Insert(doc)
	}
	log.Println("updating", doc.GetID(), doc.GetClass(), doc.GetContent())
	// increase db document version
	version := doc.GetVersion() + 1
	doc.SetVersion(version)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := n.meshDbClient.Database(n.meshDbName).Collection(doc.GetClass())
	filter := bson.M{"_id": doc.GetID()}
	_, err := collection.ReplaceOne(ctx, filter, doc)
	return err
}

func (n *NodeServiceMongoDB) FindNodeById(className string, ID primitive.ObjectID) (Node, error) {
	return findOne(className, bson.M{"_id": ID}, n.meshDbClient, n.meshDbName)
}

func (n *NodeServiceMongoDB) FindNodesFromIDList(className string, nodeIdList []primitive.ObjectID) []Node {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := n.meshDbClient.Database(n.meshDbName).Collection(className)
	findIn := bson.M{"$in": nodeIdList}
	filter := bson.M{"_id": findIn}
	findOptions := options.Find()
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	return mapNodes(cursor, ctx, className, int64(len(nodeIdList)))
}

func (n *NodeServiceMongoDB) FindNodesByClassName(className string) ([]Node, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := n.meshDbClient.Database(n.meshDbName).Collection(className)
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

func (n *NodeServiceMongoDB) FindNodeByProperty(className string, property string, value string) (Node, error) {
	return findOne(className, bson.M{property: value}, n.meshDbClient, n.meshDbName)
}

func (n *NodeServiceMongoDB) StoreBlob(file, filename string) (ID primitive.ObjectID, fileSize int64, retErr error) {
	data, err := os.ReadFile(file)
	fmt.Println("Got databytes", len(data), filename)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}
	bucket, err := gridfs.NewBucket(n.meshDbClient.Database(n.blobDbName), n.blobBucketOpts)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}

	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}
	fileDbId := uploadStream.FileID

	defer func(uploadStream *gridfs.UploadStream) {
		err = uploadStream.Close()
		if err != nil {
			retErr = fmt.Errorf("upload stream could not closed: %w, %v", err, retErr)
		}
	}(uploadStream)

	size, err := uploadStream.Write(data)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}
	log.Println("Write file to DB was successful. Wrote image:", fileDbId, ", File size:", fileSize)

	return fileDbId.(primitive.ObjectID), int64(size), nil
}

func (n *NodeServiceMongoDB) RetrieveBlobByName(fileNameInDb string, downloadPath string) error {
	// return writeToFile(bson.M{"filename": fileNameInDb}, downloadPath)
	collection := *n.blobBucketOpts.Name + ".files"
	fsFiles := n.meshDbClient.Database(n.blobDbName).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Search for blob node.
	filter := bson.M{"filename": fileNameInDb}
	var results bson.M
	if err := fsFiles.FindOne(ctx, filter).Decode(&results); err != nil {
		return fmt.Errorf("file %s not found in blob db: %w", fileNameInDb, err)
	}

	// Load blob data from gridfs bucket.
	bucket, _ := gridfs.NewBucket(n.meshDbClient.Database(n.blobDbName), n.blobBucketOpts)
	var buf bytes.Buffer
	writtenBufBytes, err := bucket.DownloadToStream(results["_id"], &buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File size to download:", writtenBufBytes)
	return os.WriteFile(downloadPath, buf.Bytes(), 0600)
}

func findOne(className string, filter bson.M, client *mongo.Client, dbName string) (Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database(dbName).Collection(className)

	creator := GetTypeConverter(className)
	node := creator()

	err := collection.FindOne(ctx, filter).Decode(*node)
	if err != nil {
		return nil, errors.New(ErrorDocumentNotFound)
	}
	var n Node
	n = *node
	contentConverter := GetContentConverter(className)
	content := contentConverter(n.GetContent().(primitive.D).Map())
	n.SetContent(content)
	return n, nil
}

func mapNodes(cursor *mongo.Cursor, ctx context.Context, className string, initialLength int64) []Node {
	resultList := make([]Node, initialLength)
	for cursor.Next(ctx) {
		creator := GetTypeConverter(className)
		node := creator()
		err := cursor.Decode(*node)
		contentNode := *node
		contentConverter := GetContentConverter(className)
		content := contentConverter(contentNode.GetContent().(primitive.D).Map())
		contentNode.SetContent(content)
		if err != nil {
			log.Fatal(err)
		}
		resultList = append(resultList, contentNode)
	}
	return resultList
}