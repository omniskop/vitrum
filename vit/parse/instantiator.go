package parse

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
	if e.err == nil {
		return "no error"
	}
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
	vit.InstantiateComponent = InstantiateComponent
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
func interpret(document VitDocument, id string, globalCtx *vit.GlobalContext) ([]vit.Component, error) {
	fileCtx := vit.NewFileContext(globalCtx)
	for _, imp := range document.Imports {
		if len(imp.File) != 0 {
			// file import
			return nil, genericErrorf(imp.Position, "not yet implemented")
		} else if len(imp.Namespace) != 0 {
			// namespace import
			lib, err := ResolveLibrary(imp.Namespace)
			if err != nil {
				return nil, ParseError{imp.Position, err}
			}
			for _, name := range lib.ComponentNames() {
				fileCtx.KnownComponents.Set(name, &LibraryInstantiator{lib, name})
			}
		} else {
			return nil, genericErrorf(imp.Position, "incomplete namespace")
		}
	}

	var instances []vit.Component
	for _, comp := range document.Components {
		instance, err := instantiateCustomComponent(comp, id, document.Name, fileCtx)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// instantiateCustomComponent creates a component described by a componentDefinition and wraps it in a Custom component with the given id.
func instantiateCustomComponent(def *vit.ComponentDefinition, id string, name string, fileCtx *vit.FileContext) (vit.Component, error) {
	comp, err := InstantiateComponent(def, fileCtx)
	if err != nil {
		return nil, err
	}

	cst := vit.NewCustom(id, name, comp)

	return cst, nil
}

// InstantiateComponent creates a component described by a componentDefinition.
func InstantiateComponent(def *vit.ComponentDefinition, fileCtx *vit.FileContext) (vit.Component, error) {
	src, ok := fileCtx.Get(def.BaseName)
	if !ok {
		// TODO: improve context for error; either here or upstream
		return nil, unknownComponentError{def.BaseName}
	}
	instance, err := src.Instantiate(def.ID, fileCtx.Global)
	if err != nil {
		return nil, componentError{src, err}
	}

	err = populateComponent(instance, def, fileCtx)
	if err != nil {
		return instance, componentError{src, err}
	}

	fileCtx.RegisterComponent(def.ID, instance)
	// TODO: figure out where the components will be unregistered again

	return instance, nil
}

// populateComponent takes a fresh component instance as well as it's definition and populates all attributes and children with their correct values.
func populateComponent(instance vit.Component, def *vit.ComponentDefinition, fileCtx *vit.FileContext) error {
	for _, enum := range def.Enumerations {
		if !instance.DefineEnum(enum) {
			return genericErrorf(*enum.Position, "enum %q already defined", enum.Name)
		}
	}

	for _, method := range def.Methods {
		if !instance.DefineMethod(method.CopyInContext(fileCtx)) {
			return genericErrorf(*method.AsyncFunction.Position, "method %q already defined", method.Name)
		}
	}

	for _, prop := range def.Properties {
		if prop.VitType != "" {
			// this defines a new property
			if err := instance.DefineProperty(prop, fileCtx); err != nil {
				return err
			}
			// instance.SetProperty(prop.identifier[0], prop.expression)
		} else if len(prop.Identifier) == 1 {
			// simple property assignment
			var err error
			if len(prop.Components) == 0 {
				err = instance.SetPropertyCode(prop.Identifier[0], vit.Code{Code: prop.Expression, Position: prop.ValuePos, FileCtx: fileCtx})
				if err != nil {
					// if this property doesn't exist, check if it's an event
					if ev, ok := instance.Event(prop.Identifier[0]); ok {
						// register event listener
						l := ev.CreateListener(vit.Code{Code: prop.Expression, Position: prop.ValuePos, FileCtx: fileCtx})
						instance.RootC().AddListenerFunction(l)
						err = nil
					}
				}
			} else if len(prop.Components) == 1 {
				err = instance.SetProperty(prop.Identifier[0], vit.ComponentDefinitionInContext{prop.Components[0], fileCtx})
			} else {
				err = instance.SetProperty(prop.Identifier[0], vit.ComponentDefinitionListInContext{prop.Components, fileCtx})
			}
			if err != nil {
				return genericError{prop.Pos, err}
			}
		} else {
			// assign property with qualifier
			// TODO: make this universal?
			if prop.Identifier[0] == "anchors" {
				v, ok := instance.Property(prop.Identifier[0])
				if !ok {
					return genericErrorf(prop.Pos, "unknown property %q of component %q", prop.Identifier[0], def.BaseName)
				}
				anchors, ok := v.(*vit.AnchorsValue)
				if !ok {
					return genericErrorf(prop.Pos, "cannot assign to non group-property %q of component %q", prop.Identifier[0], def.BaseName)
				}

				ok = anchors.SetPropertyCode(prop.Identifier[1], vit.Code{Code: prop.Expression, Position: prop.ValuePos, FileCtx: fileCtx})
				if !ok {
					return genericErrorf(prop.Pos, "unknown property %q of component %q", strings.Join(prop.Identifier, "."), def.BaseName)
				}
			} else {
				v, ok := instance.Property(prop.Identifier[0])
				if !ok {
					return genericErrorf(prop.Pos, "unknown property %q of component %q", prop.Identifier[0], def.BaseName)
				}
				switch v := v.(type) {
				case *vit.GroupValue:
					// set property of group value
					err := v.SetCodeOf(prop.Identifier[1], vit.Code{Code: prop.Expression, Position: prop.ValuePos, FileCtx: fileCtx})
					if err != nil {
						return genericErrorf(prop.Pos, "group-property %q of component %q: %w", prop.Identifier[0], def.BaseName, err)
					}
				case *vit.ComponentRefValue:
					// add listener to an event of another component
					event, ok := v.Component().RootC().Event(prop.Identifier[1])
					if !ok {
						return genericErrorf(prop.Pos, "unknown event %q of component %q", prop.Identifier[1], prop.Identifier[0])
					}
					l := event.CreateListener(vit.Code{Code: prop.Expression, Position: prop.ValuePos, FileCtx: fileCtx})
					instance.RootC().AddListenerFunction(l)
				default:
					return genericErrorf(prop.Pos, "cannot assign %q on property %q of component %q", prop.Identifier[1], prop.Identifier[0], def.BaseName)
				}
			}
		}
	}

	for _, childDef := range def.Children {
		childInstance, err := InstantiateComponent(childDef, fileCtx)
		if err != nil {
			return err
		}
		instance.AddChild(childInstance)
	}

	return nil
}
