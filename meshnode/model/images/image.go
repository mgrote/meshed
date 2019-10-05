package images

import (
	"log"
	"meshed/meshnode"
	"meshed/meshnode/mesh"
	"meshed/meshnode/model"
	"meshed/meshnode/model/categories"
	"meshed/meshnode/model/users"
	"os"
)

const ClassName = "image"

func ImageNodeType() mesh.NodeType {
	return meshnode.NewNodeType([]string{categories.ClassName, users.ClassName}, ClassName)
}

type Image struct {
	Title    string	`json:"title"`
	SubTitle string	`json:"subtitle"`
	Filename string	`json:"filename"`
	Path     string	`json:"path"`
	Size     int64	`json:"size"`
}

func init() {
	log.Println("image init called")
	model.RegisterType("image", func() *mesh.MeshNode {
		node := meshnode.NewNodeWithContent(ImageNodeType(), Image{})
		return &node
	})
}

func NewNode(title string, filename string) mesh.MeshNode {
	image := Image{
		Title:    title,
		Filename: filename,
	}
	return getNode(image)
}

func NewCheckedNode(title string, pathToFile string) (mesh.MeshNode, bool) {
	if fi, err := os.Stat(pathToFile); err == nil {
		image := Image{
			Title:    title,
			Path:     pathToFile,
			Size:     fi.Size(),
			Filename: fi.Name(),
		}
		return getNode(image), true
	}
	return nil, false
}

func getNode(image Image) mesh.MeshNode {
	node := meshnode.NewNodeWithContent(ImageNodeType(), image)
	node.Save()
	return node
}

func GetImage(m mesh.MeshNode) Image {
	image, ok := m.GetContent().(Image)
	if !ok {
		log.Fatal("could not convert content from ", m)
	}
	return image
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