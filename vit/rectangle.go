package vit

import (
	"fmt"
)

type Rectangle struct {
	Item
	id string

	color *ColorValue
}

func NewRectangle(id string, scope ComponentContainer) *Rectangle {
	return &Rectangle{
		Item:  *NewItem(id, scope),
		id:    id,
		color: NewColorValue("", nil),
	}
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle{%s}", r.id)
}

func (r *Rectangle) Property(key string) (Value, bool) {
	switch key {
	case "color":
		return r.color, true
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

func (r *Rectangle) MustProperty(key string) Value {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Rectangle) SetProperty(key string, value interface{}, position *PositionRange) bool {
	switch key {
	case "color":
		r.color.Expression.ChangeCode(value.(string), position)
	// case "width":
	// 	r.width.Expression.ChangeCode(value.(string))
	// case "stuff":
	// 	r.stuff.Expression.ChangeCode(value.(string))
	// case "anchors":
	// 	r.anchors = value.(Anchors)
	default:
		return r.Item.SetProperty(key, value, position)
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

func (r *Rectangle) UpdateExpressions() (int, ErrorGroup) {
	var sum int
	var errs ErrorGroup
	if r.color.ShouldEvaluate() {
		sum++
		err := r.color.Update(r)
		if err != nil {
			errs.Add(newExpressionError("Rectangle", "color", r.id, r.color.Expression, err))
		}
	}

	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	for name, prop := range r.Item.properties {
		if prop.ShouldEvaluate() {
			sum++
			err := prop.Update(r)
			if err != nil {
				errs.Add(newExpressionError("Rectangle", name, r.id, *prop.GetExpression(), err))
			}
		}
	}
	n, err := r.Item.UpdateExpressions()
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (r *Rectangle) ID() string {
	return r.id
}

func (r *Rectangle) finish() error {
	for _, props := range r.properties {
		if alias, ok := props.(*AliasValue); ok {
			err := alias.Update(r)
			if err != nil {
				return fmt.Errorf("alias error: %w", err)
			}
		}
	}

	for _, child := range r.children {
		err := child.finish()
		if err != nil {
			return err
		}
	}

	return nil
}
