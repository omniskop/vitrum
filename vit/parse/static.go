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
		for i := 0; i < len(doc.components); i++ {
			for c := 0; c < len(doc.components[i].properties); c++ {
				prop := &doc.components[i].properties[c]
				if !prop.static || prop.staticValue != nil {
					continue
				}
				val, err := script.RunContained(prop.expression, staticVariableResolver{documents, doc.components[i]})
				if err != nil {
					fmt.Println("1>", err)
					continue
				}
				prop.staticValue = val
			}
		}
		documents.Global[componentName] = &DocumentInstantiator{doc}
	}
}

type staticVariableResolver struct {
	documents vit.ComponentContainer
	component *componentDefinition
}

func (r staticVariableResolver) ResolveVariable(name string) (interface{}, bool) {
	if comp, ok := r.documents.Get(name); ok {
		return comp, true
	}

	for _, prop := range r.component.properties {
		if prop.static && len(prop.identifier) == 1 && prop.identifier[0] == name {
			if prop.staticValue != nil {
				return prop.staticValue, true
			}
			val, err := script.RunContained(prop.expression, r)
			if err != nil {
				fmt.Println("2>", err)
				return nil, false
			}
			prop.staticValue = val
			return val, true
		}
	}

	for _, enum := range r.component.enumerations {
		if enum.Name == name {
			// enum fullfills the script.VariableSource interface and can be returned
			return enum, true
		}
	}

	return nil, false
}
