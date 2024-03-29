package std

import (
	"fmt"

	vit "github.com/omniskop/vitrum/vit"
)

type Container struct {
	Item
	id string

	content vit.ComponentDefListValue

	children []vit.Component
}

func NewContainer(id string, context *vit.FileContext) *Container {
	return &Container{
		Item:    *NewItem(id, context),
		id:      id,
		content: *vit.NewEmptyComponentDefListValue(),
	}
}

func (r *Container) String() string {
	return fmt.Sprintf("Container(%s)", r.id)
}

func (r *Container) Property(key string) (vit.Value, bool) {
	switch key {
	case "content":
		return &r.content, true
	default:
		return r.Item.Property(key)
	}
}

func (r *Container) MustProperty(key string) vit.Value {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Container) SetProperty(key string, value interface{}) error {
	switch key {
	case "content":
		err := r.content.SetValue(value.([]vit.ComponentDefinition))
		if err != nil {
			return vit.NewPropertyError("container", key, r.id, err)
		}
	default:
		return r.Item.SetProperty(key, value)
	}
	return nil
}

func (r *Container) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "content":
		r.content.SetComponentDefinitions(nil)
	default:
		return r.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (r *Container) Event(name string) (vit.Listenable, bool) {
	return r.Item.Event(name)
}

func (r *Container) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "content":
		return r.content, true
	default:
		return r.Item.ResolveVariable(key)
	}
}

func (r *Container) AddChild(child vit.Component) {
	child.SetParent(r)
	r.AddChildButKeepParent(child)
}

func (r *Container) Children() []vit.Component {
	return r.Item.Children()
}

func (r *Container) UpdateExpressions(context vit.Component) (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup
	if context == nil {
		context = r
	}

	if changed, _ := r.content.Update(context); changed {
		sum++
		for _, child := range r.children {
			r.RootC().RemoveChild(child)
		}
		r.children = r.children[:]
		for _, def := range r.content.ComponentDefinitions() {
			comp, err := r.RootC().InstantiateInScope(def)
			if err != nil {
				errs.Add(err)
				continue
			}
			r.AddChild(comp)
			r.children = append(r.children, comp)
		}
	}

	n, err := r.Item.UpdateExpressions(context)
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (r *Container) ID() string {
	return r.id
}

func (r *Container) Finish() error {
	return r.RootC().FinishInContext(r)
}
