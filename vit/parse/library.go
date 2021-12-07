package parse

import (
	"fmt"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

// A Library describes defines one or more components that can be used in other files
type Library interface {
	ComponentNames() []string
	NewComponent(string, string, vit.ComponentResolver) (vit.Component, bool)
}

// resolveLibraryImport takes a library identifier and returns the corresponding library or an error if the identifier is unknown.
// Currently this is hardcoded but should be made dynamic in the future.
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
