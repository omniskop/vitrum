package vit

import (
	"fmt"
	"image/color"

	"github.com/omniskop/vitrum/vit/vcolor"
)

type ColorValue struct {
	Expression
	Value color.Color
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

func (v *ColorValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		return err
	}

	switch result := val.(type) {
	case string:
		c, err := vcolor.String(result)
		if err != nil {
			v.Value = color.Black
			return err
		}
		v.Value = c
	default:
		return fmt.Errorf("color expression evaluated to %T, expected string", result)
	}

	return nil
}

func (v *ColorValue) GetValue() interface{} {
	return v.Value
}
