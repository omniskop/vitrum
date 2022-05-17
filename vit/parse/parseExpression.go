package parse

import (
	"strings"

	"github.com/omniskop/vitrum/vit"
)

type propertyValueType int

const (
	valueTypeComponent propertyValueType = iota
	valueTypeList
	valueTypeExpression
)

// a propertyValue is either a component definition, a list of more componentValues or a JavaScript expression
type propertyValue struct {
	valueType propertyValueType
	component *vit.ComponentDefinition
	list      []propertyValue
}

func newComponentValue(compDef *vit.ComponentDefinition) propertyValue {
	return propertyValue{
		valueType: valueTypeComponent,
		component: compDef,
	}
}

func newListValue(list []propertyValue) propertyValue {
	return propertyValue{
		valueType: valueTypeList,
		list:      list,
	}
}

func newExpressionValue() propertyValue {
	return propertyValue{
		valueType: valueTypeExpression,
	}
}

func (v propertyValue) String() string {
	switch v.valueType {
	case valueTypeComponent:
		return v.component.String()
	case valueTypeList:
		var parts []string
		for _, v := range v.list {
			parts = append(parts, v.String())
		}
		return "[ " + strings.Join(parts, ", ") + " ]"
	case valueTypeExpression:
		return "<code>"
	default:
		return ""
	}
}

// onlyContainsExpressions returns true if this propertyValue only contains expressions or lists of expressions.
// If it returns false it will also return the component definition that it does contain.
func (v propertyValue) onlyContainsExpressions() (bool, *vit.ComponentDefinition) {
	if v.valueType == valueTypeExpression {
		return true, nil
	} else if v.valueType == valueTypeComponent {
		return false, v.component
	} else {
		for _, v := range v.list {
			if yes, comp := v.onlyContainsExpressions(); !yes {
				return false, comp
			}
		}
		return true, nil
	}
}

// componentList returns component definitions that are inside this list.
// Only valid if this propertyValue is of type list and contains only component definitions.
func (v propertyValue) componentList() []*vit.ComponentDefinition {
	var components []*vit.ComponentDefinition
	for _, v := range v.list {
		components = append(components, v.component)
	}
	return components
}

// parsePropertyValueFromExpression analyses the given expression token to figure out if it might contain special values.
// Instead of actual JavaScript code an expression might contain a component definition or a list of those.
// If it does it will create a local lexer instance to further process the code as it is not intended for the JavaScript interpreter.
func parsePropertyValueFromExpression(expressionToken token) (value propertyValue, err error) {
	defer func() {
		if r := recover(); r != nil {
			value = newExpressionValue()
			err = nil
		}
	}()

	// TODO: potentially optimize this by performing a quick and simple check on the string wether or not this can be a component definition of a list

	lexer := NewLexerAtPosition(strings.NewReader(expressionToken.literal), expressionToken.position.Start())
	tokens := NewTokenBuffer(lexer.Lex)

	value, err = parsePropertyValueFromTokens(tokens)
	if err != nil {
		return propertyValue{}, err
	}

	// validation step to make sure the value is in a form that is currently supported
	isExpression, err := validatePropertyValue(value, expressionToken)
	if err != nil {
		return propertyValue{}, err
	}
	if isExpression {
		return newExpressionValue(), nil
	}

	return value, nil
}

// parsePropertyValueFromTokens parses the given tokens to figure out what type of propertyValue this is.
func parsePropertyValueFromTokens(tokens *tokenBuffer) (propertyValue, error) {
nextToken:
	t := tokens.next()
	switch t.tokenType {
	case tokenNewline:
		goto nextToken // skip newlines

	case tokenIdentifier: // start of a component
		ignoreTokens(tokens, tokenNewline)
		_, err := expectToken(tokens.next, tokenLeftBrace) // opening brace
		if err != nil {
			return newExpressionValue(), nil
		}
		compDef, err := parseComponent(t.literal, tokens) // actual component content
		if err != nil {
			return propertyValue{}, err // javascript
		}
		return newComponentValue(compDef), nil

	case tokenLeftBracket: // opening of a list
		var values []propertyValue
	parseNextItem:
		// parse list item
		content, err := parsePropertyValueFromTokens(tokens)
		if err != nil {
			return propertyValue{}, err // javascript
		}
		values = append(values, content)   // save the item
		ignoreTokens(tokens, tokenNewline) // skip newlines

		t = tokens.next()
		if t.tokenType == tokenComma {
			goto parseNextItem // another items follows
		} else if t.tokenType == tokenRightBracket {
			return newListValue(values), nil // the list is finished
		} else {
			return propertyValue{}, unexpectedToken(t, tokenRightBracket, tokenComma)
		}

	default:
		return newExpressionValue(), nil
	}
}

// validatePropertyValue checks if the given propertyValue is in a form that is currently supported.
// It will specifically check lists.
// The only allowed lists are:
// - a single dimension list that only contains component definitions
// - an arbitrary list that only contains other lists and expressions
// All other lists will be blocked.
func validatePropertyValue(value propertyValue, expressionToken token) (isExpression bool, err error) {
	switch value.valueType {
	case valueTypeComponent:
		return false, nil
	case valueTypeList:
		var consistentType propertyValueType = valueTypeList
		for _, v := range value.list {
			if v.valueType == valueTypeExpression {
				if consistentType == valueTypeComponent {
					return false, parseErrorf(expressionToken.position, "lists with a mix between JavaScript code and component definitions are currently not cupported")
				}
				consistentType = valueTypeExpression
			} else if v.valueType == valueTypeComponent {
				if consistentType == valueTypeExpression {
					return false, parseErrorf(expressionToken.position, "lists with a mix between JavaScript code and component definitions are currently not cupported")
				}
				consistentType = valueTypeComponent
			} else if v.valueType == valueTypeList {
				if yes, comp := v.onlyContainsExpressions(); yes {
					return true, nil
				} else {
					return true, parseErrorf(comp.Pos, "component definitions in nested lists are currently not supported")
				}
			}
		}
		return false, nil
	case valueTypeExpression:
		return true, nil
	}
	return true, nil
}
