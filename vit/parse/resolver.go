package parse

import (
	"fmt"

	"github.com/omniskop/vitrum/vit"
)

type Resolver interface {
	Resolve(names ...string) (interface{}, bool)
}

// DocumentVariableResolver will resolve variables that are global to a document.
// That can include components from the following sources:
// - Components inside an import folder
// - Components in the same folder as this document
// - Components that have been imported directly
type DocumentVariableResolver struct {
	variables map[string]interface{}
}

func NewDocumentVariableResolver(doc VitDocument) *DocumentVariableResolver {
	variables := make(map[string]interface{})

	for _, imp := range doc.imports {
		name := imp.identifier
		if name == "" {
			name = imp.namespace[len(imp.namespace)-1]
		}
		variables[name] = imp
	}

	if len(doc.components) == 1 {
		variables[doc.name] = doc.components[0]
	}

	return &DocumentVariableResolver{
		variables: variables,
	}
}

func (r *DocumentVariableResolver) Resolve(names ...string) (interface{}, bool) {
	if len(names) == 0 {
		return nil, false
	}
	if val, ok := r.variables[names[0]]; ok {
		if len(names) == 1 {
			return val, true
		}
		if res, ok := val.(Resolver); ok {
			return res.Resolve(names[1:]...)
		}
	}
	return nil, false
}

// DocumentInstantiator implements the AbstractComponent interface for instantiating vit files.
type DocumentInstantiator struct {
	doc VitDocument
}

var _ vit.AbstractComponent = (*DocumentInstantiator)(nil) // make sure that DocumentInstantiator implements the AbstractComponent interface

func (i *DocumentInstantiator) Instantiate(id string, components vit.ComponentResolver) (vit.Component, error) {
	// TODO: use id
	comp, err := Interpret(i.doc, components)
	if err != nil {
		return nil, err
	}
	return comp[0], nil
}

func (i *DocumentInstantiator) ResolveVariable(name string) (interface{}, bool) {
	enumerations := i.doc.components[0].enumerations
	for _, enum := range enumerations {
		if enum.Embedded {
			// If this enum is embedded we immediately search if it contains this variable. If it doesn't that's not an issue
			v, ok := enum.ResolveVariable(name)
			if ok {
				return v, true
			}
		}
		if enum.Name == name {
			return enum, true
		}
	}
	return nil, false
}

type LibraryInstantiator struct {
	library       Library
	componentName string
}

var _ vit.AbstractComponent = (*LibraryInstantiator)(nil) // make sure that DocumentInstantiator implements the AbstractComponent interface

func (i *LibraryInstantiator) Instantiate(id string, components vit.ComponentResolver) (vit.Component, error) {
	c, ok := i.library.NewComponent(i.componentName, id, components)
	if !ok {
		return nil, fmt.Errorf("unknown component %q", i.componentName)
	}
	return c, nil
}

func (i *LibraryInstantiator) ResolveVariable(name string) (interface{}, bool) {
	panic("LibraryInstantiator.ResolveVariable() not implemented")
	return nil, false
}
