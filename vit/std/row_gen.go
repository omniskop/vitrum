// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
	parse "github.com/omniskop/vitrum/vit/parse"
)

func newFileContextForRow(globalCtx *vit.GlobalContext) (*vit.FileContext, error) {
	return vit.NewFileContext(globalCtx), nil
}

type Row struct {
	*Item
	id string

	topPadding    vit.OptionalValue[*vit.FloatValue]
	rightPadding  vit.OptionalValue[*vit.FloatValue]
	bottomPadding vit.OptionalValue[*vit.FloatValue]
	leftPadding   vit.OptionalValue[*vit.FloatValue]
	padding       vit.FloatValue
	spacing       vit.FloatValue
	childLayouts  vit.LayoutList
}

// newRowInGlobal creates an appropriate file context for the component and then returns a new Row instance.
// The returned error will only be set if a library import that is required by the component fails.
func newRowInGlobal(id string, globalCtx *vit.GlobalContext, thisLibrary parse.Library) (*Row, error) {
	fileCtx, err := newFileContextForRow(globalCtx)
	if err != nil {
		return nil, err
	}
	parse.AddLibraryToContainer(thisLibrary, &fileCtx.KnownComponents)
	return NewRow(id, fileCtx), nil
}
func NewRow(id string, context *vit.FileContext) *Row {
	r := &Row{
		Item:          NewItem("", context),
		id:            id,
		topPadding:    *vit.NewOptionalValue(vit.NewFloatValueFromCode(vit.Code{FileCtx: context, Code: "0", Position: nil})),
		rightPadding:  *vit.NewOptionalValue(vit.NewFloatValueFromCode(vit.Code{FileCtx: context, Code: "0", Position: nil})),
		bottomPadding: *vit.NewOptionalValue(vit.NewFloatValueFromCode(vit.Code{FileCtx: context, Code: "0", Position: nil})),
		leftPadding:   *vit.NewOptionalValue(vit.NewFloatValueFromCode(vit.Code{FileCtx: context, Code: "0", Position: nil})),
		padding:       *vit.NewFloatValueFromCode(vit.Code{FileCtx: context, Code: "0", Position: nil}),
		spacing:       *vit.NewFloatValueFromCode(vit.Code{FileCtx: context, Code: "0", Position: nil}),
		childLayouts:  make(vit.LayoutList),
	}
	// property assignments on embedded components
	// register listeners for when a property changes
	r.topPadding.AddDependent(vit.FuncDep(r.recalculateLayout))
	r.rightPadding.AddDependent(vit.FuncDep(r.recalculateLayout))
	r.bottomPadding.AddDependent(vit.FuncDep(r.recalculateLayout))
	r.leftPadding.AddDependent(vit.FuncDep(r.recalculateLayout))
	r.padding.AddDependent(vit.FuncDep(r.recalculateLayout))
	r.spacing.AddDependent(vit.FuncDep(r.recalculateLayout))
	// register event listeners
	// register enumerations
	// add child components

	context.RegisterComponent("", r)

	return r
}

func (r *Row) String() string {
	return fmt.Sprintf("Row(%s)", r.id)
}

func (r *Row) Property(key string) (vit.Value, bool) {
	switch key {
	case "topPadding":
		return &r.topPadding, true
	case "rightPadding":
		return &r.rightPadding, true
	case "bottomPadding":
		return &r.bottomPadding, true
	case "leftPadding":
		return &r.leftPadding, true
	case "padding":
		return &r.padding, true
	case "spacing":
		return &r.spacing, true
	default:
		return r.Item.Property(key)
	}
}

func (r *Row) MustProperty(key string) vit.Value {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Row) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "topPadding":
		err = r.topPadding.SetValue(value)
	case "rightPadding":
		err = r.rightPadding.SetValue(value)
	case "bottomPadding":
		err = r.bottomPadding.SetValue(value)
	case "leftPadding":
		err = r.leftPadding.SetValue(value)
	case "padding":
		err = r.padding.SetValue(value)
	case "spacing":
		err = r.spacing.SetValue(value)
	default:
		return r.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("Row", key, r.id, err)
	}
	return nil
}

func (r *Row) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "topPadding":
		r.topPadding.SetCode(code)
	case "rightPadding":
		r.rightPadding.SetCode(code)
	case "bottomPadding":
		r.bottomPadding.SetCode(code)
	case "leftPadding":
		r.leftPadding.SetCode(code)
	case "padding":
		r.padding.SetCode(code)
	case "spacing":
		r.spacing.SetCode(code)
	default:
		return r.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (r *Row) Event(name string) (vit.Listenable, bool) {
	switch name {
	default:
		return r.Item.Event(name)
	}
}

func (r *Row) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "topPadding":
		return &r.topPadding, true
	case "rightPadding":
		return &r.rightPadding, true
	case "bottomPadding":
		return &r.bottomPadding, true
	case "leftPadding":
		return &r.leftPadding, true
	case "padding":
		return &r.padding, true
	case "spacing":
		return &r.spacing, true
	default:
		return r.Item.ResolveVariable(key)
	}
}

func (r *Row) AddChild(child vit.Component) {
	defer r.childWasAdded(child)
	child.SetParent(r)
	r.AddChildButKeepParent(child)
}

func (r *Row) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	defer r.childWasAdded(addThis)
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

func (r *Row) UpdateExpressions(context vit.Component) (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if context == nil {
		context = r
	}
	// properties
	if changed, err := r.topPadding.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "topPadding", r.id, err))
		}
	}
	if changed, err := r.rightPadding.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "rightPadding", r.id, err))
		}
	}
	if changed, err := r.bottomPadding.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "bottomPadding", r.id, err))
		}
	}
	if changed, err := r.leftPadding.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "leftPadding", r.id, err))
		}
	}
	if changed, err := r.padding.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "padding", r.id, err))
		}
	}
	if changed, err := r.spacing.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "spacing", r.id, err))
		}
	}

	// methods

	n, err := r.Item.UpdateExpressions(context)
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (r *Row) As(target *vit.Component) bool {
	if _, ok := (*target).(*Row); ok {
		*target = r
		return true
	}
	return r.Item.As(target)
}

func (r *Row) ID() string {
	return r.id
}

func (r *Row) Finish() error {
	return r.RootC().FinishInContext(r)
}
