package gui

import (
	"fmt"

	vit "github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/parse"
)

//go:generate go build -o gencmd github.com/omniskop/vitrum/vit/generator/gencmd
//go:generate ./gencmd -i WindowComponent.vit -o windowComponent_gen.go -p github.com/omniskop/vitrum/gui
//go:generate rm ./gencmd

func init() {
	parse.RegisterLibrary("GUI", GUILib{})
}

type GUILib struct{}

func (l GUILib) ComponentNames() []string {
	return []string{"Window"}
}

func (l GUILib) NewComponent(name string, id string, globalCtx *vit.GlobalContext) (vit.Component, bool) {
	var comp vit.Component
	var err error
	switch name {
	case "Window":
		comp, err = newWindowComponentInGlobal(id, globalCtx)
	default:
		return nil, false
	}
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return comp, true
}

func (l GUILib) StaticAttribute(componentName string, attributeName string) (interface{}, bool) {
	switch componentName {

	}
	return nil, false
}
