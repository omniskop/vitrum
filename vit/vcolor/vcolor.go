package vcolor

import (
	"encoding/hex"
	"fmt"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

// CSS color names from "Basic color keywords" https://www.w3.org/TR/css-color-3/#html4
// plus some extra
var ColorNames = map[string]color.Color{
	"transparent": color.Transparent,
	"opaque":      color.Opaque,

	"black":   color.Black,
	"silver":  color.RGBA{192, 192, 192, 255},
	"gray":    color.RGBA{128, 128, 128, 255},
	"grey":    color.RGBA{128, 128, 128, 255},
	"white":   color.White,
	"maroon":  color.RGBA{128, 0, 0, 255},
	"red":     color.RGBA{255, 0, 0, 255},
	"purple":  color.RGBA{128, 0, 128, 255},
	"fuchsia": color.RGBA{255, 0, 255, 255},
	"green":   color.RGBA{0, 128, 0, 255},
	"lime":    color.RGBA{0, 255, 0, 255},
	"olive":   color.RGBA{128, 128, 0, 255},
	"yellow":  color.RGBA{255, 255, 0, 255},
	"navy":    color.RGBA{0, 0, 128, 255},
	"blue":    color.RGBA{0, 0, 255, 255},
	"teal":    color.RGBA{0, 128, 128, 255},
	"aqua":    color.RGBA{0, 255, 255, 255},
}

func StringRGBA(input string) (r, g, b, a byte, err error) {
	if input[0] == '#' {
		// hex rgb(a) color
		switch len(input) {
		case 4:
			var out [3]byte
			_, err := hex.Decode(out[0:3], []byte{input[1], input[1], input[2], input[2], input[3], input[3]})
			if err != nil {
				return 0, 0, 0, 0, err
			}
			return out[0], out[1], out[2], 255, nil
		case 7:
			out, err := hex.DecodeString(input[1:])
			if err != nil {
				return 0, 0, 0, 0, err
			}
			return out[0], out[1], out[2], 255, nil
		case 9:
			out, err := hex.DecodeString(input[1:])
			if err != nil {
				return 0, 0, 0, 0, err
			}
			return out[0], out[1], out[2], out[3], nil
		default:
			return 0, 0, 0, 0, fmt.Errorf("invalid color code %q", input)
		}
	} else {
		// css color name
		if c, ok := ColorNames[input]; ok {
			r, g, b, a := c.RGBA()
			return byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8), nil
		}
		return 0, 0, 0, 0, fmt.Errorf("unknown color %q", input)
	}
}

func String(input string) (color.Color, error) {
	if input[0] != '#' {
		// doing this through StringRGBA would cause unnecessary conversions
		if c, ok := ColorNames[input]; ok {
			return c, nil
		}
		return nil, fmt.Errorf("unknown color %q", input)
	}

	r, g, b, a, err := StringRGBA(input)
	if err != nil {
		return nil, err
	}
	return color.RGBA{r, g, b, a}, nil
}

func ToColorful(in color.Color) (colorful.Color, byte) {
	r, g, b, a := in.RGBA()
	out := colorful.Color{R: float64(r) / 65535, G: float64(g) / 65535, B: float64(b) / 65535}
	return out, byte(a >> 8)
}

func RGBAToHex(r, g, b, a byte) string {
	return "#" + hex.EncodeToString([]byte{r, g, b, a})
}

func ColorfulToHex(c colorful.Color, a byte) string {
	return fmt.Sprintf("%s%02x", c.Hex(), a)
}
