package script

import (
	"fmt"
	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
)

var runtime *goja.Runtime
var runtimeMux sync.Mutex

func Setup() {
	runtime = goja.New()
	runtime.SetParserOptions(parser.WithDisableSourceMaps)
	runtime.Set("Vit", builtinFunctions)
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
		return nil, err
	}
	return val.Export(), nil
}

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

type Variable struct {
	Identifier []string
	Value      interface{}
}

type VariableSource interface {
	ResolveVariable(string) (interface{}, bool)
}

type VariableBridge struct {
	Source VariableSource
}

func (b *VariableBridge) Get(key string) goja.Value {
	val, ok := b.Source.ResolveVariable(key)
	if !ok {
		fmt.Printf("[VariableBridge] get %q: undefined\n", key)
		return goja.Undefined()
	}
	if subBridge, ok := val.(VariableBridge); ok {
		fmt.Printf("[VariableBridge] get %q: (abstract) component\n", key)
		return runtime.NewDynamicObject(&subBridge)
	}
	fmt.Printf("[VariableBridge] get %q: (%T) %v\n", key, val, val)
	return runtime.ToValue(val)
}

func (b *VariableBridge) Set(key string, value goja.Value) bool {
	fmt.Printf("[VariableBridge] set %q: %v (%T)", key, value, value)
	return true
}

func (b *VariableBridge) Has(key string) bool {
	fmt.Println("Has", key)
	_, ok := b.Source.ResolveVariable(key)
	return ok
}

func (b *VariableBridge) Delete(key string) bool {
	fmt.Println("Delete", key)
	return false
}

func (b *VariableBridge) Keys() []string {
	fmt.Println("Keys")
	return nil
}
