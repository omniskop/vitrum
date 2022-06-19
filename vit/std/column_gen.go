// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
)

type Column struct {
	Item
	id string

	topPadding    vit.OptionalValue[*vit.FloatValue]
	rightPadding  vit.OptionalValue[*vit.FloatValue]
	bottomPadding vit.OptionalValue[*vit.FloatValue]
	leftPadding   vit.OptionalValue[*vit.FloatValue]
	padding       vit.FloatValue
	spacing       vit.FloatValue
	childLayouts  vit.LayoutList
}

func NewColumn(id string, scope vit.ComponentContainer) *Column {
	return &Column{
		Item:          *NewItem(id, scope),
		id:            id,
		topPadding:    *vit.NewOptionalValue(vit.NewFloatValue("", nil)),
		rightPadding:  *vit.NewOptionalValue(vit.NewFloatValue("", nil)),
		bottomPadding: *vit.NewOptionalValue(vit.NewFloatValue("", nil)),
		leftPadding:   *vit.NewOptionalValue(vit.NewFloatValue("", nil)),
		padding:       *vit.NewFloatValue("", nil),
		spacing:       *vit.NewFloatValue("", nil),
		childLayouts:  make(vit.LayoutList),
	}
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

func (c *Column) SetProperty(key string, value interface{}, position *vit.PositionRange) bool {
	switch key {
	case "topPadding":
		c.topPadding.ChangeCode(value.(string), position)
	case "rightPadding":
		c.rightPadding.ChangeCode(value.(string), position)
	case "bottomPadding":
		c.bottomPadding.ChangeCode(value.(string), position)
	case "leftPadding":
		c.leftPadding.ChangeCode(value.(string), position)
	case "padding":
		c.padding.ChangeCode(value.(string), position)
	case "spacing":
		c.spacing.ChangeCode(value.(string), position)
	default:
		return c.Item.SetProperty(key, value, position)
	}
	return true
}

func (c *Column) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case c.id:
		return c, true
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

func (c *Column) UpdateExpressions() (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if c.topPadding.ShouldEvaluate() {
		sum++
		err := c.topPadding.Update(c)
		if err != nil {
			errs.Add(vit.NewExpressionError("Column", "topPadding", c.id, *c.topPadding.GetExpression(), err))
		}
		c.recalculateLayout(c.topPadding)
	}
	if c.rightPadding.ShouldEvaluate() {
		sum++
		err := c.rightPadding.Update(c)
		if err != nil {
			errs.Add(vit.NewExpressionError("Column", "rightPadding", c.id, *c.rightPadding.GetExpression(), err))
		}
		c.recalculateLayout(c.rightPadding)
	}
	if c.bottomPadding.ShouldEvaluate() {
		sum++
		err := c.bottomPadding.Update(c)
		if err != nil {
			errs.Add(vit.NewExpressionError("Column", "bottomPadding", c.id, *c.bottomPadding.GetExpression(), err))
		}
		c.recalculateLayout(c.bottomPadding)
	}
	if c.leftPadding.ShouldEvaluate() {
		sum++
		err := c.leftPadding.Update(c)
		if err != nil {
			errs.Add(vit.NewExpressionError("Column", "leftPadding", c.id, *c.leftPadding.GetExpression(), err))
		}
		c.recalculateLayout(c.leftPadding)
	}
	if c.padding.ShouldEvaluate() {
		sum++
		err := c.padding.Update(c)
		if err != nil {
			errs.Add(vit.NewExpressionError("Column", "padding", c.id, *c.padding.GetExpression(), err))
		}
		c.recalculateLayout(c.padding)
	}
	if c.spacing.ShouldEvaluate() {
		sum++
		err := c.spacing.Update(c)
		if err != nil {
			errs.Add(vit.NewExpressionError("Column", "spacing", c.id, *c.spacing.GetExpression(), err))
		}
		c.recalculateLayout(c.spacing)
	}
	if c.childLayouts.ShouldEvaluate() {
		sum++
		err := c.childLayouts.Update(c)
		if err != nil {
			errs.Add(vit.NewExpressionError("Column", "childLayouts", c.id, *c.childLayouts.GetExpression(), err))
		}
		c.recalculateLayout(c.childLayouts)
	}

	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	n, err := c.UpdatePropertiesInContext(c)
	sum += n
	errs.AddGroup(err)
	n, err = c.Item.UpdateExpressions()
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
