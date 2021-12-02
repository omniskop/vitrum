package vit

import (
	"fmt"

	"github.com/omniskop/vitrum/vit/script"
)

/*
	Eine Wert kann entweder eine Konstante sein oder eine expression.
*/

type Value interface {
	Update(context Component) error
	GetValue() interface{}
	MakeDirty([]*Expression)
	GetExpression() *Expression
	AddDependent(*Expression)
	RemoveDependent(*Expression)
	ShouldEvaluate() bool
	Err() error
}

type IntValue struct {
	Expression
	Value int
}

func NewEmptyIntValue() *IntValue {
	return &IntValue{
		Expression: *NewExpression("0"),
	}
}

func NewIntValue(expression string) *IntValue {
	return &IntValue{
		Expression: *NewExpression(expression),
	}
}

func (v *IntValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		return err
	}
	castVal, ok := castInt(val)
	if !ok {
		return fmt.Errorf("expression did not evaluate to expected type int but %T instead", val)
	}
	v.Value = int(castVal)
	return nil
}

func (c *IntValue) GetValue() interface{} {
	return c.Value
}

type Expression struct {
	code         string
	dirty        bool
	dependencies map[Value]bool
	dependents   map[*Expression]bool
	program      script.Script
	err          error
}

func NewExpression(code string) *Expression {
	prog, err := script.NewScript("expression", code)
	if err != nil {
		panic(err)
	}
	return &Expression{
		code:         code,
		dirty:        true,
		dependencies: make(map[Value]bool),
		dependents:   make(map[*Expression]bool),
		program:      prog,
		err:          nil,
	}
}

func (e *Expression) Evaluate(context Component) (interface{}, error) {
	fmt.Printf("[expression] evaluating %q\n", e.code)
	collector := NewAccessCollector(context)
	val, err := e.program.Run(collector)
	variables := collector.GetReadValues()
	// fmt.Printf("[expression] expression %q read from:\n", e.code)
	for _, variable := range variables {
		// fmt.Printf("\t%v\n", variable.GetExpression().code)
		if _, ok := e.dependencies[variable]; !ok {
			e.dependencies[variable] = true
			variable.AddDependent(e)
		}
	}
	variables = collector.GetWrittenValues()
	// fmt.Printf("[expression] expression %q wrote to:\n", e.code)
	for _, variable := range variables {
		// fmt.Printf("\t%v\n", variable.GetExpression().code)
		if _, ok := e.dependents[variable.GetExpression()]; !ok {
			e.dependents[variable.GetExpression()] = true
		}
	}
	e.dirty = false
	if err != nil {
		return nil, fmt.Errorf("expression %q failed: %v", e.code, err)
	}
	return val, nil
}

// ShouldEvaluate returns true if this expression should be reevaluated because any dependencies have changed
func (e *Expression) ShouldEvaluate() bool {
	return e.dirty
}

func (e *Expression) MakeDirty(stack []*Expression) {
	for _, exp := range stack {
		if exp == e {
			panic("circular dependency detected")
		}
	}
	e.dirty = true
	e.NotifyDependents(append(stack, e))
}

func (e *Expression) ChangeCode(code string) {
	e.program, e.err = script.NewScript("expression", code)
	if e.err != nil {
		fmt.Printf("[expression] code error %q: %v\n", code, e.err)
		return
	}
	e.code = code
	e.ClearDependencies()
	e.dirty = true
	e.NotifyDependents(nil)
}

func (e *Expression) ClearDependencies() {
	for exp := range e.dependencies {
		exp.RemoveDependent(e)
	}
	e.dependencies = make(map[Value]bool)
}

// NotifyDependents informs all expressions that depend on the result of this expression that they will need to be reevaluated
func (e *Expression) NotifyDependents(stack []*Expression) {
	for exp := range e.dependents {
		exp.MakeDirty(stack)
	}
}

func (e *Expression) IsConstant() bool {
	return len(e.dependencies) == 0
}

func (e *Expression) GetExpression() *Expression {
	return e
}

func (e *Expression) AddDependent(exp *Expression) {
	e.dependents[exp] = true
}

func (e *Expression) RemoveDependent(exp *Expression) {
	delete(e.dependents, exp)
}

func (e *Expression) Err() error {
	return e.err
}

type AccessCollector struct {
	context       Component
	readValues    *map[Value]bool
	writtenValues *map[Value]bool
}

func NewAccessCollector(context Component) *AccessCollector {
	r := make(map[Value]bool)
	w := make(map[Value]bool)
	return &AccessCollector{
		context:       context,
		readValues:    &r,
		writtenValues: &w,
	}
}

func (c *AccessCollector) SubContext(context Component) *AccessCollector {
	return &AccessCollector{
		context:       context,
		readValues:    c.readValues,
		writtenValues: c.writtenValues,
	}
}

func (c *AccessCollector) ResolveVariable(key string) (interface{}, bool) {
	variable, ok := c.context.ResolveVariable(key)
	if !ok {
		return nil, false
	}

	switch actual := variable.(type) {
	case Component:
		return script.VariableBridge{Source: c.SubContext(actual)}, true
	case Value:
		(*c.readValues)[actual] = true
		return actual.GetValue(), true
	case IntValue:
		var intf Value = &actual
		(*c.readValues)[intf] = true
		return actual.Value, true
	default:
		panic(fmt.Errorf("resolved variable %q to unhandled type %T", key, actual))
	}
}

func (c *AccessCollector) GetReadValues() []Value {
	result := make([]Value, 0, len(*c.readValues))
	for val := range *c.readValues {
		result = append(result, val)
	}
	return result
}

func (c *AccessCollector) GetWrittenValues() []Value {
	result := make([]Value, 0, len(*c.writtenValues))
	for val := range *c.writtenValues {
		result = append(result, val)
	}
	return result
}

func castInt(val interface{}) (int, bool) {
	switch n := val.(type) {
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	default:
		return 0, false
	}
}
