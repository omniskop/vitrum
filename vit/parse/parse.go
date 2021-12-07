package parse

import (
	"errors"
	"fmt"

	"github.com/omniskop/vitrum/vit"
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

var keywords = map[string]bool{
	"property": true,
	"default":  true,
	"required": true,
	"readonly": true,
	"enum":     true,
	"embedded": true,
}

// A tokenSource is a function that returns tokens
type tokenSource func() token

// unitType describes the type of a code unit that has been parsed
type unitType int

const (
	unitTypeNil          unitType = iota // no valid unit was parsed
	unitTypeEOF                          // the file has ended
	unitTypeComponent                    // a component definition has been parsed
	unitTypeComponentEnd                 // a component has ended
	unitTypeProperty                     // a property has been parsed
	unitTypeEnum                         // an enum has been parsed
)

// String returns the name of a unitType as a string
func (uType unitType) String() string {
	switch uType {
	case unitTypeEOF:
		return "end of file"
	case unitTypeComponent:
		return "component"
	case unitTypeComponentEnd:
		return "end of component"
	case unitTypeProperty:
		return "property"
	case unitTypeEnum:
		return "enum"
	default:
		return "unknown unit"
	}
}

// Parse takes a tokenBuffer from a single vit file and returns a parsed document
func Parse(tokens *tokenBuffer) (file *VitDocument, err error) {
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

	file = new(VitDocument)

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

// parseImports parses all import statements at the beginning of a file
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

// parseSingleImport parses a single import statement
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
// It specifies the returned unit by it's type. (See 'unitType')
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
		if _, ok := keywords[lineIdentifier[0].literal]; ok {
			return parseAttributeDeclaration(lineIdentifier[0].literal, tokens)
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

// parseComponent parses the content of a component and returns a component definition.
// It takes the name of the component as parameter.
func parseComponent(identifier string, tokens *tokenBuffer) (*componentDefinition, error) {
	c := &componentDefinition{
		name: identifier,
	}

	// read units and store them correspondingly until we reach the end of the component
	for {
		unitIntf, uType, err := parseUnit(tokens)
		if err != nil {
			return c, err
		}
		switch uType {
		case unitTypeComponentEnd: // the component has ended
			return c, nil
		case unitTypeProperty: // property declaration/definition
			prop := unitIntf.(property)
			// check if this property has the identifier "id" and handle it specially
			if len(prop.identifier) == 1 && prop.identifier[0] == "id" {
				// TODO: validate that the expression is a valid id. Calculations are not allowed
				c.id = prop.expression
			} else {
				// check if the property has already been defined before
				if c.identifierIsKnown(prop.identifier) {
					return c, fmt.Errorf("identifier %v is already defined", prop.identifier)
				}
				// save it in the component
				c.properties = append(c.properties, prop)
			}
		case unitTypeEnum: // an enumeration
			enum := unitIntf.(vit.Enumeration)
			if c.identifierIsKnown([]string{enum.Name}) {
				return c, fmt.Errorf("identifier %q is already defined", enum.Name)
			}
			c.enumerations = append(c.enumerations, enum)
		case unitTypeComponent: // child component
			child := unitIntf.(*componentDefinition)
			c.children = append(c.children, child)
		default:
			return c, fmt.Errorf("unexpected %v while parsing unit", uType)
		}
	}
}

// parseAttributeDeclaration parses the declaration or a component attribute. That could for example be a property or an enumeration.
// The provided 'keyword' should be the first word that has already been read from the line. (Which has determined that this will be an attribute declaration)
func parseAttributeDeclaration(keyword string, tokens *tokenBuffer) (interface{}, unitType, error) {
	// We will collect modifiers that are listed before the actual attribute type is specified
	var modifiers []string

	// The switch comes before be read an actual token to handle the provided keyword that has been read before
	for {
		switch keyword {
		case "property":
			// this attribute is a property
			prop, err := parseProperty(tokens, modifiers)
			if err != nil {
				return nil, unitTypeNil, err
			}
			return prop, unitTypeProperty, nil
		case "enum":
			// this attribute is an enumeration
			en, err := parseEnum(tokens, modifiers)
			if err != nil {
				return nil, unitTypeNil, err
			}
			return en, unitTypeEnum, nil
		default:
			// a modifier
			for _, m := range modifiers {
				if m == keyword {
					return nil, unitTypeNil, fmt.Errorf("duplicate modifier %q", keyword)
				}
			}
			modifiers = []string{keyword}
		}

		// read the next word
		t, err := expectToken(tokens.next, tokenIdentifier)
		if err != nil {
			return nil, unitTypeNil, err
		}
		keyword = t.literal
	}
}

// parseProperty parses a property declaration/definition with the given modifiers
func parseProperty(tokens *tokenBuffer, modifiers []string) (property, error) {
	// read the type of the property
	typeToken, err := expectToken(tokens.next, tokenIdentifier)
	if err != nil {
		return property{}, err
	}

	// property name
	identifier, err := expectToken(tokens.next, tokenIdentifier)
	if err != nil {
		return property{}, err
	}

	prop := property{
		identifier: []string{identifier.literal},
		vitType:    typeToken.literal,
		readOnly:   stringSliceContains(modifiers, "readonly"),
	}

	// read the next token and determine if a value will follow
	switch tokens.next().tokenType {
	case tokenColon: // continue to read the value
	case tokenNewline:
		return prop, nil // property is finished with no value
	default:
		return property{}, unexpectedToken(tokens.next(), tokenColon, tokenNewline)
	}

	// read the value (expression that will determine the value)
	expression, err := expectToken(tokens.next, tokenExpression)
	if err != nil {
		return property{}, err
	}
	prop.expression = expression.literal

	// remove any newlines or semicolons
	_, err = expectToken(tokens.next, tokenNewline, tokenSemicolon)
	if err != nil {
		return property{}, err
	}

	return prop, nil
}

// parseEnum parses an enumeration declaration with the given modifiers
func parseEnum(tokens *tokenBuffer, modifiers []string) (vit.Enumeration, error) {
	enum := vit.Enumeration{
		Values:   make(map[string]int),
		Embedded: stringSliceContains(modifiers, "embedded"),
	}

	// name
	t, err := expectToken(tokens.next, tokenIdentifier)
	if err != nil {
		return enum, err
	}
	enum.Name = t.literal

	// expect a '{'
	_, err = expectToken(tokens.next, tokenLeftBrace)
	if err != nil {
		return enum, err
	}

	// throw away new lines
	ignoreTokens(tokens, tokenNewline)

	var nextValue = 0 // Contains the value that the next enum key will have. We start at 0
lineLoop:
	for {
		// read enum key or closing brace
		keyToken := tokens.next()
		if keyToken.tokenType == tokenRightBrace {
			break lineLoop // enum ended
		} else if keyToken.tokenType != tokenIdentifier {
			return enum, unexpectedToken(keyToken, tokenIdentifier, tokenRightBrace)
		}
		// check if the key already exists
		if _, ok := enum.Values[keyToken.literal]; ok {
			return enum, parseError{keyToken.start, fmt.Errorf("duplicate enum key %q", keyToken.literal)}
		}

		// check if there is a manual assignment
		t := tokens.next()
		if t.tokenType == tokenAssignment {
			valueToken := tokens.next()
			if valueToken.tokenType != tokenInteger {
				return enum, parseError{valueToken.start, fmt.Errorf("only integer literals can be assigned to enum keys, but found %q", keyToken.literal)}
			}
			nextValue = valueToken.IntValue()
			if nextValue < 0 {
				return enum, parseError{valueToken.start, fmt.Errorf("enum value can't be negative")}
			}
			t = tokens.next() // read the next token after the assignment
		}

		// store the value
		enum.Values[keyToken.literal] = nextValue
		nextValue++

		// check how the line ends
		switch t.tokenType {
		case tokenComma:
			ignoreTokens(tokens, tokenNewline)
		case tokenNewline:
			// a newline without a semicolon is only allowed at the end of the block, thus we expect a closing brace here
			t = tokens.next()
			if t.tokenType != tokenRightBrace {
				return enum, unexpectedToken(t, tokenRightBrace)
			}
			enum.Values[t.literal] = nextValue
			break lineLoop
		default:
			return enum, unexpectedToken(t, tokenAssignment, tokenComma, tokenNewline)
		}

	}

	return enum, nil
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

// expectToken ready the next token and checks if it is of one of the given types.
// If not it returns a unexpectedToken error.
func expectToken(nextToken tokenSource, tTypes ...tokenType) (token, error) {
	t := nextToken()

	for _, tType := range tTypes {
		if t.tokenType == tType {
			return t, nil
		}
	}

	return t, unexpectedToken(t, tTypes...)
}

// expectKeyword reads the next token and checks if it is an identifier with the given value.
// If it is not an identifier it returns an unexpectedToken error. If it is but the values don't match it returns a descriptive parseError.
func expectKeyword(nextToken tokenSource, value string) (token, error) {
	t := nextToken()

	if t.tokenType != tokenIdentifier {
		return t, unexpectedToken(t, tokenIdentifier)
	} else if t.literal != value {
		return t, parseError{pos: t.start, err: fmt.Errorf("unexpected token %v, expected keyword %q", t.literal, value)}
	}

	return t, nil
}

// literalsToStrings converts the literals of a tokens list into a string slice
func literalsToStrings(tokens []token) []string {
	strs := make([]string, len(tokens))
	for i, ident := range tokens {
		strs[i] = ident.literal
	}
	return strs
}

// stringSliceContains checks if a string slice contains a given string
func stringSliceContains(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

// unexpectedToken returns true if two string slices equal and false otherwise
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, e := range a {
		if e != b[i] {
			return false
		}
	}
	return true
}
