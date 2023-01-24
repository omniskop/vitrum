package vit

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/omniskop/vitrum/vit/script"
)

// TODO: Change all values to only update when their value actually changed (either through expression or directly)

type Value interface {
	GetValue() interface{}      // returns the current value in it's natural type
	AddDependent(Dependent)     // adds a dependent that should be notified about changes to this value
	RemoveDependent(Dependent)  // removes a dependent
	SetValue(interface{}) error // changes this value to a new one. Might return an error if the type is incorrect
	SetCode(Code)               // changes this value to the result of an expression
	Update(context Component) (bool, error)
}

type Dependent interface {
	MakeDirty([]Dependent)
}

// FunctionDependency is an adapter that implements the Dependent interface
// and simply calls the callback function when the dependency changes.
type FunctionDependency struct {
	Callback *func()
}

func FuncDep(cb func()) FunctionDependency {
	return FunctionDependency{
		Callback: &cb,
	}
}

func (d FunctionDependency) MakeDirty([]Dependent) {
	(*d.Callback)()
}

type baseValue struct {
	dependents map[Dependent]bool
}

func newBaseValue() baseValue {
	return baseValue{
		dependents: make(map[Dependent]bool),
	}
}

func (v *baseValue) AddDependent(d Dependent) {
	v.dependents[d] = true
}

func (v *baseValue) RemoveDependent(d Dependent) {
	delete(v.dependents, d)
}

func (v *baseValue) notifyDependents(stack []Dependent) {
	for d := range v.dependents {
		d.MakeDirty(stack)
	}
}

func newValueForType(vitType string, code Code) (Value, error) {
	switch vitType {
	case "int":
		return NewIntValueFromCode(code), nil
	case "float":
		return NewFloatValueFromCode(code), nil
	case "string":
		return NewStringValueFromCode(code), nil
	case "bool":
		return NewBoolValueFromCode(code), nil
	case "alias":
		return NewAliasValueFromCode(code), nil
	case "component":
		return NewComponentRefValueFromCode(code), nil
	case "var":
		return NewAnyValueFromCode(code), nil
	case "componentdef":
		return nil, fmt.Errorf("unable to create ComponentDefValue from code")
	}
	return nil, UnknownTypeError{vitType}
}

func newValueFromGo(value interface{}) (Value, error) {
	switch value := value.(type) {
	case int:
		return NewIntValue(value), nil
	case uint:
		return NewIntValue(int(value)), nil
	case int8:
		return NewIntValue(int(value)), nil
	case uint8:
		return NewIntValue(int(value)), nil
	case int16:
		return NewIntValue(int(value)), nil
	case uint16:
		return NewIntValue(int(value)), nil
	case int32:
		return NewIntValue(int(value)), nil
	case uint32:
		return NewIntValue(int(value)), nil
	case int64:
		return NewIntValue(int(value)), nil
	case uint64:
		return NewIntValue(int(value)), nil
	case float32:
		return NewFloatValue(float64(value)), nil
	case float64:
		return NewFloatValue(value), nil
	case string:
		return NewStringValue(value), nil
	case bool:
		return NewBoolValue(value), nil
	default:
		if reflect.ValueOf(value).Kind() == reflect.Func {
			return NewFunctionValue(value)
		}
		return nil, fmt.Errorf("unable to create value from %T", value)
	}
}

func vitTypeFromGo(value interface{}) (string, bool) {
	switch value.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return "int", true
	case float32, float64:
		return "float", true
	case string:
		return "string", true
	case bool:
		return "bool", true
	default:
		return "", false
	}
}

type typeError struct {
	expectedType string
	actualType   string
}

// newTypeError returns a new type error that will determine the actually received type by itself
func newTypeError(expected string, actualValue interface{}) typeError {
	return typeError{
		expectedType: expected,
		actualType:   fmt.Sprintf("%T", actualValue),
	}
}

func (e typeError) Error() string {
	return fmt.Sprintf("evaluated to %s but %s was expected", e.actualType, e.expectedType)
}

func (e typeError) Is(target error) bool {
	_, ok := target.(typeError)
	return ok
}

// ========================================= List Value ============================================

type ListValue[ElementType Value] struct {
	baseValue
	value      []ElementType
	expression *Expression
}

func NewListValueFromCode[ElementType Value](code Code) *ListValue[ElementType] {
	return &ListValue[ElementType]{
		baseValue:  newBaseValue(),
		expression: NewExpression(code),
	}
}

func NewListValue[ElementType Value](value []ElementType) *ListValue[ElementType] {
	return &ListValue[ElementType]{
		baseValue: newBaseValue(),
		value:     value,
	}
}

func NewEmptyListValue[ElementType Value]() *ListValue[ElementType] {
	return &ListValue[ElementType]{
		baseValue: newBaseValue(),
		value:     make([]ElementType, 0),
	}
}

func (v *ListValue[ElementType]) GetValue() interface{} {
	var out []interface{}
	for _, element := range v.value {
		out = append(out, element.GetValue())
	}
	return out
}

func (v *ListValue[ElementType]) Slice() []ElementType {
	return v.value
}

func (v *ListValue[ElementType]) SetValue(value interface{}) error {
	if value == nil {
		v.value = make([]ElementType, 0)
		v.notifyDependents(nil)
		return nil
	} else if slice, ok := value.([]ElementType); ok {
		v.value = slice
		v.notifyDependents(nil)
		return nil
	}
	return newTypeError(fmt.Sprintf("slice of %T", *new(ElementType)), value)
}

func (v *ListValue[ElementType]) SetSlice(slice []ElementType) {
	v.value = slice
	v.notifyDependents(nil)
}

func (v *ListValue[ElementType]) SetCode(code Code) {
	v.expression = NewExpression(code)
	v.notifyDependents([]Dependent{v.expression})
}

func (v *ListValue[ElementType]) Update(context Component) (bool, error) {
	if v.expression == nil {
		return false, nil
	}
	if !v.expression.ShouldEvaluate() {
		return false, nil
	}
	val, err := v.expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return false, nil
		}
		return false, err
	}
	castVal, ok := castList[ElementType](val)
	if !ok {
		return false, newTypeError(fmt.Sprintf("list of type %T", *new(ElementType)), val)
	}
	// TODO: check if the list actually changed?
	v.value = castVal
	v.notifyDependents(nil)
	return true, nil
}

// ========================================== Int Value ============================================

type IntValue struct {
	baseValue
	value      int
	expression *Expression
}

func NewIntValueFromCode(code Code) *IntValue {
	return &IntValue{
		baseValue:  newBaseValue(),
		value:      0,
		expression: NewExpression(code),
	}
}

func NewIntValue(value int) *IntValue {
	return &IntValue{
		baseValue: newBaseValue(),
		value:     value,
	}
}

func NewEmptyIntValue() *IntValue {
	return &IntValue{
		baseValue: newBaseValue(),
		value:     0,
	}
}

func (v *IntValue) GetValue() interface{} {
	return v.value
}

func (v *IntValue) Int() int {
	return v.value
}

func (v *IntValue) SetValue(newValue interface{}) error {
	if intVal, ok := castInt(newValue); ok {
		v.value = intVal
		v.expression = nil
		v.notifyDependents(nil) // as this is a fixed value there is no need to add ourself to the stack
		return nil
	}
	return newTypeError(fmt.Sprintf("number"), newValue)
}

func (v *IntValue) SetIntValue(newValue int) {
	v.value = newValue
	v.expression = nil
	v.notifyDependents(nil) // as this is a fixed value there is no need to add ourself to the stack
}

func (v *IntValue) SetCode(code Code) {
	v.expression = NewExpression(code)
	v.notifyDependents([]Dependent{v.expression})
}

func (v *IntValue) Update(context Component) (bool, error) {
	if v.expression == nil {
		return false, nil
	}
	if !v.expression.ShouldEvaluate() {
		return false, nil
	}
	val, err := v.expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return false, nil
		}
		return false, err
	}
	castVal, ok := castInt(val)
	if !ok {
		return false, newTypeError("number", val)
	}
	if v.value != castVal {
		v.value = castVal
		v.notifyDependents(nil)
	}
	return true, nil
}

func castInt(val interface{}) (int, bool) {
	switch n := val.(type) {
	case uint:
		return int(n), true
	case int:
		return n, true
	case uint8:
		return int(n), true
	case uint16:
		return int(n), true
	case uint32:
		return int(n), true
	case uint64:
		return int(n), true
	case int8:
		return int(n), true
	case int16:
		return int(n), true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	case float32:
		return int(n), true
	case float64:
		return int(n), true
	default:
		return 0, false
	}
}

// ========================================= Float Value ===========================================

type FloatValue struct {
	baseValue
	value      float64
	expression *Expression
}

func NewFloatValueFromCode(code Code) *FloatValue {
	return &FloatValue{
		baseValue:  newBaseValue(),
		value:      0,
		expression: NewExpression(code),
	}
}

func NewFloatValue(value float64) *FloatValue {
	return &FloatValue{
		baseValue: newBaseValue(),
		value:     value,
	}
}

func NewEmptyFloatValue() *FloatValue {
	return &FloatValue{
		baseValue: newBaseValue(),
		value:     0,
	}
}

func (v *FloatValue) GetValue() interface{} {
	return v.value
}

func (v *FloatValue) Float64() float64 {
	return v.value
}

func (v *FloatValue) SetValue(newValue interface{}) error {
	if floatVal, ok := castFloat64(newValue); ok {
		v.value = floatVal
		v.expression = nil
		v.notifyDependents(nil) // as this is a fixed value there is no need to add ourself to the stack
		return nil
	}
	return newTypeError(fmt.Sprintf("number"), newValue)
}

func (v *FloatValue) SetFloatValue(newValue float64) {
	v.value = newValue
	v.expression = nil
	v.notifyDependents(nil) // as this is a fixed value there is no need to add ourself to the stack
}

func (v *FloatValue) SetCode(code Code) {
	v.expression = NewExpression(code)
	v.notifyDependents([]Dependent{v.expression})
}

func (v *FloatValue) Update(context Component) (bool, error) {
	if v.expression == nil {
		return false, nil
	}
	if !v.expression.ShouldEvaluate() {
		return false, nil
	}
	val, err := v.expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return false, nil
		}
		return false, err
	}
	castVal, ok := castFloat64(val)
	if !ok {
		return false, newTypeError("number", val)
	}
	if v.value != castVal {
		v.value = castVal
		v.notifyDependents(nil)
	}
	return true, nil
}

func castFloat64(val interface{}) (float64, bool) {
	switch n := val.(type) {
	case float32:
		return float64(n), true
	case float64:
		return n, true
	case uint:
		return float64(n), true
	case int:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

// ======================================== String Value ===========================================

type StringValue struct {
	baseValue
	value      string
	expression *Expression
}

func NewStringValueFromCode(code Code) *StringValue {
	return &StringValue{
		baseValue:  newBaseValue(),
		value:      "",
		expression: NewExpression(code),
	}
}

func NewStringValue(value string) *StringValue {
	return &StringValue{
		baseValue: newBaseValue(),
		value:     value,
	}
}

func NewEmptyStringValue() *StringValue {
	return &StringValue{
		baseValue: newBaseValue(),
		value:     "",
	}
}

func (v *StringValue) GetValue() interface{} {
	return v.value
}

func (v *StringValue) String() string {
	return v.value
}

func (v *StringValue) SetValue(newValue interface{}) error {
	if strVal, ok := castString(newValue); ok {
		v.value = strVal
		v.expression = nil
		v.notifyDependents(nil) // as this is a fixed value there is no need to add ourself to the stack
		return nil
	}
	return newTypeError("string", newValue)
}

func (v *StringValue) SetStringValue(newValue string) {
	v.value = newValue
	v.expression = nil
	v.notifyDependents(nil)
}

func (v *StringValue) SetCode(code Code) {
	v.expression = NewExpression(code)
	v.notifyDependents([]Dependent{v.expression})
}

func (v *StringValue) Update(context Component) (bool, error) {
	if v.expression == nil {
		return false, nil
	}
	if !v.expression.ShouldEvaluate() {
		return false, nil
	}
	val, err := v.expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return false, nil
		}
		return false, err
	}
	strVal, ok := castString(val)
	if !ok {
		return false, newTypeError("string", val)
	}
	if v.value != strVal {
		v.value = strVal
		v.notifyDependents([]Dependent{v.expression})
	}
	return true, nil
}

// Position returns the PositionRange that this value was defined at.
// The boolean indicates wether a position is set.
func (v *StringValue) Position() (*PositionRange, bool) {
	if v.expression == nil {
		return nil, false
	}
	if v.expression.position == nil {
		return nil, false
	}
	return v.expression.position, true
}

func castString(value interface{}) (string, bool) {
	switch actual := value.(type) {
	case string:
		return actual, true
	case int64:
		return strconv.FormatInt(actual, 10), true
	case float64:
		return strconv.FormatFloat(actual, 'f', 64, 10), true
	}
	return "", false
}

// ========================================= Bool Value ============================================

type BoolValue struct {
	baseValue
	value      bool
	expression *Expression
}

func NewBoolValueFromCode(code Code) *BoolValue {
	return &BoolValue{
		baseValue:  newBaseValue(),
		value:      false,
		expression: NewExpression(code),
	}
}

func NewBoolValue(value bool) *BoolValue {
	return &BoolValue{
		baseValue: newBaseValue(),
		value:     value,
	}
}

func NewEmptyBoolValue() *BoolValue {
	return &BoolValue{
		baseValue: newBaseValue(),
		value:     false,
	}
}

func (v *BoolValue) GetValue() interface{} {
	return v.value
}

func (v *BoolValue) Bool() bool {
	return v.value
}

func (v *BoolValue) SetValue(newValue interface{}) error {
	if boolVal, ok := castBool(newValue); ok {
		v.value = boolVal
		v.expression = nil
		v.notifyDependents(nil) // as this is a fixed value there is no need to add ourself to the stack
		return nil
	}
	return newTypeError("boolean", newValue)
}

func (v *BoolValue) SetBoolValue(newValue bool) {
	v.expression = nil
	if v.value != newValue {
		v.value = newValue
		v.notifyDependents(nil) // as this is a fixed value there is no need to add ourself to the stack
	}
}

func (v *BoolValue) SetCode(code Code) {
	v.expression = NewExpression(code)
	v.notifyDependents([]Dependent{v.expression})
}

func (v *BoolValue) Update(context Component) (bool, error) {
	if v.expression == nil {
		return false, nil
	}
	if !v.expression.ShouldEvaluate() {
		return false, nil
	}
	val, err := v.expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return false, nil
		}
		return false, err
	}
	boolVal, ok := castBool(val)
	if !ok {
		return false, newTypeError("boolean", val)
	}
	if v.value != boolVal {
		v.value = boolVal
		v.notifyDependents(nil)
	}
	return true, nil
}

func castBool(value interface{}) (bool, bool) {
	switch actual := value.(type) {
	case bool:
		return actual, true
	case int64:
		return actual != 0, true
	case float64:
		return actual != 0, true
	case string:
		return len(actual) > 0, true
	}
	return false, false
}

// ========================================= Alias Value ===========================================

// Alias value points to a different value, potentially of another component
type AliasValue struct {
	baseValue
	expression string
	position   *PositionRange
	fileCtx    *FileContext
	other      Value
	changed    bool
}

func NewAliasValueFromCode(code Code) *AliasValue {
	return &AliasValue{
		baseValue:  newBaseValue(),
		expression: code.Code,
		position:   code.Position,
		fileCtx:    code.FileCtx,
		other:      nil,
		changed:    true,
	}
}

func (v *AliasValue) GetValue() interface{} {
	if v.other == nil {
		return nil
	}
	return v.other.GetValue()
}

func (v *AliasValue) SetValue(newValue interface{}) error {
	// TODO: Should this update the alias itself or the aliased valued?
	if v.other != nil {
		return v.other.SetValue(newValue)
	}
	return nil
}

func (v *AliasValue) SetCode(code Code) {
	v.expression = code.Code
	v.position = code.Position
	if v.other != nil {
		v.other.RemoveDependent(v)
	}
	v.other = nil
	v.changed = true
}

func (v *AliasValue) Update(context Component) (bool, error) {
	if v.other != nil {
		return false, nil
	}
	if v.expression == "" {
		// TODO: add position to error
		return false, NewExpressionError(v.expression, v.position, fmt.Errorf("expression is empty"))
	}
	parts := strings.Split(v.expression, ".")
	var comp Component = context
	var currentProperty Value
	// find component using the id in the expression
	if strings.Contains(parts[0], " ") {
		return false, NewExpressionError(v.expression, v.position, fmt.Errorf("invalid alias reference: %q", v.expression))
	}
	var ok bool
	comp, ok = v.fileCtx.IDs[parts[0]]
	if !ok {
		return false, NewExpressionError(v.expression, v.position, fmt.Errorf("unable to resolve alias reference: %q", v.expression))
	}
	parts = parts[1:]

	// find property using the remaining parts
	for _, part := range parts {
		val, ok := comp.Property(part)
		if !ok {
			return false, NewExpressionError(v.expression, v.position, fmt.Errorf("unable to resolve alias reference: %q", v.expression))
		}
		currentProperty = val
	}

	// nothing found
	if currentProperty == nil {
		return false, NewExpressionError(v.expression, v.position, fmt.Errorf("unable to resolve alias reference: %q", v.expression))
	}
	// referenced itself
	if currentProperty == v {
		return false, NewExpressionError(v.expression, v.position, fmt.Errorf("alias cannot reference itself: %q", v.expression))
	}

	v.other = currentProperty // saving this also marks the alias as updated, preventing an infinite loop in the next check

	// if we referenced another alias we need will update that as well and make sure there are no circular references
	if otherAlias, ok := currentProperty.(*AliasValue); ok {
		_, err := otherAlias.Update(comp)
		if err != nil {
			return false, NewExpressionError(v.expression, v.position, fmt.Errorf("error in nested alias update: %w", err))
		}
		if yes, chain := isAliasRecursive(v, nil); yes {
			return false, NewExpressionError(v.expression, v.position, fmt.Errorf("alias contains circular reference: %v", formatAliasChain(chain)))
		}
	}

	v.other.AddDependent(v)

	return true, nil
}

func (v *AliasValue) MakeDirty(stack []Dependent) {
	v.notifyDependents(append(stack, v))
}

func isAliasRecursive(alias *AliasValue, chain []*AliasValue) (bool, []*AliasValue) {
	if subAlias, ok := alias.other.(*AliasValue); ok {
		for _, a := range chain {
			if a == subAlias {
				return true, append(chain, alias, subAlias)
			}
		}
		return isAliasRecursive(subAlias, append(chain, alias))
	}

	return false, nil
}

func formatAliasChain(chain []*AliasValue) string {
	var steps []string
	for _, a := range chain {
		steps = append(steps, fmt.Sprintf("%q", a.expression))

	}
	return strings.Join(steps, " -> ")
}

// ========================================== Any Value ============================================

type AnyValue struct {
	baseValue
	value      interface{}
	expression *Expression
}

func NewAnyValueFromCode(code Code) *AnyValue {
	return &AnyValue{
		baseValue:  newBaseValue(),
		value:      nil,
		expression: NewExpression(code),
	}
}

func NewAnyValue(value interface{}) *AnyValue {
	return &AnyValue{
		baseValue:  newBaseValue(),
		value:      value,
		expression: nil,
	}
}

func NewEmptyAnyValue() *AnyValue {
	return &AnyValue{
		baseValue:  newBaseValue(),
		value:      nil,
		expression: nil,
	}
}

func (v *AnyValue) GetValue() interface{} {
	return v.value
}

func (v *AnyValue) SetValue(value interface{}) error {
	v.value = value
	v.expression = nil
	v.notifyDependents([]Dependent{nil}) // as this is a fixed value there is no need to add ourself to the stack
	return nil
}

func (v *AnyValue) SetCode(code Code) {
	v.expression = NewExpression(code)
	v.notifyDependents([]Dependent{v.expression})
}

func (v *AnyValue) Update(context Component) (bool, error) {
	if v.expression == nil {
		return false, nil
	}
	if !v.expression.ShouldEvaluate() {
		return false, nil
	}
	val, err := v.expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return false, nil
		}
		return false, err
	}
	if v.value != val {
		v.value = val
		v.notifyDependents(nil)
	}
	return true, nil
}

// ================================= Component Definition Value ====================================

type ComponentDefValue struct {
	baseValue
	value   *ComponentDefinition
	context *FileContext
	changed bool
	err     error
}

func NewComponentDefValue(component *ComponentDefinition, context *FileContext) *ComponentDefValue {
	return &ComponentDefValue{
		baseValue: newBaseValue(),
		value:     component,
		context:   context,
		changed:   true,
	}
}

func NewEmptyComponentDefValue() *ComponentDefValue {
	return &ComponentDefValue{
		baseValue: newBaseValue(),
		value:     nil,
		changed:   false,
	}
}

func (v *ComponentDefValue) GetValue() interface{} {
	return v.value
}

func (v *ComponentDefValue) ComponentDefinition() *ComponentDefinition {
	return v.value
}

func (v *ComponentDefValue) Context() *FileContext {
	return v.context
}

func (v *ComponentDefValue) SetValue(newValue interface{}) error {
	switch compDef := newValue.(type) {
	case ComponentDefinitionInContext:
		v.value = compDef.ComponentDefinition
		v.context = compDef.Context
		v.changed = true
		v.err = nil
		v.notifyDependents(nil)
	case *ComponentDefinition:
		v.value = compDef
		v.context = nil
		v.changed = true
		v.err = nil
		v.notifyDependents(nil)
	default:
		return newTypeError("component definition", newValue)
	}
	return nil
}

func (v *ComponentDefValue) SetComponentDefinition(component *ComponentDefinition, context *FileContext) {
	v.value = component
	v.context = context
	v.changed = true
	v.err = nil
	v.notifyDependents(nil)
}

func (v *ComponentDefValue) SetCode(code Code) {
	panic("must not call SetExpression on ComponentDefValue")
}

func (v *ComponentDefValue) Update(context Component) (bool, error) {
	changed := v.changed
	v.changed = false
	return changed, v.err
}

// =============================== Component Definition List Value =================================

type ComponentDefListValue struct {
	baseValue
	components []*ComponentDefinition
	changed    bool
	err        error
}

func NewComponentDefListValue(components []*ComponentDefinition, position *PositionRange) *ComponentDefListValue {
	return &ComponentDefListValue{
		baseValue:  newBaseValue(),
		components: components,
		changed:    true,
		err:        nil,
	}
}

func NewEmptyComponentDefListValue() *ComponentDefListValue {
	return &ComponentDefListValue{
		baseValue:  newBaseValue(),
		components: nil,
		changed:    true,
		err:        nil,
	}
}

func (v *ComponentDefListValue) GetValue() interface{} {
	return v.components
}

func (v *ComponentDefListValue) ComponentDefinitions() []*ComponentDefinition {
	return v.components
}

func (v *ComponentDefListValue) SetValue(newValue interface{}) error {
	if compDefs, ok := newValue.([]*ComponentDefinition); ok {
		v.components = compDefs
		v.changed = true
		v.err = nil
		v.notifyDependents(nil)
		return nil
	}
	return newTypeError("slice of component definitions", newValue)
}

func (v *ComponentDefListValue) SetComponentDefinitions(components []*ComponentDefinition) {
	v.components = components
	v.changed = true
	v.notifyDependents(nil)
}

func (v *ComponentDefListValue) SetCode(code Code) {
	panic("must not call SetExpression on ComponentDefListValue")
}

func (v *ComponentDefListValue) Update(context Component) (bool, error) {
	changed := v.changed
	v.changed = false
	v.notifyDependents(nil)
	return changed, v.err
}

// ======================================= Optional Value ==========================================

type OptionalValue[T Value] struct {
	baseValue
	value   T
	isSet   bool
	changed bool
}

func NewOptionalValue[T Value](v T) *OptionalValue[T] {
	return &OptionalValue[T]{
		baseValue: newBaseValue(),
		value:     v,
	}
}

func (v *OptionalValue[T]) GetValue() interface{} {
	if v.isSet {
		return v.value.GetValue()
	}
	return nil
}

func (v *OptionalValue[T]) IsSet() bool {
	return v.isSet
}

// returns the actual value weather it is set or not.
// This should only be used for reading. The wrapped value must not be written to directly through this.
func (v *OptionalValue[T]) Value() T {
	return v.value
}

func (v *OptionalValue[T]) SetValue(newValue interface{}) error {
	err := v.value.SetValue(newValue)
	if err != nil {
		return err
	}
	v.isSet = true
	v.changed = true
	v.notifyDependents(nil)
	return nil
}

func (v *OptionalValue[T]) SetCode(code Code) {
	v.value.SetCode(code)
	v.isSet = true
	v.changed = true
	v.notifyDependents([]Dependent{v})
}

func (v *OptionalValue[T]) Update(context Component) (bool, error) {
	if v.isSet {
		v.changed = false
		return v.value.Update(context)
	}
	// we keep track if the value was changed ourself because we wouldn't know otherwise if the value was unset
	changed := v.changed
	v.changed = false
	if changed {
		v.notifyDependents(nil)
	}
	return changed, nil
}

func (v *OptionalValue[T]) MakeDirty(stack []Dependent) {
	v.changed = true
	v.notifyDependents(append(stack, v))
}

// ================================== Component Reference Value ====================================

type ComponentRefValue struct {
	baseValue
	value      Component
	expression *Expression
}

func NewComponentRefValueFromCode(code Code) *ComponentRefValue {
	return &ComponentRefValue{
		baseValue:  newBaseValue(),
		value:      nil,
		expression: NewExpression(code),
	}
}

func NewComponentRefValue(comp Component) *ComponentRefValue {
	return &ComponentRefValue{
		baseValue: newBaseValue(),
		value:     comp,
	}
}

func NewEmptyComponentRefValue() *ComponentRefValue {
	return &ComponentRefValue{
		baseValue: newBaseValue(),
		value:     nil,
	}
}

func (v *ComponentRefValue) GetValue() interface{} {
	return v.value
}

func (v *ComponentRefValue) Component() Component {
	return v.value
}

func (v *ComponentRefValue) SetValue(newValue interface{}) error {
	if comp, ok := newValue.(Component); ok {
		v.value = comp
		v.expression = nil
		v.notifyDependents(nil)
		return nil
	}
	return newTypeError("component reference", newValue)
}

func (v *ComponentRefValue) SetComponent(comp Component) {
	v.value = comp
	v.expression = nil
	v.notifyDependents(nil)
}

func (v *ComponentRefValue) SetCode(code Code) {
	v.expression = NewExpression(code)
	v.notifyDependents([]Dependent{v.expression})
}

func (v *ComponentRefValue) Update(context Component) (bool, error) {
	if v.expression == nil {
		return false, nil
	}
	if !v.expression.ShouldEvaluate() {
		return false, nil
	}
	val, err := v.expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return false, nil
		}
		return false, err
	}

	bridge, ok := val.(*script.VariableBridge)
	if !ok {
		return false, newTypeError("component", val)
	}
	collector, ok := bridge.Source.(*AccessCollector)
	if !ok {
		return false, newTypeError("component", val)
	}
	component, ok := collector.context.(Component)
	if !ok {
		return false, newTypeError("component", val)
	}

	if v.value != component {
		v.value = component
		v.notifyDependents(nil)
	}
	return true, nil
}

// ========================================= Group Value ===========================================

type groupEntry struct {
	value       Value
	overwritten bool
}

type GroupValue struct {
	baseValue
	values     map[string]*groupEntry
	expression *Expression
}

// TODO: make parser smarter
// Currently the expression for a group will just be parsed in javascript but it would be much better to analyse the code ourselves.
// It would be faster und allow for much more control. (It would also make the code more consistent)
// But in order to do that we would need to access the parser.
// If that will ever change in the future I don't think we would need the work around with the groupEntry anymore.

func NewGroupValueFromCode(schema map[string]Value, code Code) *GroupValue {
	return &GroupValue{
		baseValue:  newBaseValue(),
		values:     createGroupEntries(schema),
		expression: NewExpression(code),
	}
}

func NewEmptyGroupValue(schema map[string]Value) *GroupValue {
	return &GroupValue{
		baseValue:  newBaseValue(),
		values:     createGroupEntries(schema),
		expression: nil,
	}
}

func createGroupEntries(schema map[string]Value) map[string]*groupEntry {
	entries := make(map[string]*groupEntry)
	for k, v := range schema {
		entries[k] = &groupEntry{value: v}
	}
	return entries
}

func (v *GroupValue) GetValue() interface{} {
	return v.values
}

func (v *GroupValue) Get(key string) (Value, bool) {
	value, ok := v.values[key]
	return value.value, ok
}

func (v *GroupValue) MustGet(key string) Value {
	value, ok := v.values[key]
	if !ok {
		panic(fmt.Sprintf("tried to read unknown key %q from group value", key))
	}
	return value.value
}

func (v *GroupValue) SetValue(newValue interface{}) error {
	var gErr ErrorGroup
	if valueMap, ok := newValue.(map[string]interface{}); ok {
		for key, value := range valueMap {
			if val, ok := v.values[key]; ok {
				err := val.value.SetValue(value)
				if err != nil {
					gErr.Add(err)
				}
			} else {
				gErr.Add(fmt.Errorf("unknown group key %q", key))
			}
		}
	}
	v.notifyDependents(nil)
	if gErr.Failed() {
		return gErr
	}
	return nil
}

func (v *GroupValue) SetCode(code Code) {
	v.expression = NewExpression(code)
	for _, value := range v.values {
		// disable overwrites
		// TODO: this disables all overwrites, even for properties that will not be set in this expression
		value.overwritten = false
	}
	v.notifyDependents([]Dependent{v.expression})
}

func (v *GroupValue) Update(context Component) (bool, error) {
	changed, errs := v.updateIndividualValues(context)
	defer func() {
		if changed {
			v.notifyDependents(nil)
		}
	}()

	if v.expression == nil {
		if errs.Failed() {
			return false, errs
		}
		return changed, nil
	}
	if !v.expression.ShouldEvaluate() {
		if errs.Failed() {
			return false, errs
		}
		return changed, nil
	}
	val, err := v.expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return changed, errs
		}
		errs.Add(err)
		return changed, errs
	}
	dataMap := val.(map[string]interface{})
	for key, jsValue := range dataMap {
		if entry, ok := v.values[key]; ok {
			if entry.overwritten {
				continue // don't overwrite values that were set specifically
			}
			err := entry.value.SetValue(jsValue)
			if err != nil {
				errs.Add(err)
			}
		} else {
			errs.Add(fmt.Errorf("unknown group key %q", key))
		}
	}

	if errs.Failed() {
		return true, errs
	}
	return true, nil
}

// updateIndividualValues updates all values. Even if they are not overwritten specifically.
// It returns true if at least one value was changed. It always returns an error group wether something went wrong or not.
func (v *GroupValue) updateIndividualValues(context Component) (bool, ErrorGroup) {
	var somethingChanged bool
	var errs ErrorGroup
	for _, value := range v.values {
		changed, err := value.value.Update(context)
		if err != nil {
			errs.Add(err)
		}
		if changed {
			somethingChanged = true
		}
	}
	return somethingChanged, errs
}

func (v *GroupValue) SetValueOf(name string, newValue interface{}) error {
	if value, ok := v.values[name]; ok {
		value.overwritten = true
		return value.value.SetValue(newValue)
	}
	return fmt.Errorf("unknown group key %q", name)
}

func (v *GroupValue) SetCodeOf(key string, code Code) error {
	if value, ok := v.values[key]; ok {
		value.overwritten = true
		value.value.SetCode(code)
		return nil
	}
	return fmt.Errorf("unknown group key %q", key)
}

// ======================================= Function Value ==========================================

// FunctionValue represents a function that is defined in go. JavaScript can only call it and not set or depend on it.
// Can be set from Go through the SetValue method.
type FunctionValue struct {
	fun interface{}
}

// NewFunctionValue returns a new function value containing the given function. If the parameter is not a function it returns with a TypeError.
func NewFunctionValue(value interface{}) (*FunctionValue, error) {
	// make sure this is actually a function
	if reflect.TypeOf(value).Kind() != reflect.Func {
		return nil, newTypeError("function", value)
	}
	return &FunctionValue{fun: value}, nil
}

// MustNewFunctionValue returns a new function value containing the given function. If the parameter is not a function it panics.
func MustNewFunctionValue(value interface{}) *FunctionValue {
	v, err := NewFunctionValue(value)
	if err != nil {
		panic(fmt.Errorf("called MustNewFunctionValue: %v", err))
	}
	return v
}

// NewEmptyFunctionValue returns a new function value that doesn't contain a function.
func NewEmptyFunctionValue() *FunctionValue {
	return &FunctionValue{}
}

func (v *FunctionValue) GetValue() interface{} {
	return v.fun
}

func (v *FunctionValue) AddDependent(Dependent) {}

func (v *FunctionValue) RemoveDependent(Dependent) {}

func (v *FunctionValue) SetValue(newValue interface{}) error {
	// make sure this is actually a function
	if reflect.TypeOf(newValue).Kind() != reflect.Func {
		return newTypeError("function", newValue)
	}
	v.fun = newValue
	return nil
}

func (v *FunctionValue) SetCode(Code) {}

func (v *FunctionValue) Update(context Component) (bool, error) {
	return false, nil
}
