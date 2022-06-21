// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
)

type Grid_HorizontalItemAlignment uint

const (
	Grid_HorizontalItemAlignment_AlignLeft    Grid_HorizontalItemAlignment = 0
	Grid_HorizontalItemAlignment_AlignHCenter Grid_HorizontalItemAlignment = 1
	Grid_HorizontalItemAlignment_AlignRight   Grid_HorizontalItemAlignment = 2
)

type Grid_VerticalItemAlignment uint

const (
	Grid_VerticalItemAlignment_AlignTop     Grid_VerticalItemAlignment = 0
	Grid_VerticalItemAlignment_AlignVCenter Grid_VerticalItemAlignment = 1
	Grid_VerticalItemAlignment_AlignBottom  Grid_VerticalItemAlignment = 2
)

type Grid_Flow uint

const (
	Grid_Flow_LeftToRight Grid_Flow = 0
	Grid_Flow_TopToBottom Grid_Flow = 1
)

type Grid struct {
	Item
	id string

	topPadding              vit.OptionalValue[*vit.FloatValue]
	rightPadding            vit.OptionalValue[*vit.FloatValue]
	bottomPadding           vit.OptionalValue[*vit.FloatValue]
	leftPadding             vit.OptionalValue[*vit.FloatValue]
	padding                 vit.FloatValue
	spacing                 vit.FloatValue
	columnSpacing           vit.OptionalValue[*vit.FloatValue]
	rowSpacing              vit.OptionalValue[*vit.FloatValue]
	columns                 vit.OptionalValue[*vit.IntValue]
	rows                    vit.OptionalValue[*vit.IntValue]
	horizontalItemAlignment vit.IntValue
	verticalItemAlignment   vit.IntValue
	flow                    vit.IntValue
	childLayouts            vit.LayoutList
}

func NewGrid(id string, scope vit.ComponentContainer) *Grid {
	g := &Grid{
		Item:                    *NewItem(id, scope),
		id:                      id,
		topPadding:              *vit.NewOptionalValue(vit.NewFloatValueFromExpression("0", nil)),
		rightPadding:            *vit.NewOptionalValue(vit.NewFloatValueFromExpression("0", nil)),
		bottomPadding:           *vit.NewOptionalValue(vit.NewFloatValueFromExpression("0", nil)),
		leftPadding:             *vit.NewOptionalValue(vit.NewFloatValueFromExpression("0", nil)),
		padding:                 *vit.NewFloatValueFromExpression("0", nil),
		spacing:                 *vit.NewFloatValueFromExpression("0", nil),
		columnSpacing:           *vit.NewOptionalValue(vit.NewFloatValueFromExpression("0", nil)),
		rowSpacing:              *vit.NewOptionalValue(vit.NewFloatValueFromExpression("0", nil)),
		columns:                 *vit.NewOptionalValue(vit.NewEmptyIntValue()),
		rows:                    *vit.NewOptionalValue(vit.NewEmptyIntValue()),
		horizontalItemAlignment: *vit.NewIntValueFromExpression("HorizontalItemAlignment.AlignLeft", nil),
		verticalItemAlignment:   *vit.NewIntValueFromExpression("VerticalItemAlignment.AlignTop", nil),
		flow:                    *vit.NewIntValueFromExpression("Flow.LeftToRight", nil),
		childLayouts:            make(vit.LayoutList),
	}
	g.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "HorizontalItemAlignment",
		Position: nil,
		Values:   map[string]int{"AlignRight": 2, "AlignLeft": 0, "AlignHCenter": 1},
	})
	g.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "VerticalItemAlignment",
		Position: nil,
		Values:   map[string]int{"AlignVCenter": 1, "AlignBottom": 2, "AlignTop": 0},
	})
	g.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "Flow",
		Position: nil,
		Values:   map[string]int{"TopToBottom": 1, "LeftToRight": 0},
	})
	return g
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid(%s)", g.id)
}

func (g *Grid) Property(key string) (vit.Value, bool) {
	switch key {
	case "topPadding":
		return &g.topPadding, true
	case "rightPadding":
		return &g.rightPadding, true
	case "bottomPadding":
		return &g.bottomPadding, true
	case "leftPadding":
		return &g.leftPadding, true
	case "padding":
		return &g.padding, true
	case "spacing":
		return &g.spacing, true
	case "columnSpacing":
		return &g.columnSpacing, true
	case "rowSpacing":
		return &g.rowSpacing, true
	case "columns":
		return &g.columns, true
	case "rows":
		return &g.rows, true
	case "horizontalItemAlignment":
		return &g.horizontalItemAlignment, true
	case "verticalItemAlignment":
		return &g.verticalItemAlignment, true
	case "flow":
		return &g.flow, true
	default:
		return g.Item.Property(key)
	}
}

func (g *Grid) MustProperty(key string) vit.Value {
	v, ok := g.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (g *Grid) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "topPadding":
		err = g.topPadding.SetValue(value)
	case "rightPadding":
		err = g.rightPadding.SetValue(value)
	case "bottomPadding":
		err = g.bottomPadding.SetValue(value)
	case "leftPadding":
		err = g.leftPadding.SetValue(value)
	case "padding":
		err = g.padding.SetValue(value)
	case "spacing":
		err = g.spacing.SetValue(value)
	case "columnSpacing":
		err = g.columnSpacing.SetValue(value)
	case "rowSpacing":
		err = g.rowSpacing.SetValue(value)
	case "columns":
		err = g.columns.SetValue(value)
	case "rows":
		err = g.rows.SetValue(value)
	case "horizontalItemAlignment":
		err = g.horizontalItemAlignment.SetValue(value)
	case "verticalItemAlignment":
		err = g.verticalItemAlignment.SetValue(value)
	case "flow":
		err = g.flow.SetValue(value)
	default:
		return g.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("Grid", key, g.id, err)
	}
	return nil
}

func (g *Grid) SetPropertyExpression(key string, code string, pos *vit.PositionRange) error {
	switch key {
	case "topPadding":
		g.topPadding.SetExpression(code, pos)
	case "rightPadding":
		g.rightPadding.SetExpression(code, pos)
	case "bottomPadding":
		g.bottomPadding.SetExpression(code, pos)
	case "leftPadding":
		g.leftPadding.SetExpression(code, pos)
	case "padding":
		g.padding.SetExpression(code, pos)
	case "spacing":
		g.spacing.SetExpression(code, pos)
	case "columnSpacing":
		g.columnSpacing.SetExpression(code, pos)
	case "rowSpacing":
		g.rowSpacing.SetExpression(code, pos)
	case "columns":
		g.columns.SetExpression(code, pos)
	case "rows":
		g.rows.SetExpression(code, pos)
	case "horizontalItemAlignment":
		g.horizontalItemAlignment.SetExpression(code, pos)
	case "verticalItemAlignment":
		g.verticalItemAlignment.SetExpression(code, pos)
	case "flow":
		g.flow.SetExpression(code, pos)
	default:
		return g.Item.SetPropertyExpression(key, code, pos)
	}
	return nil
}

func (g *Grid) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case g.id:
		return g, true
	case "topPadding":
		return &g.topPadding, true
	case "rightPadding":
		return &g.rightPadding, true
	case "bottomPadding":
		return &g.bottomPadding, true
	case "leftPadding":
		return &g.leftPadding, true
	case "padding":
		return &g.padding, true
	case "spacing":
		return &g.spacing, true
	case "columnSpacing":
		return &g.columnSpacing, true
	case "rowSpacing":
		return &g.rowSpacing, true
	case "columns":
		return &g.columns, true
	case "rows":
		return &g.rows, true
	case "horizontalItemAlignment":
		return &g.horizontalItemAlignment, true
	case "verticalItemAlignment":
		return &g.verticalItemAlignment, true
	case "flow":
		return &g.flow, true
	default:
		return g.Item.ResolveVariable(key)
	}
}

func (g *Grid) AddChild(child vit.Component) {
	defer g.childWasAdded(child)
	child.SetParent(g)
	g.AddChildButKeepParent(child)
}

func (g *Grid) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	defer g.childWasAdded(addThis)
	var targetType vit.Component = afterThis

	for ind, child := range g.Children() {
		if child.As(&targetType) {
			addThis.SetParent(g)
			g.AddChildAtButKeepParent(addThis, ind+1)
			return
		}
	}
	g.AddChild(addThis)
}

func (g *Grid) UpdateExpressions() (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if changed, err := g.topPadding.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "topPadding", g.id, err))
		}
		g.recalculateLayout(g.topPadding)
	}
	if changed, err := g.rightPadding.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "rightPadding", g.id, err))
		}
		g.recalculateLayout(g.rightPadding)
	}
	if changed, err := g.bottomPadding.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "bottomPadding", g.id, err))
		}
		g.recalculateLayout(g.bottomPadding)
	}
	if changed, err := g.leftPadding.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "leftPadding", g.id, err))
		}
		g.recalculateLayout(g.leftPadding)
	}
	if changed, err := g.padding.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "padding", g.id, err))
		}
		g.recalculateLayout(g.padding)
	}
	if changed, err := g.spacing.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "spacing", g.id, err))
		}
		g.recalculateLayout(g.spacing)
	}
	if changed, err := g.columnSpacing.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "columnSpacing", g.id, err))
		}
		g.recalculateLayout(g.columnSpacing)
	}
	if changed, err := g.rowSpacing.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "rowSpacing", g.id, err))
		}
		g.recalculateLayout(g.rowSpacing)
	}
	if changed, err := g.columns.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "columns", g.id, err))
		}
		g.recalculateLayout(g.columns)
	}
	if changed, err := g.rows.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "rows", g.id, err))
		}
		g.recalculateLayout(g.rows)
	}
	if changed, err := g.horizontalItemAlignment.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "horizontalItemAlignment", g.id, err))
		}
		g.recalculateLayout(g.horizontalItemAlignment)
	}
	if changed, err := g.verticalItemAlignment.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "verticalItemAlignment", g.id, err))
		}
		g.recalculateLayout(g.verticalItemAlignment)
	}
	if changed, err := g.flow.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "flow", g.id, err))
		}
		g.recalculateLayout(g.flow)
	}
	if changed, err := g.childLayouts.Update(g); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Grid", "childLayouts", g.id, err))
		}
		g.recalculateLayout(g.childLayouts)
	}

	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	n, err := g.UpdatePropertiesInContext(g)
	sum += n
	errs.AddGroup(err)
	n, err = g.Item.UpdateExpressions()
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (g *Grid) As(target *vit.Component) bool {
	if _, ok := (*target).(*Grid); ok {
		*target = g
		return true
	}
	return g.Item.As(target)
}

func (g *Grid) ID() string {
	return g.id
}

func (g *Grid) Finish() error {
	return g.RootC().FinishInContext(g)
}

func (g *Grid) staticAttribute(name string) (interface{}, bool) {
	switch name {
	case "AlignLeft":
		return uint(Grid_HorizontalItemAlignment_AlignLeft), true
	case "AlignHCenter":
		return uint(Grid_HorizontalItemAlignment_AlignHCenter), true
	case "AlignRight":
		return uint(Grid_HorizontalItemAlignment_AlignRight), true
	case "AlignTop":
		return uint(Grid_VerticalItemAlignment_AlignTop), true
	case "AlignVCenter":
		return uint(Grid_VerticalItemAlignment_AlignVCenter), true
	case "AlignBottom":
		return uint(Grid_VerticalItemAlignment_AlignBottom), true
	case "LeftToRight":
		return uint(Grid_Flow_LeftToRight), true
	case "TopToBottom":
		return uint(Grid_Flow_TopToBottom), true
	default:
		return nil, false
	}
}
