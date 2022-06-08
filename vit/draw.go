package vit

import (
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

type Rect struct {
	X1, Y1, X2, Y2 float64
}

func (r Rect) ToCanvas() canvas.Rect {
	return canvas.Rect{r.X1, r.Y1, r.X2 - r.X1, r.Y2 - r.Y1}
}

type DrawingContext struct {
	*canvas.Context
}

func Draw(comp Component) error {
	// Create new canvas of dimension 1000x1000 mm
	c := canvas.New(1000, 1000)

	// Create a canvas context used to keep drawing state
	ctx := canvas.NewContext(c)

	// move origin of the context to the top left corner
	ctx.SetCoordSystem(canvas.CartesianIV)

	comp.Draw(DrawingContext{ctx}, Rect{0, 0, 1000, 1000})

	// Rasterize the canvas and write to a PNG file with 3.2 dots-per-mm (320x320 px)
	c.WriteFile("output.png", renderers.PNG())

	return nil
}
