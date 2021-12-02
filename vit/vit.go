package vit

// Component describes a generic vit component
type Component interface {
	DefineProperty(name string, vitType string, expression string) bool
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
