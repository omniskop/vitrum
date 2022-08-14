package vit

import (
	"fmt"
	"log"

	"github.com/omniskop/vitrum/vit/script"
)

// This will be set by the parser package at initialization.
// It is ugly but we can't use circular dependencies.
// Right now the alternative would be to pass a type from the parser package to every vit component specifically for this purpose and that's not really what I want to do either.
var InstantiateComponent func(*ComponentDefinition, *FileContext) (Component, error)

// Component describes a generic vit component
type Component interface {
	DefineProperty(PropertyDefinition) error // Creates a new property. On failure it returns either a RedeclarationError or UnknownTypeError.
	DefineEnum(Enumeration) bool
	DefineMethod(Method) bool
	Property(name string) (Value, bool)               // returns the property with the given name, and a boolean indicating whether the property exists
	MustProperty(name string) Value                   // same as Property but panics if the property doesn't exist
	SetProperty(name string, value interface{}) error // sets the property with the given name to the given value
	SetPropertyCode(name string, code Code) error     // sets the property with the given name to the given expression
	Event(name string) (Listenable, bool)             // returns the event with the given name, and a boolean indicating whether the event exists
	ResolveVariable(name string) (interface{}, bool)  // searches the scope for a variable with the given name. Returns either an expression or a component. The boolean indicates wether the variable was found.
	AddChild(Component)                               // Adds the given component as a child and also set's their parent to this component
	AddChildAfter(Component, Component)
	Children() []Component                // Returns all children of this component
	SetParent(Component)                  // Sets the parent of this component to the given component
	ID() string                           // Returns the id of this component
	String() string                       // Returns a short string representation of this component
	UpdateExpressions() (int, ErrorGroup) // Recursively reevaluate all expressions that got dirty. Returns the number of reevaluated expression (includes potential failed ones)
	As(*Component) bool                   // Returns true if this component is of the same type as the given parameter. It also changes the parameter to point to this component.
	ApplyLayout(*Layout)

	Draw(DrawingContext, Rect) error
	Bounds() Rect

	RootC() *Root  // returns the root of this component
	Finish() error // Finishes the component instantiation. Should only be called by components that embed this one.
}

type FocusableComponent interface {
	Focus()
	Blur()
}

func FinishComponent(comp Component) error {
	return comp.Finish()
}

type Enumeration struct {
	Name     string
	Embedded bool
	Values   map[string]int
	Position *PositionRange
}

func (e Enumeration) ResolveVariable(name string) (interface{}, bool) {
	if value, ok := e.Values[name]; ok {
		return value, true
	}
	return nil, false
}

type Method struct {
	Name string
	AsyncFunction
}

func NewMethod(name string, code string, positon *PositionRange, fileCtx *FileContext) Method {
	return Method{
		Name:          name,
		AsyncFunction: *NewAsyncFunction(code, positon, fileCtx),
	}
}

func NewMethodFromCode(name string, code Code) Method {
	return NewMethod(name, code.Code, code.Position, code.FileCtx)
}

func (m Method) CopyInContext(fileCtx *FileContext) Method {
	return NewMethod(m.Name, m.code, m.Position, fileCtx)
}

type AbstractComponent interface {
	script.VariableSource
	Instantiate(string, *GlobalContext) (Component, error)
	Name() string
	// Static values
}

// ComponentContainer holds a list of abstract components
type ComponentContainer struct {
	Global map[string]AbstractComponent // globally defined components
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

// GlobalContext holds information about a vitrum instance
type GlobalContext struct {
	KnownComponents ComponentContainer // globally known components
	Variables       map[string]Value
	Environment     ExecutionEnvironment
}

func (c *GlobalContext) Get(name string) (AbstractComponent, bool) {
	return c.KnownComponents.Get(name)
}

func (c *GlobalContext) ResolveVariable(name string) (interface{}, bool) {
	if comp, ok := c.KnownComponents.Get(name); ok {
		return comp, true
	}
	v, ok := c.Variables[name]
	return v, ok
}

func (c *GlobalContext) SetVariable(name string, value interface{}) error {
	if _, ok := c.Variables[name]; ok {
		return c.Variables[name].SetValue(value)
	}

	v, err := newValueFromGo(value)
	if err != nil {
		return err
	}
	c.Variables[name] = v
	return nil
}

// FileContext holds information about a file.
// It also contains a reference to the global context and can be used to access things like component definitions.
type FileContext struct {
	Global          *GlobalContext       // global context
	KnownComponents ComponentContainer   // Components that are known inside the file
	IDs             map[string]Component // mapping from id's to components in the file
}

func NewFileContext(global *GlobalContext) *FileContext {
	return &FileContext{
		Global:          global,
		KnownComponents: NewComponentContainer(),
		IDs:             make(map[string]Component),
	}
}

func (ctx *FileContext) RegisterComponent(id string, comp Component) {
	if id != "" {
		ctx.IDs[id] = comp
	}
	ctx.Global.Environment.RegisterComponent(id, comp)
}

func (ctx *FileContext) UnregisterComponent(id string, comp Component) {
	if id != "" {
		delete(ctx.IDs, id)
	}
	ctx.Global.Environment.UnregisterComponent(id, comp)
}

// Get returns the component with the given name.
// The returned boolean indicates whether the component was found.
func (ctx *FileContext) Get(name string) (AbstractComponent, bool) {
	if comp, ok := ctx.KnownComponents.Get(name); ok {
		return comp, true
	}
	if comp, ok := ctx.Global.Get(name); ok {
		return comp, true
	}
	return nil, false
}

// ResolveVariable returns defined components with the given name, existing components with the given id or globally defined values.
func (ctx *FileContext) ResolveVariable(name string) (interface{}, bool) {
	if comp, ok := ctx.KnownComponents.Get(name); ok {
		return comp, true
	}
	if comp, ok := ctx.IDs[name]; ok {
		return comp, true
	}
	return ctx.Global.ResolveVariable(name)
}

func (ctx *FileContext) GetComponentByID(id string) (Component, bool) {
	if comp, ok := ctx.IDs[id]; ok {
		return comp, true
	}
	return nil, false
}

type ExecutionEnvironment interface {
	RegisterComponent(string, Component)
	UnregisterComponent(string, Component)
	RequestFocus(FocusableComponent)
	Logger() *log.Logger
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

type PropertyError struct {
	componentName string
	componentID   string
	propertyName  string
	err           error
}

func NewPropertyError(componentName string, propertyName string, componentID string, err error) PropertyError {
	return PropertyError{
		componentName: componentName,
		componentID:   componentID,
		propertyName:  propertyName,
		err:           err,
	}
}

func (e PropertyError) Error() string {
	var identifier string
	if e.componentID == "" {
		if e.componentName == "" {
			identifier = e.propertyName
		} else {
			identifier = fmt.Sprintf("%s.%s", e.componentName, e.propertyName)
		}
	} else {
		if e.componentName == "" {
			identifier = fmt.Sprintf("(%s).%s", e.componentID, e.propertyName)
		} else {
			identifier = fmt.Sprintf("%s(%s).%s", e.componentID, e.componentName, e.propertyName)
		}
	}
	return fmt.Sprintf("%s: %s", identifier, e.err.Error())
}

func (e PropertyError) Is(target error) bool {
	_, ok := target.(PropertyError)
	return ok
}

func (e PropertyError) Unwrap() error {
	return e.err
}

type UnknownPropertyError struct{}

func unknownPropErr(componentName, propertyName, componentID string) PropertyError {
	return NewPropertyError(componentName, propertyName, componentID, UnknownPropertyError{})
}

func (e UnknownPropertyError) Error() string {
	return fmt.Sprintf("unknown property")
}

func (e UnknownPropertyError) Is(target error) bool {
	_, ok := target.(UnknownPropertyError)
	return ok
}

func copyMap[K comparable, V any](original map[K]V) map[K]V {
	var output = make(map[K]V)
	for k, v := range original {
		output[k] = v
	}
	return output
}
