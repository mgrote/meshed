package meshnode

type NodeType struct {
	acceptTypes	[]string
	class		string
}

func NewNodeType(acceptTypes []string, className string) NodeType {
	return NodeType{
		acceptTypes: acceptTypes,
		class: className,
	}
}


func (t NodeType) GetClass() string {
	return t.class
}

func (t NodeType) AcceptTypes() []string {
	return t.acceptTypes
}

func (t NodeType) IsAccepting(className string) bool {
	if len(t.acceptTypes) == 0 {
		return true
	}
	for _, cn := range t.acceptTypes {
		if cn == className {
			return true
		}
	}
	return false
}
