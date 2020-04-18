package model

import (
	"log"
	"meshed/meshnode/mesh"
)

var typeregistry = make(map[string]func() *mesh.MeshNode)
var contentregistry = make(map[string]func(map[string]interface{}) interface{})

func RegisterTypeConverter(className string, creator func() *mesh.MeshNode) {
	log.Println("registering", className)
	typeregistry[className] = creator
}

func RegisterContentConverter(className string, contentConverter func(map[string]interface{}) interface{}) {
	contentregistry[className] = contentConverter
}

func GetTypeConverter(className string) func() *mesh.MeshNode {
	return typeregistry[className]
}

func GetContentConverter(className string) func(map[string]interface{}) interface{} {
	return contentregistry[className]
}

func GetTypes() []string {
	keys := make([]string, len(typeregistry))
	i := 0
	for key := range typeregistry {
		keys[i] = key
		i++
	}
	return keys
}
