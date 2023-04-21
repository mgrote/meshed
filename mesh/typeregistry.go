package mesh

import (
	"log"
)

var typeRegistry = make(map[string]func() *Node)
var contentRegistry = make(map[string]func(map[string]interface{}) interface{})

func RegisterTypeConverter(className string, creator func() *Node) {
	log.Println("registering", className)
	typeRegistry[className] = creator
}

func RegisterContentConverter(className string, contentConverter func(map[string]interface{}) interface{}) {
	contentRegistry[className] = contentConverter
}

func GetTypeConverter(className string) func() *Node {
	return typeRegistry[className]
}

func GetContentConverter(className string) func(map[string]interface{}) interface{} {
	return contentRegistry[className]
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
