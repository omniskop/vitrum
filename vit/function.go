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

func NewFunction(code string, position *PositionRange, fileCtx *FileContext) *Function {
	originalCode := code
	// convert code that is just enclosed in curly braces to a valid function
	if startsWith(code, '{') {
		code = fmt.Sprintf("function()%s", code)
		if position != nil {
			// adjust position in a way to hide the fact that we added code here
			p := position.StartColumnShifted(-11)
			position = &p
		}
	}

	// wrap code in parenthesis to make sure we handle the function like a value
	code = fmt.Sprintf("(%s)", code)
	if position != nil {
		p := position.StartColumnShifted(-1)
		position = &p
	}

	prog, err := script.NewScript("function", code)
	if err != nil {
		err = NewExpressionError(originalCode, position, err)
	}

	return &Function{
		fileCtx:  fileCtx,
		code:     originalCode,
		program:  prog,
		Position: position,
		err:      err,
	}
}

func NewFunctionFromCode(code Code) *Function {
	return NewFunction(code.Code, code.Position, code.FileCtx)
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

func NewAsyncFunction(code string, position *PositionRange, fileCtx *FileContext) *AsyncFunction {
	return &AsyncFunction{
		Function:  *NewFunction(code, position, fileCtx),
		arguments: nil,
		dirty:     false, // async functions start clean
	}
}

func NewAsyncFunctionFromCode(code Code) *AsyncFunction {
	return NewAsyncFunction(code.Code, code.Position, code.FileCtx)
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
