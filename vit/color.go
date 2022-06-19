package vit

import (
	"fmt"
	"image/color"

	"github.com/omniskop/vitrum/vit/vcolor"
)

type ColorValue struct {
	Expression
	value color.Color
}

func NewColorValue(expression string, position *PositionRange) *ColorValue {
	v := new(ColorValue)
	if expression == "" {
		v.Expression = *NewExpression(`"black"`, position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *ColorValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *ColorValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		return err
	}

	switch result := val.(type) {
	case string:
		c, err := vcolor.String(result)
		if err != nil {
			v.value = color.Black
			return err
		}
		v.value = c
	default:
		return fmt.Errorf("color expression evaluated to %T, expected string", result)
	}

	return nil
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
