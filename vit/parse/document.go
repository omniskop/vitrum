package parse

import (
	"fmt"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

// VitDocument contains everything there is to know about a parsed vit file
type VitDocument struct {
	Name       string            // Name of the file without extension. Usually the name of the component this file describes.
	Imports    []ImportStatement // all imported libraries and files
	Components []*vit.ComponentDefinition
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

// An ImportStatement can either import a module/namespace or a file
// namespaces have a version with major and minor part.
// Either namespace or file can be set, but not both.
type ImportStatement struct {
	Namespace []string // fully qualified name of the module to import
	File      string   // file path that should be imported
	Version   string   // version string for namespace imports
	Qualifier string   // optional qualifier that allows the user to refer to the import by a different name
	Position  vit.PositionRange
}

// String returns a human readable multiline string representation of the import
func (s ImportStatement) String() string {
	var out strings.Builder
	if len(s.Namespace) == 0 {
		out.WriteString(fmt.Sprintf("%q", s.File))
	} else {
		out.WriteString(strings.Join(s.Namespace, "."))
	}
	out.WriteRune('@')

	out.WriteString(s.Version)

	return out.String()
}
