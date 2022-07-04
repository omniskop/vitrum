package gui

import (
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

func (l GUILib) NewComponent(name string, id string, scope vit.ComponentContext) (vit.Component, bool) {
	switch name {
	case "Window":
		return NewWindowComponent(id, scope), true
	}
	return nil, false
}

func (l GUILib) StaticAttribute(componentName string, attributeName string) (interface{}, bool) {
	switch componentName {

	}
	return nil, false
}
