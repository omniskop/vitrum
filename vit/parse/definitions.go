package parse

import (
	"fmt"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

var dataTypes = map[string]bool{
	"bool":   true,
	"int":    true,
	"float":  true,
	"string": true,
	"color":  true,
}

var keywords = map[string]bool{
	"property": true,
	"default":  true,
	"required": true,
	"readonly": true,
	"enum":     true,
	"embedded": true,
}

type tokenSource func() token

type VitDocument struct {
	name       string
	imports    []importStatement
	components []*componentDefinition
}

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

// An importStatement can either import a namespace or a file
// namspaces have a version with major and minor part.
// Additionally an identifier can be specified.
type importStatement struct {
	namespace  []string
	file       string
	version    string
	identifier string
}

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

type unitType int

const (
	unitTypeNil unitType = iota
	unitTypeEOF
	unitTypeComponent
	unitTypeComponentEnd
	unitTypeProperty
	unitTypeEnum
)

func (uType unitType) String() string {
	switch uType {
	case unitTypeEOF:
		return "end of file"
	case unitTypeComponent:
		return "component"
	case unitTypeComponentEnd:
		return "end of component"
	case unitTypeProperty:
		return "property"
	case unitTypeEnum:
		return "enum"
	default:
		return "unknown unit"
	}
}

type componentDefinition struct {
	name       string
	id         string
	properties []property
	children   []*componentDefinition
	// signals               []signal
	// signalHandlers        []signalHandler
	// methods               []method
	// attachedProperties    []property
	// attachedSignalHandler []signalHandler
	enumerations []vit.Enumeration
}

func (d *componentDefinition) IdentifierIsKnown(identifier []string) bool {
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

type property struct {
	identifier []string
	vitType    string
	expression string
	component  *componentDefinition
}

func (p property) String() string {
	ident := strings.Join(p.identifier, ".")
	if p.component != nil {
		return fmt.Sprintf("%s (%s): %v", ident, p.vitType, p.component)
	} else {
		return fmt.Sprintf("%s (%s): %s", ident, p.vitType, p.expression)
	}
}
