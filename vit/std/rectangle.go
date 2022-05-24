package std

import (
	"github.com/omniskop/vitrum/vit"
)

func (r *Rectangle) Draw(ctx vit.DrawingContext, area vit.Rect) error {
	rect := vit.Rect{
		X1: r.left.Float64(),
		Y1: r.top.Float64(),
		X2: r.right.Float64(),
		Y2: r.bottom.Float64(),
	}

	ctx.FillColor = r.color.RGBAColor()
	ctx.DrawPath(0, 0, rect.ToCanvas().ToPath())

	r.Root.DrawChildren(ctx, rect)

	return nil
}
