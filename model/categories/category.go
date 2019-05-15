package categories

import (
	"log"
	"meshnode/mesh"
	"meshnode/meshnode"
)

func CategoryNodeType() mesh.NodeType {
	return meshnode.NodeType{
		[]string{},
		"category"}
}

type Category struct {
	Name    string
}

func NewNode(name string) mesh.MeshNode {
	category := Category{
		Name:    name,
	}
	node := meshnode.NewNodeWithContent(CategoryNodeType(), category)
	node.Save()
	return node
}

func GetCategory(m mesh.MeshNode) Category {
	user, ok := m.GetContent().(Category)
	if !ok {
		log.Fatal("could not convert content from ", m)
	}
	return user
}