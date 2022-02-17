package vit

import (
	"errors"
	"fmt"

	"github.com/omniskop/vitrum/vit/script"
)

// indicates that an expression was not evaluated fully because it read from another expression that has been marked as dirty
var unsettledDependenciesError = errors.New("unsettled dependencies")

type Expression struct {
	code         string
	dirty        bool
	dependencies map[Value]bool
	dependents   map[*Expression]bool
	program      script.Script
	position     *PositionRange
	err          error
}

func NewExpression(code string, position *PositionRange) *Expression {
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
		position:     position,
		err:          nil,
	}
}

func (e *Expression) Evaluate(context Component) (interface{}, error) {
	// fmt.Printf("[expression] evaluating %q from %v\n", e.code, e.position)
	collector := NewAccessCollector(context)
	val, err := e.program.Run(collector)
	variables := collector.GetReadValues()
	// fmt.Printf("[expression] expression %q read from expressions:\n", e.code)
	var dontStoreValue bool
	for _, variable := range variables {
		// fmt.Printf("\t%v\n", variable.GetExpression().code)
		if variable.ShouldEvaluate() {
			// fmt.Printf("[expression] this expression is dirty, we will not update out value for now\n")
			dontStoreValue = true
		}
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
	if dontStoreValue {
		return nil, unsettledDependenciesError
	}
	e.dirty = false
	if err != nil {
		return nil, err
	}
	// NOTE: Currently we don't update other expression that depend on out value.
	// This is due to the fact that this should already have happened when this expression was marked as dirty.
	// In the future this might turn out to not be sufficient and if dependents will be marked as dirty in here in the future we should
	// consider removing that part of the code in MakeDirty.
	// Cases where that change might be necessary could be:
	//   - if this expression surprisingly directly sets the value of another expression
	//   - if this expression uses volatile vallues like the current time which would evaluate differently each time without a previous call to MakeDirty.
	//     Altough that would beg the question of why this expression would've been reevaluated in the first place.
	//     And this sounds like a very special case that would need to be handled specifically anyways.
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

func (e *Expression) ChangeCode(code string, position *PositionRange) {
	e.program, e.err = script.NewScript("expression", code)
	if e.err != nil {
		fmt.Printf("[expression] code error %q: %v\r\n", code, e.err)
		return
	}
	e.position = position
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
	case Component: // reference to an existing component instance
		return c.SubContext(actual), true
	case AbstractComponent: // static component values
		return &variableConverter{actual}, true
	case Enumeration:
		return &variableConverter{actual}, true
	case Value:
		(*c.readValues)[actual] = true // mark as read
		return actual.GetValue(), true
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64, bool, string:
		return actual, true
	default:
		panic(script.Exception(fmt.Sprintf("resolved variable %q to unhandled type %T", key, actual)))
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

type variableConverter struct {
	context script.VariableSource
}

func (c *variableConverter) ResolveVariable(key string) (interface{}, bool) {
	variable, ok := c.context.ResolveVariable(key)
	if !ok {
		return nil, false
	}

	switch actual := variable.(type) {
	case AbstractComponent: // static component values
		return &variableConverter{actual}, true
	case Enumeration:
		return &variableConverter{actual}, true
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64, bool, string:
		return actual, true
	default:
		panic(script.Exception(fmt.Sprintf("resolved variable %q to unhandled type %T", key, actual)))
	}
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

func castFloat(val interface{}) (float64, bool) {
	switch n := val.(type) {
	case int64:
		return float64(n), true
	case float64:
		return float64(n), true
	default:
		return 0, false
	}
}

type ExpressionError struct {
	ComponentName string
	ComponentID   string
	PropertyName  string
	Code          string
	Position      *PositionRange
	err           error
}

func newExpressionError(componentName string, propertyName string, componentID string, expression Expression, err error) ExpressionError {
	return ExpressionError{
		ComponentName: componentName,
		ComponentID:   componentID,
		PropertyName:  propertyName,
		Code:          expression.code,
		Position:      expression.position,
		err:           err,
	}
}

func (e ExpressionError) Error() string {
	var identifier string
	if e.ComponentID == "" {
		identifier = fmt.Sprintf("%s.%s", e.ComponentName, e.PropertyName)
	} else {
		identifier = fmt.Sprintf("%s(%s).%s", e.ComponentName, e.ComponentID, e.PropertyName)
	}
	if e.Position != nil {
		return fmt.Sprintf("%v: %s: expression %q: %v", e.Position, identifier, e.Code, e.err)
	}
	return fmt.Sprintf("%s: expression %q: %v", identifier, e.Code, e.err)
}

func (e ExpressionError) Is(target error) bool {
	_, ok := target.(ExpressionError)
	return ok
}

func (e ExpressionError) Unwrap() error {
	return e.err
}
