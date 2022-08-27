package vit

import "fmt"

// Root is the base component all other components embed. It provides some basic functionality.
type Root struct {
	context        *FileContext
	parent         Component
	id             string           // id of this component. Can only be set on creation and not be changed.
	properties     map[string]Value // custom properties defined in a vit file
	enumerations   map[string]Enumeration
	methods        map[string]*Method
	eventListeners []Evaluater // functions that are defined on this component that will be triggered through external events
	onCompleted    EventAttribute[struct{}]
	children       []Component
}

// TODO: Investigate wether we could bring eventListeners into methods. They kinda do the same.

func NewRoot(id string, context *FileContext) Root {
	return Root{
		context:      context,
		id:           id,
		properties:   make(map[string]Value),
		enumerations: make(map[string]Enumeration),
		methods:      make(map[string]*Method),

		onCompleted: *NewEventAttribute[struct{}](),
	}
}

func (r *Root) String() string {
	return fmt.Sprintf("Root{%s}", r.id)
}

// DefineProperty creates a new property on the component.
// TODO: currently properties can be redefined. Make a decision on that behaviour and update the documentation accordingly (including the Component interface).
func (r *Root) DefineProperty(propDef PropertyDefinition) error {
	name := propDef.Identifier[0]
	if propDef.VitType == "componentdef" {
		if propDef.ListDimensions > 0 {
			r.properties[name] = NewComponentDefListValue(propDef.Components, &propDef.Pos)
		} else {
			r.properties[name] = NewComponentDefValue(propDef.Components[0], nil)
		}
	} else {
		value, err := newValueForType(propDef.VitType, Code{Code: propDef.Expression, Position: &propDef.Pos})
		if err != nil {
			// TODO: add more info?
			return err
		}
		r.properties[name] = value
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

func (r *Root) DefineMethod(method Method) bool {
	if _, ok := r.methods[method.Name]; ok {
		return false
	}
	r.methods[method.Name] = &method
	return true
}

func (r *Root) Property(key string) (Value, bool) {
	if key == "Root" {
		// this is a special case for now
		return NewComponentRefValue(r), true
	}
	v, ok := r.properties[key]
	if ok {
		return v, true
	}
	return nil, false
}

func (r *Root) MustProperty(key string) Value {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Root) SetProperty(key string, newValue interface{}) error {
	prop, ok := r.properties[key]
	if !ok {
		return unknownPropErr("", key, r.id)
	}
	err := prop.SetValue(newValue)
	if err != nil {
		return NewPropertyError("", key, r.id, err)
	}
	return nil
}

func (r *Root) SetPropertyCode(key string, code Code) error {
	prop, ok := r.properties[key]
	if !ok {
		return unknownPropErr("", key, r.id)
	}
	prop.SetCode(code)
	return nil
}

func (r *Root) AddListenerFunction(f Evaluater) {
	// TODO: offer a way to maybe remove these again?
	r.eventListeners = append(r.eventListeners, f)
}

func (r *Root) Event(name string) (Listenable, bool) {
	switch name {
	case "onCompleted":
		return &r.onCompleted, true
	}
	return nil, false
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

	for name, meth := range r.methods {
		if name == key {
			return meth, true
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

func (r *Root) AddChildAfter(afterThis, addThis Component) {
	var dynType Component = afterThis

	for j, child := range r.Children() {
		if child.As(&dynType) {
			addThis.SetParent(r)
			r.AddChildAtButKeepParent(addThis, j+1)
			return
		}
	}
	r.AddChild(addThis)
}

func (r *Root) RemoveChild(child Component) {
	for i, c := range r.children {
		if c == child {
			r.children = append(r.children[:i], r.children[i+1:]...)
			child.SetParent(nil)
			return
		}
	}
}

// AddChildButKeepParent adds the given component as a child but doesn't change the child's parent.
// SHOULD ONLY BE CALLED FROM OTHER COMPONENTS THAT EMBED THIS ROOT.
func (r *Root) AddChildButKeepParent(child Component) {
	r.children = append(r.children, child)
}

func (r *Root) AddChildAtButKeepParent(child Component, index int) {
	r.children = append(r.children[:index], append([]Component{child}, r.children[index:]...)...)
}

func (r *Root) SetParent(parent Component) {
	r.parent = parent
}

func (r *Root) Parent() Component {
	return r.parent
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

// UpdatePropertiesInContext updates all properties with the given component as a context.
// SHOULD ONLY BE CALLED FROM OTHER COMPONENTS THAT EMBED THIS ROOT.
func (r *Root) UpdatePropertiesInContext(context Component) (int, ErrorGroup) {
	var sum int
	var errs ErrorGroup
	for name, prop := range r.properties {
		if changed, err := prop.Update(r); changed || err != nil {
			sum++
			if err != nil {
				errs.Add(NewPropertyError("", name, r.id, err))
			}
		}
	}
	for _, f := range r.eventListeners {
		if f.ShouldEvaluate() {
			_, err := f.Evaluate(context)
			sum++
			if err != nil {
				errs.Add(NewPropertyError("", "<eventListener>", r.id, err))
			}
		}
	}
	for _, m := range r.methods {
		if m.ShouldEvaluate() {
			_, err := m.Evaluate(context)
			sum++
			if err != nil {
				errs.Add(NewPropertyError("", "<method>", r.id, err))
			}
		}
	}
	return sum, errs
}

// InstantiateInScope will instantiate the given component in the scope of this root component.
func (r *Root) InstantiateInScope(comp *ComponentDefinition) (Component, error) {
	return InstantiateComponent(comp, r.context)
}

func (r *Root) ID() string {
	return r.id
}

func (r *Root) As(target *Component) bool {
	if _, ok := (*target).(*Root); ok {
		*target = r
		return true
	}
	return false
}

func (r *Root) RootC() *Root {
	return r
}

// Finish put's the final touches on an instantiated component.
// It is guaranteed that all other surrounding components are instantiated just not necessarily finished.
// This needs to be reimplemented by each component.
func (r *Root) Finish() error {
	for _, props := range r.properties {
		if alias, ok := props.(*AliasValue); ok {
			_, err := alias.Update(r)
			if err != nil {
				return fmt.Errorf("alias error: %w", err)
			}
		}
	}

	for _, child := range r.children {
		err := child.Finish()
		if err != nil {
			return err
		}
	}

	r.onCompleted.Fire(nil)

	return nil
}

// FinishInContext put's the final touches on an instantiated component.
// It is guaranteed that all other surrounding components are instantiated just not necessarily finished.
// SHOULD ONLY BE CALLED FROM OTHER COMPONENTS THAT EMBED THIS ROOT.
func (r *Root) FinishInContext(context Component) error {
	for _, props := range r.properties {
		if alias, ok := props.(*AliasValue); ok {
			_, err := alias.Update(context)
			if err != nil {
				return fmt.Errorf("alias error: %w", err)
			}
		}
	}

	for _, child := range r.children {
		err := child.Finish()
		if err != nil {
			return err
		}
	}

	r.onCompleted.Fire(nil)

	return nil
}

func (r *Root) Draw(ctx DrawingContext, area Rect) error {
	return r.DrawChildren(ctx, area)
}

func (r *Root) DrawChildren(ctx DrawingContext, area Rect) error {

	for _, child := range r.children {
		child.Draw(ctx, area)
	}

	return nil
}

func (r *Root) ApplyLayout(*Layout) {}

func (r *Root) Bounds() Rect {
	return Rect{}
}

func (r *Root) Context() *FileContext {
	return r.context
}
