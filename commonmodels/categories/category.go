package categories

import (
	"github.com/mgrote/meshed/commonmodels"
	"github.com/mgrote/meshed/mesh"
	"log"
)

func CategoryNodeType() mesh.NodeType {
	return mesh.NewNodeType([]string{}, commonmodels.CategoryType)
}

type Category struct {
	Name string `json:"name"`
}

func init() {
	log.Println("category init called")
	mesh.RegisterTypeConverter(commonmodels.CategoryType, func() *mesh.Node {
		node := mesh.NewNodeWithContent(CategoryNodeType(), Category{})
		return &node
	})
	mesh.RegisterContentConverter(commonmodels.CategoryType, GetFromMap)
}

func NewNode(name string) mesh.Node {
	category := Category{
		Name: name,
	}
	node := mesh.NewNodeWithContent(CategoryNodeType(), category)
	node.Save()
	return node
}

func GetCategory(m mesh.Node) Category {
	category, ok := m.GetContent().(Category)
	if !ok {
		log.Fatal("could not convert content from ", m)
	}
	return category
}

func GetFromMap(docmap map[string]interface{}) interface{} {
	return Category{
		Name: docmap["name"].(string),
	}
}
