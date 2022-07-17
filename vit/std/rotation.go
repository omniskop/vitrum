package std

import vit "github.com/omniskop/vitrum/vit"

func (r *Rotation) Draw(ctx vit.DrawingContext, area vit.Rect) error {
	rect := r.Bounds()

	var hPivot float64
	var vPivot float64

	switch Rotation_HorizontalPivot(r.horizontalPivot.Int()) {
	case Rotation_HorizontalPivot_PivotLeft:
		hPivot = rect.Left()
	case Rotation_HorizontalPivot_PivotHCenter:
		hPivot = rect.CenterX()
	case Rotation_HorizontalPivot_PivotRight:
		hPivot = rect.Right()
	}

	switch Rotation_VerticalPivot(r.verticalPivot.Int()) {
	case Rotation_VerticalPivot_PivotTop:
		vPivot = rect.Top()
	case Rotation_VerticalPivot_PivotVCenter:
		vPivot = rect.CenterY()
	case Rotation_VerticalPivot_PivorBottom:
		vPivot = rect.Bottom()
	}

	ctx.RotateAbout(r.degrees.Float64(), hPivot, vPivot)
	defer ctx.RotateAbout(-r.degrees.Float64(), hPivot, vPivot)

	return r.Root.DrawChildren(ctx, rect)
}
