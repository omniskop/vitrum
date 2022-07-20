package std

import (
	"math"
	"testing"

	vit "github.com/omniskop/vitrum/vit"
)

var testCases = []struct {
	mode Image_FillMode
	in   vit.Rect
	out  vit.Rect
}{
	{
		mode: Image_FillMode_Fill, // check scale up
		in:   vit.NewRect(0, 0, 50, 60),
		out:  vit.NewRect(100, 100, 100, 100),
	}, {
		mode: Image_FillMode_Fill, // check scale down
		in:   vit.NewRect(0, 0, 150, 180),
		out:  vit.NewRect(100, 100, 100, 100),
	},
	{
		mode: Image_FillMode_Fit, // check scale up of square image
		in:   vit.NewRect(0, 0, 50, 50),
		out:  vit.NewRect(100, 100, 100, 100),
	}, {
		mode: Image_FillMode_Fit, // check scale up with landscape ratio
		in:   vit.NewRect(0, 0, 50, 25),
		out:  vit.NewRect(100, 125, 100, 50),
	}, {
		mode: Image_FillMode_Fit, // check scale up with portrait ratio
		in:   vit.NewRect(0, 0, 25, 50),
		out:  vit.NewRect(125, 100, 50, 100),
	}, {
		mode: Image_FillMode_Fit, // check scale down of square image
		in:   vit.NewRect(0, 0, 150, 150),
		out:  vit.NewRect(100, 100, 100, 100),
	}, {
		mode: Image_FillMode_Fit, // check scale down with landscape ratio
		in:   vit.NewRect(0, 0, 150, 75),
		out:  vit.NewRect(100, 125, 100, 50),
	}, {
		mode: Image_FillMode_Fit, // check scale down with portrait ratio
		in:   vit.NewRect(0, 0, 75, 150),
		out:  vit.NewRect(125, 100, 50, 100),
	},
	{
		mode: Image_FillMode_PreferUnchanged, // check no change
		in:   vit.NewRect(0, 0, 50, 60),
		out:  vit.NewRect(100+(100-50)/2, 100+(100-60)/2, 50, 60),
	}, {
		mode: Image_FillMode_PreferUnchanged, // check scale down
		in:   vit.NewRect(0, 0, 150, 75),
		out:  vit.NewRect(100, 125, 100, 50),
	}}

func TestFill(t *testing.T) {
	var space = vit.NewRect(100, 100, 100, 100)
	var image Image
	for _, tc := range testCases {
		out := image.fill(space, tc.in, tc.mode)
		if !rectanglesEqual(out, tc.out) {
			t.Errorf("fill %v in %v using mode %v => got %v but wanted %v", tc.in, space, tc.mode, out, tc.out)
		}
	}
}

// rectanglesEqual returns true if both rectangle are considered equal enough.
func rectanglesEqual(a, b vit.Rect) bool {
	return math.Abs(a.X1-b.X1) < 0.00000001 &&
		math.Abs(a.Y1-b.Y1) < 0.00000001 &&
		math.Abs(a.X2-b.X2) < 0.00000001 &&
		math.Abs(a.Y2-b.Y2) < 0.00000001
}
