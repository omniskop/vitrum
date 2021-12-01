package vit

import "fmt"

type Item struct {
	Root

	width   IntValue
	stuff   IntValue
	anchors Anchors
}

func NewItem(id string) *Item {
	return &Item{
		Root:  NewRoot(id),
		width: NewIntValue(),
		stuff: NewIntValue(),
	}
}

func (i *Item) String() string {
	return fmt.Sprintf("Item{%s}", i.id)
}

func (i *Item) Property(key string) (interface{}, bool) {
	switch key {
	case "width":
		return i.width.Value, true
	case "stuff":
		return i.stuff.Value, true
	case "anchors":
		return &i.anchors, true
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

func (i *Item) SetProperty(key string, value interface{}) bool {
	fmt.Printf("[Item] set %q: %v\n", key, value)
	switch key {
	case "width":
		i.width.Expression.ChangeCode(value.(string))
	case "stuff":
		i.stuff.Expression.ChangeCode(value.(string))
	case "anchors":
		i.anchors = value.(Anchors)
	default:
		return i.Root.SetProperty(key, value)
	}
	return true
}

func (i *Item) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case i.id:
		return i, true
	case "width":
		return i.width, true
	case "stuff":
		return i.stuff, true
	}

	return i.Root.ResolveVariable(key)
}

func (i *Item) AddChild(child Component) {
	child.SetParent(i)
	i.children = append(i.children, child)
}

func (i *Item) UpdateExpressions() (int, error) {
	var sum int
	if i.width.ShouldEvaluate() {
		sum++
		err := i.width.Update(i)
		if err != nil {
			return sum, fmt.Errorf("evaluating 'Item.width': %w", err)
		}
	}
	if i.stuff.ShouldEvaluate() {
		sum++
		err := i.stuff.Update(i)
		if err != nil {
			return sum, fmt.Errorf("evaluating 'Item.stuff': %w", err)
		}
	}
	n, err := i.Root.UpdateExpressions()
	sum += n
	return sum, err
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
