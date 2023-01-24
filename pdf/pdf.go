package pdf

import (
	"fmt"

	vit "github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/parse"
)

//go:generate go build -o gencmd github.com/omniskop/vitrum/vit/generator/gencmd
//go:generate ./gencmd -i DocumentComponent.vit -o documentComponent_gen.go -p github.com/omniskop/vitrum/pdf
//go:generate ./gencmd -i PageComponent.vit -o pageComponent_gen.go -p github.com/omniskop/vitrum/pdf
//go:generate rm ./gencmd

func init() {
	parse.RegisterLibrary("PDF", PDFLib{})
}

type PDFLib struct{}

func (l PDFLib) ComponentNames() []string {
	return []string{"Document", "Page"}
}

func (l PDFLib) NewComponent(name string, id string, globalCtx *vit.GlobalContext) (vit.Component, bool) {
	var comp vit.Component
	var err error
	switch name {
	case "Document":
		comp, err = newDocumentComponentInGlobal(id, globalCtx, l)
	case "Page":
		comp, err = newPageComponentInGlobal(id, globalCtx, l)
	default:
		return nil, false
	}
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return comp, true
}

func (l PDFLib) StaticAttribute(componentName string, attributeName string) (interface{}, bool) {
	switch componentName {
	case "Page":
		return (*PageComponent)(nil).staticAttribute(attributeName)
	}
	return nil, false
}
