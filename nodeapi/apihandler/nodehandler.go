package apihandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"meshed/meshnode/dbclient"
	"meshed/meshnode/model"
	"net/http"
	"reflect"
)

const TypeName = "typename"
const NodeID = "nodeid"

// List existing Entriypoints (existing node types)
// curl localhost:8001/listtypes
func ListNodeTypes(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	model.GetTypes()
	if err := json.NewEncoder(writer).Encode(model.GetTypes()); err != nil {
		log.Fatal("handler.TodoIndex: error while encoding types")
	}
}

// List all nodes of an type
// curl localhost:8001/nodes/category
func ListNodes(writer http.ResponseWriter, request *http.Request) {
	requestVars := mux.Vars(request)
	if typeName, err := requestVars[TypeName]; !err {
		log.Println("Could not find any type from request", typeName)
		writeNotFound(writer)
	} else {
		nodes, success := dbclient.FindAllByClassName(typeName)
		if success {
			writer.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(writer).Encode(nodes); err != nil {
				log.Fatal("Error while encoding respose")
			}
		} else {
			writeNotFound(writer)
		}
	}
}

// Show one node from this type with this ID
// curl localhost:8001/nodes/category/5cfe56a4eb825f1c8ed6e248
func ShowNode(writer http.ResponseWriter, request *http.Request) {
	requestVars := mux.Vars(request)
	if typeName, err := requestVars[TypeName]; !err {
		log.Println("Could not find any type from request")
		writeNotFound(writer)
	} else if nodeid, err := requestVars[NodeID]; !err {
		log.Println("Could not find any id from request")
		writeNotFound(writer)
	} else {
		writer.WriteHeader(http.StatusOK)
		id, err := primitive.ObjectIDFromHex(nodeid)
		if err != nil {
			log.Fatal("Could not get ObjectID from", nodeid)
		}
		node := dbclient.FindById(typeName, id)
		log.Println("got node", node.GetContent(), reflect.TypeOf(node.GetContent()), reflect.TypeOf(node))
		if err := json.NewEncoder(writer).Encode(node); err != nil {
			log.Fatal("Error while encoding respose")
		}
	}
}

func writeNotFound(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(writer).Encode(JSONError{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		log.Fatal("Error while encoding respose")
	}
}
