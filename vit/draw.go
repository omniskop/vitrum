package vit

import (
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

type Rect struct {
	X1, Y1, X2, Y2 float64
}

func NewRect(x, y, w, h float64) Rect {
	return Rect{x, y, x + w, y + h}
}

func (r Rect) ToCanvas() canvas.Rect {
	return canvas.Rect{r.X1, r.Y1, r.X2 - r.X1, r.Y2 - r.Y1}
}

// MovedX returns a copy of the rectangle moved by the given amount on the x-axis.
func (r Rect) MovedX(x float64) Rect {
	return Rect{r.X1 + x, r.Y1, r.X2 + x, r.Y2}
}

// MovedY returns a copy of the rectangle moved by the given amount on the y-axis.
func (r Rect) MovedY(y float64) Rect {
	return Rect{r.X1, r.Y1 + y, r.X2, r.Y2 + y}
}

func (r Rect) Width() float64 {
	return r.X2 - r.X1
}

func (r Rect) Height() float64 {
	return r.Y2 - r.Y1
}

func (r Rect) CenterX() float64 {
	return (r.X1 + r.X2) / 2
}

func (r Rect) CenterY() float64 {
	return (r.Y1 + r.Y2) / 2
}

func (r Rect) Top() float64 {
	return r.Y1
}

func (r Rect) Right() float64 {
	return r.X2
}

func (r Rect) Bottom() float64 {
	return r.Y2
}

func (r Rect) Left() float64 {
	return r.X1
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
