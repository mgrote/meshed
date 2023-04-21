package mongodb

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mgrote/meshed/configurations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var GridMongoClient *mongo.Client
var gridDbConfig *configurations.DbConfig
var bucketOpts *options.BucketOptions

func InitFileStreamDbClientWithConfig(configFileName string) error {
	var err error
	gridDbConfig, err = configurations.ReadDbConfig(configFileName)
	if err != nil {
		return fmt.Errorf("init file stream database client configuration: %w", err)
	}
	GridMongoClient, err = InitDbConnection(gridDbConfig)
	if err != nil {
		return fmt.Errorf("init file stream database client: %w", err)
	}
	bucketOpts = options.GridFSBucket().SetName(gridDbConfig.BlobBucketName)
	return nil
}

func UploadFile(file, filename string) (ID primitive.ObjectID, fileSize int64, retErr error) {

	data, err := os.ReadFile(file)
	fmt.Println("Got databytes", len(data), filename)
	if err != nil {
		return primitive.NilObjectID, 0, err
	}
	bucket, err := gridfs.NewBucket(GridMongoClient.Database(gridDbConfig.BlobDbName), bucketOpts)
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

func DownloadFileByName(fileNameInDb string, downloadPath string) error {
	return writeToFile(bson.M{"filename": fileNameInDb}, downloadPath)
}

func DownloadFileById(id primitive.ObjectID, downloadPath string) error {
	return writeToFile(bson.M{"_id": id}, downloadPath)
}

func writeToFile(filter bson.M, downloadPath string) error {
	fsFiles := GridMongoClient.Database(gridDbConfig.BlobDbName).Collection(gridDbConfig.BlobBucketName + ".files")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results bson.M
	err := fsFiles.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		log.Fatal("File not found in image db", err)
	}
	log.Println("found requested image to stream", results)
	bucket, _ := gridfs.NewBucket(GridMongoClient.Database(gridDbConfig.BlobDbName), bucketOpts)
	var buf bytes.Buffer
	writtenBufBytes, err := bucket.DownloadToStream(results["_id"], &buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File size to download:", writtenBufBytes)
	return os.WriteFile(downloadPath, buf.Bytes(), 0600)
}
