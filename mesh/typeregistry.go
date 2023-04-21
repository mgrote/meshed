package mesh

import (
	"log"
)

var typeRegistry = make(map[string]func() *Node)
var contentRegistry = make(map[string]func(map[string]interface{}) interface{})

func RegisterTypeConverter(typeName string, creator func() *Node) {
	log.Println("registering", typeName)
	typeRegistry[typeName] = creator
}

func RegisterContentConverter(typeName string, contentConverter func(map[string]interface{}) interface{}) {
	contentRegistry[typeName] = contentConverter
}

func GetTypeConverter(typeName string) func() *Node {
	return typeRegistry[typeName]
}

func GetContentConverter(typeName string) func(map[string]interface{}) interface{} {
	return contentRegistry[typeName]
}

func GetTypes() []string {
	keys := make([]string, len(typeRegistry))
	i := 0
	for key := range typeRegistry {
		keys[i] = key
		i++
	}
	return keys
}
