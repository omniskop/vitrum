package parse

import (
	"fmt"

	"github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/script"
)

// evaluates all static expressions in all documents
func evaluateStaticExpressions(documents vit.ComponentContainer) {
	for componentName, abstract := range documents.Global {
		docInst, ok := abstract.(*DocumentInstantiator)
		if !ok {
			continue
		}
		doc := docInst.doc
		for i := 0; i < len(doc.Components); i++ {
			for c := 0; c < len(doc.Components[i].Properties); c++ {
				prop := &doc.Components[i].Properties[c]
				if !prop.Static || prop.StaticValue != nil {
					continue
				}
				val, err := script.RunContained(prop.Expression, staticVariableResolver{documents, doc.Components[i]})
				if err != nil {
					fmt.Println("1>", err)
					continue
				}
				prop.StaticValue = val
			}
		}
		documents.Global[componentName] = &DocumentInstantiator{doc}
	}
}

type staticVariableResolver struct {
	documents vit.ComponentContainer
	component *vit.ComponentDefinition
}

func (r staticVariableResolver) ResolveVariable(name string) (interface{}, bool) {
	if comp, ok := r.documents.Get(name); ok {
		return comp, true
	}

	for _, prop := range r.component.Properties {
		if prop.Static && len(prop.Identifier) == 1 && prop.Identifier[0] == name {
			if prop.StaticValue != nil {
				return prop.StaticValue, true
			}
			val, err := script.RunContained(prop.Expression, r)
			if err != nil {
				fmt.Println("2>", err)
				return nil, false
			}
			prop.StaticValue = val
			return val, true
		}
	}

	for _, enum := range r.component.Enumerations {
		if enum.Name == name {
			// enum fulfills the script.VariableSource interface and can be returned
			return enum, true
		}
	}

	return nil, false
}
