package std

import vit "github.com/omniskop/vitrum/vit"

func (c *Column) getTopPadding() float64 {
	if c.topPadding.IsSet() {
		return c.topPadding.ActualValue().Float64()
	}
	return c.padding.Float64()
}

func (c *Column) getRightPadding() float64 {
	if c.rightPadding.IsSet() {
		return c.rightPadding.ActualValue().Float64()
	}
	return c.padding.Float64()
}

func (c *Column) getBottomPadding() float64 {
	if c.bottomPadding.IsSet() {
		return c.bottomPadding.ActualValue().Float64()
	}
	return c.padding.Float64()
}

func (c *Column) getLeftPadding() float64 {
	if c.leftPadding.IsSet() {
		return c.leftPadding.ActualValue().Float64()
	}
	return c.padding.Float64()
}

func (c *Column) CalculateSize() (float64, float64) {
	var totalWidth float64
	var totalHeight float64 = c.getTopPadding() + c.getBottomPadding()
	for _, child := range c.Children() {
		width := child.MustProperty("width").GetValue().(float64)
		height := child.MustProperty("height").GetValue().(float64)
		if width == 0 || height == 0 {
			continue
		}

		totalHeight += height + c.spacing.Float64()
		if width > totalWidth {
			totalWidth = width
		}
	}
	if len(c.Children()) > 0 {
		totalHeight -= c.spacing.Float64()
	}
	totalWidth += c.getLeftPadding() + c.getRightPadding()
	return totalWidth, totalHeight
}

// Recalculate Layout of all child components.
// The first parameter will usually be the property that resulted in the update but it is not used and should not be relied upon.
func (c *Column) recalculateLayout(interface{}) {
	var x float64 = c.left.Float64() + c.getLeftPadding()
	var y float64 = c.top.Float64() + c.getTopPadding()
	for _, child := range c.Children() {
		width := child.MustProperty("width").GetValue().(float64)
		height := child.MustProperty("height").GetValue().(float64)
		if width == 0 || height == 0 {
			continue
		}

		xCopy := x
		yCopy := y
		c.childLayouts[child].SetPosition(nil, &yCopy)
		c.childLayouts[child].SetPreferredPosition(&xCopy, nil)
		y += height + c.spacing.Float64()
	}
	c.layouting(c.CalculateSize())
}

func (c *Column) createNewChildLayout(child vit.Component) *vit.Layout {
	l := &vit.Layout{}
	c.childLayouts[child] = l
	l.SetTargetSize(nil, nil)
	return l
}

func (c *Column) childWasAdded(child vit.Component) {
	child.ApplyLayout(c.createNewChildLayout(child))
	c.recalculateLayout(nil)
}
