package controls

import (
	"fmt"

	"github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/parse"
)

//go:generate go build -o gencmd github.com/omniskop/vitrum/vit/generator/gencmd
//go:generate ./gencmd -i Button.vit -o button_gen.go -p github.com/omniskop/vitrum/controls
//go:generate ./gencmd -i TextField.vit -o textField_gen.go -p github.com/omniskop/vitrum/controls
//go:generate rm ./gencmd

func init() {
	parse.RegisterLibrary("Controls", ControlsLib{})
}

type ControlsLib struct {
}

func (l ControlsLib) ComponentNames() []string {
	return []string{"Button", "TextField"}
}

func (l ControlsLib) NewComponent(name string, id string, globalCtx *vit.GlobalContext) (vit.Component, bool) {
	var comp vit.Component
	var err error
	switch name {
	case "Button":
		comp, err = newButtonInGlobal(id, globalCtx, l)
	case "TextField":
		comp, err = newTextFieldInGlobal(id, globalCtx, l)
	default:
		return nil, false
	}
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return comp, true
}

func (l ControlsLib) StaticAttribute(componentName string, attributeName string) (interface{}, bool) {
	switch componentName {
	}
	return nil, false
}
