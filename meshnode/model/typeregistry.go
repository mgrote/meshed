package model

import (
	"log"
	"meshed/meshnode/mesh"
)

var registry = make(map[string]func()*mesh.MeshNode)

func RegisterType(className string, creator func()*mesh.MeshNode) {
	log.Println("registering", className)
	registry[className] = creator
}

func GetType(className string) func()*mesh.MeshNode {
	return registry[className]
}

func GetTypes() []string {
	keys := make([]string, len(registry))
	i := 0
	for key := range registry {
		keys[i] = key
		i++
	}
	return keys
}


