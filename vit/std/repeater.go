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

func (r *Repeater) evaluateInternals(interface{}) error {
	if r.delegate.GetValue() == nil {
		// No delegate available to instantiate
		// delete existing items
		for _, item := range r.items {
			r.RemoveChild(item.Component)
		}
		return nil
	}

	model, err := r.interpretModel()
	if err != nil {
		return err
	}

	compDef := r.delegate.ComponentDefinition()

	// TODO: be smart about which items actually need to be recreated

	// remove existing items
	for _, item := range r.items {
		r.RemoveChild(item.Component)
	}

	// create items
	r.items = r.items[:]
	for _, m := range model {
		instance, err := r.InstantiateInScope(compDef)
		if err != nil {
			return err
		}
		r.items = append(r.items, RepeaterItem{
			Component: instance,
			Key:       m.key,
			Value:     m.value,
		})
		if r.Parent() != nil {
			// TODO: i think this parent check could be handled better.
			//   maybe even top most components could have a sort of 'fake' parent?

			// TODO: the components should be inserted immediately after the repeater itself
			r.Parent().AddChildAfter(r, instance)
		}
		instance.UpdateExpressions()
	}

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
	switch value := r.model.GetValue().(type) {
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
