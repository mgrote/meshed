package mesh

import (
	"fmt"
	"log"
	"time"
)

type NodeType interface {
	TypeName() string
	AcceptTypes() []string
	Accepting(string) bool
}

type Node interface {
	GetID() string
	SetID(id interface{})
	AddChild(Node) error
	RemoveChild(Node) error
	GetChildren(string) []Node
	GetChildrenIn(...string) []Node
	AddParent(Node) error
	RemoveParent(Node) error
	GetParents(string) []Node
	HasChild(string, string) bool
	HasParent(string, string) bool
	RemoveAllNodes(string)
	GetNodes(string) []Node
	GetTypeName() string
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
	//ID       primitive.ObjectID              `json:"id" bson:"_id,omitempty"`
	ID       string              `json:"id"`
	Version  uint16              `json:"version"`
	TypeName string              `json:"name"`
	nodeType NodeType            //`json:"-"`
	Content  interface{}         `json:"content"`
	History  history             `json:"history"`
	Children map[string][]Node   `json:"-" bson:"-"`
	ChShadow map[string][]string `json:"chshadow"`
	Parents  map[string][]Node   `json:"-" bson:"-"`
	PtShadow map[string][]string `json:"ptshadow"`
}

type history struct {
	Created  time.Time `json:"created"`
	Changed  time.Time `json:"changed"`
	LastUser string    `json:"lastuser"`
}

func NewNode(t NodeType) Node {
	n := node{
		TypeName: t.TypeName(),
		nodeType: t,
		History:  newHistory(),
		Parents:  make(map[string][]Node),
		Children: make(map[string][]Node),
		ChShadow: make(map[string][]string),
		PtShadow: make(map[string][]string),
	}
	return &n
}

func NewNodeWithContent(t NodeType, c interface{}) Node {
	n := node{
		TypeName: t.TypeName(),
		nodeType: t,
		History:  newHistory(),
		Content:  c,
		Parents:  make(map[string][]Node),
		Children: make(map[string][]Node),
		ChShadow: make(map[string][]string),
		PtShadow: make(map[string][]string),
	}
	return &n
}

// func (n *node) GetID() primitive.ObjectID {
func (n *node) GetID() string {
	return n.ID
}

func (n *node) SetID(ident interface{}) {
	//if id, ok := ident.(primitive.ObjectID); ok {
	if id, ok := ident.(string); ok {
		n.ID = id
	}
	if id, ok := ident.(fmt.Stringer); ok {
		n.ID = id.String()
		log.Println("set id to", n.ID, n.GetContent())
	} else {
		log.Fatal("error convert ID to string", ident)
	}
}

func (n *node) AddChild(cn Node) error {
	if !n.nodeType.Accepting(cn.GetTypeName()) {
		return fmt.Errorf("could not add child, parent type %s is not accepted by child type %s",
			n.nodeType.TypeName(), cn.GetTypeName())
	}

	log.Println("add child", cn.GetID(), "to", n.GetID())
	if _, ok := n.Children[cn.GetTypeName()]; !ok {
		n.Children[cn.GetTypeName()] = make([]Node, 0)
		n.ChShadow[cn.GetTypeName()] = make([]string, 0)
	}
	n.Children[cn.GetTypeName()] = append(n.Children[cn.GetTypeName()], cn)
	n.ChShadow[cn.GetTypeName()] = append(n.ChShadow[cn.GetTypeName()], cn.GetID())
	if err := NodeService.Save(n); err != nil {
		return fmt.Errorf("%s, %s could not add child %s, %s: %w",
			n.GetTypeName(), n.GetID(), cn.GetTypeName(), cn.GetID(), err)
	}

	if !cn.HasParent(n.GetTypeName(), n.GetID()) {
		if err := cn.AddParent(n); err != nil {
			return fmt.Errorf("could not add parent %s, %s to child %s, %s: %w",
				n.GetTypeName(), n.GetID(), cn.GetTypeName(), cn.GetID(), err)
		}
	}
	return nil
}

func (n *node) RemoveChild(cn Node) error {
	if _, ok := n.Children[cn.GetTypeName()]; !ok {
		return nil
	}

	if pos, oki := containsId(n.ChShadow[cn.GetTypeName()], cn.GetID()); oki {
		if isl, ok := removeIdFromPosition(n.ChShadow[cn.GetTypeName()], pos); ok {
			n.ChShadow[cn.GetTypeName()] = isl
		} else {
			delete(n.ChShadow, cn.GetTypeName())
		}
	}
	if pos, oki := containsNode(n.Children[cn.GetTypeName()], cn); oki {
		if nsl, ok := removeNodeFromPosition(n.Children[cn.GetTypeName()], pos); ok {
			n.Children[cn.GetTypeName()] = nsl
		} else {
			delete(n.Children, cn.GetTypeName())
		}
	}
	if err := NodeService.Save(n); err != nil {
		return fmt.Errorf("%s, %s could not remove child %s, %s: %w",
			n.GetTypeName(), n.GetID(), cn.GetTypeName(), cn.GetID(), err)
	}

	if cn.HasParent(n.GetTypeName(), n.GetID()) {
		if err := cn.RemoveParent(n); err != nil {
			return fmt.Errorf("could not remove parent %s, %s from child %s, %s: %w",
				n.GetTypeName(), n.GetID(), cn.GetTypeName(), cn.GetID(), err)
		}
	}
	return nil
}

func (n *node) GetChildren(typeName string) []Node {
	// check if parent nodes completely loaded
	if len(n.ChShadow[typeName]) != len(n.Children[typeName]) {
		n.Children[typeName] = checkMissingReferences(n.ChShadow[typeName], n.Children[typeName], typeName)
	}
	return n.Children[typeName]
}

func (n *node) GetChildrenIn(typeNames ...string) []Node {
	// check if parent nodes completely loaded
	var childrenIn []Node
	for _, typeName := range typeNames {
		childrenIn = append(childrenIn, n.GetChildren(typeName)...)
	}
	return childrenIn
}

func (n *node) AddParent(pn Node) error {
	if !n.nodeType.Accepting(pn.GetTypeName()) {
		return fmt.Errorf("could not add parent, child type %s is not accepted by parent type %s",
			n.nodeType.TypeName(), pn.GetTypeName())
	}

	log.Println("add parent", pn.GetID(), "to", n.GetID())
	if _, ok := n.Children[pn.GetTypeName()]; !ok {
		n.Parents[pn.GetTypeName()] = make([]Node, 0)
		n.PtShadow[pn.GetTypeName()] = make([]string, 0)
	}
	n.Parents[pn.GetTypeName()] = append(n.Children[pn.GetTypeName()], pn)
	n.PtShadow[pn.GetTypeName()] = append(n.ChShadow[pn.GetTypeName()], pn.GetID())

	if err := NodeService.Save(n); err != nil {
		return fmt.Errorf("%s, %s could not add parent %s, %s: %w",
			n.GetTypeName(), n.GetID(), pn.GetTypeName(), pn.GetID(), err)
	}

	if !pn.HasChild(n.GetTypeName(), n.GetID()) {
		log.Println("add also to child", n.GetID(), "to", pn.GetID())
		if err := pn.AddChild(n); err != nil {
			return fmt.Errorf("could not add child %s, %s to parent %s, %s: %w",
				n.GetTypeName(), n.GetID(), pn.GetTypeName(), pn.GetID(), err)
		}
	}
	return nil
}

func (n *node) RemoveParent(pn Node) error {
	if _, ok := n.Parents[pn.GetTypeName()]; !ok {
		return nil
	}
	if pos, oki := containsId(n.PtShadow[pn.GetTypeName()], pn.GetID()); oki {
		if isl, ok := removeIdFromPosition(n.PtShadow[pn.GetTypeName()], pos); ok {
			n.PtShadow[pn.GetTypeName()] = isl
		} else {
			delete(n.PtShadow, pn.GetTypeName())
		}
	}
	if pos, oki := containsNode(n.Parents[pn.GetTypeName()], pn); oki {
		if nsl, ok := removeNodeFromPosition(n.Parents[pn.GetTypeName()], pos); ok {
			n.Parents[pn.GetTypeName()] = nsl
		} else {
			delete(n.Parents, pn.GetTypeName())
		}
	}

	if err := NodeService.Save(n); err != nil {
		return fmt.Errorf("%s, %s could not remove parent %s, %s: %w",
			n.GetTypeName(), n.GetID(), pn.GetTypeName(), pn.GetID(), err)
	}

	if pn.HasChild(n.GetTypeName(), n.GetID()) {
		if err := pn.RemoveChild(n); err != nil {
			return fmt.Errorf("could not remove child %s, %s from parent %s, %s: %w",
				n.GetTypeName(), n.GetID(), pn.GetTypeName(), pn.GetID(), err)
		}
	}
	return nil
}

func (n *node) GetParents(typeName string) []Node {
	if len(n.ChShadow[typeName]) != len(n.Children[typeName]) {
		n.Parents[typeName] = checkMissingReferences(n.PtShadow[typeName], n.Parents[typeName], typeName)
	}
	return n.Parents[typeName]
}

func (n *node) HasChild(typeName string, id string) bool {
	if ids, ok := n.ChShadow[typeName]; ok {
		_, oki := containsId(ids, id)
		return oki
	}
	return false
}

func (n *node) HasParent(typeName string, id string) bool {
	if ids, ok := n.PtShadow[typeName]; ok {
		_, oki := containsId(ids, id)
		return oki
	}
	return false
}

func (n *node) RemoveAllNodes(typeName string) {
	delete(n.Children, typeName)
	delete(n.Parents, typeName)
}

func (n *node) GetNodes(typeName string) []Node {
	return append(n.Children[typeName], n.Parents[typeName]...)
}

func (n *node) GetTypeName() string {
	return n.TypeName
}

func (n *node) Type() NodeType {
	return n.nodeType
}

func (n *node) SetType(t NodeType) {
	n.nodeType = t
	n.TypeName = t.TypeName()
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
	return NodeService.Save(n)
}

func (n *node) SaveContent(content interface{}) error {
	n.SetContent(content)
	return NodeService.Save(n)
}

func containsId(ids []string, id string) (int, bool) {
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
func checkMissingReferences(shadowIds []string, refs []Node, typeName string) []Node {
	var nodeIdList = make([]string, len(shadowIds))
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
	refs = append(refs, NodeService.FindNodesFromIDList(typeName, nodeIdList)...)
	return refs
}

func removeIdFromPosition(ids []string, pos int) ([]string, bool) {
	if len(ids) == 1 {
		return nil, false
	}
	copy(ids[pos:], ids[pos+1:])
	ids[len(ids)-1] = "" // prevents memory leak?
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
