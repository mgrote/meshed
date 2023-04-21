package mesh

import (
	"fmt"
	"github.com/mgrote/meshed/mesh/mongodb"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeType interface {
	Class() string
	AcceptTypes() []string
	Accepting(string) bool
}

type Node interface {
	GetID() primitive.ObjectID
	SetID(id interface{})
	AddChild(Node) error
	RemoveChild(Node) error
	GetChildren(string) []Node
	AddParent(Node) error
	RemoveParent(Node) error
	GetParents(string) []Node
	HasChild(string, primitive.ObjectID) bool
	HasParent(string, primitive.ObjectID) bool
	RemoveAllNodes(string)
	GetNodes(string) []Node
	GetClass() string
	Type() NodeType
	SetType(NodeType)
	SetContent(interface{})
	GetContent() interface{}
	GetVersion() uint16
	SetVersion(uint16)
	Save() error
	SaveContent(interface{}) error
}

type node struct {
	ID       primitive.ObjectID              `json:"id" bson:"_id,omitempty"`
	Version  uint16                          `json:"version"`
	Class    string                          `json:"class"`
	nodeType NodeType                        //`json:"-"`
	Content  interface{}                     `json:"content"`
	History  history                         `json:"history"`
	Children map[string][]Node               `json:"-" bson:"-"`
	ChShadow map[string][]primitive.ObjectID `json:"chshadow"`
	Parents  map[string][]Node               `json:"-" bson:"-"`
	PtShadow map[string][]primitive.ObjectID `json:"ptshadow"`
}

type history struct {
	Created  time.Time `json:"created"`
	Changed  time.Time `json:"changed"`
	LastUser string    `json:"lastuser"`
}

func NewNode(t NodeType) Node {
	n := node{
		Class:    t.Class(),
		nodeType: t,
		History:  newHistory(),
		Parents:  make(map[string][]Node),
		Children: make(map[string][]Node),
		ChShadow: make(map[string][]primitive.ObjectID),
		PtShadow: make(map[string][]primitive.ObjectID),
	}
	return &n
}

func NewNodeWithContent(t NodeType, c interface{}) Node {
	n := node{
		Class:    t.Class(),
		nodeType: t,
		History:  newHistory(),
		Content:  c,
		Parents:  make(map[string][]Node),
		Children: make(map[string][]Node),
		ChShadow: make(map[string][]primitive.ObjectID),
		PtShadow: make(map[string][]primitive.ObjectID),
	}
	return &n
}

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

func (n *node) AddChild(cn Node) error {
	if !n.nodeType.Accepting(cn.GetClass()) {
		return fmt.Errorf("could not add child, parent type %s is not accepted by child type %s",
			n.nodeType.Class(), cn.GetClass())
	}

	log.Println("add child", cn.GetID(), "to", n.GetID())
	if _, ok := n.Children[cn.GetClass()]; !ok {
		n.Children[cn.GetClass()] = make([]Node, 0)
		n.ChShadow[cn.GetClass()] = make([]primitive.ObjectID, 0)
	}
	n.Children[cn.GetClass()] = append(n.Children[cn.GetClass()], cn)
	n.ChShadow[cn.GetClass()] = append(n.ChShadow[cn.GetClass()], cn.GetID())
	if err := Service.Save(n); err != nil {
		return fmt.Errorf("%s, %s could not add child %s, %s: %w",
			n.GetClass(), n.GetID().String(), cn.GetClass(), cn.GetID().String(), err)
	}

	if !cn.HasParent(n.GetClass(), n.GetID()) {
		if err := cn.AddParent(n); err != nil {
			return fmt.Errorf("could not add parent %s, %s to child %s, %s: %w",
				n.GetClass(), n.GetID().String(), cn.GetClass(), cn.GetID().String(), err)
		}
	}
	return nil
}

func (n *node) RemoveChild(cn Node) error {
	if _, ok := n.Children[cn.GetClass()]; !ok {
		return nil
	}

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
	if err := mongodb.Save(n); err != nil {
		return fmt.Errorf("%s, %s could not remove child %s, %s: %w",
			n.GetClass(), n.GetID().String(), cn.GetClass(), cn.GetID().String(), err)
	}

	if cn.HasParent(n.GetClass(), n.GetID()) {
		if err := cn.RemoveParent(n); err != nil {
			return fmt.Errorf("could not remove parent %s, %s from child %s, %s: %w",
				n.GetClass(), n.GetID().String(), cn.GetClass(), cn.GetID().String(), err)
		}
	}
	return nil
}

func (n *node) GetChildren(className string) []Node {
	// check if parent nodes completely loaded
	if len(n.ChShadow[className]) != len(n.Children[className]) {
		n.Children[className] = checkMissingReferences(n.ChShadow[className], n.Children[className], className)
	}
	return n.Children[className]
}

func (n *node) AddParent(pn Node) error {
	if !n.nodeType.Accepting(pn.GetClass()) {
		return fmt.Errorf("could not add parent, child type %s is not accepted by parent type %s",
			n.nodeType.Class(), pn.GetClass())
	}

	log.Println("add parent", pn.GetID(), "to", n.GetID())
	if _, ok := n.Children[pn.GetClass()]; !ok {
		n.Parents[pn.GetClass()] = make([]Node, 0)
		n.PtShadow[pn.GetClass()] = make([]primitive.ObjectID, 0)
	}
	n.Parents[pn.GetClass()] = append(n.Children[pn.GetClass()], pn)
	n.PtShadow[pn.GetClass()] = append(n.ChShadow[pn.GetClass()], pn.GetID())

	if err := mongodb.Save(n); err != nil {
		return fmt.Errorf("%s, %s could not add parent %s, %s: %w",
			n.GetClass(), n.GetID().String(), pn.GetClass(), pn.GetID().String(), err)
	}

	if !pn.HasChild(n.GetClass(), n.GetID()) {
		log.Println("add also to child", n.GetID(), "to", pn.GetID())
		if err := pn.AddChild(n); err != nil {
			return fmt.Errorf("could not add child %s, %s to parent %s, %s: %w",
				n.GetClass(), n.GetID().String(), pn.GetClass(), pn.GetID().String(), err)
		}
	}
	return nil
}

func (n *node) RemoveParent(pn Node) error {
	if _, ok := n.Parents[pn.GetClass()]; !ok {
		return nil
	}
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

	if err := mongodb.Save(n); err != nil {
		return fmt.Errorf("%s, %s could not remove parent %s, %s: %w",
			n.GetClass(), n.GetID().String(), pn.GetClass(), pn.GetID().String(), err)
	}

	if pn.HasChild(n.GetClass(), n.GetID()) {
		if err := pn.RemoveChild(n); err != nil {
			return fmt.Errorf("could not remove child %s, %s from parent %s, %s: %w",
				n.GetClass(), n.GetID().String(), pn.GetClass(), pn.GetID().String(), err)
		}
	}
	return nil
}

func (n *node) GetParents(className string) []Node {
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

func (n *node) GetNodes(className string) []Node {
	return append(n.Children[className], n.Parents[className]...)
}

func (n *node) GetClass() string {
	return n.Class
}

func (n *node) Type() NodeType {
	return n.nodeType
}

func (n *node) SetType(t NodeType) {
	n.nodeType = t
	n.Class = t.Class()
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

func (n *node) Save() error {
	return mongodb.Save(n)
}

func (n *node) SaveContent(content interface{}) error {
	n.SetContent(content)
	return mongodb.Save(n)
}

func containsId(ids []primitive.ObjectID, id primitive.ObjectID) (int, bool) {
	for it, slid := range ids {
		if slid == id {
			return it, true
		}
	}
	return -1, false
}

func containsNode(ids []Node, n Node) (int, bool) {
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
func checkMissingReferences(shadowIds []primitive.ObjectID, refs []Node, className string) []Node {
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
	refs = append(refs, mongodb.FindFromIDList(className, nodeIdList)...)
	return refs
}

func removeIdFromPosition(ids []primitive.ObjectID, pos int) ([]primitive.ObjectID, bool) {
	if len(ids) == 1 {
		return nil, false
	}
	copy(ids[pos:], ids[pos+1:])
	ids[len(ids)-1] = primitive.NilObjectID // prevents memory leak?
	ids = ids[:len(ids)-1]
	return ids, true
}

func removeNodeFromPosition(nodes []Node, pos int) ([]Node, bool) {
	if len(nodes) == 1 {
		return nil, false
	}
	copy(nodes[pos:], nodes[pos+1:])
	nodes[len(nodes)-1] = nil // prevents memory leak
	nodes = nodes[:len(nodes)-1]
	return nodes, true
}

func newHistory() history {
	h := history{
		Created: time.Now(),
		Changed: time.Now(),
	}
	return h
}
