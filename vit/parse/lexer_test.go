package parse

import (
	"errors"
	"fmt"
	"strings"
	"testing"
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
		output: []interface{}{token{tokenIdentifier, "hello", positionRange{"", 1, 1, 1, 5}}, tokenEOF},
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
		output: []interface{}{token{tokenIdentifier, "hello", positionRange{"", 1, 1, 1, 5}}},
	}, {
		input:  " hello ",
		output: []interface{}{token{tokenIdentifier, "hello", positionRange{"", 1, 2, 1, 6}}},
	}, {
		input:  "{hello}",
		output: []interface{}{token{tokenLeftBrace, "", positionRange{"", 1, 1, 1, 1}}, token{tokenIdentifier, "hello", positionRange{"", 1, 2, 1, 6}}, token{tokenRightBrace, "", positionRange{"", 1, 7, 1, 7}}},
	}, {
		input:  `[one, two]`,
		output: []interface{}{token{tokenLeftBracket, "", positionRange{"", 1, 1, 1, 1}}, token{tokenIdentifier, "one", positionRange{"", 1, 2, 1, 4}}, token{tokenComma, "", positionRange{"", 1, 5, 1, 5}}, token{tokenIdentifier, "two", positionRange{"", 1, 7, 1, 9}}, token{tokenRightBracket, "", positionRange{"", 1, 10, 1, 10}}},
	}, {
		input:  `Rectangle { color: "red" }`,
		output: []interface{}{token{tokenIdentifier, "Rectangle", positionRange{"", 1, 1, 1, 9}}, token{tokenLeftBrace, "", positionRange{"", 1, 11, 1, 11}}, token{tokenIdentifier, "color", positionRange{"", 1, 13, 1, 17}}, token{tokenColon, "", positionRange{"", 1, 18, 1, 18}}, token{tokenExpression, `"red" `, positionRange{"", 1, 20, 1, 25}}, token{tokenRightBrace, "", positionRange{"", 1, 26, 1, 26}}},
	}, {
		input:  "5",
		output: []interface{}{token{tokenInteger, "5", positionRange{"", 1, 1, 1, 1}}},
	}, {
		input:  "5test: 12.3",
		output: []interface{}{token{tokenInteger, "5", positionRange{"", 1, 1, 1, 1}}, token{tokenIdentifier, "test", positionRange{"", 1, 2, 1, 5}}, token{tokenColon, "", positionRange{"", 1, 6, 1, 6}}, token{tokenExpression, "12.3", positionRange{"", 1, 8, 1, 11}}},
	}, {
		input:  "12.34",
		output: []interface{}{token{tokenFloat, "12.34", positionRange{"", 1, 1, 1, 5}}},
	}, {
		input:  "one: 1\n    two: 2",
		output: []interface{}{token{tokenIdentifier, "one", positionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", positionRange{"", 1, 4, 1, 4}}, token{tokenExpression, "1", positionRange{"", 1, 6, 1, 6}}, token{tokenNewline, "", positionRange{"", 1, 7, 1, 7}}, token{tokenIdentifier, "two", positionRange{"", 2, 5, 2, 7}}, token{tokenColon, "", positionRange{"", 2, 8, 2, 8}}, token{tokenExpression, "2", positionRange{"", 2, 10, 2, 10}}},
	}, {
		input:  "import QtQuick.Controls 5.15",
		output: []interface{}{token{tokenIdentifier, "import", positionRange{"", 1, 1, 1, 6}}, token{tokenIdentifier, "QtQuick", positionRange{"", 1, 8, 1, 14}}, token{tokenPeriod, "", positionRange{"", 1, 15, 1, 15}}, token{tokenIdentifier, "Controls", positionRange{"", 1, 16, 1, 23}}, token{tokenFloat, "5.15", positionRange{"", 1, 25, 1, 28}}},
	}, {
		input:  `property color nextColor: "blue"`,
		output: []interface{}{token{tokenIdentifier, "property", positionRange{"", 1, 1, 1, 8}}, token{tokenIdentifier, "color", positionRange{"", 1, 10, 1, 14}}, token{tokenIdentifier, "nextColor", positionRange{"", 1, 16, 1, 24}}, token{tokenColon, "", positionRange{"", 1, 25, 1, 25}}, token{tokenExpression, `"blue"`, positionRange{"", 1, 27, 1, 32}}},
	}, {
		input:  `:"test"`,
		output: []interface{}{token{tokenColon, "", positionRange{"", 1, 1, 1, 1}}, token{tokenExpression, `"test"`, positionRange{"", 1, 2, 1, 7}}},
	}, {
		input:  ":'\n;\\'\"test'",
		output: []interface{}{token{tokenColon, "", positionRange{"", 1, 1, 1, 1}}, token{tokenExpression, "'\n;\\'\"test'", positionRange{"", 1, 2, 2, 9}}},
	}, {
		input:  ": {\none\ntwo\n}",
		output: []interface{}{token{tokenColon, "", positionRange{"", 1, 1, 1, 1}}, token{tokenExpression, "{\none\ntwo\n}", positionRange{"", 1, 3, 4, 1}}},
	}, {
		input:  "one //two",
		output: []interface{}{token{tokenIdentifier, "one", positionRange{"", 1, 1, 1, 3}}},
	}, {
		input:  "one//two\nthree",
		output: []interface{}{token{tokenIdentifier, `one`, positionRange{"", 1, 1, 1, 3}}, token{tokenNewline, "", positionRange{"", 1, 9, 1, 9}}, token{tokenIdentifier, `three`, positionRange{"", 2, 1, 2, 5}}},
	}, {
		input:  "one/*two\nthree*/four",
		output: []interface{}{token{tokenIdentifier, `one`, positionRange{"", 1, 1, 1, 3}}, token{tokenIdentifier, `four`, positionRange{"", 2, 8, 2, 11}}},
	}, {
		input:  "one: /*stuff*/5",
		output: []interface{}{token{tokenIdentifier, "one", positionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", positionRange{"", 1, 4, 1, 4}}, token{tokenExpression, "5", positionRange{"", 1, 15, 1, 15}}},
	}, {
		input:  "one: // stuff",
		output: []interface{}{token{tokenIdentifier, "one", positionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", positionRange{"", 1, 4, 1, 4}}, LexError{position{"", 1, 6}, "unexpected token: '//'"}},
	}, {
		input:  "one: two//stuff",
		output: []interface{}{token{tokenIdentifier, "one", positionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", positionRange{"", 1, 4, 1, 4}}, token{tokenExpression, "two", positionRange{"", 1, 6, 1, 8}}},
	}, {
		input:  "one: two/*stuff*/three",
		output: []interface{}{token{tokenIdentifier, "one", positionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", positionRange{"", 1, 4, 1, 4}}, token{tokenExpression, "twothree", positionRange{"", 1, 6, 1, 22}}},
	}, {
		input:  "one: /*#*/ Item{",
		output: []interface{}{token{tokenIdentifier, "one", positionRange{"", 1, 1, 1, 3}}, token{tokenColon, "", positionRange{"", 1, 4, 1, 4}}, token{tokenIdentifier, "Item", positionRange{"", 1, 12, 1, 15}}, token{tokenLeftBrace, "", positionRange{"", 1, 16, 1, 16}}},
	}, {
		input:  "one/*two",
		output: []interface{}{token{tokenIdentifier, `one`, positionRange{"", 1, 1, 1, 3}}, LexError{}},
	}, {
		input:  ":`a${\n5}b`",
		output: []interface{}{token{tokenColon, "", positionRange{"", 1, 1, 1, 1}}, token{tokenExpression, "`a${\n5}b`", positionRange{"", 1, 2, 2, 4}}},
	}, {
		input: `Rectangle {width: 100; height: 100; gradient: Gradient { GradientStop { position: 0.0; color: "yellow" }; GradientStop { position: 1.0; color: "green" } }}`,
		output: []interface{}{
			token{tokenIdentifier, "Rectangle", positionRange{"", 1, 1, 1, 9}},        // Rectangle
			token{tokenLeftBrace, "", positionRange{"", 1, 11, 1, 11}},                // {
			token{tokenIdentifier, "width", positionRange{"", 1, 12, 1, 16}},          // width
			token{tokenColon, "", positionRange{"", 1, 17, 1, 17}},                    // :
			token{tokenExpression, "100", positionRange{"", 1, 19, 1, 21}},            // 100
			token{tokenSemicolon, "", positionRange{"", 1, 22, 1, 22}},                // ;
			token{tokenIdentifier, "height", positionRange{"", 1, 24, 1, 29}},         // height
			token{tokenColon, "", positionRange{"", 1, 30, 1, 30}},                    // :
			token{tokenExpression, "100", positionRange{"", 1, 32, 1, 34}},            // 100
			token{tokenSemicolon, "", positionRange{"", 1, 35, 1, 35}},                // ;
			token{tokenIdentifier, "gradient", positionRange{"", 1, 37, 1, 44}},       // gradient
			token{tokenColon, "", positionRange{"", 1, 45, 1, 45}},                    // :
			token{tokenIdentifier, "Gradient", positionRange{"", 1, 47, 1, 54}},       // Gradient
			token{tokenLeftBrace, "", positionRange{"", 1, 56, 1, 56}},                // {
			token{tokenIdentifier, "GradientStop", positionRange{"", 1, 58, 1, 69}},   // GradientStop
			token{tokenLeftBrace, "", positionRange{"", 1, 71, 1, 71}},                // {
			token{tokenIdentifier, "position", positionRange{"", 1, 73, 1, 80}},       // position
			token{tokenColon, "", positionRange{"", 1, 81, 1, 81}},                    // :
			token{tokenExpression, "0.0", positionRange{"", 1, 83, 1, 85}},            // 0.0
			token{tokenSemicolon, "", positionRange{"", 1, 86, 1, 86}},                // ;
			token{tokenIdentifier, "color", positionRange{"", 1, 88, 1, 92}},          // color
			token{tokenColon, "", positionRange{"", 1, 93, 1, 93}},                    // :
			token{tokenExpression, `"yellow" `, positionRange{"", 1, 95, 1, 103}},     // "yellow"
			token{tokenRightBrace, "", positionRange{"", 1, 104, 1, 104}},             // }
			token{tokenSemicolon, "", positionRange{"", 1, 105, 1, 105}},              // ;
			token{tokenIdentifier, "GradientStop", positionRange{"", 1, 107, 1, 118}}, // GradientStop
			token{tokenLeftBrace, "", positionRange{"", 1, 120, 1, 120}},              // {
			token{tokenIdentifier, "position", positionRange{"", 1, 122, 1, 129}},     // position
			token{tokenColon, "", positionRange{"", 1, 130, 1, 130}},                  // :
			token{tokenExpression, "1.0", positionRange{"", 1, 132, 1, 134}},          // 1.0
			token{tokenSemicolon, "", positionRange{"", 1, 135, 1, 135}},              // ;
			token{tokenIdentifier, "color", positionRange{"", 1, 137, 1, 141}},        // color
			token{tokenColon, "", positionRange{"", 1, 142, 1, 142}},                  // :
			token{tokenExpression, `"green" `, positionRange{"", 1, 144, 1, 151}},     // "green"
			token{tokenRightBrace, "", positionRange{"", 1, 152, 1, 152}},             // }
			token{tokenRightBrace, "", positionRange{"", 1, 154, 1, 154}},             // }
			token{tokenRightBrace, "", positionRange{"", 1, 155, 1, 155}},             // }
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
		output: token{tokenInteger, "1", positionRange{"", 1, 1, 1, 1}},
	}, {
		input:  "1.0",
		output: token{tokenFloat, "1.0", positionRange{"", 1, 1, 1, 3}},
	}, {
		input:  ".5",
		output: token{tokenFloat, ".5", positionRange{"", 1, 1, 1, 2}},
	}, {
		input:  "1.",
		output: token{tokenFloat, "1.", positionRange{"", 1, 1, 1, 2}},
	}, {
		input:  "1.0.3",
		output: token{tokenIdentifier, "1.0.3", positionRange{"", 1, 1, 1, 5}},
	}, {
		input:  "01.20",
		output: token{tokenFloat, "01.20", positionRange{"", 1, 1, 1, 5}},
	}}
	for _, test := range tests {
		l := NewLexer(strings.NewReader(test.input), "")
		tok, err := l.scanNumber()

		if ok, msg := checkLexResult(tok, err, test.output); !ok {
			t.Logf("input %q:", test.input)
			t.Error(msg)
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
