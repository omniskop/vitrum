package vit

import (
	"fmt"

	"github.com/omniskop/vitrum/vit/script"
)

// Component describes a generic vit component
type Component interface {
	DefineProperty(name string, vitType string, expression string, position *PositionRange) error // Creates a new property. On failure it returns either a RedeclarationError or UnknownTypeError.
	DefineEnum(Enumeration) bool
	Property(name string) (Value, bool)                                       // returns the property with the given name, and a boolean indicating whether the property exists
	MustProperty(name string) Value                                           // same as Property but panics if the property doesn't exist
	SetProperty(name string, value interface{}, position *PositionRange) bool // sets the property with the given name to the given value and returns a boolean indicating whether the property exists
	ResolveVariable(name string) (interface{}, bool)                          // searches the scope for a variable with the given name. Returns either an expression or a component. The boolean indicates wether the variable was found.
	ResolveID(id string) (Component, bool)                                    // Recursively searches the children for a component with the given id. It does not check itself, only it's children!
	AddChild(Component)                                                       // Adds the given component as a child and also set's their parent to this component
	Children() []Component                                                    // Returns all children of this component
	SetParent(Component)                                                      // Sets the parent of this component to the given component
	ID() string                                                               // Returns the id of this component
	String() string                                                           // Returns a short string representation of this component
	UpdateExpressions() (int, ErrorGroup)                                     // Recursively reevaluate all expressions that got dirty. Returns the number of reevaluated expression (includes potential failed ones)
	As(*Component) bool                                                       // Returns true if this component is of the same type as the given parameter. It also changes the parameter to point to this component.

	RootC() *Root  // returns the root of this component
	Finish() error // Finishes the component instantiation. Should only be called by components that embed this one.
}

func FinishComponent(comp Component) error {
	return comp.Finish()
}

type Enumeration struct {
	Name     string
	Embedded bool
	Values   map[string]int
	Position PositionRange
}

func (e Enumeration) ResolveVariable(name string) (interface{}, bool) {
	if value, ok := e.Values[name]; ok {
		return value, true
	}
	return nil, false
}

type AbstractComponent interface {
	script.VariableSource
	Instantiate(string, ComponentContainer) (Component, error)
	Name() string
	// Static values
}

// ComponentContainer holds a list of abstract components
type ComponentContainer struct {
	Global map[string]AbstractComponent // globally defined componets
	Local  map[string]AbstractComponent // components specific to the current document
}

func NewComponentContainer() ComponentContainer {
	return ComponentContainer{
		Global: make(map[string]AbstractComponent),
		Local:  make(map[string]AbstractComponent),
	}
}

// Returns a new ComponentContainer carrying over the global components from this one.
func (c ComponentContainer) JustGlobal() ComponentContainer {
	return ComponentContainer{
		Global: c.Global,
		Local:  make(map[string]AbstractComponent),
	}
}

// Returns a new ComponentContainer using this local components as global ones.
func (c ComponentContainer) ToGlobal() ComponentContainer {
	return ComponentContainer{
		Global: c.Local,
		Local:  make(map[string]AbstractComponent),
	}
}

func (c ComponentContainer) Set(name string, comp AbstractComponent) {
	c.Local[name] = comp
}

func (c ComponentContainer) Get(names string) (AbstractComponent, bool) {
	src, ok := c.Local[names]
	if !ok {
		src, ok = c.Global[names]
	}
	return src, ok
}

// ErrorGroup contains a list of multiple error and may be used whenever multiple errors may occur without the need to fail immediately.
// To check if an error actually occurred use the method 'Failed'.
type ErrorGroup struct {
	Errors []error
}

// Add the error to the list. If err is nil it won't be added.
func (e *ErrorGroup) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// AddGroup adds all errors of another group to this one. It doesn't matter if the other group is empty or not.
func (e *ErrorGroup) AddGroup(group ErrorGroup) {
	if !group.Failed() {
		return
	}
	e.Errors = append(e.Errors, group.Errors...)
}

// Failed returns true if the group contains at least one error.
func (e *ErrorGroup) Failed() bool {
	return len(e.Errors) > 0
}

// Error implements the error interface. It does not actually return any of the errors itself, but just a short information about the amount of errors.
func (e ErrorGroup) Error() string {
	if !e.Failed() {
		return "no errors"
	}
	return fmt.Sprintf("group with %d errors", len(e.Errors))
}

func (e ErrorGroup) Is(target error) bool {
	_, ok := target.(ErrorGroup)
	return ok
}

type RedeclarationError struct {
	PropertyName       string
	PreviousDefinition PositionRange
}

func (e RedeclarationError) Error() string {
	return fmt.Sprintf("property %q is already declared. (Previous declaration at %s)", e.PropertyName, e.PreviousDefinition.String())
}

func (e RedeclarationError) Is(target error) bool {
	_, ok := target.(RedeclarationError)
	return ok
}

type UnknownTypeError struct {
	TypeName string
}

func (e UnknownTypeError) Error() string {
	return fmt.Sprintf("unknown type '%s'", e.TypeName)
}

func (e UnknownTypeError) Is(target error) bool {
	_, ok := target.(UnknownTypeError)
	return ok
}
