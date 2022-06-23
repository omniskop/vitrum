package std

import (
	"math"

	"github.com/omniskop/vitrum/vit"
	"github.com/tdewolff/canvas"
)

func (r *Rectangle) Draw(ctx vit.DrawingContext, area vit.Rect) error {
	rect := r.Bounds()

	ctx.SetFillColor(r.color.Color())

	borderWidth := float64(r.border.MustGet("width").GetValue().(int))
	if borderWidth > 0 {
		ctx.SetStrokeJoiner(canvas.MiterJoin)
		ctx.SetStrokeColor(r.border.MustGet("color").(*vit.ColorValue).Color())
		ctx.SetStrokeWidth(float64(r.border.MustGet("width").(*vit.IntValue).Int()))
		radius := math.Max(r.radius.Float64()-borderWidth/2, 0) // decrease radius. Make sure it can't get below 0
		ctx.DrawPath(
			rect.X1+borderWidth/2,
			rect.Y1+borderWidth/2,
			canvas.RoundedRectangle(rect.Width()-borderWidth, rect.Height()-borderWidth, radius),
		)
	} else {
		ctx.DrawPath(rect.X1, rect.Y1, canvas.RoundedRectangle(rect.Width(), rect.Height(), r.radius.Float64()))
	}

	return r.Root.DrawChildren(ctx, rect)
}
