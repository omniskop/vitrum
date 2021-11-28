package parse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// LexError contains additional information about the error that occurred
type LexError struct {
	pos position
	msg string
}

func (e LexError) Error() string {
	return fmt.Sprintf("%v: %s", e.pos, e.msg)
}

func (e LexError) Is(subject error) bool {
	_, ok := subject.(LexError)
	return ok
}

type ReadError struct {
	err error
}

func (e ReadError) Error() string {
	return fmt.Sprintf("read error: %v", e.err)
}

func (e ReadError) Is(subject error) bool {
	_, ok := subject.(ReadError)
	return ok
}

func (e ReadError) Unwrap() error {
	return e.err
}

type lexer struct {
	source              *bufio.Reader
	filePath            string
	pos                 position
	previousPosition    position
	expressionFollowing bool // weather the next scanned part should be an expression
}

func NewLexer(input io.Reader, filePath string) *lexer {
	return &lexer{
		source:   bufio.NewReader(input),
		filePath: filePath,
		pos: position{
			filePath: filePath,
			line:     1,
			column:   1,
		},
	}
}

func LexAll(input io.Reader, filePath string) ([]token, error) {
	l := NewLexer(input, filePath)
	allTokens := make([]token, 0, 100) // we will just start with some capacity
	for {
		t, err := l.Lex()
		if err != nil {
			return allTokens, err
		}
		allTokens = append(allTokens, t)
		if t.tokenType == tokenEOF {
			break
		}
	}
	return allTokens, nil
}

// nextRune returns the next rune from the source as well as the position of that rune in the file.
// The error returned will always be of type ReadError.
func (l *lexer) nextRune() (rune, position, error) {
	r, _, err := l.source.ReadRune()
	if err != nil {
		return r, l.pos, ReadError{err}
	}
	p := l.pos
	l.previousPosition = l.pos
	if r == '\n' {
		l.pos.line++
		l.pos.column = 1
	} else {
		l.pos.column++
	}

	return r, p, nil
}

// unreadRune adds the last read rune back to the buffer.
// Calling this more than once might result in incorrect positions.
func (l *lexer) unreadRune() {
	l.source.UnreadRune()
	l.pos = l.previousPosition
	l.previousPosition.column-- // just an estimate, we might be missing a line break
}

// Lex returns the next token from the source.
// The error returned will either be of type LexError if the source file contained an error or of type ReadError if the source broke.
func (l *lexer) Lex() (token, error) {
	if l.expressionFollowing {
		t, err := l.scanExpression()
		if err != nil {
			return t, err
		}
		l.expressionFollowing = false
		if t.tokenType == tokenExpression && len(t.literal) == 0 {
			return t, LexError{pos: t.start, msg: "empty expression"}
		}
		return t, nil
	}

	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return token{
					tokenType: tokenEOF,
					start:     pos,
					end:       pos,
				}, nil
			}

			return token{}, err
		}

		switch {
		case unicode.IsSpace(r) && r != '\n':
			continue
		case r == '"' || r == '\'' || r == '`':
			return l.scanString(r)
		case r == '{':
			return token{
				tokenType: tokenLeftBrace,
				start:     pos, end: pos,
			}, nil
		case r == '}':
			return token{
				tokenType: tokenRightBrace,
				start:     pos, end: pos,
			}, nil
		case r == ':':
			l.expressionFollowing = true
			return token{
				tokenType: tokenColon,
				start:     pos, end: pos,
			}, nil
		case r == '.':
			return token{
				tokenType: tokenPeriod,
				start:     pos, end: pos,
			}, nil
		case r == '[':
			return token{
				tokenType: tokenLeftBracket,
				start:     pos, end: pos,
			}, nil
		case r == ']':
			return token{
				tokenType: tokenRightBracket,
				start:     pos, end: pos,
			}, nil
		case r == ',':
			return token{
				tokenType: tokenComma,
				start:     pos, end: pos,
			}, nil
		case r == ';':
			return token{
				tokenType: tokenSemicolon,
				start:     pos, end: pos,
			}, nil
		case r == '\n':
			return token{
				tokenType: tokenNewline,
				start:     pos, end: pos,
			}, nil
		case r == '/':
			r, _, err := l.nextRune()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return token{}, LexError{l.pos, fmt.Sprintf("unexpected symbol %q", string(r))}
				}
				return token{}, err
			}
			if r == '/' || r == '*' {
				err = l.skipComment(r == '*')
				if err != nil {
					return token{}, err
				}
			} else {
				l.unreadRune()
				return token{}, LexError{l.pos, fmt.Sprintf("unexpected symbol %q", string(r))}
			}
			continue
		case unicode.IsLetter(r):
			return l.scanIdentifier(r)
		case unicode.IsNumber(r):
			l.unreadRune()
			return l.scanNumber()
		default:
			return token{}, LexError{l.pos, fmt.Sprintf("unexpected symbol %q", string(r))}
		}
	}
}

// scanExpression will scan the input for an expression. The returned token can either be an expression or a literal.
// The returned error is either of type LexError or ReadError.
//
// Expressions will not be dissected in detail because it is potentially JavaScript code that will be handled separately.
// We will still need to keep track of braces and strings though so we know where the expression ends.
// An expression can also be just a single line or a block expression that is surrounded by curly braces, in which case it would be a function.
//
// But it is also possible to have a component as a value and we will need to detect that and handle it differently.
// We will detect that by checking that the expression only contains spaces and/or valid runes for a literal
// followed by a left brace. This *could* also match valid JavaScript code but I don't see a reason to assume that that would ever be usefull code in a vit file.
//
// It also removes line and block comments. This is required in the case of a vit component.
// This will also be applied to JavaScript code, which should be fine I think, but it could potentially be changed if it should cause any issues.
func (l *lexer) scanExpression() (token, error) {
	t := token{
		tokenType: tokenExpression,
		start:     l.pos,
	}

startExpressionScan:

	var blockExpression bool // is this expression encapsulated in braces
	var str strings.Builder

	// read in all spaces until the expression starts and check if it is encapsulated in braces
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return t, nil // just return the empty expression
			}

			return token{}, err
		}
		if unicode.IsSpace(r) { // also matches newlines
			continue
		} else if r == '{' {
			blockExpression = true
		}
		t.start = pos
		l.unreadRune() // put whatever we just read back
		break
	}

	// Now read the actual expression.
	// A simple expression ends with a semicolon, a new line or a closing brace.
	// A block expression is only terminated by a closing brace.
	// We need to keep track of opened strings to make sure that we ignore terminating runes inside them.
	// For example it is valid for a newline character to be inside a string without terminating the expression.
	var openBraces int           // number of open braces ('{') that have yet to be closed
	var insideString bool        // weather or not we are inside a string
	var stringDelimiter rune     // which rune started the string (either ', ` or ")
	var escaped bool             // weather the next rune is escaped
	var previousPosition = l.pos // the position of the previous rune
	var potentialComment bool    // set to true if we encountered a / in the last rune
	// couldBeLiteral is true if all read runes until this point would be a valid literal.
	// We need to check this because properties can also hold components which would start with a literal.
	var couldBeLiteral = !blockExpression // only if it's not a block expression
	// If we thought that the read runes couldn't be a literal because the last read character was a / we will store the position of that slash here.
	var literalFailedBecauseOfSlashAt int = -1

	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if blockExpression {
					return t, LexError{pos, "unexpected end of file, expecting closing brace"}
				} else {
					t.end = previousPosition
					t.literal = str.String()
					return t, nil
				}
			}

			return token{}, err
		}

		if couldBeLiteral {
			if r == '{' {
				// this is not an expression but a vit component
				// we will report a literal token and leave

				l.unreadRune() // put the brace back so it can be read on the next call to Lex
				// remove potential spaces or newlines around the literal
				value := strings.Trim(str.String(), " \n")
				// value can now only contain characters appearing on the same line
				// thus we can calculate the end like this:
				end := t.start
				end.column += len(value) - 1 // -1 because we wan't the location of the last rune not the following one
				return token{
					tokenIdentifier,
					value,
					t.start,
					end,
				}, nil
			}
			// if this is not a valid literal rune, space or newline ...
			if !validLiteralRune(r, str.Len() == 0) && r != ' ' && r != '\n' {
				// ... this has no longer the potential to be a vit component
				if r == '/' {
					// if it is a / we will note the position in case that this is just the start of a comment
					literalFailedBecauseOfSlashAt = pos.column
				}
				couldBeLiteral = false
			}
		}

		// if we are inside a string ...
		if insideString {
			// ... and this rune should not be escaped and is the same rune that started the string ...
			if !escaped && r == stringDelimiter {
				// ... then we are out of the string
				insideString = false
			}
			// go to the next rune as we will not parse any special characters inside a string
			goto continueLoop
		}

		// we are not inside a string here

		// now figure out if this is a special character
		switch r {
		case '"', '\'', '`': // opening a string
			insideString = true
			stringDelimiter = r
		case '{':
			openBraces++
		case '}':
			openBraces--
			// if we closed the last brace ...
			if openBraces <= 0 {
				// ... the expression has ended
				if blockExpression {
					// if it is a block expression the closing brace is part of it
					str.WriteRune(r)
					fmt.Println(pos)
					t.end = pos
				} else {
					// As this is not a block expression we assume that this closing brace ends the component whose property we are reading.
					// Thus we will put the brace back and return here.
					t.end = previousPosition
					l.unreadRune()
				}
				t.literal = str.String()
				return t, nil
			}
		case '\n', ';':
			// inline expressions end here
			if !blockExpression {
				t.end = previousPosition
				l.unreadRune()
				t.literal = str.String()
				return t, nil
			}
		case '/':
			if potentialComment {
				// a line comment started
				t.end = previousPosition
				t.end.column-- // remove first '/'; we can just go back one rune because we now that we can't be at the start of a line
				err := l.skipComment(false)
				t.literal = str.String()
				if len(t.literal) == 0 { // should never happen
					return t, err
				}
				// cut away the first '/' at the end of the literal
				t.literal = t.literal[:len(t.literal)-1]
				return t, err // expression ended here
			}
			potentialComment = true
			goto continueLoop
		case '*':
			if potentialComment {
				// a block comment started
				t.end = previousPosition
				l.unreadRune()
				err := l.skipComment(true)
				if err != nil {
					return t, err
				}
				s := str.String()
				if len(s) > 0 { // should always be the case
					str.Reset()
					str.WriteString(s[:len(s)-1]) // cut away the '/' before the * at the end of the literal
				}
				if len(s) == 1 && !blockExpression {
					// if the only thing scanned so far was the first '/' of the comment we restart scanning of the expression like it had just started
					// this gives more accurate position of the actual expression without the comment
					goto startExpressionScan
				}
				// if we decided that the scanned text could not be a literal just because we found the preceding '/', reset couldBeLiteral here.
				if literalFailedBecauseOfSlashAt == pos.column-1 {
					couldBeLiteral = true
					literalFailedBecauseOfSlashAt = -1
				}

				// continue with the expression
				goto continueLoopWithoutRune
			}
		}

		potentialComment = false // reset; if we just set this to true we jumped over this

	continueLoop:
		str.WriteRune(r)
	continueLoopWithoutRune:
		previousPosition = pos
	}
}

// skipSpaces skips all spaces and newlines. It returns the position of the first non-space character.
// That rune will still be in the source as unreadRune will be called at the end.
func (l *lexer) skipSpaces() (position, error) {
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return pos, fmt.Errorf("unexpected end of file")
			}

			return pos, fmt.Errorf("unexpected error: %w", err)
		}
		if !unicode.IsSpace(r) { // also matches newlines
			l.unreadRune()
			return pos, nil
		}
	}
}

// scanString scans a string token that has been started with the given quotation mark. The literal will not contain the quotation marks.
// The returned error will either be of type LexError or ReadError.
func (l *lexer) scanString(quotationMark rune) (token, error) {
	t := token{
		tokenType: tokenString,
		start:     l.pos,
	}

	// NOTE: new lines are allowed inside strings

	var escaped bool
	var str strings.Builder
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return t, LexError{pos, fmt.Sprintf("unexpected end of file")}
			}
			return t, err
		}

		if escaped {
			escaped = false
			str.WriteRune(r)
			continue
		}

		switch r {
		case quotationMark:
			t.end = pos
			t.literal = str.String()
			return t, nil
		case '\\':
			escaped = true
		}

		str.WriteRune(r)
	}
}

// scanIdentifier scans an identifier token that starts with the given first rune.
// The returned error will be of type ReadError.
func (l *lexer) scanIdentifier(first rune) (token, error) {
	t := token{
		tokenType: tokenIdentifier,
		start:     l.pos,
	}
	t.start.column--

	previousPosition := t.start
	var str strings.Builder
	str.WriteRune(first)
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				t.end = previousPosition
				t.literal = str.String()
				return t, nil
			}
			return t, err
		}

		if validLiteralRune(r, str.Len() == 0) {
			str.WriteRune(r)
		} else {
			t.end = previousPosition
			t.literal = str.String()
			l.unreadRune()
			return t, nil
		}

		previousPosition = pos
	}
}

// skipComment reads all runes until the end of the comment.
func (l *lexer) skipComment(multiLineComment bool) error {
	var starFound bool // only required for multiline comments
	for {
		r, _, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if multiLineComment {
					return fmt.Errorf("unexpected end of file")
				}
				// a singe line comment can be the last thing in a file
				return nil
			}
			return fmt.Errorf("unexpected error: %w", err)
		}

		if multiLineComment {
			if r == '*' {
				starFound = true
			} else if starFound && r == '/' {
				return nil // reached end of multiline comment
			} else {
				starFound = false
			}
		} else {
			if r == '\n' {
				// we have reached the end of the line
				l.unreadRune()
				return nil
			}
		}
	}
}

// scanNumber scans a number token.
// The returned error will either be of type LexError or ReadError.
func (l *lexer) scanNumber() (token, error) {
	t := token{
		tokenType: tokenFloat,
		start:     l.pos,
	}

	// read in all runes

	var previousPosition = l.pos
	var str strings.Builder
	var isFloatingPoint bool
	var invalid bool
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				goto parseAndReturn
			}
			return token{}, err
		}

		// if r is not a number or a period...
		if !unicode.IsNumber(r) && r != '.' {
			// ... the number has ended
			l.unreadRune()
			goto parseAndReturn
		}

		if r == '.' {
			if isFloatingPoint {
				invalid = true // more than one floating point but we will continue to scan
			} else {
				isFloatingPoint = true
			}
		}

		str.WriteRune(r)
		previousPosition = pos
	}

	// now parse the read number

parseAndReturn:
	t.end = previousPosition

	stringN := str.String()
	// catch some exceptions
	if len(stringN) == 0 || stringN == "." || stringN == "+" || stringN == "-" {
		return t, LexError{t.start, fmt.Sprintf("number %q is incomplete", stringN)}
	}
	t.literal = stringN

	if invalid {
		t.tokenType = tokenIdentifier
		return t, nil
	}

	if isFloatingPoint {
		_, err := strconv.ParseFloat(stringN, 64)
		if err != nil {
			return t, LexError{t.start, fmt.Sprintf("%q is not a valid number: %v", stringN, err)}
		}
		t.tokenType = tokenFloat
	} else {
		_, err := strconv.ParseInt(stringN, 10, 64)
		if err != nil {
			return t, LexError{t.start, fmt.Sprintf("%q is not a valid number: %v", stringN, err)}
		}
		t.tokenType = tokenInteger
	}

	return t, nil
}

// validLiteralRune checks if the rune can be used in a literal
// 'first' signals that this rune would be the first of the literal which has stricter requirements
func validLiteralRune(r rune, first bool) bool {
	// a literal cannot start with a '-' or a number but they can contain it
	return unicode.IsLetter(r) || r == '_' || (!first && (r == '-' || unicode.IsNumber(r)))
}
