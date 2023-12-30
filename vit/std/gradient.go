package std

import (
	"image/color"

	"github.com/tdewolff/canvas"
)

func (r *Gradient) ConstructGradient(start, end canvas.Point) *canvas.LinearGradient {
	grad := canvas.NewLinearGradient(start, end)
	grad.Add(0, color.RGBA{0, 0, 0, 255})

	for i, child := range r.Children() {
		if i > 0 && i != len(r.Children())-1 {
			// canvas only supports gradients with 2 steps, thus we skip all but the first and last stop
			continue
		}
		stop, ok := child.(*GradientStop)
		if !ok {
			continue
		}

		grad.Add(stop.position.Float64(), stop.color.RGBAColor())
	}

	return grad
}
