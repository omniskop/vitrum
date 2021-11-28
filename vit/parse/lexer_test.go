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
		output: []interface{}{token{tokenIdentifier, "hello", position{"", 1, 1}, position{"", 1, 5}}, tokenEOF},
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
		output: []interface{}{token{tokenIdentifier, "hello", position{"", 1, 1}, position{"", 1, 5}}},
	}, {
		input:  " hello ",
		output: []interface{}{token{tokenIdentifier, "hello", position{"", 1, 2}, position{"", 1, 6}}},
	}, {
		input:  "{hello}",
		output: []interface{}{token{tokenLeftBrace, "", position{"", 1, 1}, position{"", 1, 1}}, token{tokenIdentifier, "hello", position{"", 1, 2}, position{"", 1, 6}}, token{tokenRightBrace, "", position{"", 1, 7}, position{"", 1, 7}}},
	}, {
		input:  `[one, two]`,
		output: []interface{}{token{tokenLeftBracket, "", position{"", 1, 1}, position{"", 1, 1}}, token{tokenIdentifier, "one", position{"", 1, 2}, position{"", 1, 4}}, token{tokenComma, "", position{"", 1, 5}, position{"", 1, 5}}, token{tokenIdentifier, "two", position{"", 1, 7}, position{"", 1, 9}}, token{tokenRightBracket, "", position{"", 1, 10}, position{"", 1, 10}}},
	}, {
		input:  `Rectangle { color: "red" }`,
		output: []interface{}{token{tokenIdentifier, "Rectangle", position{"", 1, 1}, position{"", 1, 9}}, token{tokenLeftBrace, "", position{"", 1, 11}, position{"", 1, 11}}, token{tokenIdentifier, "color", position{"", 1, 13}, position{"", 1, 17}}, token{tokenColon, "", position{"", 1, 18}, position{"", 1, 18}}, token{tokenExpression, `"red" `, position{"", 1, 20}, position{"", 1, 25}}, token{tokenRightBrace, "", position{"", 1, 26}, position{"", 1, 26}}},
	}, {
		input:  "5",
		output: []interface{}{token{tokenInteger, "5", position{"", 1, 1}, position{"", 1, 1}}},
	}, {
		input:  "5test: 12.3",
		output: []interface{}{token{tokenInteger, "5", position{"", 1, 1}, position{"", 1, 1}}, token{tokenIdentifier, "test", position{"", 1, 2}, position{"", 1, 5}}, token{tokenColon, "", position{"", 1, 6}, position{"", 1, 6}}, token{tokenExpression, "12.3", position{"", 1, 8}, position{"", 1, 11}}},
	}, {
		input:  "12.34",
		output: []interface{}{token{tokenFloat, "12.34", position{"", 1, 1}, position{"", 1, 5}}},
	}, {
		input:  "one: 1\n    two: 2",
		output: []interface{}{token{tokenIdentifier, "one", position{"", 1, 1}, position{"", 1, 3}}, token{tokenColon, "", position{"", 1, 4}, position{"", 1, 4}}, token{tokenExpression, "1", position{"", 1, 6}, position{"", 1, 6}}, token{tokenNewline, "", position{"", 1, 7}, position{"", 1, 7}}, token{tokenIdentifier, "two", position{"", 2, 5}, position{"", 2, 7}}, token{tokenColon, "", position{"", 2, 8}, position{"", 2, 8}}, token{tokenExpression, "2", position{"", 2, 10}, position{"", 2, 10}}},
	}, {
		input:  "import QtQuick.Controls 5.15",
		output: []interface{}{token{tokenIdentifier, "import", position{"", 1, 1}, position{"", 1, 6}}, token{tokenIdentifier, "QtQuick", position{"", 1, 8}, position{"", 1, 14}}, token{tokenPeriod, "", position{"", 1, 15}, position{"", 1, 15}}, token{tokenIdentifier, "Controls", position{"", 1, 16}, position{"", 1, 23}}, token{tokenFloat, "5.15", position{"", 1, 25}, position{"", 1, 28}}},
	}, {
		input:  `property color nextColor: "blue"`,
		output: []interface{}{token{tokenIdentifier, "property", position{"", 1, 1}, position{"", 1, 8}}, token{tokenIdentifier, "color", position{"", 1, 10}, position{"", 1, 14}}, token{tokenIdentifier, "nextColor", position{"", 1, 16}, position{"", 1, 24}}, token{tokenColon, "", position{"", 1, 25}, position{"", 1, 25}}, token{tokenExpression, `"blue"`, position{"", 1, 27}, position{"", 1, 32}}},
	}, {
		input:  `:"test"`,
		output: []interface{}{token{tokenColon, "", position{"", 1, 1}, position{"", 1, 1}}, token{tokenExpression, `"test"`, position{"", 1, 2}, position{"", 1, 7}}},
	}, {
		input:  ":'\n;\\'\"test'",
		output: []interface{}{token{tokenColon, "", position{"", 1, 1}, position{"", 1, 1}}, token{tokenExpression, "'\n;\\'\"test'", position{"", 1, 2}, position{"", 2, 9}}},
	}, {
		input:  ": {\none\ntwo\n}",
		output: []interface{}{token{tokenColon, "", position{"", 1, 1}, position{"", 1, 1}}, token{tokenExpression, "{\none\ntwo\n}", position{"", 1, 3}, position{"", 4, 1}}},
	}, {
		input:  "one //two",
		output: []interface{}{token{tokenIdentifier, "one", position{"", 1, 1}, position{"", 1, 3}}},
	}, {
		input:  "one//two\nthree",
		output: []interface{}{token{tokenIdentifier, `one`, position{"", 1, 1}, position{"", 1, 3}}, token{tokenNewline, "", position{"", 1, 9}, position{"", 1, 9}}, token{tokenIdentifier, `three`, position{"", 2, 1}, position{"", 2, 5}}},
	}, {
		input:  "one/*two\nthree*/four",
		output: []interface{}{token{tokenIdentifier, `one`, position{"", 1, 1}, position{"", 1, 3}}, token{tokenIdentifier, `four`, position{"", 2, 8}, position{"", 2, 11}}},
	}, {
		input:  "one: /*stuff*/5",
		output: []interface{}{token{tokenIdentifier, "one", position{"", 1, 1}, position{"", 1, 3}}, token{tokenColon, "", position{"", 1, 4}, position{"", 1, 4}}, token{tokenExpression, "5", position{"", 1, 15}, position{"", 1, 15}}},
	}, {
		input:  "one: // stuff",
		output: []interface{}{token{tokenIdentifier, "one", position{"", 1, 1}, position{"", 1, 3}}, token{tokenColon, "", position{"", 1, 4}, position{"", 1, 4}}, LexError{position{"", 1, 6}, "unexpected token: '//'"}},
	}, {
		input:  "one: two//stuff",
		output: []interface{}{token{tokenIdentifier, "one", position{"", 1, 1}, position{"", 1, 3}}, token{tokenColon, "", position{"", 1, 4}, position{"", 1, 4}}, token{tokenExpression, "two", position{"", 1, 6}, position{"", 1, 8}}},
	}, {
		input:  "one: two/*stuff*/three",
		output: []interface{}{token{tokenIdentifier, "one", position{"", 1, 1}, position{"", 1, 3}}, token{tokenColon, "", position{"", 1, 4}, position{"", 1, 4}}, token{tokenExpression, "twothree", position{"", 1, 6}, position{"", 1, 22}}},
	}, {
		input:  "one: /*#*/ Item{",
		output: []interface{}{token{tokenIdentifier, "one", position{"", 1, 1}, position{"", 1, 3}}, token{tokenColon, "", position{"", 1, 4}, position{"", 1, 4}}, token{tokenIdentifier, "Item", position{"", 1, 12}, position{"", 1, 15}}, token{tokenLeftBrace, "", position{"", 1, 16}, position{"", 1, 16}}},
	}, {
		input:  "one/*two",
		output: []interface{}{token{tokenIdentifier, `one`, position{"", 1, 1}, position{"", 1, 3}}, LexError{}},
	}, {
		input:  ":`a${\n5}b`",
		output: []interface{}{token{tokenColon, "", position{"", 1, 1}, position{"", 1, 1}}, token{tokenExpression, "`a${\n5}b`", position{"", 1, 2}, position{"", 2, 4}}},
	}, {
		input: `Rectangle {width: 100; height: 100; gradient: Gradient { GradientStop { position: 0.0; color: "yellow" }; GradientStop { position: 1.0; color: "green" } }}`,
		output: []interface{}{
			token{tokenIdentifier, "Rectangle", position{"", 1, 1}, position{"", 1, 9}},        // Rectangle
			token{tokenLeftBrace, "", position{"", 1, 11}, position{"", 1, 11}},                // {
			token{tokenIdentifier, "width", position{"", 1, 12}, position{"", 1, 16}},          // width
			token{tokenColon, "", position{"", 1, 17}, position{"", 1, 17}},                    // :
			token{tokenExpression, "100", position{"", 1, 19}, position{"", 1, 21}},            // 100
			token{tokenSemicolon, "", position{"", 1, 22}, position{"", 1, 22}},                // ;
			token{tokenIdentifier, "height", position{"", 1, 24}, position{"", 1, 29}},         // height
			token{tokenColon, "", position{"", 1, 30}, position{"", 1, 30}},                    // :
			token{tokenExpression, "100", position{"", 1, 32}, position{"", 1, 34}},            // 100
			token{tokenSemicolon, "", position{"", 1, 35}, position{"", 1, 35}},                // ;
			token{tokenIdentifier, "gradient", position{"", 1, 37}, position{"", 1, 44}},       // gradient
			token{tokenColon, "", position{"", 1, 45}, position{"", 1, 45}},                    // :
			token{tokenIdentifier, "Gradient", position{"", 1, 47}, position{"", 1, 54}},       // Gradient
			token{tokenLeftBrace, "", position{"", 1, 56}, position{"", 1, 56}},                // {
			token{tokenIdentifier, "GradientStop", position{"", 1, 58}, position{"", 1, 69}},   // GradientStop
			token{tokenLeftBrace, "", position{"", 1, 71}, position{"", 1, 71}},                // {
			token{tokenIdentifier, "position", position{"", 1, 73}, position{"", 1, 80}},       // position
			token{tokenColon, "", position{"", 1, 81}, position{"", 1, 81}},                    // :
			token{tokenExpression, "0.0", position{"", 1, 83}, position{"", 1, 85}},            // 0.0
			token{tokenSemicolon, "", position{"", 1, 86}, position{"", 1, 86}},                // ;
			token{tokenIdentifier, "color", position{"", 1, 88}, position{"", 1, 92}},          // color
			token{tokenColon, "", position{"", 1, 93}, position{"", 1, 93}},                    // :
			token{tokenExpression, `"yellow" `, position{"", 1, 95}, position{"", 1, 103}},     // "yellow"
			token{tokenRightBrace, "", position{"", 1, 104}, position{"", 1, 104}},             // }
			token{tokenSemicolon, "", position{"", 1, 105}, position{"", 1, 105}},              // ;
			token{tokenIdentifier, "GradientStop", position{"", 1, 107}, position{"", 1, 118}}, // GradientStop
			token{tokenLeftBrace, "", position{"", 1, 120}, position{"", 1, 120}},              // {
			token{tokenIdentifier, "position", position{"", 1, 122}, position{"", 1, 129}},     // position
			token{tokenColon, "", position{"", 1, 130}, position{"", 1, 130}},                  // :
			token{tokenExpression, "1.0", position{"", 1, 132}, position{"", 1, 134}},          // 1.0
			token{tokenSemicolon, "", position{"", 1, 135}, position{"", 1, 135}},              // ;
			token{tokenIdentifier, "color", position{"", 1, 137}, position{"", 1, 141}},        // color
			token{tokenColon, "", position{"", 1, 142}, position{"", 1, 142}},                  // :
			token{tokenExpression, `"green" `, position{"", 1, 144}, position{"", 1, 151}},     // "green"
			token{tokenRightBrace, "", position{"", 1, 152}, position{"", 1, 152}},             // }
			token{tokenRightBrace, "", position{"", 1, 154}, position{"", 1, 154}},             // }
			token{tokenRightBrace, "", position{"", 1, 155}, position{"", 1, 155}},             // }
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
		output: token{tokenInteger, "1", position{"", 1, 1}, position{"", 1, 1}},
	}, {
		input:  "1.0",
		output: token{tokenFloat, "1.0", position{"", 1, 1}, position{"", 1, 3}},
	}, {
		input:  ".5",
		output: token{tokenFloat, ".5", position{"", 1, 1}, position{"", 1, 2}},
	}, {
		input:  "1.",
		output: token{tokenFloat, "1.", position{"", 1, 1}, position{"", 1, 2}},
	}, {
		input:  "1.0.3",
		output: token{tokenIdentifier, "1.0.3", position{"", 1, 1}, position{"", 1, 5}},
	}, {
		input:  "01.20",
		output: token{tokenFloat, "01.20", position{"", 1, 1}, position{"", 1, 5}},
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
		if !tok.start.IsEqual(expected.start) {
			return false, fmt.Sprintf("token %q starts at %v but %v was expected", tok.tokenType, tok.start, expected.start)
		}
		if !tok.end.IsEqual(expected.end) {
			return false, fmt.Sprintf("token %q ends at %v but %v was expected", tok.tokenType, tok.end, expected.end)
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
