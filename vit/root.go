package vit

import "fmt"

// Root is the base component all other components embed. It provides some basic functionality.
type Root struct {
	parent     Component
	id         string                 // id of this component. Can only be set on creation and not be changed.
	properties map[string]interface{} // custom properties defined in a vit file
	children   []Component
}

func NewRoot(id string) Root {
	return Root{
		id:         id,
		properties: make(map[string]interface{}),
	}
}

func (i *Root) String() string {
	return fmt.Sprintf("Root{%s}", i.id)
}

func (i *Root) Property(key string) (interface{}, bool) {
	if key == "id" {
		return i.id, true
	}
	v, ok := i.properties[key]
	return v, ok
}

func (i *Root) MustProperty(key string) interface{} {
	v, ok := i.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (i *Root) SetProperty(key string, value interface{}) bool {
	fmt.Printf("[Root] set %q: %v\n", key, value)
	if _, ok := i.properties[key]; !ok {
		return false
	}
	i.properties[key] = value
	return true
}

// ResolveVariable; THIS NEEDS TO BE REIMPLEMENTED BY THE EMBEDDING STRUCTS TO RETURN THE CORRECT TYPE IF THE ID OF THIS COMPONENT IS REQUESTED
func (i *Root) ResolveVariable(key string) (interface{}, bool) {
	if key == "parent" {
		if i.parent == nil {
			fmt.Println("tried to access parent of root component")
			return nil, false
		}
		return i.parent, true
	}
	if key == i.id {
		return i, true
	}

	for _, child := range i.children {
		if child.ID() == key {
			return child, true
		}
		if comp, ok := child.ResolveID(key); ok {
			return comp, true
		}
	}

	if i.parent != nil {
		return i.parent.ResolveVariable(key)
	}
	return nil, false
}

func (i *Root) ResolveID(id string) (Component, bool) {
	for _, child := range i.children {
		if child.ID() == id {
			return child, true
		}
	}
	return nil, false
}

// AddChild; THIS NEEDS TO BE REIMPLEMENTED BY THE EMBEDDING STRUCTS TO SET THE CORRECT PARENT TYPE INSTEAD OF ROOT
// TODO: remove this method here, to force reimplementation?
func (i *Root) AddChild(child Component) {
	child.SetParent(i)
	i.children = append(i.children, child)
}

func (i *Root) SetParent(parent Component) {
	i.parent = parent
}

func (i *Root) Children() []Component {
	return i.children
}

func (i *Root) UpdateExpressions() (int, error) {
	var sum int
	for _, child := range i.children {
		n, err := child.UpdateExpressions()
		sum += n
		if err != nil {
			return sum, err
		}
	}
	return sum, nil
}

func (i *Root) ID() string {
	return i.id
}
