package meshnode

import (
	"log"
	"meshnode/dbclient"
	"meshnode/mesh"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type node struct {
	ID       	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Version  	uint16             `json:"version"`
	Class    	string             `json:"class"`
	nodeType 	mesh.NodeType           `json:"-"`
	Content  	interface{}        `json:"content"`
	History  	history            `json:"history"`
	Children	map[string][]mesh.MeshNode	`json:"-" bson:"-"`
	ChShadow	map[string][]primitive.ObjectID	`json:"chshadow"`
	Parents		map[string][]mesh.MeshNode	`json:"-" bson:"-"`
	PtShadow	map[string][]primitive.ObjectID	`json:"ptshadow"`
}

type history struct {
	Created  time.Time `json:"created"`
	Changed  time.Time `json:"changed"`
	LastUser string    `json:"lastuser"`
}

func NewNode(t mesh.NodeType) mesh.MeshNode {
	n := node{
		Class:    	t.GetClass(),
		nodeType: 	t,
		History:  	newHistory(),
		Parents: 	make(map[string][]mesh.MeshNode),
		Children: 	make(map[string][]mesh.MeshNode),
		ChShadow:	make(map[string][]primitive.ObjectID),
		PtShadow:	make(map[string][]primitive.ObjectID),
	}
	return &n
}

func NewNodeWithContent(t mesh.NodeType, c interface{}) mesh.MeshNode {
	n := node{
		Class:    	t.GetClass(),
		nodeType: 	t,
		History:  	newHistory(),
		Content:  	c,
		Parents: 	make(map[string][]mesh.MeshNode),
		Children: 	make(map[string][]mesh.MeshNode),
		ChShadow:	make(map[string][]primitive.ObjectID),
		PtShadow:	make(map[string][]primitive.ObjectID),
	}
	return &n
}

// interface MeshNode
func (n *node) GetID() primitive.ObjectID {
	return n.ID
}

func (n *node) SetID(ident interface{}) {
	if id, ok := ident.(primitive.ObjectID); ok {
		n.ID = id
		log.Println("set id to", n.ID, n.GetContent())
	} else {
		log.Fatal("error convert", ident)
	}
}

func (n *node) AddChild(cn mesh.MeshNode) {
	if n.nodeType.IsAccepted(cn.GetClass()) {
		log.Println("add child", cn.GetID(), "to", n.GetID())
		if _, ok := n.Children[cn.GetClass()]; !ok {
			n.Children[cn.GetClass()] = make([]mesh.MeshNode, 0)
			n.ChShadow[cn.GetClass()] = make([]primitive.ObjectID, 0)
		}
		n.Children[cn.GetClass()] = append(n.Children[cn.GetClass()], cn)
		n.ChShadow[cn.GetClass()] = append(n.ChShadow[cn.GetClass()], cn.GetID())
		dbclient.Save(n)
		if !cn.HasParent(n.GetClass(), n.GetID()) {
			log.Println("add also to parent", n.GetID(), "to", cn.GetID())
			cn.AddParent(n)
		}
	}
}

func (n *node) RemoveChild(cn mesh.MeshNode) {
	if _, ok := n.Children[cn.GetClass()]; ok {
		if pos, oki := containsId(n.ChShadow[cn.GetClass()], cn.GetID()); oki {
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
		dbclient.Save(n)
		// check parents
		if cn.HasParent(n.GetClass(), n.GetID()) {
			cn.RemoveParent(n)
		}
	}
}

func (n *node) GetChildren(className string) []mesh.MeshNode {
	// check if parent nodes completely loaded
	if len(n.ChShadow[className]) != len(n.Children[className]) {
		n.Children[className] = checkMissingReferences(n.ChShadow[className], n.Children[className], className)
	}
	return n.Children[className]
}

func (n *node) AddParent(pn mesh.MeshNode) {
	if n.nodeType.IsAccepted(pn.GetClass()) {
		log.Println("add parent", pn.GetID(),"to", n.GetID())
		if _, ok := n.Children[pn.GetClass()]; !ok {
			n.Parents[pn.GetClass()] = make([]mesh.MeshNode, 0)
			n.PtShadow[pn.GetClass()] = make([]primitive.ObjectID, 0)
		}
		n.Parents[pn.GetClass()] = append(n.Children[pn.GetClass()], pn)
		n.PtShadow[pn.GetClass()] = append(n.ChShadow[pn.GetClass()], pn.GetID())
		dbclient.Save(n)
		if !pn.HasChild(n.GetClass(), n.GetID()) {
			log.Println("add also to child", n.GetID(), "to", pn.GetID())
			pn.AddChild(n)
		}
	}
}

func (n *node) RemoveParent(pn mesh.MeshNode) {
	if _, ok := n.Parents[pn.GetClass()]; ok {
		if pos, oki := containsId(n.PtShadow[pn.GetClass()], pn.GetID()); oki {
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
		dbclient.Save(n)
		// check children
		if pn.HasChild(n.GetClass(), n.GetID()) {
			pn.RemoveChild(n)
		}
	}
}

func (n *node) GetParents(className string) []mesh.MeshNode {
	if len(n.ChShadow[className]) != len(n.Children[className]) {
		n.Parents[className] = checkMissingReferences(n.PtShadow[className], n.Parents[className], className)
	}
	return n.Parents[className]
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

func (n *node) GetNodes(className string) []mesh.MeshNode {
	return append(n.Children[className], n.Parents[className] ...)
}

func (n *node) GetClass() string {
	return n.Class
}

func (n *node) GetType() mesh.NodeType {
	return n.nodeType
}

func (n *node) SetType(t mesh.NodeType) {
	n.nodeType = t
	n.Class = t.GetClass()
}

func (n *node) SetContent(c interface{}) {
	n.Content = c
}

func (n *node) GetContent() interface{} {
	return n.Content
}

func (n *node) SetVersion(v uint16) {
	n.Version = v
}

func (n *node) GetVersion() uint16 {
	return n.Version
}

func (n *node) Save() {
	dbclient.Save(n)
}
// interface MeshNode end


func containsId(ids []primitive.ObjectID, id primitive.ObjectID) (int, bool) {
	for it, slid := range ids {
		if slid == id {
			return it, true
		}
	}
	return -1, false
}

func containsNode(ids []mesh.MeshNode, n mesh.MeshNode) (int, bool) {
	for it, slid := range ids {
		if slid == n {
			return it, true
		}
	}
	return -1, false
}

// Checks, if nodes are missing in a list with nodes with a list of given ids.
// Returns a list with already existing and formerly missing nodes.
// The length of the id list and the returned node list have to be equal.
func checkMissingReferences(shadowIds []primitive.ObjectID, refs []mesh.MeshNode, className string) []mesh.MeshNode {
	var nodeIdList = make([]primitive.ObjectID, len(shadowIds))
CheckExistingId:
	for _, cid := range shadowIds {
		for _, child := range refs {
			if child.GetID() == cid {
				// if child is already loaded skip adding id to load list
				// feels hackish
				continue CheckExistingId
			}
		}
		nodeIdList = append(nodeIdList, cid)
	}
	refs = append(refs, dbclient.FindFromIDList(className, nodeIdList)...)
	return refs
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

func removeNodeFromPosition(nodes []mesh.MeshNode, pos int) ([]mesh.MeshNode, bool) {
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