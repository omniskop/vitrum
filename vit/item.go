package vit

import "fmt"

type Item struct {
	Root

	width   IntValue
	height  IntValue
	stuff   IntValue // TODO: delete
	anchors Anchors
	x       FloatValue
	y       FloatValue
	z       FloatValue
}

func NewItem(id string, scope ComponentContainer) *Item {
	return &Item{
		Root:   NewRoot(id, scope),
		width:  *NewIntValue("", nil),
		height: *NewIntValue("", nil),
		stuff:  *NewIntValue("", nil),
		x:      *NewFloatValue("", nil),
		y:      *NewFloatValue("", nil),
		z:      *NewFloatValue("", nil),
	}
}

func (i *Item) String() string {
	return fmt.Sprintf("Item{%s}", i.id)
}

func (i *Item) Property(key string) (interface{}, bool) {
	switch key {
	case "width":
		return &i.width.Value, true
	case "height":
		return &i.height.Value, true
	case "stuff":
		return &i.stuff.Value, true
	case "anchors":
		return &i.anchors, true
	case "x":
		return &i.x.Value, true
	case "y":
		return &i.y.Value, true
	case "z":
		return &i.z.Value, true
	default:
		return i.Root.Property(key)
	}
}

func (i *Item) MustProperty(key string) interface{} {
	v, ok := i.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (i *Item) SetProperty(key string, value interface{}, position *PositionRange) bool {
	// fmt.Printf("[Item] set %q: %v\n", key, value)
	switch key {
	case "width":
		i.width.Expression.ChangeCode(value.(string), position)
	case "height":
		i.height.Expression.ChangeCode(value.(string), position)
	case "stuff":
		i.stuff.Expression.ChangeCode(value.(string), position)
	case "anchors":
		i.anchors = value.(Anchors)
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
		return i.width, true
	case "height":
		return i.height, true
	case "stuff":
		return i.stuff, true
	case "anchors":
		return &i.anchors, true
	case "x":
		return i.x, true
	case "y":
		return i.y, true
	case "z":
		return i.z, true
	}

	return i.Root.ResolveVariable(key)
}

func (i *Item) AddChild(child Component) {
	child.SetParent(i)
	i.children = append(i.children, child)
}

func (i *Item) UpdateExpressions() (int, ErrorGroup) {
	var errs ErrorGroup
	var sum int
	if i.width.ShouldEvaluate() {
		sum++
		err := i.width.Update(i)
		if err != nil {
			errs.Add(newExpressionError("Item", "width", i.width.Expression, err))
		}
	}
	if i.height.ShouldEvaluate() {
		sum++
		err := i.height.Update(i)
		if err != nil {
			errs.Add(newExpressionError("Item", "height", i.height.Expression, err))
		}
	}
	if i.stuff.ShouldEvaluate() {
		sum++
		err := i.stuff.Update(i)
		if err != nil {
			errs.Add(newExpressionError("Item", "stuff", i.stuff.Expression, err))
		}
	}
	if i.x.ShouldEvaluate() {
		sum++
		err := i.x.Update(i)
		if err != nil {
			errs.Add(newExpressionError("Item", "x", i.x.Expression, err))
		}
	}
	if i.y.ShouldEvaluate() {
		sum++
		err := i.y.Update(i)
		if err != nil {
			errs.Add(newExpressionError("Item", "y", i.y.Expression, err))
		}
	}
	if i.z.ShouldEvaluate() {
		sum++
		err := i.z.Update(i)
		if err != nil {
			errs.Add(newExpressionError("Item", "z", i.z.Expression, err))
		}
	}
	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	for name, prop := range i.Root.properties {
		if prop.ShouldEvaluate() {
			sum++
			err := prop.Update(i)
			if err != nil {
				errs.Add(newExpressionError("Item", name, *prop.GetExpression(), err))
			}
		}
	}
	n, moreErrs := i.Root.UpdateExpressions()
	errs.AddGroup(moreErrs)
	sum += n
	return sum, errs
}

func (i *Item) ID() string {
	return i.id
}

func (i *Item) Anchors() Anchors {
	return i.anchors
}

func (i *Item) SetAnchors() Anchors {
	return i.anchors
}
