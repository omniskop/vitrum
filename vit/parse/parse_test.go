package parse

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omniskop/vitrum/vit"
)

func TestParseError(t *testing.T) {
	wrapped := errors.New("test error")
	err := ParseError{
		pos: vit.PositionRange{
			FilePath:    "test.vit",
			StartLine:   1,
			StartColumn: 2,
			EndLine:     3,
			EndColumn:   4,
		},
		err: wrapped,
	}
	if err.Error() != `test.vit:1:2: test error` {
		t.Errorf(`Expected error to be 'test.vit:1:2: test error', got '%s'`, err.Error())
	}
	if !err.Is(err) {
		t.Errorf("parseError.Is is not identifying itself correctly")
	}
	if err.Is(fmt.Errorf("test")) {
		t.Errorf("parseError.Is is identifying other errors incorrectly")
	}
	if err.Unwrap() != wrapped {
		t.Errorf("parseError.Unwrap is not returning the underlying error")
	}
}

func TestUnexpectedTokenError(t *testing.T) {
	err := unexpectedTokenError{
		got: token{
			tokenType: tokenIdentifier,
			literal:   "test",
		},
	}
	if err.Error() != `unexpected identifier "test"` {
		t.Errorf(`Error string is %q but shouldn't`, err.Error())
	}

	err.expected = []tokenType{tokenInteger, tokenFloat, tokenString}
	if err.Error() != `unexpected identifier "test", expected integer, float or string` {
		t.Errorf(`Error string is %q`, err.Error())
	}

	if !err.Is(err) {
		t.Errorf("unexpectedTokenError.Is is not identifying itself correctly")
	}
	if err.Is(fmt.Errorf("test")) {
		t.Errorf("unexpectedTokenError.Is is identifying other errors incorrectly")
	}
}

func TestUnexpectedToken(t *testing.T) {
	tok := token{
		tokenType: tokenIdentifier,
		literal:   "value",
		position:  vit.PositionRange{"test", 0, 0, 0, 0},
	}
	err := unexpectedToken(tok, tokenInteger)
	if err.pos != tok.position {
		t.Errorf("unexpectedToken set 'pos' incorrectly to %+v", tok.position)
	}
	var unexpErr unexpectedTokenError
	if !errors.As(err, &unexpErr) {
		t.Errorf("the parseError returned by unexpectedToken does not contain an unexpectedTokenError")
	}
	if unexpErr.got != tok {
		t.Errorf("the 'got' field of the unexpectedTokenError created by unexpectedToken was not set correctly")
	}
	if len(unexpErr.expected) != 1 || unexpErr.expected[0] != tokenInteger {
		t.Errorf("the 'expected' field of the unexpectedTokenError created by unexpectedToken was not set correctly")
	}
}

type testTokenSource struct {
	tokens []token
	index  int
}

func (ts *testTokenSource) Next() (token, error) {
	if ts.index >= len(ts.tokens) {
		return token{}, io.EOF
	}
	tok := ts.tokens[ts.index]
	ts.index++
	return tok, nil
}

func NewTestTokenBuffer(tokens ...token) *tokenBuffer {
	src := testTokenSource{
		tokens: tokens,
	}
	return NewTokenBuffer(src.Next)
}

func TestIgnoreTokens(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok && errors.Is(err, io.EOF) {
				t.Log("read more tokens from source than it should've and thus paniced:")
				t.Error(r)
			}
			t.Fatalf("test paniced in an unecpected way: %v", r)
		}
	}()

	// the tokens used here are arbitrary
	buf := NewTestTokenBuffer(token{tokenType: tokenLess})
	tok := ignoreTokens(buf, tokenNewline)
	if tok.tokenType != tokenLess {
		t.Errorf("Expected token type %d, got %d", tokenLess, tok.tokenType)
	}
	if buf.next().tokenType != tokenLess {
		t.Errorf("ignoreTokens seems to have consumed the token it shouldn't have")
	}

	buf = NewTestTokenBuffer(token{tokenType: tokenNewline}, token{tokenType: tokenGreater}, token{tokenType: tokenInteger})
	tok = ignoreTokens(buf, tokenNewline, tokenGreater)
	if tok.tokenType != tokenInteger {
		t.Errorf("Expected token type %d, got %d", tokenInteger, tok.tokenType)
	}
	tok = buf.next()
	if tok.tokenType == tokenNewline || tok.tokenType == tokenGreater {
		t.Errorf("ignoreTokens left a token in the buffer that it should have removed: %q", tok)
	}
	if tok.tokenType != tokenInteger {
		t.Errorf("Expected token type %d, got %d", tokenInteger, tok.tokenType)
	}
}

func TestExpectToken(t *testing.T) {
	buf := NewTestTokenBuffer(token{tokenType: tokenColon}, token{tokenType: tokenInteger})
	tok, err := expectToken(buf.next, tokenColon)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if tok.tokenType != tokenColon {
		t.Errorf("Expected token type %d, got %d", tokenColon, tok.tokenType)
	}
	tok, err = expectToken(buf.next, tokenNewline)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestLiteralsToStrings(t *testing.T) {
	result := literalsToStrings([]token{{
		literal: "one",
	}, {
		literal: "two",
	}})
	if len(result) != 2 || result[0] != "one" || result[1] != "two" {
		t.Errorf("Expected %+v, got %+v", []string{"one", "two"}, result)
	}
	result = literalsToStrings([]token{})
	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %+v", result)
	}
}

// =========================================== PARSER ==============================================

var validFile = `
// comment
import One 1.23
// comment
import Two.Three 4.56 // comment

/*
    This is a cool file
*/

Item {
    id: rect
    anchors.left: parent.left + 10; // comment
    affe: /*#invalid stuff#*/ Item {
        one: 1
    }

	property bool local: true;

    Label {
        wrapMode: Text.WordWrap
        text: "What a wonderful world"
    }
}
`

var validDocumentold = &VitDocument{
	Imports: []importStatement{
		{namespace: []string{"One"}, version: "1.23"},
		{namespace: []string{"Two", "Three"}, version: "4.56"},
	},
	Components: []*vit.ComponentDefinition{
		{
			BaseName: "Item",
			ID:       "rect",
			Properties: []vit.PropertyDefinition{
				{Identifier: []string{"anchors", "left"}, Expression: "parent.left + 10"},
				{Identifier: []string{"affe"}, Components: []*vit.ComponentDefinition{{BaseName: "Item", Properties: []vit.PropertyDefinition{{Identifier: []string{"one"}, Expression: "1"}}}}},
				{Identifier: []string{"local"}, VitType: "bool", Expression: "true"},
			},
			Children: []*vit.ComponentDefinition{
				{
					BaseName: "Label",
					Properties: []vit.PropertyDefinition{
						{Identifier: []string{"wrapMode"}, Expression: "Text.WordWrap"},
						{Identifier: []string{"text"}, Expression: `"What a wonderful world"`},
					},
				},
			},
		},
	},
}

var validDocument = &VitDocument{
	Imports: []importStatement{
		{namespace: []string{"One"}, version: "1.23", position: vit.PositionRange{FilePath: "test", StartLine: 3, StartColumn: 1, EndLine: 3, EndColumn: 15}},
		{namespace: []string{"Two", "Three"}, version: "4.56", position: vit.PositionRange{FilePath: "test", StartLine: 5, StartColumn: 1, EndLine: 5, EndColumn: 21}},
	},
	Components: []*vit.ComponentDefinition{
		{
			Pos:      vit.PositionRange{FilePath: "", StartLine: 0, StartColumn: 0, EndLine: 0, EndColumn: 0},
			BaseName: "Item",
			ID:       "rect",
			Properties: []vit.PropertyDefinition{
				{
					Pos:         vit.PositionRange{FilePath: "test", StartLine: 13, StartColumn: 5, EndLine: 13, EndColumn: 34},
					Identifier:  []string{"anchors", "left"},
					Expression:  "parent.left + 10",
					Components:  nil,
					ReadOnly:    false,
					Static:      false,
					StaticValue: interface{}(nil),
				},
				{
					Pos:        vit.PositionRange{FilePath: "test", StartLine: 14, StartColumn: 5, EndLine: 16, EndColumn: 5},
					Identifier: []string{"affe"},
					Expression: "/*#invalid stuff#*/ Item {\n        one: 1\n    }",
					Components: []*vit.ComponentDefinition{{
						Pos:      vit.PositionRange{FilePath: "", StartLine: 0, StartColumn: 0, EndLine: 0, EndColumn: 0},
						BaseName: "Item",
						ID:       "",
						Properties: []vit.PropertyDefinition{
							{
								Pos:        vit.PositionRange{FilePath: "test", StartLine: 15, StartColumn: 9, EndLine: 15, EndColumn: 14},
								Identifier: []string{"one"},
								VitType:    "", Expression: "1", Components: nil,
								ReadOnly:    false,
								Static:      false,
								StaticValue: interface{}(nil),
							},
						},
						Children:     []*vit.ComponentDefinition(nil),
						Enumerations: []vit.Enumeration(nil),
					}},
					ReadOnly:    false,
					Static:      false,
					StaticValue: interface{}(nil),
				},
				{
					Pos:         vit.PositionRange{FilePath: "test", StartLine: 18, StartColumn: 2, EndLine: 18, EndColumn: 26},
					Identifier:  []string{"local"},
					VitType:     "bool",
					Expression:  "true",
					Components:  nil,
					ReadOnly:    false,
					Static:      false,
					StaticValue: interface{}(nil),
				},
			},
			Children: []*vit.ComponentDefinition{
				{
					BaseName: "Label",
					Properties: []vit.PropertyDefinition{
						{Identifier: []string{"wrapMode"}, Expression: "Text.WordWrap", Pos: vit.PositionRange{FilePath: "test", StartLine: 21, StartColumn: 9, EndLine: 21, EndColumn: 31}},
						{Identifier: []string{"text"}, Expression: `"What a wonderful world"`, Pos: vit.PositionRange{FilePath: "test", StartLine: 22, StartColumn: 9, EndLine: 22, EndColumn: 38}},
					},
				},
			},
		},
	},
}

func TestParse(t *testing.T) {
	// we lex an example file in here but we only really care about parser and not the lexer
	l := NewLexer(strings.NewReader(validFile), "test")
	buf := NewTokenBuffer(l.Lex)
	doc, err := Parse(buf)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	options := []cmp.Option{cmp.AllowUnexported(VitDocument{}, vit.ComponentDefinition{}, vit.PropertyDefinition{}, importStatement{})}
	if !cmp.Equal(validDocument, doc, options...) {
		t.Log("Parsed document deviated from expected result:")
		t.Log("- expected")
		t.Log("+ found")
		t.Log(cmp.Diff(validDocument, doc, options...))
		t.Fail()
	}
}
