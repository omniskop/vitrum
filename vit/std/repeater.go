package std

import (
	"fmt"

	vit "github.com/omniskop/vitrum/vit"
)

type RepeaterItem struct {
	Component vit.Component
	Key       interface{}
	Value     interface{}
}

// type Repeater struct {
// 	Item
// 	id string

// 	count *vit.IntValue
// 	model *vit.AnyValue

// 	items []RepeaterItem
// }

// func NewRepeater(id string, scope vit.ComponentContainer) *Repeater {
// 	return &Repeater{
// 		Item:  *NewItem(id, scope),
// 		id:    id,
// 		count: vit.NewIntValue("", nil),
// 		model: vit.NewAnyValue("", nil),
// 	}
// }

// func (r *Repeater) String() string {
// 	return fmt.Sprintf("Repeater{%s}", r.id)
// }

// func (r *Repeater) Property(key string) (vit.Value, bool) {
// 	switch key {
// 	case "count":
// 		return r.count, true
// 	case "model":
// 		return r.model, true
// 	default:
// 		return r.Item.Property(key)
// 	}
// }

// func (r *Repeater) MustProperty(key string) vit.Value {
// 	v, ok := r.Property(key)
// 	if !ok {
// 		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
// 	}
// 	return v
// }

// func (r *Repeater) SetProperty(key string, value interface{}, position *vit.PositionRange) bool {
// 	switch key {
// 	case "count":
// 		r.count.ChangeCode(value.(string), position)
// 	case "model":
// 		r.model.ChangeCode(value.(string), position)
// 	default:
// 		return r.Item.SetProperty(key, value, position)
// 	}
// 	return true
// }

// func (r *Repeater) ResolveVariable(key string) (interface{}, bool) {
// 	switch key {
// 	case r.id:
// 		return r, true
// 	case "count":
// 		return r.count, true
// 	case "model":
// 		return r.model, true
// 	default:
// 		return r.Item.ResolveVariable(key)
// 	}
// }

// func (r *Repeater) AddChild(child vit.Component) {
// 	child.SetParent(r)
// 	r.AddChildButKeepParent(child)
// }

// func (r *Repeater) UpdateExpressions() (int, vit.ErrorGroup) {
// 	var sum int
// 	var errs vit.ErrorGroup

// 	if r.count.ShouldEvaluate() {
// 		sum++
// 		err := r.count.Update(r)
// 		if err != nil {
// 			errs.Add(vit.NewExpressionError("Repeater", "count", r.id, r.count.Expression, err))
// 		}
// 		r.DoStuff()
// 	}
// 	if r.model.ShouldEvaluate() {
// 		sum++
// 		err := r.model.Update(r)
// 		if err != nil {
// 			errs.Add(vit.NewExpressionError("Repeater", "model", r.id, r.model.Expression, err))
// 		}
// 		r.DoStuff()
// 	}

// 	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
// 	n, err := r.UpdatePropertiesInContext(r)
// 	sum += n
// 	errs.AddGroup(err)
// 	n, err = r.Item.UpdateExpressions()
// 	sum += n
// 	errs.AddGroup(err)
// 	return sum, errs
// }

// func (r *Repeater) ID() string {
// 	return r.id
// }

func (r *Repeater) DoStuff() error {
	if len(r.Children()) == 0 {
		return nil
	}

	// model, err := r.interpretModel()
	// if err != nil {
	// 	return err
	// }

	// r.items = r.items[:]
	// for _, m := range model {
	// 	r.items = append(r.items, RepeaterItem{
	// 		Component: r.Children()[0],
	// 		Key:       m.key,
	// 		Value:     m.value,
	// 	})
	// }

	return nil
}

func (e *Repeater) ItemAt(index int) (vit.Component, bool) {
	// if index < 0 || index >= len(e.items) {
	// 	return nil, false
	// }

	// return e.items[index].Component, true
	return nil, false
}

func (e *Repeater) Count() int {
	// return len(e.items)
	return 0
}

type repeaterModel struct {
	key   interface{}
	value interface{}
}

// int
// {1: 1, 2: 2, 3: 3}
// []string
// {1: "eins", 2: "zwei", 3: "drei"}
// map[string]string
//

func (r *Repeater) interpretModel() ([]repeaterModel, error) {
	var out []repeaterModel
	switch value := r.model.Value.(type) {
	case int64:
		for i := int64(0); i < value; i++ {
			out = append(out, repeaterModel{i, i})
		}
	case []string:
		for i := 0; i < len(value); i++ {
			out = append(out, repeaterModel{i, value[i]})
		}
	case map[string]string:
		var i int
		for key, value := range value {
			out = append(out, repeaterModel{key, value})
			i++
		}
	default:
		return out, fmt.Errorf("unsupported model type '%T'", value)
	}
	return out, nil
}
