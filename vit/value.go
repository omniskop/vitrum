package vit

import (
	"fmt"
	"strconv"
	"strings"
)

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

func (c *IntValue) GetValue() interface{} {
	return c.Value
}

// ========================================= Float Value ===========================================

type FloatValue struct {
	Expression
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

func (v *AliasValue) AddDependent(exp *Expression) {}

func (v *AliasValue) RemoveDependent(exp *Expression) {}

func (v *AliasValue) ShouldEvaluate() bool {
	return v.other == nil
}

func (v *AliasValue) Err() error {
	return nil
}
