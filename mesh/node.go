package mesh

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type node struct {
	ID       	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Version  	uint16             `json:"version"`
	Class    	string             `json:"class"`
	nodeType 	NodeType           `json:"-"`
	Content  	interface{}        `json:"content"`
	History  	history            `json:"history"`
	Children	map[string][]*node	`json:"_"`
	ChShadow	map[string][]primitive.ObjectID	`json:"chshadow"`
	Parents		map[string][]*node	`json:"_"`
	PtShadow	map[string][]primitive.ObjectID	`json:"ptshadow"`
}

type history struct {
	Created  time.Time `json:"created"`
	Changed  time.Time `json:"changed"`
	LastUser string    `json:"lastuser"`
}

func NewNode(t NodeType) *node {
	n := node{
		Class:    t.GetName(),
		nodeType: t,
		History:  newHistory(),
	}
	return &n
}

func NewNodeWithContent(t NodeType, c interface{}) *node {
	n := node{
		Class:    t.GetName(),
		nodeType: t,
		History:  newHistory(),
		Content:  c,
	}
	return &n
}

// interface Node
func (n *node) GetId() primitive.ObjectID {
	return n.ID
}

func (n *node) AddChild(cn *node) {
	if n.nodeType.IsAccepted(cn.GetClass()) {
		if _, ok := n.Children[cn.GetClass()]; !ok {
			n.Children[cn.GetClass()] = make([]*node, 1)
			n.ChShadow[cn.GetClass()] = make([]primitive.ObjectID, 1)
		}
		n.Children[cn.GetClass()] = append(n.Children[cn.GetClass()], cn)
		n.ChShadow[cn.GetClass()] = append(n.ChShadow[cn.GetClass()], cn.GetId())
		if cn.HasParent(n.GetClass(), n.GetId()) {
			cn.AddParent(n)
		}
	}
}

func (n *node) RemoveChild(cn *node) {
	if _, ok := n.Children[cn.GetClass()]; ok {
		if pos, oki := containsId(n.ChShadow[cn.GetClass()], cn.GetId()); oki {
			if isl, ok := removeIdFromPosition(n.ChShadow[cn.GetClass()], pos); ok {
				n.ChShadow[cn.GetClass()] = isl
			} else {
				delete(n.ChShadow, cn.GetClass())
			}
		}
		if pos, oki := containsNode(n.Children[cn.GetClass()], cn); oki {
			if nsl, ok := removeNodeFromPosition(n.Children[cn.GetClass()], pos); ok {
				n.Children[cn.GetClass()] = nsl
			} else {
				delete(n.Children, cn.GetClass())
			}
		}
		// check parents
		if cn.HasParent(n.GetClass(), n.GetId()) {
			cn.RemoveParent(n)
		}
	}
}

func (n *node) AddParent(pn *node) {
	if n.nodeType.IsAccepted(pn.GetClass()) {
		if _, ok := n.Children[pn.GetClass()]; !ok {
			n.Parents[pn.GetClass()] = make([]*node, 1)
			n.PtShadow[pn.GetClass()] = make([]primitive.ObjectID, 1)
		}
		n.Parents[pn.GetClass()] = append(n.Children[pn.GetClass()], pn)
		n.PtShadow[pn.GetClass()] = append(n.ChShadow[pn.GetClass()], pn.GetId())
		if pn.HasChild(n.GetClass(), n.GetId()) {
			pn.AddChild(n)
		}
	}
}

func (n *node) RemoveParent(pn *node) {
	if _, ok := n.Parents[pn.GetClass()]; ok {
		if pos, oki := containsId(n.PtShadow[pn.GetClass()], pn.GetId()); oki {
			if isl, ok := removeIdFromPosition(n.PtShadow[pn.GetClass()], pos); ok {
				n.PtShadow[pn.GetClass()] = isl
			} else {
				delete(n.PtShadow, pn.GetClass())
			}
		}
		if pos, oki := containsNode(n.Parents[pn.GetClass()], pn); oki {
			if nsl, ok := removeNodeFromPosition(n.Parents[pn.GetClass()], pos); ok {
				n.Parents[pn.GetClass()] = nsl
			} else {
				delete(n.Parents, pn.GetClass())
			}
		}
		// check children
		if pn.HasChild(n.GetClass(), n.GetId()) {
			pn.RemoveChild(n)
		}
	}
}

func (n *node) HasChild(className string, id primitive.ObjectID) bool {
	if ids, ok := n.ChShadow[className]; ok {
		_, oki := containsId(ids, id)
		return oki
	}
	return false
}

func (n *node) HasParent(className string, id primitive.ObjectID) bool {
	if ids, ok := n.PtShadow[className]; ok {
		_, oki := containsId(ids, id)
		return oki
	}
	return false
}

func (n *node) RemoveAllNodes(className string) {
	delete(n.Children, className)
	delete(n.Parents, className)
}

func (n *node) GetNodes(className string) []*node {
	return n.Children[className]
}

func (n *node) GetClass() string {
	return n.Class
}

func (n *node) GetType() NodeType {
	return n.nodeType
}

func (n *node) SetType(t NodeType) {
	n.nodeType = t
	n.Class = t.GetName()
}

func (n *node) SetContent(c interface{}) {
	n.Content = c
}

func (n *node) GetContent() interface{} {
	return n.Content
}
// interface Node end


func containsId(ids []primitive.ObjectID, id primitive.ObjectID) (int, bool) {
	for it, slid := range ids {
		if slid == id {
			return it, true
		}
	}
	return -1, false
}

func containsNode(ids []*node, n *node) (int, bool) {
	for it, slid := range ids {
		if slid == n {
			return it, true
		}
	}
	return -1, false
}

// https://github.com/golang/go/wiki/SliceTricks
func removeIdFromPosition(ids []primitive.ObjectID, pos int) ([]primitive.ObjectID, bool) {
	if len(ids) == 1 {
		return nil, false
	}
	copy(ids[pos:], ids[pos+1:])
	ids[len(ids) - 1] = primitive.NilObjectID // prevents memory leak?
	ids = ids[:len(ids) - 1]
	return ids, true
}

func removeNodeFromPosition(nodes []*node, pos int) ([]*node, bool) {
	if len(nodes) == 1 {
		return nil, false
	}
	copy(nodes[pos:], nodes[pos+1:])
	nodes[len(nodes) - 1] = nil // prevents memory leak
	nodes = nodes[:len(nodes) - 1]
	return nodes, true
}

func newHistory() history {
	h := history{
		Created: time.Now(),
		Changed: time.Now(),
	}
	return h
}