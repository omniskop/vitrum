package vit

type custom struct {
	Component
	RootComponent *Root
	id            string
	name          string
}

func NewCustom(id string, name string, parent Component) *custom {
	return &custom{
		Component:     parent,
		RootComponent: parent.RootC(),
		id:            id,
		name:          name,
	}
}

func (c *custom) ResolveVariable(key string) (interface{}, bool) {
	if key == c.id {
		return c, true
	}

	return c.Component.ResolveVariable(key)
}

func (c *custom) AddChild(child Component) {
	child.SetParent(c)
	c.RootComponent.children = append(c.RootComponent.children, child)
}

func (c *custom) UpdateExpressions() (int, ErrorGroup) {
	var errs ErrorGroup
	var sum int
	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	for name, prop := range c.RootComponent.properties {
		if expr, ok := prop.(ExpressionValue); ok && expr.ShouldEvaluate() {
			sum++
			err := expr.Update(c)
			if err != nil {
				errs.Add(NewExpressionError(c.name, name, c.id, *expr.GetExpression(), err))
			}
		}
	}

	s, e := c.Component.UpdateExpressions()
	sum += s
	errs.AddGroup(e)

	return sum, errs
}

func (c *custom) ID() string {
	return c.id
}

func (c *custom) As(target *Component) bool {
	if _, ok := (*target).(*custom); ok {
		*target = c
		return true
	}
	return c.Component.As(target)
}
