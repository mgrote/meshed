package categories

import (
	"log"
	"meshed/meshnode"
	"meshed/meshnode/mesh"
	"meshed/meshnode/model"
)

const ClassName = "category"

func CategoryNodeType() mesh.NodeType {
	return meshnode.NewNodeType([]string{}, ClassName)
}

type Category struct {
	Name    string	`json:"name"`
}

func init() {
	log.Println("category init called")
	model.RegisterType(ClassName, func() *mesh.MeshNode {
		node := meshnode.NewNodeWithContent(CategoryNodeType(), Category{})
		return &node
		//return meshnode.NewNodeWithContent(CategoryNodeType(), Category{})
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
