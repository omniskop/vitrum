package script

import (
	"errors"
	"fmt"
	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
)

var runtime *goja.Runtime
var runtimeMux sync.Mutex

func init() {
	Setup()
}

func Setup() {
	runtime = goja.New()
	runtime.SetParserOptions(parser.WithDisableSourceMaps)
	runtime.Set("Vit", builtinFunctions)
	runtime.Set("console", globalConsole)
}

type Script struct {
	compiled *goja.Program
}

func NewScript(name string, src string) (Script, error) {
	comp, err := goja.Compile(name, src, true)
	if err != nil {
		return Script{}, err
	}
	return Script{
		compiled: comp,
	}, nil
}

func (s *Script) Run(variables VariableSource) (interface{}, error) {
	runtimeMux.Lock()
	defer runtimeMux.Unlock()

	bridgeObj := runtime.NewDynamicObject(&VariableBridge{variables})
	// TODO: figure out if this has to be deleted from the runtime manually afterwards
	global := runtime.GlobalObject()
	global.SetPrototype(bridgeObj)
	defer global.SetPrototype(nil)
	val, err := runtime.RunProgram(s.compiled)
	if err != nil {
		if exception, ok := err.(*goja.Exception); ok {
			return nil, errors.New(exception.Value().String())
		}
		return nil, err
	}
	return val.Export(), nil
}

// Call executes the script in the assumption that it defines a function.
// It will then call this function using the provided arguments.
// The value returned by that function will in turn be returned by this method.
func (s *Script) Call(variables VariableSource, arguments ...interface{}) (interface{}, error) {
	runtimeMux.Lock()
	defer runtimeMux.Unlock()

	bridgeObj := runtime.NewDynamicObject(&VariableBridge{variables})
	// TODO: figure out if this has to be deleted from the runtime manually afterwards
	global := runtime.GlobalObject()
	global.SetPrototype(bridgeObj)
	defer global.SetPrototype(nil)
	val, err := runtime.RunProgram(s.compiled)
	if err != nil {
		if exception, ok := err.(*goja.Exception); ok {
			return nil, errors.New(exception.Value().String())
		}
		return nil, err
	}

	f, ok := goja.AssertFunction(val)
	if !ok {
		return nil, errors.New("not a function")
	}
	var argumentValues = make([]goja.Value, len(arguments))
	for i, arg := range arguments {
		// TODO: figure out if these values have to be deleted form the runtime manually afterwards
		argumentValues[i] = runtime.ToValue(arg)
	}
	result, err := f(goja.Undefined(), argumentValues...)
	if err != nil {
		return nil, err
	}
	return result.Export(), nil
}

// Run executes the given code with the provided variable context.
// The resulting value and a potential error is returned.
func Run(code string, variables VariableSource) (interface{}, error) {
	runtimeMux.Lock()
	defer runtimeMux.Unlock()

	bridgeObj := runtime.NewDynamicObject(&VariableBridge{variables})
	// TODO: figure out if this has to be deleted from the runtime manually afterwards
	global := runtime.GlobalObject()
	global.SetPrototype(bridgeObj)
	defer global.SetPrototype(nil)
	val, err := runtime.RunString(code)
	if err != nil {
		return nil, err
	}
	return val.Export(), nil
}

// RunContainer executes the given code in a separate runtime with the provided variable context.
// It behaves like Run but it supports multi threading by using it's own separate runtime for each call.
func RunContained(code string, variables VariableSource) (interface{}, error) {
	// NOTE: If the recreation of the runtime should at some point become a speed issue, it could maybe be consideret to put the runtime in a pool.
	bridgeObj := runtime.NewDynamicObject(&VariableBridge{variables})
	runtime := goja.New()
	runtime.SetParserOptions(parser.WithDisableSourceMaps)
	runtime.Set("Vit", builtinFunctions)
	runtime.Set("console", globalConsole)
	global := runtime.GlobalObject()
	global.SetPrototype(bridgeObj)
	val, err := runtime.RunString(code)
	if err != nil {
		return nil, err
	}
	return val.Export(), nil
}

func Exception(msg string) goja.Value {
	return runtime.ToValue(msg)
}

type VariableSource interface {
	ResolveVariable(string) (interface{}, bool)
}

type Variable interface {
	GetValue() interface{}
	SetValue(interface{}) error
}

type VariableBridge struct {
	Source VariableSource
}

func (b *VariableBridge) Get(key string) goja.Value {
	// fmt.Printf("[VariableBridge] get %q\n", key)
	val, ok := b.Source.ResolveVariable(key)
	if !ok {
		// fmt.Printf("[VariableBridge] get %q: undefined\n", key)
		// returning undefined here would be a better fit for JavaScript, but I think failing here will give better error messages
		// return goja.Undefined()
		if key == "id" {
			panic(Exception(fmt.Sprintf("id can't be used in an expression")))
		}
		panic(Exception(fmt.Sprintf("undefined variable %q", key)))
	}
	switch actual := val.(type) {
	case VariableBridge:
		// TODO: is this case still needed? Can this still happen?
		// fmt.Printf("[VariableBridge] get %q: (abstract) component\n", key)
		return runtime.NewDynamicObject(&actual)
	case VariableSource:
		// fmt.Printf("[VariableBridge] get %q: dynamic object\n", key)
		return runtime.NewDynamicObject(&VariableBridge{actual})
	case Variable:
		return runtime.ToValue(actual.GetValue())
	}
	return runtime.ToValue(val)
}

func (b *VariableBridge) Set(key string, value goja.Value) bool {
	// fmt.Printf("[VariableBridge] set %q: %v (%T)\n", key, value, value)
	val, ok := b.Source.ResolveVariable(key)
	if !ok {
		panic(Exception(fmt.Sprintf("undefined variable %q", key)))
	}
	variable, ok := val.(Variable)
	if !ok {
		panic(Exception(fmt.Sprintf("variable %q is not settable", key)))
	}
	err := variable.SetValue(value.Export())
	if err != nil {
		panic(runtime.NewTypeError(fmt.Sprintf("cannot assign %q to %q: %s", value.Export(), key, err)))
	}
	return true
}

func (b *VariableBridge) Has(key string) bool {
	// fmt.Println("Has", key)
	_, ok := b.Source.ResolveVariable(key)
	return ok
}

func (b *VariableBridge) Delete(key string) bool {
	// fmt.Println("Delete", key)
	return false
}

func (b *VariableBridge) Keys() []string {
	// fmt.Println("Keys")
	return nil
}
