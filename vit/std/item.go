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

	layout *vit.Layout
}

func NewItem(id string, scope vit.ComponentContainer) *Item {
	i := &Item{
		Root:             vit.NewRoot(id, scope),
		id:               id,
		width:            *vit.NewFloatValue("", nil),
		height:           *vit.NewFloatValue("", nil),
		anchors:          *vit.NewAnchors("", nil),
		x:                *vit.NewFloatValue("", nil),
		y:                *vit.NewFloatValue("", nil),
		z:                *vit.NewFloatValue("", nil),
		left:             *vit.NewAnchorLineValue(),
		horizontalCenter: *vit.NewAnchorLineValue(),
		right:            *vit.NewAnchorLineValue(),
		top:              *vit.NewAnchorLineValue(),
		verticalCenter:   *vit.NewAnchorLineValue(),
		bottom:           *vit.NewAnchorLineValue(),
	}
	i.anchors.OnChange = func() { i.layouting(0, 0) }
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

func (i *Item) SetProperty(key string, value interface{}, position *vit.PositionRange) bool {
	// fmt.Printf("[Item] set %q: %v\n", key, value)
	switch key {
	case "width":
		i.width.Expression.ChangeCode(value.(string), position)
	case "height":
		i.height.Expression.ChangeCode(value.(string), position)
	case "anchors":
		panic("not implemented")
		// i.anchors = value.(vit.ObjectValue)
	case "x":
		i.x.Expression.ChangeCode(value.(string), position)
	case "y":
		i.y.Expression.ChangeCode(value.(string), position)
	case "z":
		i.z.Expression.ChangeCode(value.(string), position)
	default:
		return i.Root.SetProperty(key, value, position)
	}
	return true
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
	if i.width.ShouldEvaluate() {
		sum++
		err := i.width.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "width", i.id, i.width.Expression, err))
		}
		w := i.width.Float64()
		h := i.height.Float64()
		i.layout.SetTargetSize(&w, &h)
	}
	if i.height.ShouldEvaluate() {
		sum++
		err := i.height.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "height", i.id, i.height.Expression, err))
		}
		w := i.width.Float64()
		h := i.height.Float64()
		i.layout.SetTargetSize(&w, &h)
	}
	if i.x.ShouldEvaluate() {
		sum++
		err := i.x.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "x", i.id, i.x.Expression, err))
		}
	}
	if i.y.ShouldEvaluate() {
		sum++
		err := i.y.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "y", i.id, i.y.Expression, err))
		}
	}
	if i.z.ShouldEvaluate() {
		sum++
		err := i.z.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "z", i.id, i.z.Expression, err))
		}
	}
	if i.left.ShouldEvaluate() {
		sum++
		err := i.left.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "left", i.id, vit.Expression{}, err))
		}
	}
	if i.horizontalCenter.ShouldEvaluate() {
		sum++
		err := i.horizontalCenter.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "horizontalCenter", i.id, vit.Expression{}, err))
		}
	}
	if i.right.ShouldEvaluate() {
		sum++
		err := i.right.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "right", i.id, vit.Expression{}, err))
		}
	}
	if i.top.ShouldEvaluate() {
		sum++
		err := i.top.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "top", i.id, vit.Expression{}, err))
		}
	}
	if i.verticalCenter.ShouldEvaluate() {
		sum++
		err := i.verticalCenter.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "verticalCenter", i.id, vit.Expression{}, err))
		}
	}
	if i.bottom.ShouldEvaluate() {
		sum++
		err := i.bottom.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "bottom", i.id, vit.Expression{}, err))
		}
	}

	if i.layout.PositionChanged() {
		i.layouting(0, 0)
		sum++
	}
	if i.layout.SizeChanged() {
		i.layouting(0, 0)
		sum++
	}

	n, err := i.anchors.UpdateExpressions(i)
	if n > 0 {
		i.layouting(0, 0)
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

func (i *Item) layouting(autoWidth, autoHeight float64) {
	var width = i.width.Float64()
	var height = i.height.Float64()
	if width == 0 {
		width = autoWidth
	}
	if height == 0 {
		height = autoHeight
	}
	var didHorizintal bool
	var didPreferredHorizontal bool
	var didVertical bool
	var didPreferredVertical bool

	if i.layout != nil {
		if w, ok := i.layout.GetWidth(); ok {
			width = w
		}
		if h, ok := i.layout.GetHeight(); ok {
			height = h
		}

		if i.layout.PositionChanged() {
			i.layout.AckPositionChange()
			if x, ok := i.layout.GetX(); ok {
				i.left.SetAbsolute(x)
				i.horizontalCenter.SetAbsolute(x + width/2)
				i.right.SetAbsolute(x + width)
				didHorizintal = true
			} else if x, ok := i.layout.GetPreferredX(); ok {
				i.left.SetAbsolute(x)
				i.horizontalCenter.SetAbsolute(x + width/2)
				i.right.SetAbsolute(x + width)
				didPreferredHorizontal = true
			}
			if y, ok := i.layout.GetY(); ok {
				i.top.SetAbsolute(y)
				i.verticalCenter.SetAbsolute(y + height/2)
				i.bottom.SetAbsolute(y + height)
				didVertical = true
			} else if y, ok := i.layout.GetPreferredY(); ok {
				i.top.SetAbsolute(y)
				i.verticalCenter.SetAbsolute(y + height/2)
				i.bottom.SetAbsolute(y + height)
				didPreferredVertical = true
			}
		}
	}

	if !(didHorizintal || didVertical) && i.anchors.Fill.GetValue() != nil {
		i.left.AssignTo(i.anchors.Fill.Component(), vit.AnchorLeft)
		i.horizontalCenter.AssignTo(i.anchors.Fill.Component(), vit.AnchorHorizontalCenter)
		i.right.AssignTo(i.anchors.Fill.Component(), vit.AnchorRight)
		i.top.AssignTo(i.anchors.Fill.Component(), vit.AnchorTop)
		i.verticalCenter.AssignTo(i.anchors.Fill.Component(), vit.AnchorVerticalCenter)
		i.bottom.AssignTo(i.anchors.Fill.Component(), vit.AnchorBottom)
		if i.anchors.LeftMargin.IsSet() {
			i.left.SetOffset(i.anchors.LeftMargin.Value.Float64())
		} else {
			i.left.SetOffset(0)
		}
		i.horizontalCenter.SetOffset(0)
		if i.anchors.RightMargin.IsSet() {
			i.right.SetOffset(-i.anchors.RightMargin.Value.Float64())
		} else {
			i.right.SetOffset(0)
		}
		if i.anchors.TopMargin.IsSet() {
			i.top.SetOffset(i.anchors.TopMargin.Value.Float64())
		} else {
			i.top.SetOffset(0)
		}
		i.verticalCenter.SetOffset(0)
		if i.anchors.BottomMargin.IsSet() {
			i.bottom.SetOffset(-i.anchors.BottomMargin.Value.Float64())
		} else {
			i.bottom.SetOffset(0)
		}
		return
	}
	if !(didHorizintal || didVertical) && i.anchors.CenterIn.GetValue() != nil {
		i.left.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorHorizontalCenter)
		i.left.SetOffset(-width / 2)
		i.horizontalCenter.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorHorizontalCenter)
		i.right.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorHorizontalCenter)
		i.right.SetOffset(width / 2)
		i.top.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorVerticalCenter)
		i.top.SetOffset(-height / 2)
		i.verticalCenter.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorVerticalCenter)
		i.bottom.AssignTo(i.anchors.CenterIn.Component(), vit.AnchorVerticalCenter)
		i.bottom.SetOffset(height / 2)
		i.left.SetOffset(0)
		if i.anchors.HorizontalCenterOffset.IsSet() {
			i.horizontalCenter.SetOffset(i.anchors.HorizontalCenterOffset.Value.Float64())
		} else {
			i.horizontalCenter.SetOffset(0)
		}
		i.right.SetOffset(0)
		i.top.SetOffset(0)
		if i.anchors.VerticalCenterOffset.IsSet() {
			i.verticalCenter.SetOffset(i.anchors.VerticalCenterOffset.Value.Float64())
		} else {
			i.verticalCenter.SetOffset(0)
		}
		i.bottom.SetOffset(0)
		return
	}

	if !didHorizintal && i.anchors.Left.IsSet() {
		left := i.anchors.Left.GetValue().(float64)
		if i.anchors.LeftMargin.IsSet() {
			left += i.anchors.LeftMargin.GetValue().(float64)
		}
		i.left.SetAbsolute(left)
		i.horizontalCenter.SetAbsolute(left + width/2)
		i.horizontalCenter.SetOffset(0)
		i.right.SetAbsolute(left + width)
		i.right.SetOffset(0)
		didHorizintal = true
	}

	if !didHorizintal && i.anchors.HorizontalCenter.IsSet() {
		horizontalCenter := i.anchors.HorizontalCenter.GetValue().(float64)
		if i.anchors.HorizontalCenterOffset.IsSet() {
			horizontalCenter += i.anchors.HorizontalCenterOffset.GetValue().(float64)
		}
		i.left.SetAbsolute(horizontalCenter - width/2)
		i.horizontalCenter.SetAbsolute(horizontalCenter)
		i.right.SetAbsolute(horizontalCenter + width/2)
		didHorizintal = true
	}

	if !didHorizintal && i.anchors.Right.IsSet() {
		right := i.anchors.Right.GetValue().(float64)
		if i.anchors.RightMargin.IsSet() {
			right -= i.anchors.RightMargin.GetValue().(float64)
		}
		i.left.SetAbsolute(right - width)
		i.horizontalCenter.SetAbsolute(right - width/2)
		i.right.SetAbsolute(right)
		didHorizintal = true
	}

	if !didVertical && i.anchors.Top.IsSet() {
		top := i.anchors.Top.GetValue().(float64)
		if i.anchors.TopMargin.IsSet() {
			top += i.anchors.TopMargin.GetValue().(float64)
		}
		i.top.SetAbsolute(top)
		i.verticalCenter.SetAbsolute(top + height/2)
		i.bottom.SetAbsolute(top + height)
		didVertical = true
	}

	if !didVertical && i.anchors.VerticalCenter.IsSet() {
		verticalCenter := i.anchors.VerticalCenter.GetValue().(float64)
		if i.anchors.VerticalCenterOffset.IsSet() {
			verticalCenter += i.anchors.VerticalCenterOffset.GetValue().(float64)
		}
		i.top.SetAbsolute(verticalCenter - height/2)
		i.verticalCenter.SetAbsolute(verticalCenter)
		i.bottom.SetAbsolute(verticalCenter + height/2)
		didVertical = true
	}

	if !didVertical && i.anchors.Bottom.IsSet() {
		bottom := i.anchors.Bottom.GetValue().(float64)
		if i.anchors.BottomMargin.IsSet() {
			bottom -= i.anchors.BottomMargin.GetValue().(float64)
		}
		i.top.SetAbsolute(bottom - height)
		i.verticalCenter.SetAbsolute(bottom - height/2)
		i.bottom.SetAbsolute(bottom)
		didVertical = true
	}

	if !didHorizintal && !didPreferredHorizontal {
		i.left.AssignTo(i.Parent(), vit.AnchorLeft)
		i.left.SetOffset(i.x.Float64())
		i.horizontalCenter.AssignTo(i.Parent(), vit.AnchorLeft)
		i.horizontalCenter.SetOffset(i.x.Float64() + width/2)
		i.right.AssignTo(i.Parent(), vit.AnchorLeft)
		i.right.SetOffset(i.x.Float64() + width)
	}
	if !didVertical && !didPreferredVertical {
		i.top.AssignTo(i.Parent(), vit.AnchorTop)
		i.top.SetOffset(i.y.Float64())
		i.verticalCenter.AssignTo(i.Parent(), vit.AnchorTop)
		i.verticalCenter.SetOffset(i.y.Float64() + height/2)
		i.bottom.AssignTo(i.Parent(), vit.AnchorLeft)
		i.bottom.SetOffset(i.y.Float64() + height)
	}
}

func (i *Item) setAllOffsets() {
	if i.anchors.LeftMargin.IsSet() {
		i.left.SetOffset(i.anchors.LeftMargin.Value.Float64())
	} else {
		i.left.SetOffset(0)
	}
	if i.anchors.HorizontalCenterOffset.IsSet() {
		i.horizontalCenter.SetOffset(i.anchors.HorizontalCenterOffset.Value.Float64())
	} else {
		i.horizontalCenter.SetOffset(0)
	}
	if i.anchors.RightMargin.IsSet() {
		i.right.SetOffset(i.anchors.RightMargin.Value.Float64())
	} else {
		i.right.SetOffset(0)
	}
	if i.anchors.TopMargin.IsSet() {
		i.top.SetOffset(i.anchors.TopMargin.Value.Float64())
	} else {
		i.top.SetOffset(0)
	}
	if i.anchors.VerticalCenterOffset.IsSet() {
		i.verticalCenter.SetOffset(i.anchors.VerticalCenterOffset.Value.Float64())
	} else {
		i.verticalCenter.SetOffset(0)
	}
	if i.anchors.BottomMargin.IsSet() {
		i.bottom.SetOffset(i.anchors.BottomMargin.Value.Float64())
	} else {
		i.bottom.SetOffset(0)
	}
}

func (i *Item) Draw(ctx vit.DrawingContext, area vit.Rect) error {
	i.DrawChildren(ctx, area)
	return nil
}

func (i *Item) ApplyLayout(l *vit.Layout) {
	i.layout = l
}
