package parse

import (
	"fmt"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

// VitDocument contains everything there is to know about a parsed vit file
type VitDocument struct {
	name       string            // Name of the file without extension. Usually the name of the component this file describes.
	imports    []importStatement // all imported libraries and files
	components []*componentDefinition
}

// String creates a human readable string representation of the vit document
func (d VitDocument) String() string {
	var out strings.Builder
	out.WriteString("{\r\n\tImports:")
	for _, imp := range d.imports {
		out.WriteString(fmt.Sprintf(" %v", imp))
	}

	out.WriteString("\r\n\tComponents: \r\n")

	for _, comp := range d.components {
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

// componentDefinition contains everything about a component that is defined in a vit file
type componentDefinition struct {
	name         string                 // name of the instantiated component
	id           string                 // custom id of the component
	properties   []property             // all explicitly defined or declared properties
	children     []*componentDefinition // child components
	enumerations []vit.Enumeration      // all explicitly defined enumerations
}

// identifierIsKnown returns true of the given identifier has already been defined
func (d *componentDefinition) identifierIsKnown(identifier []string) bool {
	if len(identifier) == 0 {
		return false
	}
	for _, p := range d.properties {
		if stringSlicesEqual(p.identifier, identifier) {
			return true
		}
	}
	for _, e := range d.enumerations {
		if e.Name == identifier[0] {
			return true
		}
	}

	return false
}

// String returns a human readable string representation of the component definition
func (d *componentDefinition) String() string {
	if d == nil {
		return "<nil>"
	}
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s { ", d.name))

	for _, p := range d.properties {
		s.WriteString(fmt.Sprintf("%s, ", p.String()))
	}

	if len(d.children) == 0 {

	} else if len(d.children) == 1 {
		s.WriteString(fmt.Sprintf("%d child", len(d.children)))
	} else {
		s.WriteString(fmt.Sprintf("%d children", len(d.children)))
	}

	// TODO: add other fields

	s.WriteString(" }")

	return s.String()
}

// property contains everything about a defined or declared property
type property struct {
	identifier []string             // Identifier of this property. This will usually be only one value but can contain multiple parts for example with 'Anchors.fill'
	vitType    string               // data type of the property in vit terms, not go
	expression string               // Expression string that defines the property. Can be empty.
	component  *componentDefinition // Only set if this properties type is a component.
	readOnly   bool                 // Readonly properties are statically defined on the component itself and cannot be changed directly. They will however be recalculated if one of the expressions dependencies should change.
}

// String returns a human readable string representation of the property
func (p property) String() string {
	ident := strings.Join(p.identifier, ".")
	if p.component != nil {
		return fmt.Sprintf("%s (%s): %v", ident, p.vitType, p.component)
	}

	return fmt.Sprintf("%s (%s): %s", ident, p.vitType, p.expression)
}
