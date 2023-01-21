package std

import (
	vit "github.com/omniskop/vitrum/vit"
)

func (c *Column) getTopPadding() float64 {
	if c.topPadding.IsSet() {
		return c.topPadding.Value().Float64()
	}
	return c.padding.Float64()
}

func (c *Column) getRightPadding() float64 {
	if c.rightPadding.IsSet() {
		return c.rightPadding.Value().Float64()
	}
	return c.padding.Float64()
}

func (c *Column) getBottomPadding() float64 {
	if c.bottomPadding.IsSet() {
		return c.bottomPadding.Value().Float64()
	}
	return c.padding.Float64()
}

func (c *Column) getLeftPadding() float64 {
	if c.leftPadding.IsSet() {
		return c.leftPadding.Value().Float64()
	}
	return c.padding.Float64()
}

func (c *Column) ContentSize() (float64, float64) {
	var totalWidth float64
	var totalHeight float64 = c.getTopPadding() + c.getBottomPadding()
	for _, child := range c.Children() {
		bounds := child.Bounds()
		if bounds.Width() == 0 || bounds.Height() == 0 {
			continue
		}

		totalHeight += bounds.Height() + c.spacing.Float64()
		if bounds.Width() > totalWidth {
			totalWidth = bounds.Width()
		}
	}
	if len(c.Children()) > 0 {
		totalHeight -= c.spacing.Float64()
	}
	totalWidth += c.getLeftPadding() + c.getRightPadding()
	return totalWidth, totalHeight
}

// Recalculate Layout of all child components.
func (c *Column) recalculateLayout() {
	var x float64 = c.left.Float64() + c.getLeftPadding()
	var y float64 = c.top.Float64() + c.getTopPadding()
	for _, child := range c.Children() {
		bounds := child.Bounds()
		if bounds.Width() == 0 || bounds.Height() == 0 {
			continue
		}

		xCopy := x
		yCopy := y
		c.childLayouts[child].SetPosition(nil, &yCopy)
		c.childLayouts[child].SetPreferredPosition(&xCopy, nil)
		y += bounds.Height() + c.spacing.Float64()
	}
	c.childLayouts.Update(c) // acknowledge all changes
	c.SetContentSize(c.ContentSize())
}

func (c *Column) createNewChildLayout(child vit.Component) *vit.Layout {
	l := vit.NewLayout()
	c.childLayouts[child] = l
	l.SetTargetSize(nil, nil)
	l.AddDependent(vit.FuncDep(c.recalculateLayout))
	return l
}

func (c *Column) childWasAdded(child vit.Component) {
	child.ApplyLayout(c.createNewChildLayout(child))
	c.recalculateLayout()
}
