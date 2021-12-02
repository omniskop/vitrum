package parse

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

type Library interface {
	ComponentNames() []string
	NewComponent(string, string) (vit.Component, bool)
}

func ParseFile(fileName string) (*VitDocument, error) {
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

	return doc, nil
}

func DoMagic(fileName string) (vit.Component, error) {
	doc, err := ParseFile(fileName)
	if err != nil {
		return nil, err
	}

	fmt.Println("main file parsed...")

	entries, err := os.ReadDir(path.Dir(fileName))
	if err != nil {
		return nil, err
	}
	var documents = make(map[string]VitDocument)
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".vit") {
			doc, err := ParseFile(path.Join(path.Dir(fileName), entry.Name()))
			if err != nil {
				return nil, err
			}
			documents[strings.TrimSuffix(entry.Name(), ".vit")] = *doc
		}
	}

	fmt.Println("now interpreting...")

	components, err := Interpret(*doc, documents)
	if err != nil {
		return nil, err
	}

	return components[0], nil
}

func Interpret(document VitDocument, neighbors map[string]VitDocument) ([]vit.Component, error) {
	componentIndex := make(map[string]Library)

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
				componentIndex[name] = lib
			}
		} else {
			return nil, fmt.Errorf("incomplete namespace")
		}
	}

	for name, doc := range neighbors {
		componentIndex[name] = &standaloneDocument{name, doc, neighbors}
	}

	var instances []vit.Component

	for _, comp := range document.components {
		instance, err := instantiateComponent(componentIndex, comp)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	fmt.Println("components constructured")
	fmt.Println("evaluating expressions...")

evaluateExpressions:
	n, err := CheckForReevaluation(instances)
	if err != nil {
		return nil, err
	}
	if n > 0 {
		fmt.Printf("evaluated %d expressions\n", n)
		goto evaluateExpressions
	}

	return instances, nil
}

func instantiateComponent(componentIndex map[string]Library, def *componentDefinition) (vit.Component, error) {
	lib, ok := componentIndex[def.name]
	if !ok {
		return nil, fmt.Errorf("unknown component %q", def.name)
	}
	instance, ok := lib.NewComponent(def.name, def.id)
	if !ok {
		return nil, fmt.Errorf("component %q could not be instantiated", def.name)
	}

	err := populateComponent(componentIndex, instance, def)
	if err != nil {
		return instance, err
	}

	return instance, nil
}

func populateComponent(componentIndex map[string]Library, instance vit.Component, def *componentDefinition) error {
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
		childInstance, err := instantiateComponent(componentIndex, childDef)
		if err != nil {
			return err
		}
		instance.AddChild(childInstance)
	}

	return nil
}

func resolveLibraryImport(namespace []string) (Library, error) {
	if len(namespace) == 0 {
		return nil, fmt.Errorf("empty namespace")
	}
	switch namespace[0] {
	case "Vit":
		if len(namespace) == 1 {
			return vit.StdLib{}, nil
		}
	}

	return nil, fmt.Errorf("unknown library %q", strings.Join(namespace, "."))
}

func CheckForReevaluation(components []vit.Component) (int, error) {
	var sum int
	for _, c := range components {
		n, err := c.UpdateExpressions()
		sum += n
		if err != nil {
			return sum, err
		}
	}
	return sum, nil
}

type standaloneDocument struct {
	name      string
	doc       VitDocument
	neighbors map[string]VitDocument
}

func (d *standaloneDocument) ComponentNames() []string {
	return []string{d.name}
}

func (d *standaloneDocument) NewComponent(name, id string) (vit.Component, bool) {
	if name != d.name {
		return nil, false
	}

	components, err := Interpret(d.doc, d.neighbors)
	if err != nil {
		fmt.Printf("standalone document: %v\n", err)
		return nil, false
	}
	return components[0], true
}
