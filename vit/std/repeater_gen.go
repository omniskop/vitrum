// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
)

type Repeater struct {
	Item
	id string

	count *vit.IntValue
	model *vit.AnyValue
	items *vit.ListValue[*vit.ComponentValue]
}

func NewRepeater(id string, scope vit.ComponentContainer) *Repeater {
	return &Repeater{
		Item:  *NewItem(id, scope),
		id:    id,
		count: vit.NewIntValue("", nil),
		model: vit.NewAnyValue("", nil),
		items: vit.NewListValue[*vit.ComponentValue]("", nil),
	}
}

func (r *Repeater) String() string {
	return fmt.Sprintf("Repeater(%s)", r.id)
}

func (r *Repeater) Property(key string) (vit.Value, bool) {
	switch key {
	case "count":
		return r.count, true
	case "model":
		return r.model, true
	case "items":
		return r.items, true
	default:
		return r.Item.Property(key)
	}
}

func (r *Repeater) MustProperty(key string) vit.Value {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Repeater) SetProperty(key string, value interface{}, position *vit.PositionRange) bool {
	switch key {
	case "count":
		r.count.ChangeCode(value.(string), position)
	case "model":
		r.model.ChangeCode(value.(string), position)
	case "items":
		r.items.ChangeCode(value.(string), position)
	default:
		return r.Item.SetProperty(key, value, position)
	}
	return true
}

func (r *Repeater) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case r.id:
		return r, true
	case "count":
		return r.count, true
	case "model":
		return r.model, true
	case "items":
		return r.items, true
	default:
		return r.Item.ResolveVariable(key)
	}
}

func (r *Repeater) AddChild(child vit.Component) {
	child.SetParent(r)
	r.AddChildButKeepParent(child)
}

func (r *Repeater) UpdateExpressions() (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if r.count.ShouldEvaluate() {
		sum++
		err := r.count.Update(r)
		if err != nil {
			errs.Add(vit.NewExpressionError("Repeater", "count", r.id, r.count.Expression, err))
		}
	}
	if r.model.ShouldEvaluate() {
		sum++
		err := r.model.Update(r)
		if err != nil {
			errs.Add(vit.NewExpressionError("Repeater", "model", r.id, r.model.Expression, err))
		}
	}
	if r.items.ShouldEvaluate() {
		sum++
		err := r.items.Update(r)
		if err != nil {
			errs.Add(vit.NewExpressionError("Repeater", "items", r.id, r.items.Expression, err))
		}
	}

	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	n, err := r.UpdatePropertiesInContext(r)
	sum += n
	errs.AddGroup(err)
	n, err = r.Item.UpdateExpressions()
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (r *Repeater) As(target *vit.Component) bool {
	if _, ok := (*target).(*Repeater); ok {
		*target = r
		return true
	}
	return r.Item.As(target)
}

func (r *Repeater) ID() string {
	return r.id
}

func (r *Repeater) Finish() error {
	return r.RootC().FinishInContext(r)
}