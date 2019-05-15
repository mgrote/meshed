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
	for _, slid := range t.acceptTypes {
		if slid == className {
			return true
		}
	}
	return false
}