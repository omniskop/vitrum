package parse

import (
	"errors"
	"fmt"

	"github.com/omniskop/vitrum/vit"
)

// ParseError describes an error that occurred during parsing. Tt contains the position in the file where the error occurred
type ParseError struct {
	pos vit.PositionRange
	err error
}

func parseErrorf(p vit.PositionRange, format string, args ...interface{}) ParseError {
	return ParseError{
		pos: p,
		err: fmt.Errorf(format, args...),
	}
}

func (e ParseError) Error() string {
	return fmt.Sprintf("%v: %v", e.pos, e.err)
}

func (e ParseError) Report() string {
	return e.pos.Report()
}

func (e ParseError) Is(subject error) bool {
	_, ok := subject.(ParseError)
	return ok
}

func (e ParseError) Unwrap() error {
	return e.err
}

// an unexpected token error describes that a token of a specific type (or multiple) was expected but a different one was found.
type unexpectedTokenError struct {
	got      token       // the found token
	expected []tokenType // the expected types
}

// unexpectedToken creates an unexpectedTokenError
func unexpectedToken(got token, expected ...tokenType) ParseError {
	return ParseError{
		pos: got.position,
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

// a list of some keywords that are used to detect component attributes
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

// Parse takes a tokenBuffer from a single vit file and returns a parsed document.
// The returned error will always be of type LexError, ReadError or ParseError.
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
		parsedUnit, err := parseUnit(tokens)
		if err != nil {
			return nil, err
		}

		switch parsedUnit.kind {
		case unitTypeEOF:
			break scanComponents
		case unitTypeComponent:
			component := parsedUnit.value.(*componentDefinition)
			file.components = append(file.components, component)
		default:
			return nil, parseErrorf(parsedUnit.position, "unexpected %v in global scope", parsedUnit.kind)
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
		t := tokens.next()

		imp, err := parseSingleImport(tokens)
		imp.position = vit.CombineRanges(imp.position, t.position)
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
	imp.position = t.position
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
	imp.position.SetEnd(t.position.End())

	_, err = expectToken(tokens.next, tokenNewline, tokenSemicolon)
	if err != nil {
		return imp, err
	}

	return imp, nil
}

// parseUnit parses a semantic unit of the file.
// It specifies the returned unit by it's type. (See 'unitType')
func parseUnit(tokens *tokenBuffer) (unit, error) {
	ignoreTokens(tokens, tokenNewline, tokenSemicolon)

	var lineIdentifier []token
	var startingPosition vit.Position

	// scan identifier
scanLineIdentifier:
	switch t := tokens.next(); t.tokenType {
	case tokenRightBrace:
		return componentEndUnit(t.position.Start()), nil // end of component
	case tokenIdentifier:
		// part of the line identifier
		lineIdentifier = append(lineIdentifier, t)
	case tokenEOF:
		return eofUnit(t.position.Start()), nil
	default:
		return nilUnit(), unexpectedToken(t, tokenIdentifier)
	}

	// check if the scanned identifier is a keyword
	if len(lineIdentifier) == 1 {
		startingPosition = lineIdentifier[0].position.Start()
		if _, ok := keywords[lineIdentifier[0].literal]; ok {
			return parseAttributeDeclaration(lineIdentifier[0], tokens)
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
			return componentUnit(component.position, component), err
		} else {
			return nilUnit(), ParseError{lineIdentifier[1].position, fmt.Errorf("qualified identifier is not allowed for components")}
		}
	case tokenColon:
		// property
		t := tokens.next()
		switch t.tokenType {
		case tokenIdentifier:
			// value of the property is a component
			_, err := expectToken(tokens.next, tokenLeftBrace)
			if err != nil {
				return nilUnit(), err
			}
			component, err := parseComponent(t.literal, tokens)
			return propertyUnit(component.position, property{
				position:   vit.NewRangeFromStartToEnd(startingPosition, t.position.End()),
				identifier: literalsToStrings(lineIdentifier),
				component:  component,
			}), err
		case tokenExpression:
			// value of the property is set by an expression
			return propertyUnit(t.position, property{
				position:   vit.NewRangeFromStartToEnd(startingPosition, t.position.End()),
				identifier: literalsToStrings(lineIdentifier),
				expression: t.literal,
			}), nil
		default:
			return nilUnit(), unexpectedToken(t, tokenIdentifier, tokenExpression)
		}
	default:
		return nilUnit(), unexpectedToken(t, tokenPeriod, tokenLeftBrace, tokenColon, tokenIdentifier)
	}

	return nilUnit(), nil
}

// parseComponent parses the content of a component and returns a component definition.
// It takes the name of the component as parameter.
func parseComponent(identifier string, tokens *tokenBuffer) (*componentDefinition, error) {
	c := &componentDefinition{
		name: identifier,
	}

	// read units and store them correspondingly until we reach the end of the component
	for {
		parsedUnit, err := parseUnit(tokens)
		if err != nil {
			return c, err
		}
		switch parsedUnit.kind {
		case unitTypeComponentEnd: // the component has ended
			return c, nil
		case unitTypeProperty: // property declaration/definition
			prop := parsedUnit.value.(property)
			// check if this property has the identifier "id" and handle it specially
			if len(prop.identifier) == 1 && prop.identifier[0] == "id" {
				// TODO: validate that the expression is a valid id. Calculations are not allowed
				c.id = prop.expression
			} else {
				// check if the property has already been defined before
				if c.identifierIsKnown(prop.identifier) {
					// TODO: theoretically we could get and show the position of the previous declaration here
					return c, parseErrorf(prop.position, "identifier %v is already defined", prop.identifier)
				}
				// save it in the component
				c.properties = append(c.properties, prop)
			}
		case unitTypeEnum: // an enumeration
			enum := parsedUnit.value.(vit.Enumeration)
			if c.identifierIsKnown([]string{enum.Name}) {
				return c, parseErrorf(parsedUnit.position, "identifier %q is already defined", enum.Name)
			}
			c.enumerations = append(c.enumerations, enum)
		case unitTypeComponent: // child component
			child := parsedUnit.value.(*componentDefinition)
			c.children = append(c.children, child)
		default:
			return c, parseErrorf(parsedUnit.position, "unexpected %v while parsing unit", parsedUnit.kind)
		}
	}
}

// parseAttributeDeclaration parses the declaration or a component attribute. That could for example be a property or an enumeration.
// The provided token should be the first word that has already been read from the line. (Which has determined that this will be an attribute declaration)
func parseAttributeDeclaration(t token, tokens *tokenBuffer) (unit, error) {
	// We will collect modifiers that are listed before the actual attribute type is specified
	var modifiers []string
	var err error
	var startingPosition = t.position.Start()

	// The switch comes before be read an actual token to handle the provided keyword that has been read before
	for {
		switch t.literal {
		case "property":
			// this attribute is a property
			prop, err := parseProperty(tokens, modifiers, startingPosition)
			if err != nil {
				return nilUnit(), err
			}
			return propertyUnit(prop.position, prop), nil
		case "enum":
			// this attribute is an enumeration
			en, err := parseEnum(tokens, modifiers, startingPosition)
			if err != nil {
				return nilUnit(), err
			}
			return enumUnit(en.Position, en), nil
		default:
			// a modifier
			for _, m := range modifiers {
				if m == t.literal {
					return nilUnit(), parseErrorf(t.position, "duplicate modifier %q", t.literal)
				}
			}
			modifiers = []string{t.literal}
		}

		// read the next word
		t, err = expectToken(tokens.next, tokenIdentifier)
		if err != nil {
			return nilUnit(), err
		}
	}
}

// parseProperty parses a property declaration/definition with the given modifiers
func parseProperty(tokens *tokenBuffer, modifiers []string, startingPosition vit.Position) (property, error) {
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
		position:   vit.NewRangeFromStartToEnd(startingPosition, identifier.position.End()),
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
	prop.position.SetEnd(expression.position.End())

	// remove any newlines or semicolons
	_, err = expectToken(tokens.next, tokenNewline, tokenSemicolon)
	if err != nil {
		return property{}, err
	}

	return prop, nil
}

// parseEnum parses an enumeration declaration with the given modifiers
func parseEnum(tokens *tokenBuffer, modifiers []string, startingPosition vit.Position) (vit.Enumeration, error) {
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
			enum.Position = vit.NewRangeFromStartToEnd(startingPosition, t.position.End())
			break lineLoop // enum ended
		} else if keyToken.tokenType != tokenIdentifier {
			return enum, unexpectedToken(keyToken, tokenIdentifier, tokenRightBrace)
		}
		// check if the key already exists
		if _, ok := enum.Values[keyToken.literal]; ok {
			return enum, parseErrorf(keyToken.position, "duplicate enum key %q", keyToken.literal)
		}

		// check if there is a manual assignment
		t := tokens.next()
		if t.tokenType == tokenAssignment {
			valueToken := tokens.next()
			if valueToken.tokenType != tokenInteger {
				return enum, parseErrorf(valueToken.position, "only integer literals can be assigned to enum keys, but found %q", keyToken.literal)
			}
			nextValue = valueToken.IntValue()
			if nextValue < 0 {
				return enum, parseErrorf(valueToken.position, "enum value can't be negative")
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
			enum.Position = vit.NewRangeFromStartToEnd(startingPosition, t.position.End())
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
// If not it returns a unexpectedToken error wrapped in a parseError.
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
		return t, parseErrorf(t.position, "unexpected token %v, expected keyword %q", t.literal, value)
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
