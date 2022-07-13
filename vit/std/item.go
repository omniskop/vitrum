package std

import (
	"fmt"

	vit "github.com/omniskop/vitrum/vit"
)

type Item struct {
	vit.Root
	id string

	width            vit.FloatValue
	height           vit.FloatValue
	anchors          vit.AnchorsValue
	x                vit.FloatValue
	y                vit.FloatValue
	z                vit.FloatValue
	left             vit.AnchorLineValue
	horizontalCenter vit.AnchorLineValue
	right            vit.AnchorLineValue
	top              vit.AnchorLineValue
	verticalCenter   vit.AnchorLineValue
	bottom           vit.AnchorLineValue

	contentWidth  float64
	contentHeight float64

	layout *vit.Layout
}

func NewItem(id string, context *vit.FileContext) *Item {
	i := &Item{
		Root:             vit.NewRoot(id, context),
		id:               id,
		width:            *vit.NewEmptyFloatValue(),
		height:           *vit.NewEmptyFloatValue(),
		anchors:          *vit.NewAnchors(),
		x:                *vit.NewEmptyFloatValue(),
		y:                *vit.NewEmptyFloatValue(),
		z:                *vit.NewEmptyFloatValue(),
		left:             *vit.NewAnchorLineValue(),
		horizontalCenter: *vit.NewAnchorLineValue(),
		right:            *vit.NewAnchorLineValue(),
		top:              *vit.NewAnchorLineValue(),
		verticalCenter:   *vit.NewAnchorLineValue(),
		bottom:           *vit.NewAnchorLineValue(),
	}
	i.x.AddDependent(vit.FuncDep(i.layouting))
	i.y.AddDependent(vit.FuncDep(i.layouting))
	i.z.AddDependent(vit.FuncDep(i.layouting))
	i.width.AddDependent(vit.FuncDep(i.layouting))
	i.width.AddDependent(vit.FuncDep(i.updateLayoutSize))
	i.height.AddDependent(vit.FuncDep(i.layouting))
	i.height.AddDependent(vit.FuncDep(i.updateLayoutSize))
	return i
}

func (i *Item) String() string {
	return fmt.Sprintf("Item{%s}", i.id)
}

func (i *Item) Property(key string) (vit.Value, bool) {
	switch key {
	case "width":
		return &i.width, true
	case "height":
		return &i.height, true
	case "anchors":
		return &i.anchors, true
	case "x":
		return &i.x, true
	case "y":
		return &i.y, true
	case "z":
		return &i.z, true
	case "left":
		return &i.left, true
	case "horizontalCenter":
		return &i.horizontalCenter, true
	case "right":
		return &i.right, true
	case "top":
		return &i.top, true
	case "verticalCenter":
		return &i.verticalCenter, true
	case "bottom":
		return &i.bottom, true
	default:
		return i.Root.Property(key)
	}
}

func (i *Item) MustProperty(key string) vit.Value {
	v, ok := i.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (i *Item) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "width":
		err = i.width.SetValue(value)
	case "height":
		err = i.height.SetValue(value)
	case "anchors":
		panic("not implemented")
		// i.anchors = value.(vit.ObjectValue)
	case "x":
		err = i.x.SetValue(value)
	case "y":
		err = i.y.SetValue(value)
	case "z":
		err = i.z.SetValue(value)
	default:
		return i.Root.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("item", key, i.ID(), err)
	}
	return nil
}

func (i *Item) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "width":
		i.width.SetCode(code)
	case "height":
		i.height.SetCode(code)
	case "anchors":
		panic("not implemented")
	case "x":
		i.x.SetCode(code)
	case "y":
		i.y.SetCode(code)
	case "z":
		i.z.SetCode(code)
	default:
		return i.Root.SetPropertyCode(key, code)
	}
	return nil
}

func (i *Item) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case i.id:
		return i, true
	case "width":
		return &i.width, true
	case "height":
		return &i.height, true
	case "anchors":
		return &i.anchors, true
	case "x":
		return &i.x, true
	case "y":
		return &i.y, true
	case "z":
		return &i.z, true
	case "left":
		return &i.left, true
	case "horizontalCenter":
		return &i.horizontalCenter, true
	case "right":
		return &i.right, true
	case "top":
		return &i.top, true
	case "verticalCenter":
		return &i.verticalCenter, true
	case "bottom":
		return &i.bottom, true
	default:
		return i.Root.ResolveVariable(key)
	}
}

func (i *Item) AddChild(child vit.Component) {
	child.SetParent(i)
	i.Root.AddChildButKeepParent(child)
}

func (i *Item) AddChildAfter(afterThis, addThis vit.Component) {
	var dynType vit.Component = afterThis

	for j, child := range i.Children() {
		if child.As(&dynType) {
			addThis.SetParent(i)
			i.AddChildAtButKeepParent(addThis, j+1)
			return
		}
	}
	i.AddChild(addThis)
}

func (i *Item) UpdateExpressions() (int, vit.ErrorGroup) {
	var errs vit.ErrorGroup
	var sum int
	if changed, err := i.width.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "width", i.id, err))
		}
	}
	if changed, err := i.height.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "height", i.id, err))
		}
	}
	if changed, err := i.x.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "x", i.id, err))
		}
	}
	if changed, err := i.y.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "y", i.id, err))
		}
	}
	if changed, err := i.z.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "z", i.id, err))
		}
	}
	if changed, err := i.left.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "left", i.id, err))
		}
	}
	if changed, err := i.horizontalCenter.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "horizontalCenter", i.id, err))
		}
	}
	if changed, err := i.right.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "right", i.id, err))
		}
	}
	if changed, err := i.top.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "top", i.id, err))
		}
	}
	if changed, err := i.verticalCenter.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "verticalCenter", i.id, err))
		}
	}
	if changed, err := i.bottom.Update(i); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Item", "bottom", i.id, err))
		}
	}

	if i.layout.PositionChanged() {
		i.layouting()
		sum++
	}
	if i.layout.SizeChanged() {
		i.layouting()
		sum++
	}

	n, err := i.anchors.UpdateExpressions(i)
	if n > 0 {
		i.layouting()
	}
	sum += n
	errs.AddGroup(err)

	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	n, err = i.UpdatePropertiesInContext(i)
	sum += n
	errs.AddGroup(err)
	n, err = i.Root.UpdateExpressions()
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (i *Item) ID() string {
	return i.id
}

func (i *Item) As(target *vit.Component) bool {
	if _, ok := (*target).(*Item); ok {
		*target = i
		return true
	}
	return false
}

func (i *Item) Finish() error {
	return i.RootC().FinishInContext(i)
}

func (i *Item) layouting() {
	var width = i.width.Float64()
	var height = i.height.Float64()
	if width == 0 {
		width = i.contentWidth
	}
	if height == 0 {
		height = i.contentHeight
	}
	var didSetPreferredLeft bool
	var didSetPreferredTop bool

	var didSetTop bool
	var didSetRight bool
	var didSetBottom bool
	var didSetLeft bool

	// check if the component has a forced layout
	if i.layout != nil {
		// set width and height if it's set through the layout
		if w, ok := i.layout.GetWidth(); ok {
			width = w
		}
		if h, ok := i.layout.GetHeight(); ok {
			height = h
		}

		var updateValues bool
		if i.layout.PositionChanged() {
			i.layout.AckPositionChange()
			updateValues = true
		}
		if x, ok := i.layout.GetX(); ok {
			if updateValues {
				i.left.SetAbsolute(x)
			}
			didSetLeft = true
		} else if x, ok := i.layout.GetPreferredX(); ok {
			if updateValues {
				i.left.SetAbsolute(x)
			}
			didSetPreferredLeft = true
		}
		if y, ok := i.layout.GetY(); ok {
			if updateValues {
				i.top.SetAbsolute(y)
			}
			didSetTop = true
		} else if y, ok := i.layout.GetPreferredY(); ok {
			if updateValues {
				i.top.SetAbsolute(y)
			}
			didSetPreferredTop = true
		}
	} else {
		// we only consider fill and centerIn if no layout is set

		// check anchors.fill
		if i.anchors.Fill.GetValue() != nil {
			i.left.AssignTo(i.anchors.Fill.Component(), vit.AnchorLeft)
			i.horizontalCenter.AssignTo(i.anchors.Fill.Component(), vit.AnchorHorizontalCenter)
			i.right.AssignTo(i.anchors.Fill.Component(), vit.AnchorRight)
			i.top.AssignTo(i.anchors.Fill.Component(), vit.AnchorTop)
			i.verticalCenter.AssignTo(i.anchors.Fill.Component(), vit.AnchorVerticalCenter)
			i.bottom.AssignTo(i.anchors.Fill.Component(), vit.AnchorBottom)
			i.left.SetOffset(i.anchors.CalcLeftMargin())
			i.horizontalCenter.SetOffset(0)
			i.right.SetOffset(-i.anchors.CalcRightMargin())
			i.top.SetOffset(i.anchors.CalcTopMargin())
			i.verticalCenter.SetOffset(0)
			i.bottom.SetOffset(-i.anchors.CalcBottomMargin())
			return // all is set, nothing can be overwritten anymore, so we can stop here
		}

		// check anchors.centerIn
		if i.anchors.CenterIn.GetValue() != nil {
			var hOffset float64
			var vOffset float64
			if i.anchors.HorizontalCenterOffset.IsSet() {
				hOffset = i.anchors.HorizontalCenterOffset.Value().Float64()
			}
			if i.anchors.VerticalCenterOffset.IsSet() {
				vOffset = i.anchors.VerticalCenterOffset.Value().Float64()
			}
			i.left.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorHorizontalCenter)
			i.left.SetOffset(-width/2 + hOffset)
			i.horizontalCenter.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorHorizontalCenter)
			i.horizontalCenter.SetOffset(hOffset)
			i.right.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorHorizontalCenter)
			i.right.SetOffset(width/2 + hOffset)
			i.top.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorVerticalCenter)
			i.top.SetOffset(-height/2 + vOffset)
			i.verticalCenter.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorVerticalCenter)
			i.verticalCenter.SetOffset(vOffset)
			i.bottom.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorVerticalCenter)
			i.bottom.SetOffset(height/2 + vOffset)
			return // all is set, nothing can be overwritten anymore, so we can stop here
		}
	}

	// apply anchors.left
	if !didSetLeft && i.anchors.Left.IsSet() {
		left := i.anchors.Left.GetValue().(float64)
		left += i.anchors.CalcLeftMargin()
		i.left.SetAbsolute(left)
		didSetLeft = true
	}

	// apply anchors.horizontalCenter
	if !didSetLeft && !didSetRight && i.anchors.HorizontalCenter.IsSet() {
		horizontalCenter := i.anchors.HorizontalCenter.GetValue().(float64)
		if i.anchors.HorizontalCenterOffset.IsSet() {
			horizontalCenter += i.anchors.HorizontalCenterOffset.GetValue().(float64)
		}
		i.left.SetAbsolute(horizontalCenter - width/2)
		i.horizontalCenter.SetAbsolute(horizontalCenter)
		i.right.SetAbsolute(horizontalCenter + width/2)
		didSetLeft = true
		didSetRight = true
	}

	// apply anchors.right
	if !didSetRight && i.anchors.Right.IsSet() {
		right := i.anchors.Right.GetValue().(float64)
		right -= i.anchors.CalcRightMargin()
		i.right.SetAbsolute(right)
		didSetRight = true
	}

	// apply anchors.top
	if !didSetTop && i.anchors.Top.IsSet() {
		top := i.anchors.Top.GetValue().(float64)
		top += i.anchors.CalcTopMargin()
		i.top.SetAbsolute(top)
		didSetTop = true
	}

	// apply anchors.verticalCenter
	if !didSetTop && !didSetBottom && i.anchors.VerticalCenter.IsSet() {
		verticalCenter := i.anchors.VerticalCenter.GetValue().(float64)
		if i.anchors.VerticalCenterOffset.IsSet() {
			verticalCenter += i.anchors.VerticalCenterOffset.GetValue().(float64)
		}
		i.top.SetAbsolute(verticalCenter - height/2)
		i.verticalCenter.SetAbsolute(verticalCenter)
		i.bottom.SetAbsolute(verticalCenter + height/2)
		didSetTop = true
		didSetBottom = true
	}

	// apply anchors.bottom
	if !didSetBottom && i.anchors.Bottom.IsSet() {
		bottom := i.anchors.Bottom.GetValue().(float64)
		bottom -= i.anchors.CalcBottomMargin()
		i.bottom.SetAbsolute(bottom)
		didSetBottom = true
	}

	// set left implicitly based on right
	if !didSetLeft && !didSetPreferredLeft && didSetRight {
		// if preferred left is set, we won't set left implicitly
		i.left.SetAbsolute(i.right.Float64() - width)
		didSetLeft = true
	}
	// set right implicitly based on left
	if !didSetRight && didSetLeft {
		i.right.SetAbsolute(i.left.Float64() + width)
		didSetRight = true
	}
	// set top implicitly based on bottom
	if !didSetTop && !didSetPreferredTop && didSetBottom {
		i.top.SetAbsolute(i.bottom.Float64() - height)
		didSetTop = true
	}
	// set bottom implicitly based on top
	if !didSetBottom && didSetTop {
		i.bottom.SetAbsolute(i.top.Float64() + height)
		didSetBottom = true
	}

	// set left implicitly based on position
	if !didSetLeft {
		i.left.SetAbsolute(i.x.Float64())
		didSetLeft = true
	}
	// set right implicitly based on position
	if !didSetRight {
		i.right.SetAbsolute(i.x.Float64() + width)
		didSetRight = true
	}
	// set top implicitly based on position
	if !didSetTop {
		i.top.SetAbsolute(i.y.Float64())
		didSetTop = true
	}
	// set bottom implicitly based on position
	if !didSetBottom {
		i.bottom.SetAbsolute(i.y.Float64() + height)
		didSetBottom = true
	}

	bounds := i.Bounds()
	width = bounds.Width()
	height = bounds.Height()

	// calculate horizontal and vertical center
	// this might overwrite an existing value but it should be the same anyway
	i.verticalCenter.SetAbsolute(bounds.X1 + height/2)
	i.horizontalCenter.SetAbsolute(bounds.Y1 + width/2)

	// update potential layout
	i.layout.SetTargetSize(&width, &height)
}

func (i *Item) updateLayoutSize() {
	w := i.width.Float64()
	h := i.height.Float64()
	i.layout.SetTargetSize(&w, &h)
}

func (i *Item) Draw(ctx vit.DrawingContext, area vit.Rect) error {
	i.DrawChildren(ctx, area)
	return nil
}

func (i *Item) ApplyLayout(l *vit.Layout) {
	i.layout = l
	i.layouting()
}

func (i *Item) Bounds() vit.Rect {
	return vit.Rect{
		X1: i.left.Float64(),
		Y1: i.top.Float64(),
		X2: i.right.Float64(),
		Y2: i.bottom.Float64(),
	}
}
