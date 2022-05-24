package vit

type FakeExpression struct {
	dirty        bool
	dependencies map[Value]bool
	dependents   map[Dependent]bool
	code         func(Component) (interface{}, error)
}

func (e *FakeExpression) Evaluate(context Component) (interface{}, error) {
	return e.code(context)
}

func (e *FakeExpression) ShouldEvaluate() bool {
	return e.dirty
}

func (e *FakeExpression) MakeDirty([]*Expression) {
	e.dirty = true
}

func (e *FakeExpression) ChangeCode(code string, position *PositionRange) {
}

func (e *FakeExpression) ClearDependencies() {
	for exp := range e.dependencies {
		exp.RemoveDependent(e)
	}
	e.dependencies = make(map[Value]bool)
}

func (e *FakeExpression) NotifyDependents(stack []*Expression) {
	for exp := range e.dependents {
		exp.MakeDirty(stack)
	}
}

func (e *FakeExpression) IsConstant() {

}

func (e *FakeExpression) GetExpression() *Expression {
	return nil
}

func (e *FakeExpression) AddDependent(dep Dependent) {
	e.dependents[dep] = true
}

func (e *FakeExpression) RemoveDependent(dep Dependent) {
	delete(e.dependents, dep)
}

func (e *FakeExpression) Err() error {
	return nil
}
