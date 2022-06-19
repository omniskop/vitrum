package std

import (
	vit "github.com/omniskop/vitrum/vit"
)

func (r *Row) getTopPadding() float64 {
	if r.topPadding.IsSet() {
		return r.topPadding.ActualValue().Float64()
	}
	return r.padding.Float64()
}

func (r *Row) getRightPadding() float64 {
	if r.rightPadding.IsSet() {
		return r.rightPadding.ActualValue().Float64()
	}
	return r.padding.Float64()
}

func (r *Row) getBottomPadding() float64 {
	if r.bottomPadding.IsSet() {
		return r.bottomPadding.ActualValue().Float64()
	}
	return r.padding.Float64()
}

func (r *Row) getLeftPadding() float64 {
	if r.leftPadding.IsSet() {
		return r.leftPadding.ActualValue().Float64()
	}
	return r.padding.Float64()
}

func (r *Row) ContentSize() (float64, float64) {
	var totalWidth float64 = r.getLeftPadding() + r.getRightPadding()
	var totalHeight float64
	for _, child := range r.Children() {
		bounds := child.Bounds()
		if bounds.Width() == 0 || bounds.Height() == 0 {
			continue
		}

		totalWidth += bounds.Width() + r.spacing.Float64()
		if bounds.Height() > totalHeight {
			totalHeight = bounds.Height()
		}
	}
	if len(r.Children()) > 0 {
		totalWidth -= r.spacing.Float64()
	}
	totalHeight += r.getTopPadding() + r.getBottomPadding()
	return totalWidth, totalHeight
}

// Recalculate Layout of all child components.
// The first parameter will usually be the property that resulted in the update but it is not used and should not be relied upon.
func (r *Row) recalculateLayout(interface{}) {
	var x float64 = r.left.Float64() + r.getLeftPadding()
	var y float64 = r.top.Float64() + r.getTopPadding()
	for _, child := range r.Children() {
		bounds := child.Bounds()
		if bounds.Width() == 0 || bounds.Height() == 0 {
			continue
		}

		xCopy := x
		yCopy := y
		r.childLayouts[child].SetPosition(&xCopy, nil)
		r.childLayouts[child].SetPreferredPosition(nil, &yCopy)
		x += bounds.Width() + r.spacing.Float64()
	}
	r.contentWidth, r.contentHeight = r.ContentSize()
	r.layouting(r.contentWidth, r.contentHeight)
}

func (r *Row) createNewChildLayout(child vit.Component) *vit.Layout {
	l := &vit.Layout{}
	r.childLayouts[child] = l
	l.SetTargetSize(nil, nil)
	return l
}

func (r *Row) childWasAdded(child vit.Component) {
	child.ApplyLayout(r.createNewChildLayout(child))
	r.recalculateLayout(nil)
}
