package inmemorymap

import (
	"github.com/mgrote/meshed/mesh"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeServiceMemoryMap struct {
	MemoryMap map[string]map[string]mesh.Node
}

func InitApi() error {
	mesh.NodeService = NewNodeServiceMemoryMap()
	return nil
}

func NewNodeServiceMemoryMap() *NodeServiceMemoryMap {
	return &NodeServiceMemoryMap{
		MemoryMap: make(map[string]map[string]mesh.Node),
	}
}

func (n *NodeServiceMemoryMap) Insert(doc mesh.Node) error {
	//TODO implement me
	panic("implement me")
}

func (n *NodeServiceMemoryMap) Save(doc mesh.Node) error {
	//TODO implement me
	panic("implement me")
}

func (n *NodeServiceMemoryMap) FindNodeById(typeName string, id primitive.ObjectID) (mesh.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (n *NodeServiceMemoryMap) FindNodesFromIDList(typeName string, nodeIdList []primitive.ObjectID) []mesh.Node {
	//TODO implement me
	panic("implement me")
}

func (n *NodeServiceMemoryMap) FindNodesByTypeName(typeName string) ([]mesh.Node, bool) {
	//TODO implement me
	panic("implement me")
}

func (n *NodeServiceMemoryMap) FindNodeByProperty(typeName string, property string, value string) (mesh.Node, error) {
	//TODO implement me
	panic("implement me")
}

func (n *NodeServiceMemoryMap) StoreBlob(file, filename string) (ID primitive.ObjectID, fileSize int64, retErr error) {
	//TODO implement me
	panic("implement me")
}

func (n *NodeServiceMemoryMap) RetrieveBlobByName(fileNameInDb string, downloadPath string) error {
	//TODO implement me
	panic("implement me")
}
