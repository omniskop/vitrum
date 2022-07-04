package std

import vit "github.com/omniskop/vitrum/vit"

//go:generate go build -o gencmd github.com/omniskop/vitrum/vit/generator/gencmd
//go:generate ./gencmd -i Rectangle.vit -o rectangle_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Repeater.vit -o repeater_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Row.vit -o row_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Column.vit -o column_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Grid.vit -o grid_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Text.vit -o text_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i MouseArea.vit -o mouseArea_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate rm ./gencmd

type StdLib struct {
}

func (l StdLib) ComponentNames() []string {
	return []string{"Item", "Rectangle", "Repeater", "Container", "Row", "Column", "Grid", "Text", "MouseArea"}
}

func (l StdLib) NewComponent(name string, id string, scope vit.ComponentContext) (vit.Component, bool) {
	switch name {
	case "Item":
		return NewItem(id, scope), true
	case "Rectangle":
		return NewRectangle(id, scope), true
	case "Repeater":
		return NewRepeater(id, scope), true
	case "Container":
		return NewContainer(id, scope), true
	case "Row":
		return NewRow(id, scope), true
	case "Column":
		return NewColumn(id, scope), true
	case "Grid":
		return NewGrid(id, scope), true
	case "Text":
		return NewText(id, scope), true
	case "MouseArea":
		return NewMouseArea(id, scope), true
	}
	return nil, false
}

func (l StdLib) StaticAttribute(componentName string, attributeName string) (interface{}, bool) {
	switch componentName {
	case "Grid":
		return (*Grid)(nil).staticAttribute(attributeName)
	case "Text":
		return (*Text)(nil).staticAttribute(attributeName)
	case "MouseArea":
		return (*MouseArea)(nil).staticAttribute(attributeName)
	}
	return nil, false
}
