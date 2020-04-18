package dbclient

import (
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"meshed/configuration/configurations"
	"time"
)



var GridMongoClient *mongo.Client
var gridDbConfig configurations.DbConfig
var bucketOpts *options.BucketOptions
const gridDbConfigFile = "/Users/michaelgrote/etc/go/imagestream.db.properties.ini"

func init() {
	initStreamingDatabase(gridDbConfigFile)
}

func ReinitFileStreamDbClientWithConfig(pathToConfigFile string) {
	initStreamingDatabase(pathToConfigFile)
}

func initStreamingDatabase(pathToConfigFile string) {
	gridDbConfig = configurations.ReadConfig(pathToConfigFile)
	GridMongoClient = InitDbConnection(gridDbConfig)
	bucketOpts = options.GridFSBucket().SetName(gridDbConfig.Bucketname)
}

func UploadFile(file, filename string) {

	data, err := ioutil.ReadFile(file)
	fmt.Println("Got databytes", len(data), filename)
	if err != nil {
		log.Fatal(err)
	}
	bucket, err := gridfs.NewBucket(GridMongoClient.Database(gridDbConfig.Dbname), bucketOpts)
	if err != nil {
		log.Fatal(err)
	}
	uploadStream, err := bucket.OpenUploadStream(
		filename,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Write file to DB was successful. File size: %d M\n", fileSize)
}

func DownloadFile(fileNameInDb string, downloadPath string) {
	fsFiles := GridMongoClient.Database(gridDbConfig.Dbname).Collection(gridDbConfig.Bucketname + ".files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Fatal("tralala-hihi", err)
	}
	// you can print out the result
	fmt.Println("found random image to stream", results)

	bucket, _ := gridfs.NewBucket(GridMongoClient.Database(gridDbConfig.Dbname), bucketOpts)
	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(fileNameInDb, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size to download: %v \n", dStream)
	ioutil.WriteFile(downloadPath, buf.Bytes(), 0600)
}


