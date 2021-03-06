package vfont

import (
	"image/color"

	"github.com/tdewolff/canvas"
)

type Style struct {
	Color     color.Color
	PointSize float64
	PixelSize int
	Italic    bool
	Underline bool
	Strikeout bool

	Weight Weight
}

func (s Style) canvasStyle() canvas.FontStyle {
	if s.Italic {
		return s.Weight.toCanvas() | canvas.FontItalic
	} else {
		return s.Weight.toCanvas()
	}
}

func (s Style) searchableName() string {
	if s.Italic {
		return s.Weight.searchableName() + " italic"
	}
	return s.Weight.searchableName()
}

type Weight uint

const (
	Thin       Weight = 100
	ExtraLight Weight = 200
	UltraLight Weight = 200
	Light      Weight = 300
	Normal     Weight = 400
	Regular    Weight = 400
	Medium     Weight = 500
	SemiBold   Weight = 600
	DemiBold   Weight = 600
	Bold       Weight = 700
	ExtraBold  Weight = 800
	UltraBold  Weight = 800
	Black      Weight = 900
	Heavy      Weight = 900
)

func (w Weight) toCanvas() canvas.FontStyle {
	switch w {
	case Thin:
		return canvas.FontExtraLight
	case ExtraLight:
		return canvas.FontLight
	case Light:
		return canvas.FontBook
	case Normal:
		return canvas.FontRegular
	case Medium:
		return canvas.FontMedium
	case DemiBold:
		return canvas.FontSemibold
	case Bold:
		return canvas.FontBold
	case ExtraBold:
		return canvas.FontBlack
	case Black:
		return canvas.FontExtraBlack

	default:
		return canvas.FontRegular
	}
}

// searchableName returns a string representation that will be recognized by github.com/adrg/sysfont (which is used by canvas)
func (w Weight) searchableName() string {
	switch w {
	case Thin:
		return "thin"
	case ExtraLight:
		return "extralight"
	case Light:
		return "light"
	case Normal:
		return "regular"
	case Medium:
		return "medium"
	case SemiBold:
		return "semibold"
	case Bold:
		return "bold"
	case ExtraBold:
		return "extrabold"
	case Black:
		return "heavy"

	default:
		return ""
	}
}

func LoadFontFace(familyName string, style Style) (*canvas.FontFace, error) {
	fontData := canvas.NewFontFamily(familyName)
	err := fontData.LoadLocalFont(familyName+" "+style.searchableName(), style.canvasStyle())
	if err != nil {
		return nil, err
	}

	var decorators []canvas.FontDecorator
	if style.Underline {
		decorators = append(decorators, canvas.FontUnderline)
	}
	if style.Strikeout {
		decorators = append(decorators, canvas.FontStrikethrough)
	}

	if style.PointSize == 0 {
		style.PointSize = PixelsToPoints(style.PixelSize)
	}

	return fontData.Face(
		style.PointSize,
		style.Color,
		style.canvasStyle(),
		canvas.FontNormal,
		decorators...,
	), nil
}

func PixelsToPoints(px int) float64 {
	// TODO: implement a proper conversion
	//       maybe it could be done on the drawing context instead?
	return float64(px)
	// return float64(px) / 96
}
