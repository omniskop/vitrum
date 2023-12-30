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
		return canvas.FontThin
	case ExtraLight:
		return canvas.FontExtraLight
	case Light:
		return canvas.FontLight
	case Normal:
		return canvas.FontRegular
	case Medium:
		return canvas.FontMedium
	case DemiBold:
		return canvas.FontSemiBold
	case Bold:
		return canvas.FontBold
	case ExtraBold:
		return canvas.FontExtraBold
	case Black:
		return canvas.FontBlack

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

type LoadedFont struct {
	font   *canvas.FontFamily
	styles map[Style]bool // what styles are contained in the font
}

var cachedFonts = make(map[string]LoadedFont)

func LoadFontFace(familyName string, style Style) (*canvas.FontFace, error) {
	loadedFont, ok := cachedFonts[familyName]
	if !ok {
		loadedFont.font = canvas.NewFontFamily(familyName)
		loadedFont.styles = make(map[Style]bool)
	}

	if !loadedFont.styles[style] {
		err := loadedFont.font.LoadSystemFont(familyName, style.canvasStyle())
		if err != nil {
			return nil, err
		}
		loadedFont.styles[style] = true
		cachedFonts[familyName] = loadedFont
		// TODO: clean the cache somehow
	}

	decorators := []any{
		style.Color,
		style.canvasStyle(),
		canvas.FontNormal,
	}
	if style.Underline {
		decorators = append(decorators, canvas.FontUnderline)
	}
	if style.Strikeout {
		decorators = append(decorators, canvas.FontStrikethrough)
	}

	if style.PointSize == 0 {
		style.PointSize = PixelsToPoints(style.PixelSize)
	}

	return loadedFont.font.Face(style.PointSize, decorators...), nil
}

func PixelsToPoints(px int) float64 {
	// TODO: implement a proper conversion
	//       maybe it could be done on the drawing context instead?
	return float64(px)
	// return float64(px) / 96
}
