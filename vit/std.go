package vit

type StdLib struct {
}

func (l StdLib) ComponentNames() []string {
	return []string{"Item", "Rectangle"}
}

func (l StdLib) NewComponent(name string, id string, scope ComponentContainer) (Component, bool) {
	switch name {
	case "Item":
		return NewItem(id, scope), true
	case "Rectangle":
		return NewRectangle(id, scope), true
	}
	return nil, false
}
