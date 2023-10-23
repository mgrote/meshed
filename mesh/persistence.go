package mesh

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
	FindNodeById(typeName string, id string) (Node, error)
	FindNodesFromIDList(typeName string, nodeIdList []string) []Node
	FindNodesByTypeName(typeName string) ([]Node, bool)
	FindNodeByProperty(typeName string, property string, value string) (Node, error)
	StoreBlob(file, filename string) (ID string, fileSize int64, retErr error)
	RetrieveBlobByName(fileNameInDb string, downloadPath string) error
}

var NodeService Service
