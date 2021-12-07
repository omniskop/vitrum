package vit

import (
	"github.com/omniskop/vitrum/vit/script"
)

// Component describes a generic vit component
type Component interface {
	DefineProperty(name string, vitType string, expression string) bool
	DefineEnum(Enumeration) bool
	Property(name string) (interface{}, bool)        // returns the value of the property with the given name, and a boolean indicating whether the property exists
	MustProperty(name string) interface{}            // same as Property but panics if the property doesn't exist
	SetProperty(name string, value interface{}) bool // sets the property with the given name to the given value and returns a boolean indicating whether the property exists
	ResolveVariable(name string) (interface{}, bool) // searches the scope for a variable with the given name. Returns either an expression or a component. The boolean indicates wether the variable was found.
	ResolveID(id string) (Component, bool)           // Recursively searches the children for a component with the given id. It does not check itself, only it's children!
	AddChild(Component)                              // Adds the given component as a child and also set's their parent to this component
	Children() []Component                           // Returns all children of this component
	SetParent(Component)                             // Sets the parent of this component to the given component
	ID() string                                      // Returns the id of this component
	String() string                                  // Returns a short string representation of this component
	UpdateExpressions() (int, error)                 // Recursively reevaluate all expressions that got dirty. Returns the number of reevaluated expression (includes potential failed ones)
}

type Enumeration struct {
	Name     string
	Embedded bool
	Values   map[string]int
}

func (e Enumeration) ResolveVariable(name string) (interface{}, bool) {
	if value, ok := e.Values[name]; ok {
		return value, true
	}
	return nil, false
}

type AbstractComponent interface {
	script.VariableSource
	Instantiate(string, ComponentResolver) (Component, error)
	// Static values
}

type ComponentResolver struct {
	// TODO: make these private?
	Parent     *ComponentResolver
	Components map[string]AbstractComponent
}

func NewComponentResolver(parent *ComponentResolver) ComponentResolver {
	return ComponentResolver{
		Parent:     parent,
		Components: make(map[string]AbstractComponent),
	}
}

func (r ComponentResolver) Resolve(names ...string) (AbstractComponent, bool) {
	src, ok := r.Components[names[0]]
	if !ok {
		if r.Parent != nil {
			return r.Parent.Resolve(names...)
		} else {
			return nil, false
		}
	}
	return src, ok
}
