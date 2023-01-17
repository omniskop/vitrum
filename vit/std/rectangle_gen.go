// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
	parse "github.com/omniskop/vitrum/vit/parse"
)

func newFileContextForRectangle(globalCtx *vit.GlobalContext) (*vit.FileContext, error) {
	return vit.NewFileContext(globalCtx), nil
}

type Rectangle struct {
	*Item
	id string

	color  vit.ColorValue
	radius vit.FloatValue
	border vit.GroupValue
}

// newRectangleInGlobal creates an appropriate file context for the component and then returns a new Rectangle instance.
// The returned error will only be set if a library import that is required by the component fails.
func newRectangleInGlobal(id string, globalCtx *vit.GlobalContext, thisLibrary parse.Library) (*Rectangle, error) {
	fileCtx, err := newFileContextForRectangle(globalCtx)
	if err != nil {
		return nil, err
	}
	parse.AddLibraryToContainer(thisLibrary, &fileCtx.KnownComponents)
	return NewRectangle(id, fileCtx), nil
}
func NewRectangle(id string, context *vit.FileContext) *Rectangle {
	r := &Rectangle{
		Item:   NewItem("", context),
		id:     id,
		color:  *vit.NewColorValueFromCode(vit.Code{FileCtx: context, Code: "Vit.rgb(0, 0, 0)", Position: nil}),
		radius: *vit.NewFloatValueFromCode(vit.Code{FileCtx: context, Code: "0", Position: nil}),
		border: *vit.NewEmptyGroupValue(map[string]vit.Value{
			"color": vit.NewColorValueFromCode(vit.Code{FileCtx: context, Code: "\"transparent\"", Position: nil}),
			"width": vit.NewIntValueFromCode(vit.Code{FileCtx: context, Code: "0", Position: nil}),
		}),
	}
	// property assignments on embedded components
	// register listeners for when a property changes
	// register event listeners
	// register enumerations
	// add child components

	context.RegisterComponent("", r)

	return r
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle(%s)", r.id)
}

func (r *Rectangle) Property(key string) (vit.Value, bool) {
	switch key {
	case "color":
		return &r.color, true
	case "radius":
		return &r.radius, true
	case "border":
		return &r.border, true
	default:
		return r.Item.Property(key)
	}
}

func (r *Rectangle) MustProperty(key string) vit.Value {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Rectangle) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "color":
		err = r.color.SetValue(value)
	case "radius":
		err = r.radius.SetValue(value)
	case "border":
		err = r.border.SetValue(value)
	default:
		return r.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("Rectangle", key, r.id, err)
	}
	return nil
}

func (r *Rectangle) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "color":
		r.color.SetCode(code)
	case "radius":
		r.radius.SetCode(code)
	case "border":
		r.border.SetCode(code)
	default:
		return r.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (r *Rectangle) Event(name string) (vit.Listenable, bool) {
	switch name {
	default:
		return r.Item.Event(name)
	}
}

func (r *Rectangle) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "color":
		return &r.color, true
	case "radius":
		return &r.radius, true
	case "border":
		return &r.border, true
	default:
		return r.Item.ResolveVariable(key)
	}
}

func (r *Rectangle) AddChild(child vit.Component) {
	child.SetParent(r)
	r.AddChildButKeepParent(child)
}

func (r *Rectangle) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
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

func (r *Rectangle) UpdateExpressions(context vit.Component) (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if context == nil {
		context = r
	}
	// properties
	if changed, err := r.color.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Rectangle", "color", r.id, err))
		}
	}
	if changed, err := r.radius.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Rectangle", "radius", r.id, err))
		}
	}
	if changed, err := r.border.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Rectangle", "border", r.id, err))
		}
	}

	// methods

	n, err := r.Item.UpdateExpressions(context)
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (r *Rectangle) As(target *vit.Component) bool {
	if _, ok := (*target).(*Rectangle); ok {
		*target = r
		return true
	}
	return r.Item.As(target)
}

func (r *Rectangle) ID() string {
	return r.id
}

func (r *Rectangle) Finish() error {
	return r.RootC().FinishInContext(r)
}
