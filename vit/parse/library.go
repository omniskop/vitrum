package parse

import (
	"fmt"
	"strings"

	"github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/std"
)

// A Library describes defines one or more components that can be used in other files
type Library interface {
	ComponentNames() []string
	NewComponent(string, string, vit.ComponentContext) (vit.Component, bool)
	StaticAttribute(string, string) (interface{}, bool)
}

var libraries = make(map[string]Library)

func RegisterLibrary(name string, lib Library) {
	libraries[name] = lib
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
			return std.StdLib{}, nil
		}
	case "QtQuick":
		return std.StdLib{}, nil
	case "Dark":
		return std.StdLib{}, nil
	default:
		if lib, ok := libraries[strings.Join(namespace, ".")]; ok {
			return lib, nil
		}
	}

	return nil, fmt.Errorf("unknown library %q", strings.Join(namespace, "."))
}
