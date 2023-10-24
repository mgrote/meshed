package inmemorymap

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mgrote/meshed/mesh"
	"golang.org/x/exp/maps"
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
	if doc.GetID() == nil {
		doc.SetID(uuid.New().String())
	}
	id, ok := doc.GetID().(string)
	if !ok {
		return fmt.Errorf(mesh.InvalidID)
	}
	n.MemoryMap[doc.GetTypeName()][id] = doc
	return nil
}

func (n *NodeServiceMemoryMap) Save(doc mesh.Node) error {
	return n.Insert(doc)
}

func (n *NodeServiceMemoryMap) FindNodeByID(typeName string, id interface{}) (mesh.Node, error) {
	ID, ok := id.(string)
	if !ok {
		return nil, fmt.Errorf(mesh.InvalidID)
	}
	node, ok := n.MemoryMap[typeName][ID]
	if !ok {
		return nil, fmt.Errorf(mesh.DocumentNotFound)
	}
	return node, nil
}

func (n *NodeServiceMemoryMap) FindNodesFromIDList(typeName string, nodeIdList []interface{}) []mesh.Node {
	var nodes []mesh.Node
	for _, id := range nodeIdList {
		node, err := n.FindNodeByID(typeName, id)
		if err == nil {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (n *NodeServiceMemoryMap) FindNodesByTypeName(typeName string) ([]mesh.Node, bool) {
	nodeIds, ok := n.MemoryMap[typeName]
	return maps.Values(nodeIds), ok
}

func (n *NodeServiceMemoryMap) FindNodeByProperty(typeName string, property string, value string) (mesh.Node, error) {
	for _, node := range n.MemoryMap[typeName] {
		j, err := json.Marshal(node.GetContent())
		if err != nil {
			return nil, err
		}
		var content map[string]interface{}
		err = json.Unmarshal(j, &content)
		if content[property] == value {
			return node, nil
		}
	}
	return nil, fmt.Errorf(mesh.DocumentNotFound)
}

func (n *NodeServiceMemoryMap) StoreBlob(file, filename string) (ID interface{}, fileSize int64, retErr error) {
	//TODO implement me
	panic("implement me")
}

func (n *NodeServiceMemoryMap) RetrieveBlobByName(fileNameInDb string, downloadPath string) error {
	//TODO implement me
	panic("implement me")
}
