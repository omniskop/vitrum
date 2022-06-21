package std

import (
	"fmt"

	"github.com/omniskop/vitrum/vit"
	"github.com/tdewolff/canvas"
)

func (r *Rectangle) Draw(ctx vit.DrawingContext, area vit.Rect) error {
	rect := vit.Rect{
		X1: r.left.Float64(),
		Y1: r.top.Float64(),
		X2: r.right.Float64(),
		Y2: r.bottom.Float64(),
	}

	ctx.SetFillColor(r.color.Color())
	ctx.SetStrokeJoiner(canvas.MiterJoin)
	ctx.SetStrokeColor(r.border.MustGet("color").(*vit.ColorValue).Color())
	ctx.SetStrokeWidth(float64(r.border.MustGet("width").(*vit.IntValue).Int()))
	fmt.Println(r.border.MustGet("color").(*vit.ColorValue).Color().RGBA())
	fmt.Println(r.border.MustGet("width").(*vit.IntValue).Int())
	ctx.DrawPath(0, 0, rect.ToCanvas().ToPath())

	r.Root.DrawChildren(ctx, rect)

	return nil
}
