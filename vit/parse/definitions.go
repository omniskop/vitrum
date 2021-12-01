package parse

import (
	"fmt"
	"strings"
)

var dataTypes = map[string]bool{
	"bool":   true,
	"int":    true,
	"float":  true,
	"string": true,
}

type tokenSource func() token

type VitDocument struct {
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
	unitTypeComponent
	unitTypeComponentEnd
	unitTypeProperty
	unitTypeEOF
)

func (uType unitType) String() string {
	switch uType {
	case unitTypeComponent:
		return "component"
	case unitTypeComponentEnd:
		return "end of component"
	case unitTypeProperty:
		return "property"
	case unitTypeEOF:
		return "end of file"
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
	// enums                 []enumeration
}

func (o *componentDefinition) String() string {
	if o == nil {
		return "<nil>"
	}
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s { ", o.name))

	for _, p := range o.properties {
		s.WriteString(fmt.Sprintf("%s, ", p.String()))
	}

	if len(o.children) == 0 {

	} else if len(o.children) == 1 {
		s.WriteString(fmt.Sprintf("%d child", len(o.children)))
	} else {
		s.WriteString(fmt.Sprintf("%d children", len(o.children)))
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
