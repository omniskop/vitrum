package parse

import (
	"fmt"
	"os"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

type Library interface {
	ComponentNames() []string
	NewComponent(string, string, vit.ComponentResolver) (vit.Component, bool)
}

func ParseFile(fileName string, componentName string) (*VitDocument, error) {
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
	doc.name = componentName

	return doc, nil
}

// Interpret takes the parsed document and creates the appropriate component tree
func Interpret(document VitDocument, components vit.ComponentResolver) ([]vit.Component, error) {
	allComponents := vit.NewComponentResolver(&components)

	for _, imp := range document.imports {
		if len(imp.file) != 0 {
			// file import
			return nil, fmt.Errorf("not yet implemented")
		} else if len(imp.namespace) != 0 {
			// namespace import
			lib, err := resolveLibraryImport(imp.namespace)
			if err != nil {
				return nil, err
			}
			for _, name := range lib.ComponentNames() {
				allComponents.Components[name] = &LibraryInstantiator{lib, name}
			}
		} else {
			return nil, fmt.Errorf("incomplete namespace")
		}
	}

	var instances []vit.Component
	for _, comp := range document.components {
		instance, err := instantiateComponent(comp, allComponents)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	fmt.Println("components constructured")
	fmt.Println("evaluating expressions...")

	return instances, nil
}

// instantiateComponent creates a component described by a componentDefinition.
func instantiateComponent(def *componentDefinition, components vit.ComponentResolver) (vit.Component, error) {
	src, ok := components.Resolve(def.name)
	if !ok {
		return nil, fmt.Errorf("unknown component %q", def.name)
	}
	instance, err := src.Instantiate(def.id, components)
	if err != nil {
		return nil, fmt.Errorf("component %q could not be instantiated: %v", def.name, err)
	}

	err = populateComponent(instance, def, components)
	if err != nil {
		return instance, err
	}

	return instance, nil
}

// populateComponent takes a fresh component instance as well as it's definition and populates all attributes and children with their correct values.
func populateComponent(instance vit.Component, def *componentDefinition, components vit.ComponentResolver) error {
	for _, enum := range def.enumerations {
		if !instance.DefineEnum(enum) {
			return fmt.Errorf("enum %q already defined", enum.Name)
		}
	}

	for _, prop := range def.properties {
		exp := vit.NewExpression(prop.expression)
		if prop.vitType != "" {
			// this defines a new property
			if ok := instance.DefineProperty(prop.identifier[0], prop.vitType, prop.expression); !ok {
				return fmt.Errorf("property %q is already defined", prop.identifier[0])
			}
			// instance.SetProperty(prop.identifier[0], prop.expression)
		} else if len(prop.identifier) == 1 {
			// simple property assignment
			if ok := instance.SetProperty(prop.identifier[0], prop.expression); !ok {
				return fmt.Errorf("unknown property %q of component %q", prop.identifier[0], def.name)
			}
		} else {
			// assign property with qualifier
			if prop.identifier[0] == "anchors" {
				a, _ := instance.Property("anchors")
				a.(*vit.Anchors).SetProperty(prop.identifier[1], exp)
			}
		}
	}

	for _, childDef := range def.children {
		childInstance, err := instantiateComponent(childDef, components)
		if err != nil {
			return err
		}
		instance.AddChild(childInstance)
	}

	return nil
}

// resolveLibraryImport takes a library identifier and returns the corresponding library or an error if the identifier is unknown.
func resolveLibraryImport(namespace []string) (Library, error) {
	if len(namespace) == 0 {
		return nil, fmt.Errorf("empty namespace")
	}
	switch namespace[0] {
	case "Vit":
		if len(namespace) == 1 {
			return vit.StdLib{}, nil
		}
	case "QtQuick":
		return vit.StdLib{}, nil
	case "Dark":
		return vit.StdLib{}, nil
	}

	return nil, fmt.Errorf("unknown library %q", strings.Join(namespace, "."))
}

type standaloneDocument struct {
	name       string
	doc        VitDocument
	components vit.ComponentResolver
}

func (d *standaloneDocument) ComponentNames() []string {
	return []string{d.name}
}

func (d *standaloneDocument) NewComponent(name, id string) (vit.Component, bool) {
	if name != d.name {
		return nil, false
	}

	components, err := Interpret(d.doc, d.components)
	if err != nil {
		fmt.Printf("standalone document: %v\n", err)
		return nil, false
	}
	return components[0], true
}
