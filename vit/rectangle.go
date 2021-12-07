package vit

import (
	"fmt"
)

type Rectangle struct {
	Item

	color ColorValue
}

func NewRectangle(id string, scope ComponentResolver) *Rectangle {
	return &Rectangle{
		Item: *NewItem(id, scope),
	}
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle{%s}", r.id)
}

func (r *Rectangle) Property(key string) (interface{}, bool) {
	switch key {
	case "color":
		return r.color.Value, true
	// case "width":
	// 	return r.width.Value, true
	// case "stuff":
	// 	return r.stuff.Value, true
	// case "anchors":
	// 	return &r.anchors, true
	default:
		return r.Item.Property(key)
	}
}

func (r *Rectangle) MustProperty(key string) interface{} {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Rectangle) SetProperty(key string, value interface{}) bool {
	fmt.Printf("[Rectangle] set %q: %v\n", key, value)
	switch key {
	case "color":
		r.color.Expression.ChangeCode(value.(string))
	// case "width":
	// 	r.width.Expression.ChangeCode(value.(string))
	// case "stuff":
	// 	r.stuff.Expression.ChangeCode(value.(string))
	// case "anchors":
	// 	r.anchors = value.(Anchors)
	default:
		return r.Item.SetProperty(key, value)
	}
	return true
}

func (r *Rectangle) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case r.id:
		return r, true
	case "color":
		return r.color.Value, true
		// case "width":
		// 	return r.width, true
		// case "stuff":
		// 	return r.stuff, true
	}

	return r.Item.ResolveVariable(key)
}

func (r *Rectangle) AddChild(child Component) {
	child.SetParent(r)
	r.children = append(r.children, child)
}

func (r *Rectangle) UpdateExpressions() (int, error) {
	var sum int
	if r.color.ShouldEvaluate() {
		sum++
		err := r.color.Update(r)
		if err != nil {
			return sum, fmt.Errorf("evaluating 'Rectangle.color' as %q: %w", r.color.Expression.code, err)
		}
	}
	// if r.stuff.ShouldEvaluate() {
	// 	sum++
	// 	err := r.stuff.Update(r)
	// 	if err != nil {
	// 		return sum, fmt.Errorf("evaluating 'Rectangle.stuff': %w", err)
	// 	}
	// }

	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	for name, prop := range r.Item.properties {
		if prop.ShouldEvaluate() {
			sum++
			err := prop.Update(r)
			if err != nil {
				return sum, fmt.Errorf("evaluating custom property %q: %w", name, err)
			}
		}
	}
	n, err := r.Item.UpdateExpressions()
	sum += n
	return sum, err
}

func (r *Rectangle) ID() string {
	return r.id
}
