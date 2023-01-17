package std

import (
	"fmt"

	vit "github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/parse"
)

//go:generate go build -o gencmd github.com/omniskop/vitrum/vit/generator/gencmd
//go:generate ./gencmd -i Rectangle.vit -o rectangle_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Repeater.vit -o repeater_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Row.vit -o row_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Column.vit -o column_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Grid.vit -o grid_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Text.vit -o text_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i MouseArea.vit -o mouseArea_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i KeyArea.vit -o keyArea_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Rotation.vit -o rotation_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Image.vit -o image_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate rm ./gencmd

func init() {
	parse.RegisterLibrary("Vit", StdLib{})
}

type StdLib struct {
}

func (l StdLib) ComponentNames() []string {
	return []string{"Item", "Rectangle", "Repeater", "Container", "Row", "Column", "Grid", "Text", "MouseArea", "KeyArea", "Rotation", "Image"}
}

func (l StdLib) NewComponent(name string, id string, globalCtx *vit.GlobalContext) (vit.Component, bool) {
	var comp vit.Component
	var err error
	switch name {
	case "Item":
		var fileCtx = vit.NewFileContext(globalCtx)
		return NewItem(id, fileCtx), true
	case "Rectangle":
		comp, err = newRectangleInGlobal(id, globalCtx, l)
	case "Repeater":
		comp, err = newRepeaterInGlobal(id, globalCtx, l)
	case "Container":
		var fileCtx = vit.NewFileContext(globalCtx)
		return NewContainer(id, fileCtx), true
	case "Row":
		comp, err = newRowInGlobal(id, globalCtx, l)
	case "Column":
		comp, err = newColumnInGlobal(id, globalCtx, l)
	case "Grid":
		comp, err = newGridInGlobal(id, globalCtx, l)
	case "Text":
		comp, err = newTextInGlobal(id, globalCtx, l)
	case "MouseArea":
		comp, err = newMouseAreaInGlobal(id, globalCtx, l)
	case "KeyArea":
		comp, err = newKeyAreaInGlobal(id, globalCtx, l)
	case "Rotation":
		comp, err = newRotationInGlobal(id, globalCtx, l)
	case "Image":
		comp, err = newImageInGlobal(id, globalCtx, l)
	default:
		return nil, false
	}
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return comp, true
}

func (l StdLib) StaticAttribute(componentName string, attributeName string) (interface{}, bool) {
	switch componentName {
	case "Grid":
		return (*Grid)(nil).staticAttribute(attributeName)
	case "Text":
		return (*Text)(nil).staticAttribute(attributeName)
	case "MouseArea":
		return (*MouseArea)(nil).staticAttribute(attributeName)
	case "Rotation":
		return (*Rotation)(nil).staticAttribute(attributeName)
	case "Image":
		return (*Image)(nil).staticAttribute(attributeName)
	}
	return nil, false
}
