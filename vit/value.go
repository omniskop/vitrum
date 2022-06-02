package vit

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/omniskop/vitrum/vit/script"
)

type Value interface {
	GetValue() interface{}
	SetFromProperty(PropertyDefinition)
	AddDependent(Dependent)
	RemoveDependent(Dependent)
}

type ExpressionValue interface {
	Value

	ShouldEvaluate() bool
	Update(context Component) error
	MakeDirty([]*Expression)
	GetExpression() *Expression
	Err() error
}

type Dependent interface {
	MakeDirty([]*Expression)
}

func ValueConstructorForType(vitType string, value interface{}, position *PositionRange) (Value, error) {
	switch vitType {
	case "string":
		return NewStringValue(value.(string), position), nil
	case "int":
		return NewIntValue(value.(string), position), nil
	case "float":
		return NewFloatValue(value.(string), position), nil
	case "bool":
		return NewBoolValue(value.(string), position), nil
	case "var":
		return NewAnyValue(value.(string), position), nil
	case "component":
		return NewComponentDefValue(value.(*ComponentDefinition), position), nil
	}
	return nil, UnknownTypeError{vitType}
}

type ChangeMonitor[T Value] struct {
	handlers []func(T)
}

func (m *ChangeMonitor[T]) OnChange(handler func(T)) {
	m.handlers = append(m.handlers, handler)
}

func (m *ChangeMonitor[T]) triggerChange(v T) {
	for _, handler := range m.handlers {
		handler(v)
	}
}

// ========================================= List Value ============================================

type ListValue[ElementType Value] struct {
	Expression
	Value []ElementType
}

func NewListValue[ElementType Value](expression string, position *PositionRange) *ListValue[ElementType] {
	v := new(ListValue[ElementType])
	if expression == "" {
		v.Expression = *NewExpression("[]", position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (c *ListValue[ElementType]) SetFromProperty(prop PropertyDefinition) {
	c.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *ListValue[ElementType]) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	castVal, ok := castList[ElementType](val)
	if !ok {
		return fmt.Errorf("did not evaluate to expected type list but %T instead", val)
	}
	v.Value = castVal
	return nil
}

func (v *ListValue[ElementType]) GetValue() interface{} {
	var out []interface{}
	for _, element := range v.Value {
		out = append(out, element.GetValue())
	}
	return out
}

// ========================================== Int Value ============================================

type IntValue struct {
	Expression
	Value int
}

func NewIntValue(expression string, position *PositionRange) *IntValue {
	v := new(IntValue)
	if expression == "" {
		v.Expression = *NewExpression("0", position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *IntValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *IntValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	castVal, ok := castInt(val)
	if !ok {
		return fmt.Errorf("did not evaluate to expected type int but %T instead", val)
	}
	v.Value = int(castVal)
	return nil
}

func (v *IntValue) GetValue() interface{} {
	return v.Value
}

func (v *IntValue) Int() int {
	return v.Value
}

// ========================================= Float Value ===========================================

type FloatValue struct {
	Expression
	ChangeMonitor[*FloatValue]
	Value float64
}

func NewFloatValue(expression string, position *PositionRange) *FloatValue {
	v := new(FloatValue)
	if expression == "" {
		v.Expression = *NewExpression("0", position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *FloatValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *FloatValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	castVal, ok := castFloat(val)
	if !ok {
		return fmt.Errorf("did not evaluate to expected type int but %T instead", val)
	}
	v.Value = float64(castVal)
	return nil
}

func (c *FloatValue) GetValue() interface{} {
	return c.Value
}

func (c *FloatValue) Float64() float64 {
	return c.Value
}

// ======================================== String Value ===========================================

type StringValue struct {
	Expression
	Value string
}

func NewStringValue(expression string, position *PositionRange) *StringValue {
	v := new(StringValue)
	if expression == "" {
		v.Expression = *NewExpression(`""`, position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *StringValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *StringValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	var ok bool
	v.Value, ok = convertJSValueToString(val)
	if !ok {
		return fmt.Errorf("did not evaluate to expected type string but %T instead", val)
	}
	return nil
}

func (c *StringValue) GetValue() interface{} {
	return c.Value
}

func convertJSValueToString(value interface{}) (string, bool) {
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
	Expression
	Value bool
}

func NewBoolValue(expression string, position *PositionRange) *BoolValue {
	v := new(BoolValue)
	if expression == "" {
		v.Expression = *NewExpression("false", position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *BoolValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *BoolValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	var ok bool
	v.Value, ok = convertJSValueToBool(val)
	if !ok {
		return fmt.Errorf("did not evaluate to expected type string but %T instead", val)
	}
	return nil
}

func (c *BoolValue) GetValue() interface{} {
	return c.Value
}

func convertJSValueToBool(value interface{}) (bool, bool) {
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
	Expression string
	Position   *PositionRange
	other      Value
}

func NewAliasValue(expression string, position *PositionRange) *AliasValue {
	v := new(AliasValue)
	v.Position = position
	v.Expression = expression
	return v
}

func (v *AliasValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression = prop.Expression
}

func (v *AliasValue) Update(context Component) error {
	if v.other != nil {
		return nil
	}
	if v.Expression == "" {
		// TODO: add position to error
		return fmt.Errorf("alias reference is empty")
	}
	parts := strings.Split(v.Expression, ".")
	var currentComponent Component = context
	var currentProperty Value
	// find component using the id's listed in the expression
	for {
		part := parts[0]

		if strings.Contains(part, " ") {
			return fmt.Errorf("invalid alias reference: %q", v.Expression)
		}

		if currentComponent.ID() == part {
			parts = parts[1:]
			continue // no change
		}
		if childComp, ok := currentComponent.ResolveID(part); ok {
			currentComponent = childComp
			parts = parts[1:]
			continue
		}
		break
	}
	// find property using the remaining parts
	for _, part := range parts {
		val, ok := currentComponent.Property(part)
		if !ok {
			return fmt.Errorf("unable to resolve alias reference: %q", v.Expression)
		}
		currentProperty = val
	}

	// nothing found
	if currentProperty == nil {
		return fmt.Errorf("unable to resolve alias reference: %q", v.Expression)
	}
	// referenced itself
	if currentProperty == v {
		return fmt.Errorf("alias cannot reference itself: %q", v.Expression)
	}

	v.other = currentProperty // saving this also marks the alias as updated, preventing an infinite loop in the next check

	// if we referenced another alias we need will update that as well and make sure there are no circular references
	if otherAlias, ok := currentProperty.(*AliasValue); ok {
		err := otherAlias.Update(currentComponent)
		if err != nil {
			return fmt.Errorf("error in nested alias update: %w", err)
		}
		if yes, chain := isAliasRecursive(v, nil); yes {
			return fmt.Errorf("alias contains circular reference: %v", formatAliasChain(chain))
		}
	}

	return nil
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
		steps = append(steps, fmt.Sprintf("%q", a.Expression))

	}
	return strings.Join(steps, " -> ")
}

func (v *AliasValue) GetValue() interface{} {
	if v.other == nil {
		return nil
	}
	return v.other.GetValue()
}

func (v *AliasValue) MakeDirty(stack []*Expression) {}

func (v *AliasValue) GetExpression() *Expression {
	return NewExpression(v.Expression, v.Position)
}

func (v *AliasValue) AddDependent(exp Dependent) {
	if v.other == nil {
		v.other.AddDependent(exp)
	}
}

func (v *AliasValue) RemoveDependent(exp Dependent) {
	if v.other == nil {
		v.other.RemoveDependent(exp)
	}
}

func (v *AliasValue) ShouldEvaluate() bool {
	expr, ok := v.other.(ExpressionValue) // if v.other is nil ok will be false
	return ok && expr.ShouldEvaluate()
}

func (v *AliasValue) Err() error {
	return nil
}

// ========================================== Any Value ============================================

type AnyValue struct {
	Expression
	Value interface{}
}

func NewAnyValue(expression string, position *PositionRange) *AnyValue {
	v := new(AnyValue)
	if expression == "" {
		v.Expression = *NewExpression("null", position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *AnyValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *AnyValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	v.Value = val
	return nil
}

func (c *AnyValue) GetValue() interface{} {
	return c.Value
}

// ================================= Component Definition Value ====================================

type ComponentDefValue struct {
	Value   *ComponentDefinition
	Changed bool
	err     error
}

func NewComponentDefValue(component *ComponentDefinition, position *PositionRange) *ComponentDefValue {
	return &ComponentDefValue{
		Value:   component,
		Changed: true,
	}
}

func (v *ComponentDefValue) ChangeComponent(component *ComponentDefinition) {
	v.Value = component
	v.Changed = true
	v.err = nil
}

func (v *ComponentDefValue) SetFromProperty(prop PropertyDefinition) {
	if len(prop.Components) == 0 {
		v.Value = nil
		v.err = nil
	} else if len(prop.Components) == 1 {
		v.Value = prop.Components[0]
		v.err = nil
	} else {
		v.Value = prop.Components[0]
		v.err = fmt.Errorf("cannot assign multiple components to a single component value at %s", prop.Pos.String())
	}
	v.Changed = true
}

func (v *ComponentDefValue) Update(context Component) error {
	v.Changed = false
	return v.err
}

func (v *ComponentDefValue) GetValue() interface{} {
	return v.Value
}

func (v *ComponentDefValue) MakeDirty(stack []*Expression) {
	v.Changed = true
}

func (v *ComponentDefValue) GetExpression() *Expression {
	return NewExpression("", nil)
}

func (v *ComponentDefValue) AddDependent(exp Dependent) {}

func (v *ComponentDefValue) RemoveDependent(exp Dependent) {}

func (v *ComponentDefValue) ShouldEvaluate() bool {
	return v.Changed
}

func (v *ComponentDefValue) Err() error {
	return nil
}

// =============================== Component Definition List Value =================================

type ComponentDefListValue struct {
	StaticBaseValue
	components []*ComponentDefinition
}

func NewComponentDefListValue(components []*ComponentDefinition, position *PositionRange) *ComponentDefListValue {
	return &ComponentDefListValue{
		StaticBaseValue: *NewStaticBaseValue(),
		components:      components,
	}
}

func (v *ComponentDefListValue) SetFromProperty(prop PropertyDefinition) {
	v.components = prop.Components
}

func (v *ComponentDefListValue) GetValue() interface{} {
	return v.components
}

func (v *ComponentDefListValue) ChangeComponents(components []*ComponentDefinition) {
	v.components = components
	v.Changed = true
}

// ====================================== Static Base Value ========================================

type StaticBaseValue struct {
	Changed      bool
	dependencies map[Value]bool
	dependents   map[Dependent]bool
}

func NewStaticBaseValue() *StaticBaseValue {
	return &StaticBaseValue{
		Changed:      true,
		dependencies: make(map[Value]bool),
		dependents:   make(map[Dependent]bool),
	}
}

func (v *StaticBaseValue) SetFromProperty(prop PropertyDefinition) {
	v.Changed = true
}

func (v *StaticBaseValue) Update(context Component) error {
	v.Changed = false
	return nil
}

func (v *StaticBaseValue) MakeDirty(stack []*Expression) {
	v.Changed = true
	// TODO: check the stack for circular dependencies
	for exp := range v.dependents {
		exp.MakeDirty(stack)
	}
}

func (v *StaticBaseValue) GetExpression() *Expression {
	return NewExpression("", nil)
}

func (v *StaticBaseValue) AddDependent(dep Dependent) {
	v.dependents[dep] = true
}

func (v *StaticBaseValue) RemoveDependent(exp Dependent) {
	delete(v.dependents, exp)
}

func (v *StaticBaseValue) ShouldEvaluate() bool {
	return v.Changed
}

func (v *StaticBaseValue) Err() error {
	return nil
}

type StaticListValue[ElementType Value] struct {
	StaticBaseValue
	Items []ElementType
}

func NewStaticListValue[ElementType Value](items []ElementType, position *PositionRange) *StaticListValue[ElementType] {
	return &StaticListValue[ElementType]{
		StaticBaseValue: *NewStaticBaseValue(),
		Items:           items,
	}
}

func (v *StaticListValue[ElementType]) GetValue() interface{} {
	return v.Items
}

func (v *StaticListValue[ElementType]) Set(value []ElementType) {
	v.Items = value
	v.Changed = true
}

// ======================================= Optional Value ==========================================

type OptionalValue[T Value] struct {
	ChangeMonitor[T]
	Value T
	isSet bool
}

func NewOptionalValue[T Value](v T) *OptionalValue[T] {
	return &OptionalValue[T]{
		Value: v,
	}
}

func (v *OptionalValue[T]) IsSet() bool {
	return v.isSet
}

func (v *OptionalValue[T]) SetFromProperty(prop PropertyDefinition) {
	v.Value.SetFromProperty(prop)
	v.isSet = true
}

func (v *OptionalValue[T]) Update(context Component) error {
	if expr, ok := Value(v.Value).(ExpressionValue); ok {
		return expr.Update(context)
	}
	return nil
}

func (v *OptionalValue[T]) GetValue() interface{} {
	if v.isSet {
		return v.Value.GetValue()
	}
	return nil
}

func (v *OptionalValue[T]) MakeDirty(stack []*Expression) {
	if expr, ok := Value(v.Value).(ExpressionValue); ok {
		expr.MakeDirty(stack)
	}
}

func (v *OptionalValue[T]) GetExpression() *Expression {
	// TODO: this won't change isSet
	if expr, ok := Value(v.Value).(ExpressionValue); ok {
		return expr.GetExpression()
	}
	return nil
}

func (v *OptionalValue[T]) ChangeCode(code string, position *PositionRange) {
	if expr, ok := Value(v.Value).(ExpressionValue); ok {
		expr.GetExpression().ChangeCode(code, position)
	}
	v.isSet = true
}

func (v *OptionalValue[T]) AddDependent(exp Dependent) {
	v.Value.AddDependent(exp)
}

func (v *OptionalValue[T]) RemoveDependent(exp Dependent) {
	v.Value.RemoveDependent(exp)
}

func (v *OptionalValue[T]) ShouldEvaluate() bool {
	expr, ok := Value(v.Value).(ExpressionValue)
	return ok && expr.ShouldEvaluate()
}

func (v *OptionalValue[T]) Err() error {
	if expr, ok := Value(v.Value).(ExpressionValue); ok {
		return expr.Err()
	}
	return nil
}

// ================================== Component Reference Value ====================================

type ComponentRefValue struct {
	Expression
	Value Component
}

func NewComponentRefValue(expression string, position *PositionRange) *IntValue {
	v := new(IntValue)
	if expression == "" {
		v.Expression = *NewExpression("", position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *ComponentRefValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *ComponentRefValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	// TODO: maybe check if casts are valid?
	v.Value = val.(*script.VariableBridge).Source.(*AccessCollector).context
	return nil
}

func (v *ComponentRefValue) GetValue() interface{} {
	return v.Value
}

func (v *ComponentRefValue) Component() Component {
	return v.Value
}

// ===================================== Static Float Value ========================================

type StaticFloatValue struct {
	StaticBaseValue
	value float64
}

func NewStaticFloatValue() *StaticFloatValue {
	return &StaticFloatValue{
		StaticBaseValue: *NewStaticBaseValue(),
	}
}

func (v *StaticFloatValue) GetValue() interface{} {
	return v.value
}

func (v *StaticFloatValue) Float64() float64 {
	return v.value
}

func (v *StaticFloatValue) Set(newValue float64) {
	v.value = newValue
}
