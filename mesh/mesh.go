package mesh

import "go.mongodb.org/mongo-driver/bson/primitive"

type NodeType interface {
	GetName() 		string
	AcceptChild() 	[]string
	IsAccepted(string) bool
}

type Node interface {
	GetId() primitive.ObjectID
	AddChild(*node)
	RemoveChild(*node)
	AddParent(*node)
	RemoveParent(*node)
	HasChild(string, primitive.ObjectID) bool
	HasParent(string, primitive.ObjectID) bool
	RemoveAllNodes(string)
	GetNodes(string) []node
	GetClass() string
	GetType() NodeType
	SetType(NodeType)
	SetContent(interface{})
	GetContent() interface{}
}