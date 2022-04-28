package parse

import (
	"fmt"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

// VitDocument contains everything there is to know about a parsed vit file
type VitDocument struct {
	Name       string            // Name of the file without extension. Usually the name of the component this file describes.
	Imports    []importStatement // all imported libraries and files
	Components []*ComponentDefinition
}

// String creates a human readable string representation of the vit document
func (d VitDocument) String() string {
	var out strings.Builder
	out.WriteString("{\r\n\tImports:")
	for _, imp := range d.Imports {
		out.WriteString(fmt.Sprintf(" %v", imp))
	}

	out.WriteString("\r\n\tComponents: \r\n")

	for _, comp := range d.Components {
		out.WriteString(fmt.Sprintf("\t\t%v\r\n", comp))
	}

	out.WriteString("}")

	return out.String()
}

// An importStatement can either import a module/namespace or a file
// namspaces have a version with major and minor part.
// Either namespace or file can be set, but not both.
type importStatement struct {
	namespace []string // fully qualified name of the module to import
	file      string   // file path that should be imported
	version   string   // version string for namespace imports
	qualifier string   // optional qualifier that allows the user to refer to the import by a different name
	position  vit.PositionRange
}

// String returns a human readable multiline string representation of the import
func (s importStatement) String() string {
	var out strings.Builder
	if len(s.namespace) == 0 {
		out.WriteString(fmt.Sprintf("%q", s.file))
	} else {
		out.WriteString(strings.Join(s.namespace, "."))
	}
	out.WriteRune('@')

	out.WriteString(s.version)

	return out.String()
}

// ComponentDefinition contains everything about a component that is defined in a vit file
type ComponentDefinition struct {
	position     vit.PositionRange
	BaseName     string                 // name of the instantiated base component
	ID           string                 // custom id of the component
	Properties   []Property             // all explicitly defined or declared properties
	Children     []*ComponentDefinition // child components
	Enumerations []vit.Enumeration      // all explicitly defined enumerations
}

// identifierIsKnown returns true of the given identifier has already been defined
func (d *ComponentDefinition) identifierIsKnown(identifier []string) bool {
	if len(identifier) == 0 {
		return false
	}
	for _, p := range d.Properties {
		if stringSlicesEqual(p.Identifier, identifier) {
			return true
		}
	}
	for _, e := range d.Enumerations {
		if e.Name == identifier[0] {
			return true
		}
	}

	return false
}

// String returns a human readable string representation of the component definition
func (d *ComponentDefinition) String() string {
	if d == nil {
		return "<nil>"
	}
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s { ", d.BaseName))

	for _, p := range d.Properties {
		s.WriteString(fmt.Sprintf("%s, ", p.String()))
	}

	if len(d.Children) == 0 {

	} else if len(d.Children) == 1 {
		s.WriteString(fmt.Sprintf("%d child", len(d.Children)))
	} else {
		s.WriteString(fmt.Sprintf("%d children", len(d.Children)))
	}

	// TODO: add other fields

	s.WriteString(" }")

	return s.String()
}

type ComponentInstantiator struct {
	definition *ComponentDefinition
	context    vit.ComponentContainer
}

var _ vit.Instantiator = (*ComponentInstantiator)(nil)

func (i ComponentInstantiator) Instantiate() (vit.Component, error) {
	return instantiateComponent(i.definition, i.context)
}

// Property contains everything about a defined or declared Property
type Property struct {
	position    vit.PositionRange    // position of the property declaration
	Identifier  []string             // Identifier of this property. This will usually be only one value but can contain multiple parts for example with 'Anchors.fill'
	VitType     string               // data type of the property in vit terms, not go
	Expression  string               // Expression string that defines the property. Can be empty.
	Component   *ComponentDefinition // Only set if this properties type is a component.
	ReadOnly    bool                 // Readonly properties are statically defined on the component itself and cannot be changed directly. They will however be recalculated if one of the expressions dependencies should change.
	Static      bool                 // Static properties are defined on the component itself. They will only be evaluated once when the component is loaded and are constant from that point on.
	staticValue interface{}          // The evaluated value of a static property.
}

// String returns a human readable string representation of the property
func (p Property) String() string {
	ident := strings.Join(p.Identifier, ".")
	if p.Component != nil {
		return fmt.Sprintf("%s (%s): %v", ident, p.VitType, p.Component)
	}

	return fmt.Sprintf("%s (%s): %s", ident, p.VitType, p.Expression)
}

func (p Property) Position() vit.PositionRange {
	return p.position
}
