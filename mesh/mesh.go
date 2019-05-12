package mesh

import "go.mongodb.org/mongo-driver/bson/primitive"

type NodeType interface {
	GetClass() 		string
	AcceptTypes() 	[]string
	IsAccepted(string) bool
}

type MeshNode interface {
	GetID() primitive.ObjectID
	SetID(id interface{})
	AddChild(MeshNode)
	RemoveChild(MeshNode)
	AddParent(MeshNode)
	RemoveParent(MeshNode)
	HasChild(string, primitive.ObjectID) bool
	HasParent(string, primitive.ObjectID) bool
	RemoveAllNodes(string)
	GetNodes(string) []MeshNode
	GetClass() string
	GetType() NodeType
	SetType(NodeType)
	SetContent(interface{})
	GetContent() interface{}
	GetVersion() uint16
	SetVersion(uint16)
}