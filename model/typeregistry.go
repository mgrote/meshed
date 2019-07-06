package model

import (
	"log"
	"meshnode/mesh"
)

var registry = make(map[string]func()mesh.MeshNode)

func RegisterType(className string, creator func()mesh.MeshNode) {
	log.Println("registering", className)
	registry[className] = creator
}

func GetType(className string) func()mesh.MeshNode {
	return registry[className]
}


