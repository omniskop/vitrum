package pdf

import (
	vit "github.com/omniskop/vitrum/vit"
	"github.com/tdewolff/canvas"
)

func (p *PageComponent) Draw(ctx vit.DrawingContext, area vit.Rect) error {
	rect := p.Bounds()

	ctx.SetFillColor(p.color.Color())

	ctx.DrawPath(rect.X1, rect.Y1, canvas.Rectangle(rect.Width(), rect.Height()))

	return p.Root.DrawChildren(ctx, rect)
}

func (p *PageComponent) sizeChanged() {
	p.SetContentSize(pageSize(PageComponent_Format(p.format.Int()), PageComponent_Orientation(p.orientation.Int())))
}

func pageSize(format PageComponent_Format, orientation PageComponent_Orientation) (float64, float64) {
	var w, h float64
	// At some point it would probably be worth putting these into a map
	switch format {
	case PageComponent_Format_A0:
		w, h = 841, 1189
	case PageComponent_Format_A1:
		w, h = 594, 841
	case PageComponent_Format_A2:
		w, h = 420, 594
	case PageComponent_Format_A3:
		w, h = 297, 420
	case PageComponent_Format_A4:
		w, h = 210, 297
	case PageComponent_Format_A5:
		w, h = 148, 210
	case PageComponent_Format_A6:
		w, h = 105, 148
	case PageComponent_Format_A7:
		w, h = 74, 105
	case PageComponent_Format_A8:
		w, h = 52, 74
	case PageComponent_Format_A9:
		w, h = 37, 52
	case PageComponent_Format_A10:
		w, h = 26, 37
	}

	if orientation == PageComponent_Orientation_Portrait {
		return w, h
	} else {
		return h, w
	}
}
