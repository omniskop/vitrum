// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
	parse "github.com/omniskop/vitrum/vit/parse"
)

func newFileContextForRotation(globalCtx *vit.GlobalContext) (*vit.FileContext, error) {
	return vit.NewFileContext(globalCtx), nil
}

type Rotation_HorizontalPivot uint

const (
	Rotation_HorizontalPivot_PivotLeft    Rotation_HorizontalPivot = 0
	Rotation_HorizontalPivot_PivotHCenter Rotation_HorizontalPivot = 1
	Rotation_HorizontalPivot_PivotRight   Rotation_HorizontalPivot = 2
)

func (enum Rotation_HorizontalPivot) String() string {
	switch enum {
	case Rotation_HorizontalPivot_PivotLeft:
		return "PivotLeft"
	case Rotation_HorizontalPivot_PivotHCenter:
		return "PivotHCenter"
	case Rotation_HorizontalPivot_PivotRight:
		return "PivotRight"
	default:
		return "<unknownHorizontalPivot>"
	}
}

type Rotation_VerticalPivot uint

const (
	Rotation_VerticalPivot_PivotTop     Rotation_VerticalPivot = 0
	Rotation_VerticalPivot_PivotVCenter Rotation_VerticalPivot = 1
	Rotation_VerticalPivot_PivorBottom  Rotation_VerticalPivot = 2
)

func (enum Rotation_VerticalPivot) String() string {
	switch enum {
	case Rotation_VerticalPivot_PivotTop:
		return "PivotTop"
	case Rotation_VerticalPivot_PivotVCenter:
		return "PivotVCenter"
	case Rotation_VerticalPivot_PivorBottom:
		return "PivorBottom"
	default:
		return "<unknownVerticalPivot>"
	}
}

type Rotation struct {
	*Item
	id string

	horizontalPivot vit.IntValue
	verticalPivot   vit.IntValue
	degrees         vit.FloatValue
}

// newRotationInGlobal creates an appropriate file context for the component and then returns a new Rotation instance.
// The returned error will only be set if a library import that is required by the component fails.
func newRotationInGlobal(id string, globalCtx *vit.GlobalContext, thisLibrary parse.Library) (*Rotation, error) {
	fileCtx, err := newFileContextForRotation(globalCtx)
	if err != nil {
		return nil, err
	}
	parse.AddLibraryToContainer(thisLibrary, &fileCtx.KnownComponents)
	return NewRotation(id, fileCtx), nil
}
func NewRotation(id string, context *vit.FileContext) *Rotation {
	r := &Rotation{
		Item:            NewItem("", context),
		id:              id,
		horizontalPivot: *vit.NewEmptyIntValue(),
		verticalPivot:   *vit.NewEmptyIntValue(),
		degrees:         *vit.NewEmptyFloatValue(),
	}
	// property assignments on embedded components
	// register listeners for when a property changes
	// register event listeners
	// register enumerations
	r.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "HorizontalPivot",
		Position: nil,
		Values:   map[string]int{"PivotLeft": 0, "PivotHCenter": 1, "PivotRight": 2},
	})
	r.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "VerticalPivot",
		Position: nil,
		Values:   map[string]int{"PivotTop": 0, "PivotVCenter": 1, "PivorBottom": 2},
	})
	// add child components

	context.RegisterComponent("", r)

	return r
}

func (r *Rotation) String() string {
	return fmt.Sprintf("Rotation(%s)", r.id)
}

func (r *Rotation) Property(key string) (vit.Value, bool) {
	switch key {
	case "horizontalPivot":
		return &r.horizontalPivot, true
	case "verticalPivot":
		return &r.verticalPivot, true
	case "degrees":
		return &r.degrees, true
	default:
		return r.Item.Property(key)
	}
}

func (r *Rotation) MustProperty(key string) vit.Value {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Rotation) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "horizontalPivot":
		err = r.horizontalPivot.SetValue(value)
	case "verticalPivot":
		err = r.verticalPivot.SetValue(value)
	case "degrees":
		err = r.degrees.SetValue(value)
	default:
		return r.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("Rotation", key, r.id, err)
	}
	return nil
}

func (r *Rotation) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "horizontalPivot":
		r.horizontalPivot.SetCode(code)
	case "verticalPivot":
		r.verticalPivot.SetCode(code)
	case "degrees":
		r.degrees.SetCode(code)
	default:
		return r.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (r *Rotation) Event(name string) (vit.Listenable, bool) {
	switch name {
	default:
		return r.Item.Event(name)
	}
}

func (r *Rotation) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "horizontalPivot":
		return &r.horizontalPivot, true
	case "verticalPivot":
		return &r.verticalPivot, true
	case "degrees":
		return &r.degrees, true
	default:
		return r.Item.ResolveVariable(key)
	}
}

func (r *Rotation) AddChild(child vit.Component) {
	child.SetParent(r)
	r.AddChildButKeepParent(child)
}

func (r *Rotation) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	var targetType vit.Component = afterThis

	for ind, child := range r.Children() {
		if child.As(&targetType) {
			addThis.SetParent(r)
			r.AddChildAtButKeepParent(addThis, ind+1)
			return
		}
	}
	r.AddChild(addThis)
}

func (r *Rotation) UpdateExpressions(context vit.Component) (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if context == nil {
		context = r
	}
	// properties
	if changed, err := r.horizontalPivot.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Rotation", "horizontalPivot", r.id, err))
		}
	}
	if changed, err := r.verticalPivot.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Rotation", "verticalPivot", r.id, err))
		}
	}
	if changed, err := r.degrees.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Rotation", "degrees", r.id, err))
		}
	}

	// methods

	n, err := r.Item.UpdateExpressions(context)
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (r *Rotation) As(target *vit.Component) bool {
	if _, ok := (*target).(*Rotation); ok {
		*target = r
		return true
	}
	return r.Item.As(target)
}

func (r *Rotation) ID() string {
	return r.id
}

func (r *Rotation) Finish() error {
	return r.RootC().FinishInContext(r)
}

func (r *Rotation) staticAttribute(name string) (interface{}, bool) {
	switch name {
	case "PivotLeft":
		return uint(Rotation_HorizontalPivot_PivotLeft), true
	case "PivotHCenter":
		return uint(Rotation_HorizontalPivot_PivotHCenter), true
	case "PivotRight":
		return uint(Rotation_HorizontalPivot_PivotRight), true
	case "PivotTop":
		return uint(Rotation_VerticalPivot_PivotTop), true
	case "PivotVCenter":
		return uint(Rotation_VerticalPivot_PivotVCenter), true
	case "PivorBottom":
		return uint(Rotation_VerticalPivot_PivorBottom), true
	default:
		return nil, false
	}
}
