package dbclient

import (
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func InitFileStreamDbClientWithConfig(configFileName string) {
	gridDbConfig = configurations.ReadDbConfig(configFileName)
	GridMongoClient = InitDbConnection(gridDbConfig)
	bucketOpts = options.GridFSBucket().SetName(gridDbConfig.Bucketname)
}

func UploadFile(file, filename string) (primitive.ObjectID, int64, error) {

	data, err := ioutil.ReadFile(file)
	fmt.Println("Got databytes", len(data), filename)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}
	bucket, err := gridfs.NewBucket(GridMongoClient.Database(gridDbConfig.Dbname), bucketOpts)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}

	uploadStream, err := bucket.OpenUploadStream(
		filename,
	)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}
	fileDbId := uploadStream.FileID

	defer uploadStream.Close()
	fileSize, err := uploadStream.Write(data)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}
	log.Println("Write file to DB was successful. Wrote image:", fileDbId, ", File size:", fileSize)

	return fileDbId.(primitive.ObjectID), int64(fileSize), nil
}

func DownloadFileByName(fileNameInDb string, downloadPath string) {
	writeToFile(bson.M{"filename": fileNameInDb}, downloadPath)
}

func DownloadFileById(id primitive.ObjectID, downloadPath string) {
	writeToFile(bson.M{"_id": id}, downloadPath)
}

func writeToFile(filter bson.M, downloadPath string) {
	fsFiles := GridMongoClient.Database(gridDbConfig.Dbname).Collection(gridDbConfig.Bucketname + ".files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	err := fsFiles.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		log.Fatal("File not found in image db", err)
	}
	log.Println("found requested image to stream", results)
	bucket, _ := gridfs.NewBucket(GridMongoClient.Database(gridDbConfig.Dbname), bucketOpts)
	var buf bytes.Buffer
	writtenBufBytes, err := bucket.DownloadToStream(results["_id"], &buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File size to download:", writtenBufBytes)
	ioutil.WriteFile(downloadPath, buf.Bytes(), 0600)
}
