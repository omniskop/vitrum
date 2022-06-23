package vit

import (
	"fmt"
)

/*
	Layout hierarchy: (from most important to least)
	- Fill through 'fill'
	- Centering through 'centerIn'
	- Positioning through 'left', 'right', 'top', 'bottom', 'horizontalCenter', 'verticalCenter'
	- Sizing through 'width', 'height' and positioning through 'x', 'y' and 'z'

	Margins:
	- are applicable to 'fill', but not 'centerIn'
	- A specific margin (through 'leftMargin', 'rightMargin', 'topMargin', 'bottomMargin') overrides the general 'margins'
*/

type AnchorLine byte

const (
	AnchorsNone AnchorLine = iota
	AnchorLeft
	AnchorRight
	AnchorTop
	AnchorBottom
	AnchorHorizontalCenter
	AnchorVerticalCenter
)

func (l AnchorLine) PropertyName() string {
	switch l {
	case AnchorsNone:
		fmt.Println("AnchorLine zero value is being used!")
		return "<AnchorsNone>"
	case AnchorLeft:
		return "left"
	case AnchorRight:
		return "right"
	case AnchorTop:
		return "top"
	case AnchorBottom:
		return "bottom"
	case AnchorHorizontalCenter:
		return "horizontalCenter"
	case AnchorVerticalCenter:
		return "verticalCenter"
	default:
		fmt.Println("Unknown AnchorLine")
		return "<Unknown AnchorLine>"
	}
}

type AnchorsValue struct {
	Expression
	OnChange               func()
	AlignWhenCentered      BoolValue
	BaselineOffset         OptionalValue[*FloatValue]
	Baseline               OptionalValue[*FloatValue]
	Bottom                 OptionalValue[*FloatValue]
	BottomMargin           OptionalValue[*FloatValue]
	CenterIn               ComponentRefValue
	Fill                   ComponentRefValue
	HorizontalCenter       OptionalValue[*FloatValue]
	HorizontalCenterOffset OptionalValue[*FloatValue]
	Left                   OptionalValue[*FloatValue]
	LeftMargin             OptionalValue[*FloatValue]
	Margins                OptionalValue[*FloatValue]
	Right                  OptionalValue[*FloatValue]
	RightMargin            OptionalValue[*FloatValue]
	Top                    OptionalValue[*FloatValue]
	TopMargin              OptionalValue[*FloatValue]
	VerticalCenter         OptionalValue[*FloatValue]
	VerticalCenterOffset   OptionalValue[*FloatValue]
}

func NewAnchors() *AnchorsValue {
	return &AnchorsValue{
		AlignWhenCentered:      *NewBoolValue(true),
		Baseline:               *NewOptionalValue(NewEmptyFloatValue()),
		BaselineOffset:         *NewOptionalValue(NewEmptyFloatValue()),
		Bottom:                 *NewOptionalValue(NewEmptyFloatValue()),
		BottomMargin:           *NewOptionalValue(NewEmptyFloatValue()),
		CenterIn:               *NewEmptyComponentRefValue(),
		Fill:                   *NewEmptyComponentRefValue(),
		HorizontalCenter:       *NewOptionalValue(NewEmptyFloatValue()),
		HorizontalCenterOffset: *NewOptionalValue(NewEmptyFloatValue()),
		Left:                   *NewOptionalValue(NewEmptyFloatValue()),
		LeftMargin:             *NewOptionalValue(NewEmptyFloatValue()),
		Margins:                *NewOptionalValue(NewEmptyFloatValue()),
		Right:                  *NewOptionalValue(NewEmptyFloatValue()),
		RightMargin:            *NewOptionalValue(NewEmptyFloatValue()),
		Top:                    *NewOptionalValue(NewEmptyFloatValue()),
		TopMargin:              *NewOptionalValue(NewEmptyFloatValue()),
		VerticalCenter:         *NewOptionalValue(NewEmptyFloatValue()),
		VerticalCenterOffset:   *NewOptionalValue(NewEmptyFloatValue()),
	}
}

func (a *AnchorsValue) Property(key string) (interface{}, bool) {
	switch key {
	case "alignWhenCentered":
		return a.AlignWhenCentered, true
	case "baseline":
		return a.Baseline, true
	case "baselineOffset":
		return a.BaselineOffset, true
	case "bottom":
		return a.Bottom, true
	case "bottomMargin":
		return a.BottomMargin, true
	case "centerIn":
		return a.CenterIn, true
	case "fill":
		return a.Fill, true
	case "horizonzalCenter":
		return a.HorizontalCenter, true
	case "horizonzalCenterOffset":
		return a.HorizontalCenterOffset, true
	case "left":
		return a.Left, true
	case "leftMargin":
		return a.LeftMargin, true
	case "margins":
		return a.Margins, true
	case "right":
		return a.Right, true
	case "rightMargin":
		return a.RightMargin, true
	case "top":
		return a.Top, true
	case "topMargin":
		return a.TopMargin, true
	case "verticalCenter":
		return a.VerticalCenter, true
	case "verticalCenterOffset":
		return a.VerticalCenterOffset, true
	}
	return nil, false
}

func (a *AnchorsValue) SetProperty(key string, expression string, position *PositionRange) bool {
	fmt.Printf("[anchors] setting property %q\n", key)
	switch key {
	case "alignWhenCentered":
		a.AlignWhenCentered.SetExpression(expression, position)
	case "baseline":
		a.Baseline.SetExpression(expression, position)
	case "baselineOffset":
		a.BaselineOffset.SetExpression(expression, position)
	case "bottom":
		a.Bottom.SetExpression(expression, position)
	case "bottomMargin":
		a.BottomMargin.SetExpression(expression, position)
	case "centerIn":
		a.CenterIn.SetExpression(expression, position)
	case "fill":
		a.Fill.SetExpression(expression, position)
	case "horizonzalCenter":
		a.HorizontalCenter.SetExpression(expression, position)
	case "horizonzalCenterOffset":
		a.HorizontalCenterOffset.SetExpression(expression, position)
	case "left":
		a.Left.SetExpression(expression, position)
	case "leftMargin":
		a.LeftMargin.SetExpression(expression, position)
	case "margins":
		a.Margins.SetExpression(expression, position)
	case "right":
		a.Right.SetExpression(expression, position)
	case "rightMargin":
		a.RightMargin.SetExpression(expression, position)
	case "top":
		a.Top.SetExpression(expression, position)
	case "topMargin":
		a.TopMargin.SetExpression(expression, position)
	case "verticalCenter":
		a.VerticalCenter.SetExpression(expression, position)
	case "verticalCenterOffset":
		a.VerticalCenterOffset.SetExpression(expression, position)
	default:
		return false
	}
	// if a.OnChange != nil {
	// 	a.OnChange()
	// }
	return true
}

func (v *AnchorsValue) UpdateExpressions(context Component) (int, ErrorGroup) {
	var errs ErrorGroup
	var sum int
	if changed, err := v.AlignWhenCentered.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.alignWhenCentered", context.ID(), err))
		}
	}
	if changed, err := v.Baseline.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.baseline", context.ID(), err))
		}
	}
	if changed, err := v.BaselineOffset.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.baselineOffset", context.ID(), err))
		}
	}
	if changed, err := v.Bottom.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.bottom", context.ID(), err))
		}
	}
	if changed, err := v.BottomMargin.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.bottomMargin", context.ID(), err))
		}
	}
	if changed, err := v.CenterIn.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.centerIn", context.ID(), err))
		}
	}
	if changed, err := v.Fill.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.fill", context.ID(), err))
		}
	}
	if changed, err := v.HorizontalCenter.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.horizonzalCenter", context.ID(), err))
		}
	}
	if changed, err := v.HorizontalCenterOffset.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.horizonzalCenterOffset", context.ID(), err))
		}
	}
	if changed, err := v.Left.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.left", context.ID(), err))
		}
	}
	if changed, err := v.LeftMargin.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.leftMargin", context.ID(), err))
		}
	}
	if changed, err := v.Margins.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.margins", context.ID(), err))
		}
	}
	if changed, err := v.Right.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.right", context.ID(), err))
		}
	}
	if changed, err := v.RightMargin.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.rightMargin", context.ID(), err))
		}
	}
	if changed, err := v.Top.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.top", context.ID(), err))
		}
	}
	if changed, err := v.TopMargin.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.topMargin", context.ID(), err))
		}
	}
	if changed, err := v.VerticalCenter.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.verticalCenter", context.ID(), err))
		}
	}
	if changed, err := v.VerticalCenterOffset.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(NewPropertyError("Item", "anchors.verticalCenterOffset", context.ID(), err))
		}
	}

	return sum, errs
}

func (v *AnchorsValue) CalcTopMargin() float64 {
	if v.TopMargin.IsSet() {
		return v.TopMargin.Value().GetValue().(float64)
	}
	if v.Margins.IsSet() {
		return v.Margins.Value().GetValue().(float64)
	}
	return 0
}

func (v *AnchorsValue) CalcRightMargin() float64 {
	if v.RightMargin.IsSet() {
		return v.RightMargin.Value().GetValue().(float64)
	}
	if v.Margins.IsSet() {
		return v.Margins.Value().GetValue().(float64)
	}
	return 0
}

func (v *AnchorsValue) CalcBottomMargin() float64 {
	if v.BottomMargin.IsSet() {
		return v.BottomMargin.Value().GetValue().(float64)
	}
	if v.Margins.IsSet() {
		return v.Margins.Value().GetValue().(float64)
	}
	return 0
}

func (v *AnchorsValue) CalcLeftMargin() float64 {
	if v.LeftMargin.IsSet() {
		return v.LeftMargin.Value().GetValue().(float64)
	}
	if v.Margins.IsSet() {
		return v.Margins.Value().GetValue().(float64)
	}
	return 0
}

func (v *AnchorsValue) SetFromProperty(def PropertyDefinition) {
	panic("not implemented")
}

func (v *AnchorsValue) Update(context Component) (bool, error) {
	panic("not implemented")
}

func (v *AnchorsValue) SetValue(value interface{}) error {
	panic("not implemented")
	return nil
}

func (v *AnchorsValue) GetValue() interface{} {
	panic("not implemented")
	return nil
}

func (v *AnchorsValue) AddDependent(dep Dependent) {
	panic("not implemented")
}

func (v *AnchorsValue) RemoveDependent(dep Dependent) {
	panic("not implemented")
}

func (v *AnchorsValue) SetExpression(code string, position *PositionRange) {
	panic("not implemented")
}

// ====================================== Anchor Line Value ========================================

type AnchorLineValue struct {
	baseValue
	changed bool
	source  *AnchorLineValue
	offset  float64
}

var _ Value = &AnchorLineValue{}

func NewAnchorLineValue() *AnchorLineValue {
	return &AnchorLineValue{
		baseValue: newBaseValue(),
		changed:   true,
	}
}

func (v *AnchorLineValue) GetValue() interface{} {
	return v.Float64()
}

func (v *AnchorLineValue) Float64() float64 {
	if v.source == nil {
		return v.offset
	}
	return v.source.Float64() + v.offset
}

func (v *AnchorLineValue) SetFromProperty(def PropertyDefinition) {
	panic("not implemented")
}

func (v *AnchorLineValue) SetValue(newValue interface{}) error {
	floatVal, ok := castFloat64(newValue)
	if ok {
		v.SetAbsolute(floatVal)
		return nil
	}
	return newTypeError("number", newValue)
}

func (v *AnchorLineValue) SetOffset(offset float64) {
	if v.offset != offset {
		v.offset = offset
		v.changed = true
	}
}

func (v *AnchorLineValue) SetAbsolute(value float64) {
	if v.source != nil {
		v.source.RemoveDependent(v)
		v.source = nil
		v.changed = true
	}
	if v.offset != value {
		v.offset = value
		v.changed = true
	}
}

func (v *AnchorLineValue) SetExpression(code string, pos *PositionRange) {
	panic("not implemented")
}

func (v *AnchorLineValue) AssignTo(comp Component, lineType AnchorLine) {
	if v.source != nil {
		v.source.RemoveDependent(v)
		v.source = nil
		v.changed = true
	}
	if comp == nil {
		return
	}
	sourceValue, ok := comp.Property(lineType.PropertyName())
	if !ok {
		fmt.Println("tried to assign anchor line to component without anchors")
		return
	}
	v.source = sourceValue.(*AnchorLineValue)
	v.source.AddDependent(v)
	v.offset = 0
	v.changed = true
}

func (v *AnchorLineValue) IsAbsolute() bool {
	return v.source == nil
}

func (v *AnchorLineValue) Offset() float64 {
	return v.offset
}

func (v *AnchorLineValue) Update(context Component) (bool, error) {
	changed := false
	v.changed = false
	return changed, nil
}

func (v *AnchorLineValue) MakeDirty(stack []Dependent) {
	v.changed = true
	v.notifyDependents(append(stack, v))
}
