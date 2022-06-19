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

func NewAnchors(expression string, position *PositionRange) *AnchorsValue {
	return &AnchorsValue{
		AlignWhenCentered: *NewBoolValue("true", nil),
		Baseline:          *NewOptionalValue(NewFloatValue("", nil)),
		BaselineOffset:    *NewOptionalValue(NewFloatValue("", nil)),
		Bottom:            *NewOptionalValue(NewFloatValue("", nil)),
		BottomMargin:      *NewOptionalValue(NewFloatValue("", nil)),
		// CenterIn:               *NewComponentValue("", nil),
		// Fill:                   *NewComponentValue("", nil),
		HorizontalCenter:       *NewOptionalValue(NewFloatValue("", nil)),
		HorizontalCenterOffset: *NewOptionalValue(NewFloatValue("", nil)),
		Left:                   *NewOptionalValue(NewFloatValue("", nil)),
		LeftMargin:             *NewOptionalValue(NewFloatValue("", nil)),
		Margins:                *NewOptionalValue(NewFloatValue("", nil)),
		Right:                  *NewOptionalValue(NewFloatValue("", nil)),
		RightMargin:            *NewOptionalValue(NewFloatValue("", nil)),
		Top:                    *NewOptionalValue(NewFloatValue("", nil)),
		TopMargin:              *NewOptionalValue(NewFloatValue("", nil)),
		VerticalCenter:         *NewOptionalValue(NewFloatValue("", nil)),
		VerticalCenterOffset:   *NewOptionalValue(NewFloatValue("", nil)),
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
	switch key {
	case "alignWhenCentered":
		a.AlignWhenCentered.ChangeCode(expression, position)
	case "baseline":
		a.Baseline.ChangeCode(expression, position)
	case "baselineOffset":
		a.BaselineOffset.ChangeCode(expression, position)
	case "bottom":
		a.Bottom.ChangeCode(expression, position)
	case "bottomMargin":
		a.BottomMargin.ChangeCode(expression, position)
	case "centerIn":
		a.CenterIn.ChangeCode(expression, position)
	case "fill":
		a.Fill.ChangeCode(expression, position)
	case "horizonzalCenter":
		a.HorizontalCenter.ChangeCode(expression, position)
	case "horizonzalCenterOffset":
		a.HorizontalCenterOffset.ChangeCode(expression, position)
	case "left":
		a.Left.ChangeCode(expression, position)
	case "leftMargin":
		a.LeftMargin.ChangeCode(expression, position)
	case "margins":
		a.Margins.ChangeCode(expression, position)
	case "right":
		a.Right.ChangeCode(expression, position)
	case "rightMargin":
		a.RightMargin.ChangeCode(expression, position)
	case "top":
		a.Top.ChangeCode(expression, position)
	case "topMargin":
		a.TopMargin.ChangeCode(expression, position)
	case "verticalCenter":
		a.VerticalCenter.ChangeCode(expression, position)
	case "verticalCenterOffset":
		a.VerticalCenterOffset.ChangeCode(expression, position)
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
	if v.AlignWhenCentered.ShouldEvaluate() {
		sum++
		err := v.AlignWhenCentered.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.alignWhenCentered", context.ID(), v.AlignWhenCentered.Expression, err))
		}
	}
	if v.Baseline.ShouldEvaluate() {
		sum++
		err := v.Baseline.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.baseline", context.ID(), v.Baseline.ActualValue().Expression, err))
		}
	}
	if v.BaselineOffset.ShouldEvaluate() {
		sum++
		err := v.BaselineOffset.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.baselineOffset", context.ID(), v.BaselineOffset.ActualValue().Expression, err))
		}
	}
	if v.Bottom.ShouldEvaluate() {
		sum++
		err := v.Bottom.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.bottom", context.ID(), v.Bottom.ActualValue().Expression, err))
		}
	}
	if v.BottomMargin.ShouldEvaluate() {
		sum++
		err := v.BottomMargin.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.bottomMargin", context.ID(), v.BottomMargin.ActualValue().Expression, err))
		}
	}
	if v.CenterIn.ShouldEvaluate() {
		sum++
		err := v.CenterIn.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.centerIn", context.ID(), v.CenterIn.Expression, err))
		}
	}
	if v.Fill.ShouldEvaluate() {
		sum++
		err := v.Fill.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.fill", context.ID(), v.Fill.Expression, err))
		}
	}
	if v.HorizontalCenter.ShouldEvaluate() {
		sum++
		err := v.HorizontalCenter.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.horizonzalCenter", context.ID(), v.HorizontalCenter.ActualValue().Expression, err))
		}
	}
	if v.HorizontalCenterOffset.ShouldEvaluate() {
		sum++
		err := v.HorizontalCenterOffset.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.horizonzalCenterOffset", context.ID(), v.HorizontalCenterOffset.ActualValue().Expression, err))
		}
	}
	if v.Left.ShouldEvaluate() {
		sum++
		err := v.Left.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.left", context.ID(), v.Left.ActualValue().Expression, err))
		}
	}
	if v.LeftMargin.ShouldEvaluate() {
		sum++
		err := v.LeftMargin.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.leftMargin", context.ID(), v.LeftMargin.ActualValue().Expression, err))
		}
	}
	if v.Margins.ShouldEvaluate() {
		sum++
		err := v.Margins.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.margins", context.ID(), v.Margins.ActualValue().Expression, err))
		}
	}
	if v.Right.ShouldEvaluate() {
		sum++
		err := v.Right.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.right", context.ID(), v.Right.ActualValue().Expression, err))
		}
	}
	if v.RightMargin.ShouldEvaluate() {
		sum++
		err := v.RightMargin.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.rightMargin", context.ID(), v.RightMargin.ActualValue().Expression, err))
		}
	}
	if v.Top.ShouldEvaluate() {
		sum++
		err := v.Top.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.top", context.ID(), v.Top.ActualValue().Expression, err))
		}
	}
	if v.TopMargin.ShouldEvaluate() {
		sum++
		err := v.TopMargin.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.topMargin", context.ID(), v.TopMargin.ActualValue().Expression, err))
		}
	}
	if v.VerticalCenter.ShouldEvaluate() {
		sum++
		err := v.VerticalCenter.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.verticalCenter", context.ID(), v.VerticalCenter.ActualValue().Expression, err))
		}
	}
	if v.VerticalCenterOffset.ShouldEvaluate() {
		sum++
		err := v.VerticalCenterOffset.Update(context)
		if err != nil {
			errs.Add(NewExpressionError("Item", "anchors.verticalCenterOffset", context.ID(), v.VerticalCenterOffset.ActualValue().Expression, err))
		}
	}

	return sum, errs
}

func (v *AnchorsValue) SetFromProperty(def PropertyDefinition) {
	panic("not implemented")
}

func (v *AnchorsValue) Update(context Component) error {
	panic("not implemented")
}

func (v *AnchorsValue) GetValue() interface{} {
	panic("not implemented")
	return nil
}

func (v *AnchorsValue) MakeDirty([]*Expression) {
	panic("not implemented")
}

func (v *AnchorsValue) GetExpression() *Expression {
	panic("not implemented")
	return nil
}

func (v *AnchorsValue) AddDependent(dep Dependent) {
	panic("not implemented")
}

func (v *AnchorsValue) RemoveDependent(dep Dependent) {
	panic("not implemented")
}

func (v *AnchorsValue) ShouldEvaluate() bool {
	panic("not implemented")
	return false
}

func (v *AnchorsValue) Err() error {
	panic("not implemented")
	return nil
}

// ====================================== Anchor Line Value ========================================

type AnchorLineValue struct {
	StaticBaseValue
	Source *AnchorLineValue
	offset float64
}

var _ Value = &AnchorLineValue{}

func NewAnchorLineValue() *AnchorLineValue {
	return &AnchorLineValue{
		StaticBaseValue: *NewStaticBaseValue(),
	}
}

func (v *AnchorLineValue) AssignTo(comp Component, lineType AnchorLine) {
	if v.Source != nil {
		v.Source.RemoveDependent(v)
		v.Source = nil
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
	v.Source = sourceValue.(*AnchorLineValue)
	v.Source.AddDependent(v)
	v.offset = 0
	v.changed = true
}

func (v *AnchorLineValue) SetOffset(offset float64) {
	if v.offset != offset {
		v.offset = offset
		v.changed = true
	}
}

func (v *AnchorLineValue) SetAbsolute(value float64) {
	if v.Source != nil {
		v.Source.RemoveDependent(v)
		v.Source = nil
		v.changed = true
	}
	if v.offset != value {
		v.offset = value
		v.changed = true
	}
}

func (v *AnchorLineValue) IsAbsolute() bool {
	return v.Source == nil
}

func (v *AnchorLineValue) Offset() float64 {
	return v.offset
}

func (v *AnchorLineValue) SetFromProperty(def PropertyDefinition) {
	panic("not implemented")
}

func (v *AnchorLineValue) GetValue() interface{} {
	return v.Float64()
}

func (v *AnchorLineValue) Update(context Component) error {
	v.StaticBaseValue.Update(context)
	// if v.Source != nil {
	// 	v.offset = v.Source.offset
	// }
	return nil
}

func (v *AnchorLineValue) Float64() float64 {
	if v.Source == nil {
		return v.offset
	}
	return v.Source.Float64() + v.offset
}

func (v *AnchorLineValue) GetExpression() *Expression {
	panic("not implemented")
	return nil
}

func (v *AnchorLineValue) Err() error {
	return nil
}
