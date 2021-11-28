package parse

import (
	"errors"
	"fmt"
)

type parseError struct {
	pos position
	err error
}

func (e parseError) Error() string {
	return fmt.Sprintf("%v: %v", e.pos, e.err)
}

func (e parseError) Is(subject error) bool {
	_, ok := subject.(parseError)
	return ok
}

func (e parseError) Unwrap() error {
	return e.err
}

type unexpectedTokenError struct {
	got      token
	expected []tokenType
}

func unexpectedToken(got token, expected ...tokenType) parseError {
	return parseError{
		pos: got.start,
		err: unexpectedTokenError{
			got:      got,
			expected: expected,
		},
	}
}

func (e unexpectedTokenError) Error() string {
	if len(e.expected) == 0 {
		return fmt.Sprintf("unexpected %v", e.got)
	} else {
		return fmt.Sprintf("unexpected %v, expected %v", e.got, joinTokenTypes(e.expected))
	}
}

func (e unexpectedTokenError) Is(subject error) bool {
	_, ok := subject.(unexpectedTokenError)
	return ok
}

func Parse(tokens *tokenBuffer) (file *vitDocument, err error) {
	defer func() {
		// To simplify reading tokens from the lexer any errors it encounters will be thrown as a panic.
		// We will catch them here and return them properly.
		if r := recover(); r != nil {
			panicErr, ok := r.(error)
			if ok {
				var re ReadError
				var le LexError
				if errors.As(panicErr, &re) {
					// LexerReadError, the source couldn't be read
					err = re
				} else if errors.As(panicErr, &le) {
					// LexError, the source contained a lexical error
					err = le
				} else {
					// if you end up here due to a panic check the next but one entry in the stack trace for the actual error location
					panic(r)
				}
			} else {
				panic(r)
			}
		}
	}()

	file = new(vitDocument)

	file.imports, err = parseImports(tokens)
	if err != nil {
		return nil, err
	}

	// scan all components
scanComponents:
	for {
		unit, uType, err := parseUnit(tokens)
		if err != nil {
			return nil, err
		}

		switch uType {
		case unitTypeEOF:
			break scanComponents
		case unitTypeComponent:
			component := unit.(*componentDefinition)
			file.components = append(file.components, component)
		default:
			return nil, fmt.Errorf("unexpected %v in global scope", uType)
		}
	}

	return file, nil
}

func parseImports(tokens *tokenBuffer) ([]importStatement, error) {
	statements := make([]importStatement, 0)
	for {
		ignoreTokens(tokens, tokenNewline)

		// expect "import" literal
		nextToken := tokens.peek()
		if nextToken.tokenType != tokenIdentifier || nextToken.literal != "import" {
			// we reached the end of import statements
			return statements, nil
		}
		tokens.next()

		imp, err := parseSingleImport(tokens)
		if err != nil {
			return statements, err
		}
		statements = append(statements, imp)
	}
}

func parseSingleImport(tokens *tokenBuffer) (importStatement, error) {
	var imp importStatement

	// parse imported namespace or filepath
	var namespaceImport bool
scanAgain:
	t := tokens.next()
	if t.tokenType == tokenIdentifier {
		namespaceImport = true
		imp.namespace = append(imp.namespace, t.literal)
		nextToken := tokens.peek()
		if nextToken.tokenType == tokenPeriod {
			tokens.next()
			goto scanAgain
		} else if nextToken.tokenType != tokenFloat && nextToken.tokenType != tokenInteger && nextToken.tokenType != tokenIdentifier {
			return imp, unexpectedToken(nextToken, tokenFloat, tokenInteger, tokenIdentifier)
		}

	} else if t.tokenType == tokenString {
		if namespaceImport {
			return imp, unexpectedToken(t, tokenIdentifier)
		}
		imp.file = t.literal
	} else {
		return imp, unexpectedToken(t, tokenIdentifier, tokenString)
	}

	// parse version
	t, err := expectToken(tokens.next, tokenInteger, tokenFloat, tokenIdentifier)
	if err != nil {
		return imp, err
	}
	imp.version = t.literal

	_, err = expectToken(tokens.next, tokenNewline, tokenSemicolon)
	if err != nil {
		return imp, err
	}

	return imp, nil
}

// parseUnit parses a semantic unit of the file.
// That could be a single parameter, or a whole component.
func parseUnit(tokens *tokenBuffer) (interface{}, unitType, error) {
	ignoreTokens(tokens, tokenNewline, tokenSemicolon)

	var lineIdentifier []token

	// scan identifier
scanLineIdentifier:
	switch t := tokens.next(); t.tokenType {
	case tokenRightBrace:
		return nil, unitTypeComponentEnd, nil // end of component
	case tokenIdentifier:
		// part of the line identifier
		lineIdentifier = append(lineIdentifier, t)
	case tokenEOF:
		return nil, unitTypeEOF, nil
	default:
		return nil, unitTypeNil, unexpectedToken(t, tokenIdentifier)
	}

	// check if the scanned identifier is a keyword
	if len(lineIdentifier) == 1 {
		switch lineIdentifier[0].literal {
		case "property":
			property, err := parseProperty(tokens)
			if err != nil {
				return nil, unitTypeNil, err
			}
			return property, unitTypeProperty, nil
		}
	}

	// find out what this line is about
	switch t := tokens.next(); t.tokenType {
	case tokenPeriod:
		// we have a qualified line identifier like "one.two"
		goto scanLineIdentifier // scan the next identifier
		// TODO: think about removing this goto and just scanning the next identifier here and jumping to the start of this switch again
		//       because by jumping all the way up again we run a bunch of code that would only be applicable for the first identifier of the line
	case tokenLeftBrace:
		// new component definition
		if len(lineIdentifier) == 0 {
			// TODO: i think this can be valid; but it would not be possible to occur right now
		} else if len(lineIdentifier) == 1 {
			component, err := parseComponent(lineIdentifier[0].literal, tokens)
			return component, unitTypeComponent, err
		} else {
			return nil, unitTypeNil, parseError{lineIdentifier[1].start, fmt.Errorf("qualified identifier is not allowed for components")}
		}
	case tokenColon:
		// property
		t := tokens.next()
		switch t.tokenType {
		case tokenIdentifier:
			// value of the property is a component
			_, err := expectToken(tokens.next, tokenLeftBrace)
			if err != nil {
				return nil, unitTypeNil, err
			}
			component, err := parseComponent(t.literal, tokens)
			return property{
				identifier: literalsToStrings(lineIdentifier),
				component:  component,
			}, unitTypeProperty, err
		case tokenExpression:
			// value of the property is set by an expression
			return property{
				identifier: literalsToStrings(lineIdentifier),
				expression: t.literal,
			}, unitTypeProperty, nil
		default:
			return nil, unitTypeNil, unexpectedToken(t, tokenIdentifier, tokenExpression)
		}
	default:
		return nil, unitTypeNil, unexpectedToken(t, tokenPeriod, tokenLeftBrace, tokenColon, tokenIdentifier)
	}

	return nil, unitTypeNil, nil
}

func parseComponent(identifier string, tokens *tokenBuffer) (*componentDefinition, error) {
	c := &componentDefinition{
		name: identifier,
	}

	for {
		unitIntf, uType, err := parseUnit(tokens)
		if err != nil {
			return c, err
		}
		switch uType {
		case unitTypeComponentEnd:
			return c, nil
		case unitTypeProperty:
			prop := unitIntf.(property)
			c.properties = append(c.properties, prop)
		case unitTypeComponent:
			child := unitIntf.(*componentDefinition)
			c.children = append(c.children, child)
		default:
			return c, fmt.Errorf("unexpected %v while parsing unit", uType)
		}

	}
}

func parseProperty(tokens *tokenBuffer) (property, error) {
	typeToken, err := expectToken(tokens.next, tokenIdentifier)
	if err != nil {
		return property{}, err
	}
	if !dataTypes[typeToken.literal] {
		return property{}, parseError{typeToken.start, fmt.Errorf("unknown data type %q", typeToken.literal)}
	}

	identifier, err := expectToken(tokens.next, tokenIdentifier)
	if err != nil {
		return property{}, err
	}

	_, err = expectToken(tokens.next, tokenColon)
	if err != nil {
		return property{}, err
	}

	expression, err := expectToken(tokens.next, tokenExpression)
	if err != nil {
		return property{}, err
	}

	_, err = expectToken(tokens.next, tokenNewline, tokenSemicolon)
	if err != nil {
		return property{}, err
	}

	return property{
		identifier: []string{identifier.literal},
		vitType:    typeToken.literal,
		expression: expression.literal,
	}, nil
}

// ignoreTokens consumes all tokens of the given types.
// The first not matching token will be returned but is also still present in the token source.
func ignoreTokens(tokens *tokenBuffer, tTypes ...tokenType) token {
start:
	t := tokens.peek()

	for _, tType := range tTypes {
		if t.tokenType == tType {
			tokens.next()
			goto start
		}
	}
	return t
}

func expectToken(nextToken tokenSource, tTypes ...tokenType) (token, error) {
	t := nextToken()

	for _, tType := range tTypes {
		if t.tokenType == tType {
			return t, nil
		}
	}

	return t, unexpectedToken(t, tTypes...)
}

func expectKeyword(nextToken tokenSource, value string) (token, error) {
	t := nextToken()

	if t.tokenType != tokenIdentifier {
		return t, unexpectedToken(t, tokenIdentifier)
	} else if t.literal != value {
		return t, parseError{pos: t.start, err: fmt.Errorf("unexpected token %v, expected keyword %q", t.literal, value)}
	}

	return t, nil
}

func literalsToStrings(tokens []token) []string {
	strs := make([]string, len(tokens))
	for i, ident := range tokens {
		strs[i] = ident.literal
	}
	return strs
}
