package mesh

type Type struct {
	acceptTypes []string
	class       string
}

func NewNodeType(acceptTypes []string, className string) NodeType {
	return Type{
		acceptTypes: acceptTypes,
		class:       className,
	}
}

func (t Type) Class() string {
	return t.class
}

func (t Type) AcceptTypes() []string {
	return t.acceptTypes
}

func (t Type) Accepting(className string) bool {
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
