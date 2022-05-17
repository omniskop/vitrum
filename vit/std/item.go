package std

import (
	"fmt"

	vit "github.com/omniskop/vitrum/vit"
)

type Item struct {
	vit.Root
	id string

	width   vit.IntValue
	height  vit.IntValue
	stuff   vit.IntValue // TODO: delete
	anchors vit.Anchors
	x       vit.FloatValue
	y       vit.FloatValue
	z       vit.FloatValue
}

func NewItem(id string, scope vit.ComponentContainer) *Item {
	return &Item{
		Root:   vit.NewRoot(id, scope),
		id:     id,
		width:  *vit.NewIntValue("", nil),
		height: *vit.NewIntValue("", nil),
		stuff:  *vit.NewIntValue("", nil),
		x:      *vit.NewFloatValue("", nil),
		y:      *vit.NewFloatValue("", nil),
		z:      *vit.NewFloatValue("", nil),
	}
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
	case "stuff":
		return &i.stuff, true
	// TODO: fix this
	// case "vit.Anchors":
	// 	return &i.vit.Anchors, true
	case "x":
		return &i.x, true
	case "y":
		return &i.y, true
	case "z":
		return &i.z, true
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
	case "stuff":
		i.stuff.Expression.ChangeCode(value.(string), position)
	case "vit.Anchors":
		i.anchors = value.(vit.Anchors)
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
	case "vit.Anchors":
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

func (i *Item) AddChild(child vit.Component) {
	child.SetParent(i)
	i.Root.AddChildButKeepParent(child)
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
	}
	if i.height.ShouldEvaluate() {
		sum++
		err := i.height.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "height", i.id, i.height.Expression, err))
		}
	}
	if i.stuff.ShouldEvaluate() {
		sum++
		err := i.stuff.Update(i)
		if err != nil {
			errs.Add(vit.NewExpressionError("Item", "stuff", i.id, i.stuff.Expression, err))
		}
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
	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	n, err := i.UpdatePropertiesInContext(i)
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

func (i *Item) Anchors() vit.Anchors {
	return i.anchors
}

func (i *Item) SetAnchors() vit.Anchors {
	return i.anchors
}

func (i *Item) Finish() error {
	return i.RootC().FinishInContext(i)
}
