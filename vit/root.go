package vit

import "fmt"

// Root is the base component all other components embed. It provides some basic functionality.
type Root struct {
	scope        ComponentContainer
	parent       Component
	id           string           // id of this component. Can only be set on creation and not be changed.
	properties   map[string]Value // custom properties defined in a vit file
	enumerations map[string]Enumeration
	children     []Component
}

func NewRoot(id string, scope ComponentContainer) Root {
	return Root{
		scope:        scope,
		id:           id,
		properties:   make(map[string]Value),
		enumerations: make(map[string]Enumeration),
	}
}

func (r *Root) String() string {
	return fmt.Sprintf("Root{%s}", r.id)
}

// DefineProperty creates a new property on the component.
// On failure it returns either a RedeclarationError or UnknownTypeError.
// TODO: currently properties can be redefined. Make a decision on that behaviour and update the documentation accordingly (inclusing the Component interface).
func (r *Root) DefineProperty(name string, vitType string, expression string, position *PositionRange) error {
	switch vitType {
	case "int":
		if expression == "" {
			r.properties[name] = NewIntValue("", position)
		} else {
			r.properties[name] = NewIntValue(expression, position)
		}
	case "float":
		if expression == "" {
			r.properties[name] = NewFloatValue("", position)
		} else {
			r.properties[name] = NewFloatValue(expression, position)
		}
	case "string":
		if expression == "" {
			r.properties[name] = NewStringValue("", position)
		} else {
			r.properties[name] = NewStringValue(expression, position)
		}
	case "alias":
		r.properties[name] = NewAliasValue(expression, position)
	default:
		if _, ok := r.enumerations[vitType]; ok {
			if expression == "" {
				r.properties[name] = NewIntValue("", position)
			} else {
				r.properties[name] = NewIntValue(expression, position)
			}
			return nil
		}
		return UnknownTypeError{TypeName: vitType}
	}
	return nil
}

func (r *Root) DefineEnum(enum Enumeration) bool {
	if _, ok := r.enumerations[enum.Name]; ok {
		return false
	}
	r.enumerations[enum.Name] = enum
	return true
}

func (r *Root) Property(key string) (interface{}, bool) {
	v, ok := r.properties[key]
	if ok {
		return v.GetValue(), true
	}
	return nil, false
}

func (r *Root) InternalProperty(key string) (Value, bool) {
	v, ok := r.properties[key]
	if ok {
		return v, true
	}
	return nil, false
}

func (r *Root) MustProperty(key string) interface{} {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Root) SetProperty(key string, value interface{}, position *PositionRange) bool {
	if _, ok := r.properties[key]; !ok {
		return false
	}
	r.properties[key].GetExpression().ChangeCode(value.(string), position)
	return true
}

// ResolveVariable; THIS NEEDS TO BE REIMPLEMENTED BY THE EMBEDDING STRUCTS TO RETURN THE CORRECT TYPE IF THE ID OF THIS COMPONENT IS REQUESTED
func (r *Root) ResolveVariable(key string) (interface{}, bool) {
	if key == "parent" {
		if r.parent == nil {
			fmt.Println("tried to access parent of root component")
			return nil, false
		}
		return r.parent, true
	}
	if key == r.id {
		return r, true
	}

	// check components in scope
	abs, ok := r.scope.Get(key)
	if ok {
		return abs, true
	}

	for name, enum := range r.enumerations {
		if name == key {
			return enum, true
		}
	}

	for name, prop := range r.properties {
		if name == key {
			return prop, true
		}
	}

	// for children we only check id's and not properties
	for _, child := range r.children {
		if child.ID() == key {
			return child, true
		}
		if comp, ok := child.ResolveID(key); ok {
			return comp, true
		}
	}

	return nil, false
}

func (r *Root) ResolveID(id string) (Component, bool) {
	for _, child := range r.children {
		if child.ID() == id {
			return child, true
		}
	}
	return nil, false
}

// AddChild; THIS NEEDS TO BE REIMPLEMENTED BY THE EMBEDDING STRUCTS TO SET THE CORRECT PARENT TYPE INSTEAD OF ROOT
// TODO: remove this method here, to force reimplementation?
func (r *Root) AddChild(child Component) {
	child.SetParent(r)
	r.children = append(r.children, child)
}

func (r *Root) SetParent(parent Component) {
	r.parent = parent
}

func (r *Root) Children() []Component {
	return r.children
}

func (r *Root) UpdateExpressions() (int, ErrorGroup) {
	var sum int
	var errors ErrorGroup
	for _, child := range r.children {
		n, err := child.UpdateExpressions()
		sum += n
		errors.AddGroup(err)
	}
	return sum, errors
}

func (r *Root) ID() string {
	return r.id
}

func (r *Root) root() *Root {
	return r
}

// finish put's the final touches on an instantiated component.
// It is quaranteed that all other surrounding components are instantiated just not necessarily finished.
// This needs to be reimplemented by each component.
func (r *Root) finish() error {
	for _, props := range r.properties {
		if alias, ok := props.(*AliasValue); ok {
			err := alias.Update(r)
			if err != nil {
				return fmt.Errorf("alias error: %w", err)
			}
		}
	}

	for _, child := range r.children {
		err := child.finish()
		if err != nil {
			return err
		}
	}
	return nil
}
