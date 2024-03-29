// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
	parse "github.com/omniskop/vitrum/vit/parse"
)

func newFileContextForColumn(globalCtx *vit.GlobalContext) (*vit.FileContext, error) {
	return vit.NewFileContext(globalCtx), nil
}

type Column struct {
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

// newColumnInGlobal creates an appropriate file context for the component and then returns a new Column instance.
// The returned error will only be set if a library import that is required by the component fails.
func newColumnInGlobal(id string, globalCtx *vit.GlobalContext, thisLibrary parse.Library) (*Column, error) {
	fileCtx, err := newFileContextForColumn(globalCtx)
	if err != nil {
		return nil, err
	}
	parse.AddLibraryToContainer(thisLibrary, &fileCtx.KnownComponents)
	return NewColumn(id, fileCtx), nil
}
func NewColumn(id string, context *vit.FileContext) *Column {
	c := &Column{
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
	c.topPadding.AddDependent(vit.FuncDep(c.recalculateLayout))
	c.rightPadding.AddDependent(vit.FuncDep(c.recalculateLayout))
	c.bottomPadding.AddDependent(vit.FuncDep(c.recalculateLayout))
	c.leftPadding.AddDependent(vit.FuncDep(c.recalculateLayout))
	c.padding.AddDependent(vit.FuncDep(c.recalculateLayout))
	c.spacing.AddDependent(vit.FuncDep(c.recalculateLayout))
	c.Item.AddBoundsDependency(vit.FuncDep(c.recalculateLayout))
	// register event listeners
	// register enumerations
	// add child components

	context.RegisterComponent("", c)

	return c
}

func (c *Column) String() string {
	return fmt.Sprintf("Column(%s)", c.id)
}

func (c *Column) Property(key string) (vit.Value, bool) {
	switch key {
	case "topPadding":
		return &c.topPadding, true
	case "rightPadding":
		return &c.rightPadding, true
	case "bottomPadding":
		return &c.bottomPadding, true
	case "leftPadding":
		return &c.leftPadding, true
	case "padding":
		return &c.padding, true
	case "spacing":
		return &c.spacing, true
	default:
		return c.Item.Property(key)
	}
}

func (c *Column) MustProperty(key string) vit.Value {
	v, ok := c.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (c *Column) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "topPadding":
		err = c.topPadding.SetValue(value)
	case "rightPadding":
		err = c.rightPadding.SetValue(value)
	case "bottomPadding":
		err = c.bottomPadding.SetValue(value)
	case "leftPadding":
		err = c.leftPadding.SetValue(value)
	case "padding":
		err = c.padding.SetValue(value)
	case "spacing":
		err = c.spacing.SetValue(value)
	default:
		return c.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("Column", key, c.id, err)
	}
	return nil
}

func (c *Column) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "topPadding":
		c.topPadding.SetCode(code)
	case "rightPadding":
		c.rightPadding.SetCode(code)
	case "bottomPadding":
		c.bottomPadding.SetCode(code)
	case "leftPadding":
		c.leftPadding.SetCode(code)
	case "padding":
		c.padding.SetCode(code)
	case "spacing":
		c.spacing.SetCode(code)
	default:
		return c.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (c *Column) Event(name string) (vit.Listenable, bool) {
	switch name {
	default:
		return c.Item.Event(name)
	}
}

func (c *Column) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "topPadding":
		return &c.topPadding, true
	case "rightPadding":
		return &c.rightPadding, true
	case "bottomPadding":
		return &c.bottomPadding, true
	case "leftPadding":
		return &c.leftPadding, true
	case "padding":
		return &c.padding, true
	case "spacing":
		return &c.spacing, true
	default:
		return c.Item.ResolveVariable(key)
	}
}

func (c *Column) AddChild(child vit.Component) {
	defer c.childWasAdded(child)
	child.SetParent(c)
	c.AddChildButKeepParent(child)
}

func (c *Column) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	defer c.childWasAdded(addThis)
	var targetType vit.Component = afterThis

	for ind, child := range c.Children() {
		if child.As(&targetType) {
			addThis.SetParent(c)
			c.AddChildAtButKeepParent(addThis, ind+1)
			return
		}
	}
	c.AddChild(addThis)
}

func (c *Column) UpdateExpressions(context vit.Component) (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if context == nil {
		context = c
	}
	// properties
	if changed, err := c.topPadding.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Column", "topPadding", c.id, err))
		}
	}
	if changed, err := c.rightPadding.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Column", "rightPadding", c.id, err))
		}
	}
	if changed, err := c.bottomPadding.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Column", "bottomPadding", c.id, err))
		}
	}
	if changed, err := c.leftPadding.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Column", "leftPadding", c.id, err))
		}
	}
	if changed, err := c.padding.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Column", "padding", c.id, err))
		}
	}
	if changed, err := c.spacing.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Column", "spacing", c.id, err))
		}
	}

	// methods

	n, err := c.Item.UpdateExpressions(context)
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (c *Column) As(target *vit.Component) bool {
	if _, ok := (*target).(*Column); ok {
		*target = c
		return true
	}
	return c.Item.As(target)
}

func (c *Column) ID() string {
	return c.id
}

func (c *Column) Finish() error {
	return c.RootC().FinishInContext(c)
}
