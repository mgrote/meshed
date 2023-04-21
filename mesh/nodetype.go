package mesh

type Type struct {
	acceptTypes []string
	name        string
}

func NewNodeType(acceptTypes []string, typeName string) NodeType {
	return Type{
		acceptTypes: acceptTypes,
		name:        typeName,
	}
}

func (t Type) TypeName() string {
	return t.name
}

func (t Type) AcceptTypes() []string {
	return t.acceptTypes
}

func (t Type) Accepting(typeName string) bool {
	if len(t.acceptTypes) == 0 {
		return true
	}
	for _, cn := range t.acceptTypes {
		if cn == typeName {
			return true
		}
	}
	return false
}
