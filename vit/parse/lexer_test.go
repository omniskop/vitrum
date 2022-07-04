package parse

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/omniskop/vitrum/vit"
)

func TestLexAll(t *testing.T) {
	tests := []struct {
		input  string
		output []interface{}
	}{{
		input:  "",
		output: []interface{}{tokenEOF},
	}, {
		input:  "hello",
		output: []interface{}{token{tokenIdentifier, "hello", vit.PositionRange{"", 1, 1, 1, 5}}, tokenEOF},
	}, {
		input:  "5 'test",
		output: []interface{}{LexError{}},
	}}
	for _, test := range tests {
		tokens, err := LexAll(strings.NewReader(test.input), "")
		if err != nil {
			if ok, msg := checkLexResult(token{}, err, test.output[0]); !ok {
				t.Logf("input %q:", test.input)
				t.Error(msg)
			}
			continue
		}
		if len(tokens) != len(test.output) {
			t.Logf("input %q:", test.input)
			t.Errorf("returned %d tokens but expected %d", len(tokens), len(test.output))
		}
		for i, token := range tokens {
			if ok, msg := checkLexResult(token, err, test.output[i]); !ok {
				t.Logf("input %q:", test.input)
				t.Error(msg)
			}
		}
	}
}

func TestLex(t *testing.T) {
	tests := []struct {
		input       string
		output      []interface{}
		runOnlyThis bool // set this to true to only run this test
	}{{
		input:  "",
		output: []interface{}{tokenEOF},
	}, {
		input:  "hello",
		output: []interface{}{token{tokenIdentifier, "hello", vit.PositionRange{"", 1, 1, 1, 5}}},
	}, {
		input:  " hello ",
		output: []interface{}{token{tokenIdentifier, "hello", vit.PositionRange{"", 1, 2, 1, 6}}},
	}, {
		input:  "{hello}",
		output: []interface{}{token{tokenLeftBrace, "", vit.PositionRange{"", 1, 1, 1, 1}}, token{tokenIdentifier, "hello", vit.PositionRange{"", 1, 2, 1, 6}}, token{tokenRightBrace, "", vit.PositionRange{"", 1, 7, 1, 7}}},
	}, {
		input:  `[one, two]`,
		output: []interface{}{token{tokenLeftBracket, "", vit.PositionRange{"", 1, 1, 1, 1}}, token{tokenIdentifier, "one", vit.PositionRange{"", 1, 2, 1, 4}}, token{tokenComma, "", vit.PositionRange{"", 1, 5, 1, 5}}, token{tokenIdentifier, "two", vit.PositionRange{"", 1, 7, 1, 9}}, token{tokenRightBracket, "", vit.PositionRange{"", 1, 10, 1, 10}}},
	}, {
		input:  `Rectangle { color: "red" }`,
		output: []interface{}{token{tokenIdentifier, "Rectangle", vit.PositionRange{"", 1, 1, 1, 9}}, token{tokenLeftBrace, "", vit.PositionRange{"", 1, 11, 1, 11}}, token{tokenIdentifier, "color", vit.PositionRange{"", 1, 13, 1, 17}}, token{tokenColon, "", vit.PositionRange{"", 1, 18, 1, 18}}, token{tokenExpression, `"red" `, vit.PositionRange{"", 1, 20, 1, 25}}, token{tokenRightBrace, "", vit.PositionRange{"", 1, 26, 1, 26}}},
	}, {
		input:  "0xx5",
		output: []interface{}{LexError{}, token{tokenIdentifier, "x5", vit.PositionRange{"", 1, 3, 1, 4}}},
	}, {
		input:  "1x1",
		output: []interface{}{token{tokenInteger, "1", vit.PositionRange{"", 1, 1, 1, 1}}, token{tokenIdentifier, "x1", vit.PositionRange{"", 1, 2, 1, 3}}},
	}, {
		input:  "5test: 12.3",
		output: []interface{}{token{tokenInteger, "5", vit.PositionRange{"", 1, 1, 1, 1}}, token{tokenIdentifier, "test", vit.PositionRange{"", 1, 2, 1, 5}}, token{tokenColon, "", vit.PositionRange{"", 1, 6, 1, 6}}, token{tokenExpression, "12.3", vit.PositionRange{"", 1, 8, 1, 11}}},
	}, {
		input:  "one: 1\n    two: 2",
		output: []interface{}{token{tokenIdentifier, "one", vit.PositionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", vit.PositionRange{"", 1, 4, 1, 4}}, token{tokenExpression, "1", vit.PositionRange{"", 1, 6, 1, 6}}, token{tokenNewline, "", vit.PositionRange{"", 1, 7, 1, 7}}, token{tokenIdentifier, "two", vit.PositionRange{"", 2, 5, 2, 7}}, token{tokenColon, "", vit.PositionRange{"", 2, 8, 2, 8}}, token{tokenExpression, "2", vit.PositionRange{"", 2, 10, 2, 10}}},
	}, {
		input:  "import QtQuick.Controls 5.15",
		output: []interface{}{token{tokenIdentifier, "import", vit.PositionRange{"", 1, 1, 1, 6}}, token{tokenIdentifier, "QtQuick", vit.PositionRange{"", 1, 8, 1, 14}}, token{tokenPeriod, "", vit.PositionRange{"", 1, 15, 1, 15}}, token{tokenIdentifier, "Controls", vit.PositionRange{"", 1, 16, 1, 23}}, token{tokenFloat, "5.15", vit.PositionRange{"", 1, 25, 1, 28}}},
	}, {
		input:  `property color nextColor: "blue"`,
		output: []interface{}{token{tokenIdentifier, "property", vit.PositionRange{"", 1, 1, 1, 8}}, token{tokenIdentifier, "color", vit.PositionRange{"", 1, 10, 1, 14}}, token{tokenIdentifier, "nextColor", vit.PositionRange{"", 1, 16, 1, 24}}, token{tokenColon, "", vit.PositionRange{"", 1, 25, 1, 25}}, token{tokenExpression, `"blue"`, vit.PositionRange{"", 1, 27, 1, 32}}},
	}, {
		input:  `:"test"`,
		output: []interface{}{token{tokenColon, "", vit.PositionRange{"", 1, 1, 1, 1}}, token{tokenExpression, `"test"`, vit.PositionRange{"", 1, 2, 1, 7}}},
	}, {
		input:  ":'\n;\\'\"test'",
		output: []interface{}{token{tokenColon, "", vit.PositionRange{"", 1, 1, 1, 1}}, token{tokenExpression, "'\n;\\'\"test'", vit.PositionRange{"", 1, 2, 2, 9}}},
	}, {
		input:  ": {\none\ntwo\n}",
		output: []interface{}{token{tokenColon, "", vit.PositionRange{"", 1, 1, 1, 1}}, token{tokenExpression, "{\none\ntwo\n}", vit.PositionRange{"", 1, 3, 4, 1}}},
	}, {
		input:  "one //two",
		output: []interface{}{token{tokenIdentifier, "one", vit.PositionRange{"", 1, 1, 1, 3}}},
	}, {
		input:  "one//two\nthree",
		output: []interface{}{token{tokenIdentifier, `one`, vit.PositionRange{"", 1, 1, 1, 3}}, token{tokenNewline, "", vit.PositionRange{"", 1, 9, 1, 9}}, token{tokenIdentifier, `three`, vit.PositionRange{"", 2, 1, 2, 5}}},
	}, {
		input:  "one/*two\nthree*/four",
		output: []interface{}{token{tokenIdentifier, `one`, vit.PositionRange{"", 1, 1, 1, 3}}, token{tokenIdentifier, `four`, vit.PositionRange{"", 2, 8, 2, 11}}},
	}, {
		input:  "one: /*stuff*/5",
		output: []interface{}{token{tokenIdentifier, "one", vit.PositionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", vit.PositionRange{"", 1, 4, 1, 4}}, token{tokenExpression, "/*stuff*/5", vit.PositionRange{"", 1, 6, 1, 15}}},
	}, {
		input:  "one: // stuff",
		output: []interface{}{token{tokenIdentifier, "one", vit.PositionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", vit.PositionRange{"", 1, 4, 1, 4}}, token{tokenExpression, "// stuff", vit.PositionRange{"", 1, 6, 1, 13}}},
	}, {
		input:  "one: two//stuff",
		output: []interface{}{token{tokenIdentifier, "one", vit.PositionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", vit.PositionRange{"", 1, 4, 1, 4}}, token{tokenExpression, "two//stuff", vit.PositionRange{"", 1, 6, 1, 15}}},
	}, {
		input:  "one: two/*stuff*/three",
		output: []interface{}{token{tokenIdentifier, "one", vit.PositionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", vit.PositionRange{"", 1, 4, 1, 4}}, token{tokenExpression, "two/*stuff*/three", vit.PositionRange{"", 1, 6, 1, 22}}},
	}, {
		input:  "one: /*#*/ Item{",
		output: []interface{}{token{tokenIdentifier, "one", vit.PositionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", vit.PositionRange{"", 1, 4, 1, 4}}, token{tokenExpression, "/*#*/ Item{", vit.PositionRange{"", 1, 6, 1, 16}}},
	}, {
		input:  ": `${`${`${/*stuff*/'hi'}`}`}`",
		output: []interface{}{token{tokenColon, "", vit.PositionRange{"", 1, 1, 1, 1}}, token{tokenExpression, "`${`${`${/*stuff*/'hi'}`}`}`", vit.PositionRange{"", 1, 3, 1, 30}}},
	}, {
		input:  ": `${//stuff\n5}`",
		output: []interface{}{token{tokenColon, "", vit.PositionRange{"", 1, 1, 1, 1}}, token{tokenExpression, "`${//stuff\n5}`", vit.PositionRange{"", 1, 3, 2, 3}}},
	}, {
		input:  "one/*two",
		output: []interface{}{token{tokenIdentifier, `one`, vit.PositionRange{"", 1, 1, 1, 3}}, LexError{}},
	}, {
		input:  ":`a${\n5}b`",
		output: []interface{}{token{tokenColon, "", vit.PositionRange{"", 1, 1, 1, 1}}, token{tokenExpression, "`a${\n5}b`", vit.PositionRange{"", 1, 2, 2, 4}}},
	}, {
		input: `Rectangle {width: 100; height: 100; gradient: Gradient { }}`,
		output: []interface{}{
			token{tokenIdentifier, "Rectangle", vit.PositionRange{"", 1, 1, 1, 9}},      // Rectangle
			token{tokenLeftBrace, "", vit.PositionRange{"", 1, 11, 1, 11}},              // {
			token{tokenIdentifier, "width", vit.PositionRange{"", 1, 12, 1, 16}},        // width
			token{tokenColon, "", vit.PositionRange{"", 1, 17, 1, 17}},                  // :
			token{tokenExpression, "100", vit.PositionRange{"", 1, 19, 1, 21}},          // 100
			token{tokenSemicolon, "", vit.PositionRange{"", 1, 22, 1, 22}},              // ;
			token{tokenIdentifier, "height", vit.PositionRange{"", 1, 24, 1, 29}},       // height
			token{tokenColon, "", vit.PositionRange{"", 1, 30, 1, 30}},                  // :
			token{tokenExpression, "100", vit.PositionRange{"", 1, 32, 1, 34}},          // 100
			token{tokenSemicolon, "", vit.PositionRange{"", 1, 35, 1, 35}},              // ;
			token{tokenIdentifier, "gradient", vit.PositionRange{"", 1, 37, 1, 44}},     // gradient
			token{tokenColon, "", vit.PositionRange{"", 1, 45, 1, 45}},                  // :
			token{tokenExpression, "Gradient { }", vit.PositionRange{"", 1, 47, 1, 58}}, // Gradient { }
			token{tokenRightBrace, "", vit.PositionRange{"", 1, 59, 1, 59}},             // }
		},
	}}
	// Check if only specific test cases should be run. Allows for easier debugging.
	var testCaseFilterUsed bool
	for _, test := range tests {
		if test.runOnlyThis {
			// fail here to make sure this will never accidentally committed
			t.Error("WARNING: SOME TEST CASES WILL NOT RUN. ONLY USE THIS FOR LOCAL TESTING.")
			testCaseFilterUsed = true
			break
		}
	}

	for _, test := range tests {
		if testCaseFilterUsed && !test.runOnlyThis {
			continue // ignore this test case if another one has runOnlyThis set and it's not this one
		}
		l := NewLexer(strings.NewReader(test.input), "")
		for i, expectedIntf := range test.output {
			tok, err := l.Lex()

			if ok, msg := checkLexResult(tok, err, expectedIntf); !ok {
				t.Logf("token %d of %q:", i, test.input)
				t.Error(msg)
			}
		}
		// We should have reached the end of the input and an EOF token should be returned.
		// If the get a different response the test didn't specify enough output values.
		tok, err := l.Lex()
		if err != nil {
			t.Logf("input %q:", test.input)
			t.Errorf("lexer ended with an error that the test didn't expect: %v", err)
		}
		if tok.tokenType != tokenEOF {
			t.Logf("input %q:", test.input)
			t.Errorf("lexer returned more tokens then expected: %v", tok)
		}
	}
}

func TestScanNumber(t *testing.T) {
	tests := []struct {
		input  string
		output interface{}
	}{{
		input:  "",
		output: LexError{},
	}, {
		input:  "1",
		output: token{tokenInteger, "1", vit.PositionRange{"", 1, 1, 1, 1}},
	}, {
		input:  "1.0",
		output: token{tokenFloat, "1.0", vit.PositionRange{"", 1, 1, 1, 3}},
	}, {
		input:  ".5",
		output: token{tokenFloat, ".5", vit.PositionRange{"", 1, 1, 1, 2}},
	}, {
		input:  "1.",
		output: token{tokenFloat, "1.", vit.PositionRange{"", 1, 1, 1, 2}},
	}, {
		input:  "1.0.3",
		output: token{tokenIdentifier, "1.0.3", vit.PositionRange{"", 1, 1, 1, 5}},
	}, {
		input:  "01.20",
		output: token{tokenFloat, "01.20", vit.PositionRange{"", 1, 1, 1, 5}},
	}, {
		input:  "5",
		output: token{tokenInteger, "5", vit.PositionRange{"", 1, 1, 1, 1}},
	}, {
		input:  "-5",
		output: token{tokenInteger, "-5", vit.PositionRange{"", 1, 1, 1, 2}},
	}, {
		input:  "0x1abcdefABCDEF",
		output: token{tokenInteger, "0x1abcdefABCDEF", vit.PositionRange{"", 1, 1, 1, 15}},
	}, {
		input:  "0b10011001",
		output: token{tokenInteger, "0b10011001", vit.PositionRange{"", 1, 1, 1, 10}},
	}, {
		input:  "12.34",
		output: token{tokenFloat, "12.34", vit.PositionRange{"", 1, 1, 1, 5}},
	}}
	for _, test := range tests {
		l := NewLexer(strings.NewReader(test.input), "")
		tok, err := l.scanNumber()

		if ok, msg := checkLexResult(tok, err, test.output); !ok {
			t.Logf("input %q:", test.input)
			t.Error(msg)
		}

		if _, ok := test.output.(LexError); ok {
			continue // skip further tests if this produced an error
		}

		// validate int and float conversion
		// these methods will panic when they fail, thus I have wrapped them in a subtest
		if tok.tokenType == tokenInteger {
			t.Run(fmt.Sprintf("input %q:", test.input), func(t *testing.T) {
				tok.IntValue()
			})
		} else if tok.tokenType == tokenFloat {
			t.Run(fmt.Sprintf("input %q:", test.input), func(t *testing.T) {
				tok.FloatValue()
			})
		}
	}
}

// checkLexResult takes a the token and error that was returned from the lexer and checks if it matches the expected result.
// If the result differ from the expected false is returned as well as an error message.
// The expected value can either be a token, a tokenType or an error.
func checkLexResult(tok token, err error, expected interface{}) (bool, string) {
	switch expected := expected.(type) {
	case token:
		if err != nil {
			return false, fmt.Sprintf("got error %q but %s was expected", err.Error(), expected.tokenType)
		}
		if tok.tokenType != expected.tokenType {
			return false, fmt.Sprintf("got %q token but %q was expected", tok.tokenType, expected.tokenType)
		}
		if tok.literal != expected.literal {
			return false, fmt.Sprintf("token %q has value %#v (%T) but %#v (%T) was expected", tok.tokenType, tok.literal, tok.literal, expected.literal, expected.literal)
		}
		if !tok.position.Start().IsEqual(expected.position.Start()) {
			return false, fmt.Sprintf("token %q starts at %v but %v was expected", tok.tokenType, tok.position.Start(), expected.position.Start())
		}
		if !tok.position.End().IsEqual(expected.position.End()) {
			return false, fmt.Sprintf("token %q ends at %v but %v was expected", tok.tokenType, tok.position.End(), expected.position.End())
		}
	case tokenType:
		if err != nil {
			return false, fmt.Sprintf("got error %q but %s was expected", err.Error(), expected)
		}
		if tok.tokenType != expected {
			return false, fmt.Sprintf("got %q token but %q was expected", tok.tokenType, expected)
		}
	case error:
		if err == nil {
			return false, fmt.Sprintf("got %q token but error %q was expected", tok.tokenType, expected.Error())
		}
		if !errors.As(err, &expected) {
			// check if the error can be converted into the expected one
			return false, fmt.Sprintf("got error %q but error %q was expected", err.Error(), expected.Error())
		}
	default:
		return false, fmt.Sprintf("unknown expected value: %T", expected)
	}

	return true, ""
}
