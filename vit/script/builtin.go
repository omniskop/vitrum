package script

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/omniskop/vitrum/vit/vcolor"
)

var builtinFunctions = map[string]interface{}{
	"rgb": func(r, g, b *int) string {
		if r == nil || g == nil || b == nil {
			panic(fmt.Errorf("Vit.rgb: too few arguments"))
		}
		return vcolor.RGBAToHex(uint8(*r), uint8(*g), uint8(*b), 255)
	},
	"rgba": func(r, g, b, a *int) string {
		if r == nil || g == nil || b == nil {
			panic(fmt.Errorf("Vit.rgba: too few arguments"))
		}
		if a == nil {
			return vcolor.RGBAToHex(uint8(*r), uint8(*g), uint8(*b), 255)
		} else {
			return vcolor.RGBAToHex(uint8(*r), uint8(*g), uint8(*b), uint8(*a))
		}
	},
	// TODO: maybe implement hsva, hsla and tint
	"darker": func(color string, f *float64) string {
		if f == nil {
			d := 2.0
			f = &d
		}
		c, err := vcolor.String(color)
		if err != nil {
			panic(fmt.Errorf("Vit.darker: %v", err))
		}
		col, a := vcolor.ToColorful(c)
		h, s, v := col.Hsv()
		v /= *f
		return vcolor.ColorfulToHex(colorful.Hsv(h, s, v), a)
	},
	"lighter": func(color string, f *float64) string {
		if f == nil {
			d := 1.5
			f = &d
		}
		c, err := vcolor.String(color)
		if err != nil {
			panic(fmt.Errorf("Vit.darker: %v", err))
		}
		col, a := vcolor.ToColorful(c)
		h, s, v := col.Hsv()
		fmt.Println(h, s, v)
		v *= *f
		return vcolor.ColorfulToHex(colorful.Hsv(h, s, v), a)
	},
}
