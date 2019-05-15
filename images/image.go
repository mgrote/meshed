package images

import (
	"log"
	"meshnode/mesh"
	"meshnode/meshnode"
	"os"
)

func ImageNodeType() mesh.NodeType {
	return meshnode.NodeType{
		[]string{"image", "category"},
		"image"}
}

type Image struct {
	Title    string
	SubTitle string
	Filename string
	Path     string
	Size     int64
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
	user, ok := m.GetContent().(Image)
	if !ok {
		log.Fatal("could not convert content from ", m)
	}
	return user
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