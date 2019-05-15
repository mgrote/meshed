package meshnode

type NodeType struct {
	acceptTypes	[]string
	class		string
}

func (t NodeType) GetClass() string {
	return t.class
}

func (t NodeType) AcceptTypes() []string {
	return t.acceptTypes
}

func (t NodeType) IsAccepted(className string) bool {
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