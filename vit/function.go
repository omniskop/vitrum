package vit

import (
	"fmt"
	"unicode/utf8"

	"github.com/omniskop/vitrum/vit/script"
)

type Function struct {
	fileCtx  *FileContext
	code     string
	program  script.Script
	Position *PositionRange
	err      error
}

func NewFunction(code Code) *Function {
	originalCode := code.Code
	// convert code that is just enclosed in curly braces to a valid function
	if startsWith(code.Code, '{') {
		code.Code = fmt.Sprintf("function()%s", code.Code)
		if code.Position != nil {
			// adjust position in a way to hide the fact that we added code here
			p := code.Position.StartColumnShifted(-10)
			code.Position = &p
		}
	}

	// wrap code in parenthesis to make sure we handle the function like a value
	code.Code = fmt.Sprintf("(%s)", code.Code)
	if code.Position != nil {
		p := code.Position.StartColumnShifted(-1)
		code.Position = &p
	}

	prog, err := script.NewScript("function", code.Code)
	if err != nil {
		err = NewExpressionError(originalCode, code.Position, err)
	}

	return &Function{
		fileCtx:  code.FileCtx,
		code:     originalCode,
		program:  prog,
		Position: code.Position,
		err:      err,
	}
}

func (f *Function) Call(context Component, args ...interface{}) (interface{}, error) {
	if f.err != nil {
		return nil, f.err
	}

	collector := NewAccessCollector(f.fileCtx, context)
	val, err := f.program.Call(collector, args...)
	if err != nil {
		return nil, NewExpressionError(f.code, f.Position, err)
	}
	return val, nil
}

func (f *Function) Code() string {
	return f.code
}

type AsyncFunction struct {
	Function
	arguments []interface{}
	dirty     bool
}

func NewAsyncFunction(code Code) *AsyncFunction {
	return &AsyncFunction{
		Function:  *NewFunction(code),
		arguments: nil,
		dirty:     false, // async functions start clean
	}
}

func (f *AsyncFunction) Notify(args ...interface{}) {
	f.dirty = true
	f.arguments = args
}

func (f *AsyncFunction) Evaluate(context Component) (interface{}, error) {
	f.dirty = false
	return f.Call(context, f.arguments...)
}

func (f *AsyncFunction) ShouldEvaluate() bool {
	return f.dirty
}

func startsWith(str string, r rune) bool {
	read, _ := utf8.DecodeRuneInString(str)
	return r == read
}

func endsWith(str string, r rune) bool {
	read, _ := utf8.DecodeLastRuneInString(str)
	return r == read
}
