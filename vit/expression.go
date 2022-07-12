package vit

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/dop251/goja"
	"github.com/omniskop/vitrum/vit/script"
)

// indicates that an expression was not evaluated fully because it read from another expression that has been marked as dirty
var unsettledDependenciesError = errors.New("unsettled dependencies")

// Expression contains JavaScript code that can be executed to get a value
type Expression struct {
	code         string
	dirty        bool
	dependencies map[Value]bool // values that are required by this expression
	program      script.Script
	position     *PositionRange
	err          error
}

func NewExpression(code string, position *PositionRange) *Expression {
	// The parenthesis around the code are needed to make sure we get the correct value from all expressions.
	// For example objects (e.g. {one: 1}) would return the number '1' instead of a map.
	prog, err := script.NewScript("expression", fmt.Sprintf("(%s)", code))
	if err != nil {
		err = NewExpressionError(code, position, err)
	}
	if position != nil {
		// adjust position in a way to hide the fact that we added the parenthesis around the code
		p := position.StartColumnShifted(-1)
		position = &p
	}
	return &Expression{
		code:         code,
		dirty:        true,
		dependencies: make(map[Value]bool),
		program:      prog,
		position:     position,
		err:          err,
	}
}

func (e *Expression) Evaluate(context Component) (interface{}, error) {
	if e.err != nil {
		return nil, e.err
	}

	// fmt.Printf("[expression] evaluating %q from %v\n", e.code, e.position)
	collector := NewAccessCollector(context)
	val, err := e.program.Run(collector)
	variables := collector.GetReadValues()
	// fmt.Printf("[expression] expression %q read from expressions:\n", e.code)
	var dontStoreValue bool
	for _, variable := range variables {
		// fmt.Printf("\t%v\n", variable.GetExpression().code)

		// TODO: figure out how to do this after the value refactor
		// if expr, ok := variable.(ExpressionValue); ok && expr.ShouldEvaluate() {
		// 	// fmt.Printf("[expression] this expression is dirty, we will not update out value for now\n")
		// 	dontStoreValue = true
		// }
		if _, ok := e.dependencies[variable]; !ok {
			e.dependencies[variable] = true
			variable.AddDependent(e)
		}
	}
	variables = collector.GetWrittenValues()
	// fmt.Printf("[expression] expression %q wrote to:\n", e.code)
	// TODO: do something with these?
	_ = variables

	if dontStoreValue {
		return nil, unsettledDependenciesError
	}
	e.dirty = false
	if err != nil {
		return nil, err
	}
	// NOTE: Currently we don't update other expression that depend on our value.
	// This is due to the fact that this should already have happened when this expression was marked as dirty.
	// In the future this might turn out to not be sufficient and if dependents will be marked as dirty in here in the future we should
	// consider removing that part of the code in MakeDirty.
	// Cases where that change might be necessary could be:
	//   - if this expression surprisingly directly sets the value of another expression
	//   - if this expression uses volatile values like the current time which would evaluate differently each time without a previous call to MakeDirty.
	//     Although that would beg the question of why this expression would've been reevaluated in the first place.
	//     And this sounds like a very special case that would need to be handled specifically anyways.
	return val, nil
}

// ShouldEvaluate returns true if this expression should be reevaluated because any dependencies have changed
func (e *Expression) ShouldEvaluate() bool {
	if e == nil {
		return false
	}
	return e.dirty
}

func (e *Expression) MakeDirty(stack []Dependent) {
	for _, exp := range stack {
		if exp == e {
			panic("circular dependency detected")
		}
	}
	e.dirty = true
}

func (e *Expression) ChangeCode(code string, position *PositionRange) {
	e.program, e.err = script.NewScript("expression", fmt.Sprintf("(%s)", code))
	if e.err != nil {
		fmt.Printf("[expression] code error %q: %v\r\n", code, e.err)
		e.err = NewExpressionError(code, position, e.err)
		return
	}
	e.position = position
	e.code = code
	e.ClearDependencies()
	e.dirty = true
}

func (e *Expression) ClearDependencies() {
	for exp := range e.dependencies {
		exp.RemoveDependent(e)
	}
	e.dependencies = make(map[Value]bool)
}

func (e *Expression) IsConstant() bool {
	return len(e.dependencies) == 0
}

func (e *Expression) GetExpression() *Expression {
	return e
}

func (e *Expression) Err() error {
	return e.err
}

type AccessCollector struct {
	context       Component // TODO: could this be a script.variableSource?
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
	case *Method:
		return actual, true
	case Listenable:
		return EventAdapter{c.context, actual}, true
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

type EventAdapter struct {
	context Component
	event   Listenable
}

func (a EventAdapter) AddEventListener(f *Method) {
	a.event.AddListenerFunction(&f.AsyncFunction)
}

func castList[ElementType Value](val interface{}) ([]ElementType, bool) {
	switch list := val.(type) {
	case []Value:
		result := make([]ElementType, len(list))
		for i, v := range list {
			result[i] = v.(ElementType)
		}
		return result, true
	case []interface{}:
		result := make([]ElementType, len(list))
		for i, v := range list {
			result[i] = v.(ElementType)
		}
		return result, true
	default:
		return nil, false
	}
}

type ExpressionError struct {
	Code     string
	Position *PositionRange
	err      error
}

var compilerErrorRegex = regexp.MustCompile(`Line (\d+):(\d+) (.*)`)

func NewExpressionError(code string, pos *PositionRange, wrappedErr error) ExpressionError {
	// check if we can get some more information from the error
	var cErr *goja.CompilerSyntaxError
	if errors.As(wrappedErr, &cErr) {
		// example:
		// "expression: Line 3:13 Unexpected identifier (and 5 more errors)"
		matches := compilerErrorRegex.FindStringSubmatch(cErr.CompilerError.Message)
		if len(matches) > 0 {
			line, err := strconv.Atoi(matches[1])
			column, err2 := strconv.Atoi(matches[2])
			if err == nil && err2 == nil {
				if pos != nil {
					if line == 1 {
						pos.StartColumn += column - 1
					} else {
						pos.StartLine += line - 1
						pos.StartColumn = column
					}
					pos.SetEnd(pos.Start())
				}
				wrappedErr = errors.New(matches[3])
			}
		}
	}

	return ExpressionError{
		Code:     code,
		Position: pos,
		err:      wrappedErr,
	}
}

func (e ExpressionError) Error() string {
	return fmt.Sprintf("expression %q: %v", e.Code, e.err)
}

func (e ExpressionError) Is(target error) bool {
	_, ok := target.(ExpressionError)
	return ok
}

func (e ExpressionError) Unwrap() error {
	return e.err
}
