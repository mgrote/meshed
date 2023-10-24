package mesh

const (
	DefaultConfigPath = "./config/mesh.db.properties.ini"
	DocumentNotFound  = "documentNotFound"
	InvalidID         = "invalidID"
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
	FindNodeByID(typeName string, id interface{}) (Node, error)
	FindNodesFromIDList(typeName string, nodeIdList []interface{}) []Node
	FindNodesByTypeName(typeName string) ([]Node, bool)
	FindNodeByProperty(typeName string, property string, value string) (Node, error)
	StoreBlob(file, filename string) (ID interface{}, fileSize int64, retErr error)
	RetrieveBlobByName(fileNameInDb string, downloadPath string) error
}

var NodeService Service
