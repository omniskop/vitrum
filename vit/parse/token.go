package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

// ========================================= TOKEN TYPE ============================================

type tokenType int

const (
	tokenEOF        tokenType = iota
	tokenIdentifier           // a single word
	tokenExpression           // an expression that describes the value of a property

	tokenLeftParenthesis  // '('
	tokenRightParenthesis // ')'
	tokenLeftBrace        // {
	tokenRightBrace       // }
	tokenLeftBracket      // [
	tokenRightBracket     // ]

	tokenInteger // the actual number value of a token with this type can be read by calling IntValue()
	tokenFloat   // the actual number value of a token with this type can be read by calling FloatValue()
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

// String returns a human readable name of the token type
func (tt tokenType) String() string {
	switch tt {
	case tokenEOF:
		return "EOF"
	case tokenIdentifier:
		return "identifier"
	case tokenExpression:
		return "expression"

	case tokenLeftParenthesis:
		return "'('"
	case tokenRightParenthesis:
		return "')'"
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

// isLiteralType returns true if the tokenType is not just a single character or operator but contains a literal with a meaningful value.
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

// ============================================ TOKEN ==============================================

// A token is a small building block of a vit file
type token struct {
	tokenType tokenType // specific type of this token
	literal   string    // only set for specific token types (see 'isLiteralType')
	position  vit.PositionRange
}

// IntValue converts the literal of this token to an integer.
// It will panic if it is called on a token whose type is not integer.
func (t token) IntValue() int {
	if t.tokenType != tokenInteger {
		panic(fmt.Errorf("token is not a number"))
	}

	// detect number base and prepare string
	var stringToParse = t.literal
	var numberBase int = 10
	if len(t.literal) >= 2 {
		switch t.literal[:2] {
		case "0x", "0X":
			numberBase = 16
			stringToParse = t.literal[2:]
		case "0b", "0B":
			numberBase = 2
			stringToParse = t.literal[2:]
		}
	}

	i, err := strconv.ParseInt(stringToParse, numberBase, 64)
	if err != nil {
		// If we get here the lexer has failed
		// as it should only return number tokens when a conversion is possible.
		panic(fmt.Errorf("lexer wrongly interpreted literal %q at %s as a number but it can't be converted: %w", t.literal, t.position.Start(), err))
	}

	return int(i)
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
		panic(fmt.Errorf("lexer wrongly interpreted literal %q at %s as a number but it can't be converted: %w", t.literal, t.position.Start(), err))
	}

	return f
}

// String returns a human readable description of the token with it's type and literal
func (t token) String() string {
	if isLiteralType(t.tokenType) {
		return fmt.Sprintf("%s %q", t.tokenType, t.literal)
	} else {
		return t.tokenType.String()
	}
}

// ======================================== TOKEN BUFFER ===========================================

// a tokenBuffer takes a token source and provides a number of methods for convenience.
// It will panic when it encounters an error while reading a token from the source.
// This is supposed to make the usage easier to use and reduce clutter in the parser.
// They should be catched at a higher level and not be exposed to the user.
// As described in lexer.Lex the error in the panic will always be of type LexError or ReadError depending on the underlying problem.
type tokenBuffer struct {
	nextToken *token // The next token to be read. Will be set after calling peek and should be returned before reading the source again
	source    func() (token, error)
}

func NewTokenBuffer(source func() (token, error)) *tokenBuffer {
	return &tokenBuffer{
		source: source,
	}
}

// next returns the next available token.
// Might panic, see 'tokenBuffer'.
func (tb *tokenBuffer) next() token {
	if tb.nextToken == nil {
		t, err := tb.source()
		if err != nil {
			panic(err)
		}
		return t
	} else {
		t := *tb.nextToken
		tb.nextToken = nil
		return t
	}
}

// peek returns the next available token without consuming it.
// All subsequent call to 'peek' without a call to 'next' in between will return the same token.
// Might panic, see 'tokenBuffer'.
func (tb *tokenBuffer) peek() token {
	if tb.nextToken == nil {
		t, err := tb.source()
		if err != nil {
			panic(err)
		}
		tb.nextToken = &t
	}
	return *tb.nextToken
}
