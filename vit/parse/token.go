package parse

import (
	"fmt"
	"strconv"
	"strings"
)

// ========================================= TOKEN TYPE ============================================

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenIdentifier
	tokenExpression

	tokenLeftBrace    // {
	tokenRightBrace   // }
	tokenLeftBracket  // [
	tokenRightBracket // ]

	tokenInteger
	tokenFloat
	tokenString

	tokenLess       // <
	tokenGreater    // >
	tokenAssignment // =

	tokenComma     // ,
	tokenColon     // :
	tokenSemicolon // ;
	tokenNewline   // \n
	tokenPeriod    // .
)

func (tt tokenType) String() string {
	switch tt {
	case tokenEOF:
		return "EOF"
	case tokenIdentifier:
		return "identifier"
	case tokenExpression:
		return "expression"

	case tokenLeftBrace:
		return "'{'"
	case tokenRightBrace:
		return "'}'"
	case tokenLeftBracket:
		return "'['"
	case tokenRightBracket:
		return "']'"

	case tokenInteger:
		return "integer"
	case tokenFloat:
		return "float"
	case tokenString:
		return "string"

	case tokenLess:
		return "'<'"
	case tokenGreater:
		return "'>'"
	case tokenAssignment:
		return "'='"

	case tokenComma:
		return "','"
	case tokenColon:
		return "':'"
	case tokenSemicolon:
		return "';'"
	case tokenNewline:
		return `newline`
	case tokenPeriod:
		return "'.'"
	default:
		return "<unknwon token>"
	}
}

// isLiteralType returns true if the tokenType is not just a single character or operator but contains a literal with meaningful value.
func isLiteralType(tt tokenType) bool {
	return tt == tokenInteger || tt == tokenFloat || tt == tokenString || tt == tokenIdentifier || tt == tokenExpression
}

// joinTokens returns a nice looking string of all stringified tokenTypes separated by a comma.
// The last item will be preceded by an or instead of a comma.
func joinTokenTypes(types []tokenType) string {
	if len(types) == 0 {
		return ""
	}

	if len(types) == 1 {
		return types[0].String()
	}

	var out strings.Builder
	out.Grow(20)
	for i := 0; i < len(types)-1; i++ {
		out.WriteString(types[i].String())
		if i < len(types)-2 {
			out.WriteString(", ")
		}
	}

	out.WriteString(" or ")
	out.WriteString(types[len(types)-1].String())

	return out.String()
}

// ========================================== POSITION =============================================

// position describes a specific position in a file
type position struct {
	filePath string
	line     int // line inside the file starting at 1
	column   int // column inside the line starting at 1 (this is pointing to the rune, not the byte)
}

func (p position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.filePath, p.line, p.column)
}

func (p position) IsEqual(o position) bool {
	return p.filePath == o.filePath && p.line == o.line && p.column == o.column
}

// ============================================ TOKEN ==============================================

type token struct {
	tokenType tokenType
	literal   string   // only set for specific token types (see 'isLiteralType')
	start     position // position of first rune
	end       position // position of last rune
}

// IntValue converts the literal of this token to an integer.
// It will panic if it is called on a token whose type is not integer.
func (t token) IntValue() int {
	if t.tokenType != tokenInteger {
		panic(fmt.Errorf("token is not a number"))
	}

	i, err := strconv.Atoi(t.literal)
	if err != nil {
		// If we get here the lexer has failed
		// as it should only return number tokens when a conversion is possible.
		panic(fmt.Errorf("lexer wrongly interpreted literal %q at %s as a number but it can't be converted: %w", t.literal, t.start, err))
	}

	return i
}

// FloatValue converts the literal of this token to a float.
// It will panic if it is called on a token whose type is not float.
func (t token) FloatValue() float64 {
	if t.tokenType != tokenFloat {
		panic(fmt.Errorf("token is not a number"))
	}

	f, err := strconv.ParseFloat(t.literal, 64)
	if err != nil {
		// If we get here the lexer has failed
		// as it should only return number tokens when a conversion is possible.
		panic(fmt.Errorf("lexer wrongly interpreted literal %q at %s as a number but it can't be converted: %w", t.literal, t.start, err))
	}

	return f
}

func (t token) String() string {
	if isLiteralType(t.tokenType) {
		return fmt.Sprintf("%s %q", t.tokenType, t.literal)
	} else {
		return t.tokenType.String()
	}
}

// ======================================== TOKEN BUFFER ===========================================

type tokenBuffer struct {
	buffer []token
	source func() (token, error)
}

func NewTokenBuffer(source func() (token, error)) *tokenBuffer {
	return &tokenBuffer{
		source: source,
	}
}

func (tb *tokenBuffer) next() token {
	if len(tb.buffer) == 0 {
		t, err := tb.source()
		if err != nil {
			panic(err)
		}
		return t
	}
	t := tb.buffer[0]
	tb.buffer = tb.buffer[1:]
	return t
}

func (tb *tokenBuffer) peek() token {
	if len(tb.buffer) == 0 {
		tb.loadOne()
	}

	return tb.buffer[0]
}

func (tb *tokenBuffer) loadOne() {
	t, err := tb.source()
	if err != nil {
		panic(err)
	}

	tb.buffer = append(tb.buffer, t)
}
