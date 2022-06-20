package vit

import (
	"image/color"

	"github.com/omniskop/vitrum/vit/vcolor"
)

type ColorValue struct {
	baseValue
	value      color.Color
	expression *Expression
}

func NewColorValue(expression string, position *PositionRange) *ColorValue {
	return &ColorValue{
		baseValue:  newBaseValue(),
		value:      color.Black,
		expression: NewExpression(expression, position),
	}
}

func NewEmptyColorValue() *ColorValue {
	return &ColorValue{
		baseValue:  newBaseValue(),
		value:      color.Black,
		expression: nil,
	}
}

func (v *ColorValue) GetValue() interface{} {
	return v.value
}

func (v *ColorValue) RGBAColor() color.RGBA {
	r, g, b, a := v.value.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}

func (v *ColorValue) Color() color.Color {
	return v.value
}

func (v *ColorValue) SetFromProperty(prop PropertyDefinition) {
	v.expression = NewExpression(prop.Expression, &prop.Pos)
	v.notifyDependents([]Dependent{v.expression})
}

func (v *ColorValue) SetValue(newValue interface{}) error {
	if colVal, ok := newValue.(color.Color); ok {
		v.value = colVal
		v.expression = nil
		v.notifyDependents(nil)
		return nil
	}
	return newTypeError("color", newValue)
}

func (v *ColorValue) SetColor(newValue color.Color) {
	v.value = newValue
	v.expression = nil
	v.notifyDependents(nil)
}

func (v *ColorValue) SetExpression(code string, pos *PositionRange) {
	v.expression = NewExpression(code, pos)
	v.notifyDependents([]Dependent{v.expression})
}

func (v *ColorValue) Update(context Component) (bool, error) {
	if v.expression == nil {
		return false, nil
	}
	if !v.expression.ShouldEvaluate() {
		return false, nil
	}
	val, err := v.expression.Evaluate(context)
	if err != nil {
		return false, err
	}

	switch result := val.(type) {
	case string:
		c, err := vcolor.String(result)
		if err != nil {
			v.value = color.Black
			return false, err
		}
		v.value = c
	default:
		return false, newTypeError("color string", result)
	}

	return true, nil
}
