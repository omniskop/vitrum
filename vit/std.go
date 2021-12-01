package vit

type StdLib struct {
}

func (l StdLib) ComponentNames() []string {
	return []string{"Item"}
}

func (l StdLib) NewComponent(name string, id string) (Component, bool) {
	switch name {
	case "Item":
		return NewItem(id), true
	}
	return nil, false
}
