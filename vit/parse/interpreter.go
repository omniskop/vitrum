package parse

import (
	"errors"
	"fmt"
	"os"

	"github.com/omniskop/vitrum/vit"
)

type componentError struct {
	component vit.AbstractComponent
	err       error
}

func (e componentError) Error() string {
	var cErr componentError
	if errors.As(e.err, &cErr) {
		return fmt.Sprintf("%s > %s", e.component.Name(), e.err) // looks nicer
	}
	return fmt.Sprintf("%s: %s", e.component.Name(), e.err)
}

func (e componentError) Is(target error) bool {
	_, ok := target.(componentError)
	return ok
}

func (e componentError) Unwrap() error {
	return e.err
}

type unknownComponentError struct {
	name string
}

func (e unknownComponentError) Error() string {
	return fmt.Sprintf("unknown component %q", e.name)
}

func (e unknownComponentError) Is(target error) bool {
	_, ok := target.(unknownComponentError)
	return ok
}

type genericError struct {
	position vit.PositionRange
	err      error
}

func genericErrorf(position vit.PositionRange, format string, args ...interface{}) error {
	return genericError{position, fmt.Errorf(format, args...)}
}

func (e genericError) Error() string {
	return e.err.Error()
}

func (e genericError) Is(target error) bool {
	_, ok := target.(genericError)
	return ok
}

func (e genericError) Unwrap() error {
	return e.err
}

func init() {
	vit.InstantiateComponent = instantiateComponent
}

// parseFile parsed a given file into a document with the given component name.
func parseFile(fileName string, componentName string) (*VitDocument, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lexer := NewLexer(file, fileName)

	doc, err := Parse(NewTokenBuffer(lexer.Lex))
	if err != nil {
		return nil, err
	}
	doc.Name = componentName

	return doc, nil
}

// interpret takes the parsed document and creates the appropriate component tree.
// The returned error will always be of type ParseError
func interpret(document VitDocument, id string, components vit.ComponentContainer) ([]vit.Component, error) {
	for _, imp := range document.Imports {
		if len(imp.file) != 0 {
			// file import
			return nil, genericErrorf(imp.position, "not yet implemented")
		} else if len(imp.namespace) != 0 {
			// namespace import
			lib, err := resolveLibraryImport(imp.namespace)
			if err != nil {
				return nil, ParseError{imp.position, err}
			}
			for _, name := range lib.ComponentNames() {
				components.Set(name, &LibraryInstantiator{lib, name})
			}
		} else {
			return nil, genericErrorf(imp.position, "incomplete namespace")
		}
	}

	var instances []vit.Component
	for _, comp := range document.Components {
		instance, err := instantiateCustomComponent(comp, id, document.Name, components)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// instantiateCustomComponent creates a component described by a componentDefinition and wraps it in a Custom component with the given id.
func instantiateCustomComponent(def *vit.ComponentDefinition, id string, name string, components vit.ComponentContainer) (vit.Component, error) {
	comp, err := instantiateComponent(def, components)
	if err != nil {
		return nil, err
	}

	cst := vit.NewCustom(id, name, comp)

	return cst, nil
}

// instantiateComponent creates a component described by a componentDefinition.
func instantiateComponent(def *vit.ComponentDefinition, components vit.ComponentContainer) (vit.Component, error) {
	src, ok := components.Get(def.BaseName)
	if !ok {
		return nil, unknownComponentError{def.BaseName}
	}
	instance, err := src.Instantiate(def.ID, components.JustGlobal())
	if err != nil {
		return nil, componentError{src, err}
	}

	err = populateComponent(instance, def, components)
	if err != nil {
		return instance, componentError{src, err}
	}

	return instance, nil
}

// populateComponent takes a fresh component instance as well as it's definition and populates all attributes and children with their correct values.
func populateComponent(instance vit.Component, def *vit.ComponentDefinition, components vit.ComponentContainer) error {
	for _, enum := range def.Enumerations {
		if !instance.DefineEnum(enum) {
			return genericErrorf(enum.Position, "enum %q already defined", enum.Name)
		}
	}

	for _, prop := range def.Properties {
		if prop.VitType != "" {
			// this defines a new property
			if err := instance.DefineProperty(prop); err != nil {
				return err
			}
			// instance.SetProperty(prop.identifier[0], prop.expression)
		} else if len(prop.Identifier) == 1 {
			// simple property assignment
			if ok := instance.SetProperty(prop.Identifier[0], prop, &prop.Pos); !ok {
				return genericErrorf(prop.Pos, "unknown property %q of component %q", prop.Identifier[0], def.BaseName)
			}
		} else {
			// assign property with qualifier
			// TODO: make this universal
			if prop.Identifier[0] == "anchors" {
				// TODO: fix this
				// exp := vit.NewExpression(prop.expression, &prop.position)
				// a, _ := instance.Property("anchors")
				// a.GetValue().(*vit.Anchors).SetProperty(prop.identifier[1], exp)
			}
		}
	}

	for _, childDef := range def.Children {
		childInstance, err := instantiateComponent(childDef, components)
		if err != nil {
			return err
		}
		instance.AddChild(childInstance)
	}

	return nil
}
