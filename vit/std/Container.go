package std

import (
	"fmt"

	vit "github.com/omniskop/vitrum/vit"
)

type Container struct {
	Item
	id string

	content *vit.StaticListValue[*vit.ComponentDefValue]

	children []vit.Component
}

func NewContainer(id string, scope vit.ComponentContainer) *Container {
	return &Container{
		Item:    *NewItem(id, scope),
		id:      id,
		content: vit.NewStaticListValue[*vit.ComponentDefValue](nil, nil),
	}
}

func (r *Container) String() string {
	return fmt.Sprintf("Container(%s)", r.id)
}

func (r *Container) Property(key string) (vit.Value, bool) {
	switch key {
	case "content":
		return r.content, true
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

func (r *Container) SetProperty(key string, value interface{}, position *vit.PositionRange) bool {
	switch key {
	case "content":
		if prop, ok := value.(vit.PropertyDefinition); ok {
			values := make([]*vit.ComponentDefValue, len(prop.Components))
			for i, component := range prop.Components {
				values[i] = vit.NewComponentDefValue(component, position)
			}
			r.content.Set(values)
		} else {
			r.content.Set(nil)
		}
	default:
		return r.Item.SetProperty(key, value, position)
	}
	return true
}

func (r *Container) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case r.id:
		return r, true
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

func (r *Container) UpdateExpressions() (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if r.content.ShouldEvaluate() {
		sum++
		r.content.Update(nil)
		for _, child := range r.children {
			r.RootC().RemoveChild(child)
		}
		r.children = r.children[:]
		for _, def := range r.content.Items {
			comp, err := r.InstantiateInScope(def.ComponentDefinition())
			if err != nil {
				errs.Add(err)
				continue
			}
			r.AddChild(comp)
			r.children = append(r.children, comp)
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

func (r *Container) ID() string {
	return r.id
}

func (r *Container) Finish() error {
	return r.RootC().FinishInContext(r)
}
