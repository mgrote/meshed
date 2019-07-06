package categories

import (
	"log"
	"meshnode/mesh"
	"meshnode/meshnode"
	"meshnode/model"
)

func CategoryNodeType() mesh.NodeType {
	return meshnode.NewNodeType([]string{}, "category")
}

type Category struct {
	Name    string	`json:"name"`
}

func init() {
	model.RegisterType("category", func() mesh.MeshNode {
		return meshnode.NewNodeWithContent(CategoryNodeType(), Category{})
	})
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
	category, ok := m.GetContent().(Category)
	if !ok {
		log.Fatal("could not convert content from ", m)
	}
	return category
}