package gui

import vit "github.com/omniskop/vitrum/vit"

//go:generate go build -o gencmd github.com/omniskop/vitrum/vit/generator/gencmd
//go:generate ./gencmd -i Window.vit -o window_gen.go -p github.com/omniskop/vitrum/gui
//go:generate rm ./gencmd

type GUILib struct {
}

func (l GUILib) ComponentNames() []string {
	return []string{"Window"}
}

func (l GUILib) NewComponent(name string, id string, scope vit.ComponentContainer) (vit.Component, bool) {
	switch name {
	case "Window":
		return NewWindow(id, scope), true
	}
	return nil, false
}

func (l GUILib) StaticAttribute(componentName string, attributeName string) (interface{}, bool) {
	switch componentName {

	}
	return nil, false
}
