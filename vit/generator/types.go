package generator

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/parse"
)

// vitTypeInfo returns statements that describe the type and constructor for a given property.
// It might return an error if the property contains incompatible tags or the type is unknown.
// For example for a property of type 'float' the returned code might look like this:
//     propType:    vit.FloatType
//     constructor: *vit.NewFloatValueFromExpression("0", nil),
func vitTypeInfo(comp *vit.ComponentDefinition, prop vit.PropertyDefinition) (propType *jen.Statement, constructor *jen.Statement, err error) {
	// handles gen-initializer and gen-type
	if init, ok := prop.Tags[initializerTag]; ok {
		constructor = jen.Id(init) // a custom initializer is provided
		// this also requires a custom initializer
		if typeString, ok := prop.Tags[typeTag]; ok {
			propType = generateCustomType(typeString)

			if prop.HasTag(optionalTag) {
				// a custom type can't be optional at the same time
				err = fmt.Errorf("property %s cannot combine %q with %q", prop.Identifier, typeTag, optionalTag)
			}

			return
		} else {
			// only a custom type but no initializer was provided
			err = fmt.Errorf("property %s has %q tag but no %q tag", prop.Identifier, initializerTag, typeTag)
			return
		}
	}

	// check if this property holds component definitions
	if len(prop.Components) == 1 {
		// this property holds a single component
		propType = generateComponentDefinition(prop.Components[0])
	} else if len(prop.Components) > 1 {
		// this property holds multiple components
		propType = jen.Op("[")
		for i, comp := range prop.Components {
			if i > 0 {
				propType.Op(",")
			}
			propType.Add(generateComponentDefinition(comp))
		}
		propType.Op("]")
	} else if prop.Expression != "" {
		// this property holds an expression
		propType = jen.Lit(prop.Expression)
	}

	// handle standard vit types
	switch prop.VitType {
	case "string":
		propType = jen.Qual(vitPackage, "StringValue")
		constructor = standardConstructor(prop, "String")
	case "int":
		propType = jen.Qual(vitPackage, "IntValue")
		constructor = standardConstructor(prop, "Int")
	case "float":
		propType = jen.Qual(vitPackage, "FloatValue")
		constructor = standardConstructor(prop, "Float")
	case "bool":
		propType = jen.Qual(vitPackage, "BoolValue")
		constructor = standardConstructor(prop, "Bool")
	case "color":
		propType = jen.Qual(vitPackage, "ColorValue")
		constructor = standardConstructor(prop, "Color")
	case "var":
		propType = jen.Qual(vitPackage, "AnyValue")
		constructor = standardConstructor(prop, "Any")
	case "component":
		propType = jen.Qual(vitPackage, "ComponentDefValue")
		constructor = jen.Op("*").Qual(vitPackage, "NewEmptyComponentDefValue").Call()
	case "group":
		propType, constructor, err = typeInfoForGroup(comp, prop)
		if err != nil {
			return
		}
	default:
		if _, ok := comp.GetEnum(prop.VitType); ok {
			propType = jen.Qual(vitPackage, "IntValue")
			constructor = standardConstructor(prop, "Int")
		} else {
			err = fmt.Errorf("property %s has unknown type %q", prop.Identifier, prop.VitType)
			return
		}
	}

	// check if this property is explicitly optional
	if prop.HasTag(optionalTag) {
		// wrap the actual type in an OptionalValue
		propType = jen.Qual(vitPackage, "OptionalValue").Types(jen.Op("*").Add(propType))
		constructor = jen.Op("*").Qual(vitPackage, "NewOptionalValue").Call(jen.Add(unpointer(constructor)))
	}

	// check if this property is a list
	if prop.ListDimensions > 0 {
		if prop.VitType == "component" {
			// TODO: should this be a componentRef instead?
			propType = jen.Qual(vitPackage, "ComponentDefListValue")
			constructor = jen.Op("*").Qual(vitPackage, "NewComponentDefListValue").Call(jen.Nil(), jen.Nil())
			return
		}

		elementType := jen.Op("*").Add(propType)
		propType = jen.Qual(vitPackage, "ListValue").Types(jen.Add(elementType))
		constructor = jen.Op("*").Qual(vitPackage, "NewListValue").Types(elementType).Call(jen.Lit(""), jen.Nil())
		return
	}

	return
}

func standardConstructor(prop vit.PropertyDefinition, typeName string) *jen.Statement {
	if prop.Expression == "" {
		return jen.Op("*").Qual(vitPackage, fmt.Sprintf("NewEmpty%sValue", typeName)).Call()
	} else {
		return jen.Op("*").Qual(vitPackage, fmt.Sprintf("New%sValueFromExpression", typeName)).Call(jen.Lit(prop.Expression), generatePositionRange(prop.Pos))
	}
}

func typeInfoForGroup(comp *vit.ComponentDefinition, prop vit.PropertyDefinition) (propType *jen.Statement, constructor *jen.Statement, err error) {
	subProps, err := parse.ParseGroupDefinition(prop.Expression, prop.Pos.Start())
	if err != nil {
		return nil, nil, err
	}

	var constructors = make([]struct {
		name        string
		constructor *jen.Statement
	}, len(subProps))

	for i, subProp := range subProps {
		_, subConstructor, err := vitTypeInfo(comp, subProp)
		if err != nil {
			return nil, nil, err
		}
		constructors[i].name = subProp.Identifier[0]
		constructors[i].constructor = unpointer(subConstructor)
	}

	propType = jen.Qual(vitPackage, "GroupValue")
	constructor = jen.Op("*").Qual(vitPackage, "NewEmptyGroupValue").Call(
		jen.Map(jen.String()).Qual(vitPackage, "Value").BlockFunc(func(g *jen.Group) {
			for _, subConstructor := range constructors {
				g.Lit(subConstructor.name).Op(":").Add(subConstructor.constructor).Op(",")
			}
		}),
	)
	return
}

// generateCustomType takes a string representation of a custom type and returns the appropriate generated code.
// It allows fully qualified types including package names. Like: "github.com/omniskop/vitrum/vit.ComponentDefinition"
// The "vit" package can be used directly as a shortcut instead of providing the full path.
func generateCustomType(typeString string) *jen.Statement {
	var propType *jen.Statement
	// check if the type contains a period, which indicates that it refers to a type in another package
	if strings.Contains(typeString, ".") {
		// split the type apart and get the package path and the actual type name
		packagePath, typeName, isPointer := splitCustomType(typeString)
		if isPointer {
			propType = jen.Op("*").Qual(packagePath, typeName)
		} else {
			propType = jen.Qual(packagePath, typeName)
		}
	} else {
		propType = jen.Id(typeString) // a custom type is provided
	}
	return propType
}

// splitCustomType takes a custom type string and splits it into the package path and the actual type name.
// The "vit" package can be used directly as a shortcut instead of providing the full path.
// If the type is a pointer it will remove the star from the beginning and the 'isPointer' boolean will be true.
func splitCustomType(str string) (packagePath string, typeName string, isPointer bool) {
	lastPeriod := strings.LastIndex(str, ".")
	if lastPeriod == -1 {
		return "", str, false
	}
	packagePath = str[:lastPeriod]
	typeName = str[lastPeriod+1:]
	if len(packagePath) > 0 && packagePath[0] == '*' {
		isPointer = true
		packagePath = packagePath[1:]
	}
	if packagePath == "vit" {
		// shortcut
		packagePath = vitPackage
	}
	return
}
