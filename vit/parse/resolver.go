package parse

import (
	"github.com/omniskop/vitrum/vit"
)

// DocumentInstantiator implements the vit.AbstractComponent interface for a vit document.
type DocumentInstantiator struct {
	doc VitDocument
}

var _ vit.AbstractComponent = (*DocumentInstantiator)(nil) // make sure that DocumentInstantiator implements the AbstractComponent interface

// Instantiate this component with the given id. The componentResolver will be used to resolve components that are needed in the instantiation.
func (i *DocumentInstantiator) Instantiate(id string, components vit.ComponentResolver) (vit.Component, error) {
	// TODO: use id
	comp, err := interpret(i.doc, components)
	if err != nil {
		return nil, err
	}
	return comp[0], nil
}

// ResolveVariable tries to find static attributes of the document's component.
// It implements the script.VariableResolver interface.
func (i *DocumentInstantiator) ResolveVariable(name string) (interface{}, bool) {
	for _, prop := range i.doc.components[0].properties {
		if !prop.static {
			continue
		}
		if len(prop.identifier) == 1 && prop.identifier[0] == name {
			return prop.staticValue, true
		}
	}

	for _, enum := range i.doc.components[0].enumerations {
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

func (i *DocumentInstantiator) Name() string {
	return i.doc.name
}

// LibraryInstantiator implements the vit.AbstractComponent interface for a specific component defined in a vit library.
type LibraryInstantiator struct {
	library       Library
	componentName string // name of a specific component in the library
}

var _ vit.AbstractComponent = (*LibraryInstantiator)(nil) // make sure that DocumentInstantiator implements the AbstractComponent interface

// Instantiate this component with the given id. The componentResolver will be used to resolve components that are needed in the instantiation.
func (i *LibraryInstantiator) Instantiate(id string, components vit.ComponentResolver) (vit.Component, error) {
	c, ok := i.library.NewComponent(i.componentName, id, components)
	if !ok {
		// if this happens the LibraryInstantiator was build incorrectly
		return nil, unknownComponentError{i.componentName}
	}
	return c, nil
}

// ResolveVariable tries to find static attributes of the libraries component.
func (i *LibraryInstantiator) ResolveVariable(name string) (interface{}, bool) {
	panic("LibraryInstantiator.ResolveVariable() not implemented")
	return nil, false
}

func (i *LibraryInstantiator) Name() string {
	return i.componentName
}
