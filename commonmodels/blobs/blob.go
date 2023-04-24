package blobs

import (
	"github.com/mgrote/meshed/commonmodels"
	"github.com/mgrote/meshed/mesh"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"time"
)

func BlobNodeType() mesh.NodeType {
	return mesh.NewNodeType([]string{commonmodels.UserType, commonmodels.CategoryType}, commonmodels.BlobType)
}

type Blob struct {
	Title       string             `json:"title"`
	SubTitle    string             `json:"subtitle"`
	Filename    string             `json:"filename"`
	Path        string             `json:"path"`
	Size        int64              `json:"size"`
	ContentType string             `json:"contenttype"`
	Data        primitive.ObjectID `json:"gfsid"`
}

type GridFsDoc struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Length     int64              `json:"length"`
	ChunkSize  int32              `json:"chunkSize"`
	UploadDate time.Time          `json:"uploadDate"`
	Filename   string             `json:"filename"`
}

func init() {
	log.Println("blob init called")
	mesh.RegisterTypeConverter(commonmodels.BlobType, func() *mesh.Node {
		node := mesh.NewNodeWithContent(BlobNodeType(), Blob{})
		return &node
	})
	mesh.RegisterContentConverter(commonmodels.BlobType, GetFromMap)
}

func NewNode(title string, filename string) mesh.Node {
	blob := Blob{
		Title:    title,
		Filename: filename,
	}
	return getNode(blob)
}

func NewGridFSBlobNode(filename string, size int64, contentType string, gfsid primitive.ObjectID) mesh.Node {
	image := Blob{
		Filename:    filename,
		Size:        size,
		ContentType: contentType,
		Data:        gfsid,
	}
	return getNode(image)
}

func NewCheckedNode(title string, pathToFile string) (mesh.Node, bool) {
	if fi, err := os.Stat(pathToFile); err == nil {
		blob := Blob{
			Title:    title,
			Path:     pathToFile,
			Size:     fi.Size(),
			Filename: fi.Name(),
		}
		return getNode(blob), true
	}
	return nil, false
}

func getNode(image Blob) mesh.Node {
	node := mesh.NewNodeWithContent(BlobNodeType(), image)
	node.Save()
	return node
}

func GetBlob(m mesh.Node) Blob {
	blob, ok := m.GetContent().(Blob)
	if !ok {
		log.Fatal("could not convert content from ", m)
	}
	return blob
}

// Exists reports whether the named file or directory exists.
func ReadableFile(file string) bool {
	if _, err := os.Stat(file); err != nil {
		// file not exists or file is not readable
		if os.IsNotExist(err) || os.IsPermission(err) {
			return false
		}
	}
	return true
}

func GetFromMap(docmap map[string]interface{}) interface{} {
	return Blob{
		Title:    docmap["title"].(string),
		SubTitle: docmap["subtitle"].(string),
		Filename: docmap["filename"].(string),
		Path:     docmap["path"].(string),
		Size:     docmap["size"].(int64),
	}
}

func GetGridFsDocFromMap(docmap map[string]interface{}) interface{} {
	return GridFsDoc{
		ID:         primitive.ObjectID{},
		Length:     docmap["length"].(int64),
		ChunkSize:  docmap["chunkSize"].(int32),
		UploadDate: time.Unix(docmap["uploadDate"].(int64), 0),
		Filename:   docmap["filename"].(string),
	}
}
