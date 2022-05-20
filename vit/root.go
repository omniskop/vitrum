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
func (r *Root) DefineProperty(propDef PropertyDefinition) error {
	name := propDef.Identifier[0]
	switch propDef.VitType {
	case "int":
		if propDef.Expression == "" {
			r.properties[name] = NewIntValue("", &propDef.Pos)
		} else {
			r.properties[name] = NewIntValue(propDef.Expression, &propDef.Pos)
		}
	case "float":
		if propDef.Expression == "" {
			r.properties[name] = NewFloatValue("", &propDef.Pos)
		} else {
			r.properties[name] = NewFloatValue(propDef.Expression, &propDef.Pos)
		}
	case "string":
		if propDef.Expression == "" {
			r.properties[name] = NewStringValue("", &propDef.Pos)
		} else {
			r.properties[name] = NewStringValue(propDef.Expression, &propDef.Pos)
		}
	case "bool":
		if propDef.Expression == "" {
			r.properties[name] = NewBoolValue("", &propDef.Pos)
		} else {
			r.properties[name] = NewBoolValue(propDef.Expression, &propDef.Pos)
		}
	case "alias":
		r.properties[name] = NewAliasValue(propDef.Expression, &propDef.Pos)
	case "component":
		r.properties[name] = NewComponentValue(propDef.Components[0], &propDef.Pos)
	case "var":
		r.properties[name] = NewAnyValue(propDef.Expression, &propDef.Pos)
	default:
		if _, ok := r.enumerations[propDef.VitType]; ok {
			if propDef.Expression == "" {
				r.properties[name] = NewIntValue("", &propDef.Pos)
			} else {
				r.properties[name] = NewIntValue(propDef.Expression, &propDef.Pos)
			}
			return nil
		}
		return UnknownTypeError{TypeName: propDef.VitType}
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

func (r *Root) Property(key string) (Value, bool) {
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

func (r *Root) SetProperty(key string, newValue interface{}, position *PositionRange) bool {
	prop, ok := r.properties[key]
	if !ok {
		return false
	}
	switch actual := newValue.(type) {
	case string:
		prop.GetExpression().ChangeCode(actual, position)
	case PropertyDefinition:
		prop.SetFromProperty(actual)
	}
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

// UpdatePropertiesInContext updates all properties with the given component as a context.
// SHOULD ONLY BE CALLED FROM OTHER COMPONENTS THAT EMBED THIS ROOT.
func (r *Root) UpdatePropertiesInContext(context Component) (int, ErrorGroup) {
	var sum int
	var errs ErrorGroup
	for name, prop := range r.properties {
		if prop.ShouldEvaluate() {
			sum++
			err := prop.Update(context)
			if err != nil {
				errs.Add(NewExpressionError("Rectangle", name, r.id, *prop.GetExpression(), err))
			}
		}
	}
	return sum, errs
}

// InstantiateInScope will instantiate the given component in the scope of this root component.
func (r *Root) InstantiateInScope(comp *ComponentDefinition) (Component, error) {
	return InstantiateComponent(comp, r.scope)
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
			err := alias.Update(r)
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
	return nil
}

// FinishInContext put's the final touches on an instantiated component.
// It is guaranteed that all other surrounding components are instantiated just not necessarily finished.
// SHOULD ONLY BE CALLED FROM OTHER COMPONENTS THAT EMBED THIS ROOT.
func (r *Root) FinishInContext(context Component) error {
	for _, props := range r.properties {
		if alias, ok := props.(*AliasValue); ok {
			err := alias.Update(context)
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
	return nil
}
