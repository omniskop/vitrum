package parse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/omniskop/vitrum/vit"
)

// LexError contains additional information about the error that occurred
type LexError struct {
	pos vit.Position
	msg string
}

func (e LexError) Error() string {
	return fmt.Sprintf("%v: %s", e.pos, e.msg)
}

func (e LexError) Is(subject error) bool {
	_, ok := subject.(LexError)
	return ok
}

func unexpectedEOF(pos vit.Position) LexError {
	return LexError{
		pos: pos,
		msg: "unexpected end of file",
	}
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

type lexer struct {
	source              *bufio.Reader
	filePath            string
	pos                 vit.Position // position of the rune that will be read next
	previousPosition    vit.Position
	expressionFollowing bool // weather the next scanned part should be an expression
}

func NewLexer(input io.Reader, filePath string) *lexer {
	return &lexer{
		source:   bufio.NewReader(input),
		filePath: filePath,
		pos: vit.Position{
			FilePath: filePath,
			Line:     1,
			Column:   1,
		},
	}
}

// NewLexerAtPosition returns a new lexer that will already start at the given position.
// This can be used to improve error messages if a lexer only parses a small portion of a bigger file.
func NewLexerAtPosition(input io.Reader, position vit.Position) *lexer {
	return &lexer{
		source:   bufio.NewReader(input),
		filePath: position.FilePath,
		pos:      position,
	}
}

// nextRune returns the next rune from the source as well as the vit.Position of that rune in the file.
// The error returned will always be of type ReadError.
func (l *lexer) nextRune() (rune, vit.Position, error) {
	r, _, err := l.source.ReadRune()
	if err != nil {
		return r, l.pos, ReadError{err}
	}
	p := l.pos
	l.previousPosition = l.pos
	// TODO: check how this handles line breaks with a carriage return
	if r == '\n' {
		l.pos.Line++
		l.pos.Column = 1
	} else {
		l.pos.Column++
	}

	return r, p, nil
}

// unreadRune adds the last read rune back to the buffer.
// Calling this more than once might result in incorrect vit.Positions.
func (l *lexer) unreadRune() {
	l.source.UnreadRune()
	l.pos = l.previousPosition
	l.previousPosition.Column-- // just an estimate, we might be missing a line break
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
			return t, LexError{pos: t.position.Start(), msg: "empty expression"}
		}
		return t, nil
	}

	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return token{
					tokenType: tokenEOF,
					position:  vit.NewRangeFromPosition(pos),
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
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == '}':
			return token{
				tokenType: tokenRightBrace,
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == ':':
			l.expressionFollowing = true
			return token{
				tokenType: tokenColon,
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == '.':
			return token{
				tokenType: tokenPeriod,
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == '[':
			return token{
				tokenType: tokenLeftBracket,
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == ']':
			return token{
				tokenType: tokenRightBracket,
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == '<':
			return token{
				tokenType: tokenLess,
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == '>':
			return token{
				tokenType: tokenGreater,
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == '=':
			return token{
				tokenType: tokenAssignment,
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == ',':
			return token{
				tokenType: tokenComma,
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == ';':
			return token{
				tokenType: tokenSemicolon,
				position:  vit.NewRangeFromPosition(pos),
			}, nil
		case r == '\n':
			return token{
				tokenType: tokenNewline,
				position:  vit.NewRangeFromPosition(pos),
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
				_, err = l.readComment(r == '*')
				if err != nil {
					return token{}, err
				}
			} else {
				l.unreadRune()
				return token{}, LexError{l.pos, fmt.Sprintf("unexpected symbol %q", string(r))}
			}
			continue
		case unicode.IsLetter(r) || r == '#':
			return l.scanIdentifier(r)
		case unicode.IsNumber(r), r == '-', r == '+':
			l.unreadRune()
			return l.scanNumber()
		default:
			return token{}, LexError{l.pos, fmt.Sprintf("unexpected symbol %q", string(r))}
		}
	}
}

// scanExpression will scan the input for an expression.
// Expressions won't be dissected in detail because it is potentially JavaScript code that will be handled separately.
func (l *lexer) scanExpression() (token, error) {
	// we ignore spaces and tabs at the start
	_, err := l.skip(' ', '\t')
	if err != nil {
		return token{}, err
	}

	t := token{
		tokenType: tokenExpression,
		position:  vit.NewRangeFromPosition(l.pos),
	}

	str, pos, err := l.readExpressionUntil('\n', ';', '}')
	if err != nil {
		return t, err
	}
	t.literal = str // theoretically we could trim spaces at the end here, but we would need to change a lot of tests
	t.position.SetEnd(pos)
	return t, nil
}

// readString will read a JavaScript string from the input.
// It should be called after the first quote (which must be passed as an argument) has been read.
func (l *lexer) readString(delimiter rune) (string, vit.Position, error) {
	var out strings.Builder
	var lastRune rune
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			return "", pos, err
		}
		out.WriteRune(r)

		if lastRune == '\\' { // if this rune was escaped
			goto continueLoop
		} else if r == delimiter {
			return out.String(), pos, nil
		} else if delimiter == '`' && lastRune == '$' && r == '{' { // only important if this is a raw string
			str, pos, err := l.readExpressionUntil('}')
			if err != nil {
				return "", pos, err
			}
			out.WriteString(str)
		}

	continueLoop:
		lastRune = r
	}
}

func (l *lexer) readExpressionUntil(stopRunes ...rune) (string, vit.Position, error) {
	var out strings.Builder
	var bracketType rune      // type of the outer most bracket we are tracking. It will contain the rune that is expected to close the bracket. (one of ')', ']' or '}' )
	var openBrackets int      // number brackets (which one is specified in bracketType) that have been opened
	var potentialComment bool // set to true if we encountered a '/' in the last rune
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if openBrackets > 0 {
					// theoretically we could return an error here, but we will just take the expression as is
					// return "", pos, LexError{pos, fmt.Sprintf("unexpected end of file, expected '%s'", string(bracketType))}
					return out.String(), l.previousPosition, nil
				} else {
					return out.String(), l.previousPosition, nil
				}
			}

			return "", pos, err
		}

		// check if we reached the end of the expression
		if openBrackets == 0 && containsRune(stopRunes, r) {
			l.unreadRune()
			return out.String(), l.previousPosition, nil
		}

		out.WriteRune(r)

		// check if a string starts here ...
		if r == '\'' || r == '"' || r == '`' {
			// ... and if so read all of it
			str, pos, err := l.readString(r)
			if err != nil {
				return "", pos, err
			}
			out.WriteString(str)
			continue
		}

		// check if this might be the start of a comment
		if potentialComment {
			if r == '/' || r == '*' {
				// a line comment started
				str, err := l.readComment(r == '*')
				if err != nil {
					return "", pos, err
				}
				out.WriteString(str)
				continue
			} else {
				potentialComment = false
			}
		} else if r == '/' {
			potentialComment = true
			continue
		}

		// if we are waiting for a bracket to close...
		if openBrackets > 0 {
			// ... and if this is the right type of bracket ...
			if bracketType == r {
				// ... mark it as closed ...
				openBrackets--
				// ... and check if it was the last one we expected
				if openBrackets == 0 {
					bracketType = 0
				}
			}
		} else {
			// we are currently not inside of a bracket so we need to check if a new one starts here
			switch r {
			case '(':
				bracketType = ')'
				openBrackets++
			case '[':
				bracketType = ']'
				openBrackets++
			case '{':
				bracketType = '}'
				openBrackets++
			}
		}
	}
}

func containsRune(set []rune, r rune) bool {
	for _, s := range set {
		if s == r {
			return true
		}
	}
	return false
}

// skipSpaceLikes skips all spaces and newlines. It returns the vit.Position of the first non-space character.
// That rune will still be in the source as unreadRune will be called at the end.
func (l *lexer) skipSpaceLikes() (vit.Position, error) {
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return pos, nil
			}

			return pos, err
		}
		if !unicode.IsSpace(r) { // also matches newlines
			l.unreadRune()
			return pos, nil
		}
	}
}

// skip skips all provided runes. It returns the vit.Position of the first non-skipped rune.
// That rune will still be in the source as unreadRune will be called at the end.
func (l *lexer) skip(runes ...rune) (vit.Position, error) {
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return pos, nil
			}

			return pos, err
		}
		if !containsRune(runes, r) {
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
		position:  vit.NewRangeFromPosition(l.pos),
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
			t.position.SetEnd(pos)
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
		position:  vit.NewRangeFromPosition(l.pos),
	}
	t.position.StartColumn--

	previousPosition := t.position.Start()
	var str strings.Builder
	str.WriteRune(first)
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				t.position.SetEnd(previousPosition)
				t.literal = str.String()
				return t, nil
			}
			return t, err
		}

		if validLiteralRune(r, str.Len() == 0) {
			str.WriteRune(r)
		} else {
			t.position.SetEnd(previousPosition)
			t.literal = str.String()
			l.unreadRune()
			return t, nil
		}

		previousPosition = pos
	}
}

// readComment reads all runes until the end of the comment.
// It returns the read comment excluding the runes that started it (as it will be called after the comment has already started).
// For multiline comments the output WILL include the '*/' and the end of the comment.
// For single line comments the output WILL NOT include the newline that ended the comment.
func (l *lexer) readComment(multiLineComment bool) (string, error) {
	var out strings.Builder
	var starFound bool // only required for multiline comments
	for {
		r, pos, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if multiLineComment {
					return "", unexpectedEOF(pos)
				}
				// a singe line comment can be the last thing in a file
				return out.String(), nil
			}
			return "", err
		}

		if multiLineComment {
			out.WriteRune(r)
			if r == '*' {
				starFound = true
			} else if starFound && r == '/' {
				break // reached end of multiline comment
			} else {
				starFound = false
			}
		} else {
			if r == '\n' {
				// we have reached the end of the line
				l.unreadRune()
				break
			}
			out.WriteRune(r) // only add the rune after we checked if it's a newline
		}
	}

	return out.String(), nil
}

// scanNumber scans a number token.
// The returned error will either be of type LexError or ReadError.
func (l *lexer) scanNumber() (token, error) {
	t := token{
		tokenType: tokenFloat,
		position:  vit.NewRangeFromPosition(l.pos),
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

		switch {
		case r == '.':
			if isFloatingPoint {
				invalid = true // more than one floating point but we will continue to scan
			} else {
				isFloatingPoint = true
			}
		case r == '+', r == '-':
			// a sign that is not at the start ends this number
			if str.Len() != 0 {
				l.unreadRune()
				goto parseAndReturn
			}
		case unicode.IsNumber(r):

		default:
			// the number has ended
			l.unreadRune()
			goto parseAndReturn
		}

		str.WriteRune(r)
		previousPosition = pos
	}

	// now parse the read number

parseAndReturn:
	t.position.SetEnd(previousPosition)

	stringN := str.String()
	// catch some exceptions
	if len(stringN) == 0 || stringN == "." || stringN == "+" || stringN == "-" {
		return t, LexError{t.position.Start(), fmt.Sprintf("number %q is incomplete", stringN)}
	}
	t.literal = stringN

	if invalid {
		t.tokenType = tokenIdentifier
		return t, nil
	}

	if isFloatingPoint {
		_, err := strconv.ParseFloat(stringN, 64)
		if err != nil {
			return t, LexError{t.position.Start(), fmt.Sprintf("%q is not a valid number: %v", stringN, err)}
		}
		t.tokenType = tokenFloat
	} else {
		_, err := strconv.ParseInt(stringN, 10, 64)
		if err != nil {
			return t, LexError{t.position.Start(), fmt.Sprintf("%q is not a valid number: %v", stringN, err)}
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
