package mesh

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DefaultConfigPath     = "./config/mesh.db.properties.ini"
	ErrorDocumentNotFound = "documentNotFound"
)

type DbConfig struct {
	MeshDbName     string
	BlobDbName     string
	BlobBucketName string
	DbURL          string
}

type Service interface {
	Insert(doc Node) error
	Save(doc Node) error
	FindNodeById(typeName string, id primitive.ObjectID) (Node, error)
	FindNodesFromIDList(typeName string, nodeIdList []primitive.ObjectID) []Node
	FindNodesByTypeName(typeName string) ([]Node, bool)
	FindNodeByProperty(typeName string, property string, value string) (Node, error)
	StoreBlob(file, filename string) (ID primitive.ObjectID, fileSize int64, retErr error)
	RetrieveBlobByName(fileNameInDb string, downloadPath string) error
}

var NodeService Service
