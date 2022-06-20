package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/parse"
)

const vitPackage = "github.com/omniskop/vitrum/vit"
const stdPackage = "github.com/omniskop/vitrum/vit/std"

// All generator specific tags
const (
	onChangeTag    = "gen-onchange"
	internalTag    = "gen-internal"
	typeTag        = "gen-type"
	initializerTag = "gen-initializer"
	privateTag     = "gen-private"
	optionalTag    = "gen-optional"
)

func GenerateFromFileAndSave(srcPath string, packageName string, dstPath string) error {
	doc, err := parseVit(srcPath)
	if err != nil {
		return fmt.Errorf("unable to parse: %v", err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("unable to create destination file: %w", err)
	}

	err = GenerateFromDocument(doc, packageName, dstFile)
	dstFile.Close()
	if err != nil {
		err2 := os.Remove(dstPath)
		if err2 != nil {
			fmt.Fprintln(os.Stderr, "unable to remove output file again:", err2)
		}
		return err
	}
	return nil
}

func Generate(src io.Reader, srcPath string, packageName string, dst io.Writer) error {
	lexer := parse.NewLexer(src, srcPath)

	doc, err := parse.Parse(parse.NewTokenBuffer(lexer.Lex))
	if err != nil {
		return err
	}
	doc.Name = getComponentName(srcPath)

	return GenerateFromDocument(doc, packageName, dst)
}

func GenerateFromFile(srcFile string, packageName string, dst io.Writer) error {
	doc, err := parseVit(srcFile)
	if err != nil {
		return err
	}

	return GenerateFromDocument(doc, packageName, dst)
}

func parseVit(srcFile string) (*parse.VitDocument, error) {
	file, err := os.Open(srcFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lexer := parse.NewLexer(file, srcFile)

	doc, err := parse.Parse(parse.NewTokenBuffer(lexer.Lex))
	if err != nil {
		return nil, err
	}
	doc.Name = getComponentName(srcFile)

	return doc, nil
}

func GenerateFromDocument(doc *parse.VitDocument, packageName string, dst io.Writer) error {
	f := jen.NewFilePath(packageName)
	f.HeaderComment("Code generated by vitrum gencmd. DO NOT EDIT.")

	for _, comp := range doc.Components {
		err := generateComponent(f, doc.Name, comp)
		if err != nil {
			return err
		}
	}

	return f.Render(dst)
}

func getComponentName(fileName string) string {
	fileName = filepath.Base(fileName)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func generateComponent(f *jen.File, compName string, comp *vit.ComponentDefinition) error {
	// TODO: Figure out from which package the base component should be imported from. Currently this is hardcoded to be the std package.

	f.Add(generateComponentEnums(compName, comp))

	properties := []jen.Code{
		jen.Qual(stdPackage, comp.BaseName),
		jen.Id("id").String(),
		jen.Line(),
	}

	// All property instantiations
	// we could use jen.Dict here but I wan't to preserve the property order
	propertyInstantiations := []jen.Code{
		jen.Line().Id(comp.BaseName).Op(":").Op("*").Qual(stdPackage, fmt.Sprintf("New%s", comp.BaseName)).Op("(").Id("id").Op(",").Id("scope").Op(")"),
		jen.Line().Id("id").Op(":").Id("id"),
	}

	// name of the variable that will hold the receiver of the components methods
	receiverName := strings.ToLower(string(compName[0]))

	// setup all properties for the struct definition as well as the instantiations
	for _, prop := range comp.Properties {
		if isInternalProperty(prop) {
			continue
		}
		propType, propConstructor, err := vitTypeInfo(comp, prop)
		if err != nil {
			return err
		}

		properties = append(properties, jen.Id(prop.Identifier[0]).Add(propType))

		// property instantiation
		propertyInstantiations = append(propertyInstantiations, jen.Line().Id(prop.Identifier[0]).Op(":").Add(propConstructor))
	}

	propertyInstantiations = append(propertyInstantiations, jen.Line())

	f.Type().Id(compName).Struct(properties...)

	// constructor
	f.Func().
		Id(fmt.Sprintf("New%s", compName)).
		Params(jen.Id("id").String(), jen.Id("scope").Qual(vitPackage, "ComponentContainer")).
		Params(jen.Op("*").Id(compName)).
		Block(
			jen.Return(jen.Op("&").Id(compName).Values(propertyInstantiations...)),
		)

	f.Line()

	// .String() string
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("String").
		Params().
		Params(jen.String()).
		Block(
			jen.Return(jen.Qual("fmt", "Sprintf").Call(jen.Lit(compName+"(%s)"), jen.Id(receiverName).Dot("id"))),
		).
		Line()

	// .Property(key string) (Value, bool)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("Property").
		Params(jen.Id("key").String()).
		Params(jen.Qual(vitPackage, "Value"), jen.Bool()).
		Block(
			jen.Switch(jen.Id("key")).Block(
				append(mapProperties(comp.Properties, func(prop vit.PropertyDefinition, propId string) jen.Code {
					if !isReadable(prop) {
						return nil // don't add unreadable properties
					}
					return jen.Case(jen.Lit(propId)).Block(
						jen.Return(jen.Op("&").Id(receiverName).Dot(propId), jen.True()),
					)
				}),
					jen.Default().Block(
						jen.Return(jen.Id(receiverName).Dot(comp.BaseName).Dot("Property").Call(jen.Id("key"))),
					),
				)...,
			),
		).
		Line()

	// .MustProperty(key string) Value
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("MustProperty").
		Params(jen.Id("key").String()).
		Params(jen.Qual(vitPackage, "Value")).
		Block(
			jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Id(receiverName).Dot("Property").Call(jen.Id("key")),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Panic(jen.Qual("fmt", "Errorf").Call(jen.Lit("MustProperty called with unknown key %q"), jen.Id("key"))),
			),
			jen.Return(jen.Id("v")),
		).
		Line()

	// .SetProperty(key string, value interface{}) error
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("SetProperty").
		Params(jen.Id("key").String(), jen.Id("value").Interface()).
		Params(jen.Error()).
		Block(
			jen.Var().Id("err").Error(),
			jen.Switch(jen.Id("key")).Block(
				append(mapProperties(comp.Properties, func(prop vit.PropertyDefinition, propId string) jen.Code {
					if !isWritable(prop) {
						return nil // don't add unwritable properties
					}
					return jen.Case(jen.Lit(propId)).Block(
						jen.Id("err").Op("=").Id(receiverName).Dot(propId).Op(".").Id("SetValue").Call(jen.Id("value")),
					)
				}),
					jen.Default().Block(
						jen.Return(jen.Id(receiverName).Dot(comp.BaseName).Dot("SetProperty").Call(jen.Id("key"), jen.Id("value"))),
					),
				)...,
			),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return().Qual(vitPackage, "NewPropertyError").Call(jen.Lit(compName), jen.Id("key"), jen.Id(receiverName).Dot("id"), jen.Id("err")),
			),
			jen.Return(jen.Nil()),
		).
		Line()

	// .SetPropertyExpression(key string, code string, pos *vit.PositionRange) error
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("SetPropertyExpression").
		Params(jen.Id("key").String(), jen.Id("code").String(), jen.Id("pos").Op("*").Qual(vitPackage, "PositionRange")).
		Params(jen.Error()).
		Block(
			jen.Switch(jen.Id("key")).Block(
				append(mapProperties(comp.Properties, func(prop vit.PropertyDefinition, propId string) jen.Code {
					if !isWritable(prop) {
						return nil // don't add unwritable properties
					}
					return jen.Case(jen.Lit(propId)).Block(
						jen.Id(receiverName).Dot(propId).Op(".").Id("SetExpression").Call(jen.Id("code"), jen.Id("pos")),
					)
				}),
					jen.Default().Block(
						jen.Return(jen.Id(receiverName).Dot(comp.BaseName).Dot("SetPropertyExpression").Call(jen.Id("key"), jen.Id("code"), jen.Id("pos"))),
					),
				)...,
			),
			jen.Return().Nil(),
		).
		Line()

	// .ResolveVariable(key string) (interface{}, bool)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("ResolveVariable").
		Params(jen.Id("key").String()).
		Params(jen.Interface(), jen.Bool()).
		Block(
			jen.Switch(jen.Id("key")).Block(
				prepend(
					jen.Case(jen.Id(receiverName).Dot("id")).Block(
						jen.Return(jen.Id(receiverName), jen.True()),
					),
					mapProperties(comp.Properties, func(prop vit.PropertyDefinition, propId string) jen.Code {
						if !isReadable(prop) {
							return nil // don't add unreadable properties
						}
						return jen.Case(jen.Lit(propId)).Block(
							jen.Return(jen.Op("&").Id(receiverName).Dot(propId), jen.True()),
						)
					}),
					jen.Default().Block(
						jen.Return(jen.Id(receiverName).Dot(comp.BaseName).Dot("ResolveVariable").Call(jen.Id("key"))),
					),
				)...,
			),
		).
		Line()

	// .AddChild(child Component)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("AddChild").
		Params(jen.Id("child").Qual(vitPackage, "Component")).
		Block(
			generateCallbackForAddedChild(comp, receiverName, "child"),
			jen.Id("child").Dot("SetParent").Call(jen.Id(receiverName)),
			jen.Id(receiverName).Dot("AddChildButKeepParent").Call(jen.Id("child")),
		).
		Line()

	// .AddChildAfter(afterThis, addThis Component)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("AddChildAfter").
		Params(jen.Id("afterThis").Qual(vitPackage, "Component"), jen.Id("addThis").Qual(vitPackage, "Component")).
		Block(
			generateCallbackForAddedChild(comp, receiverName, "addThis"),
			jen.Var().Id("targetType").Qual(vitPackage, "Component").Op("=").Id("afterThis"),
			jen.Line(),
			jen.For(jen.List(jen.Id("ind"), jen.Id("child")).Op(":=").Range().Id(receiverName).Dot("Children").Call()).Block(
				jen.If(jen.Id("child").Dot("As").Call(jen.Op("&").Id("targetType"))).Block(
					jen.Id("addThis").Dot("SetParent").Call(jen.Id(receiverName)),
					jen.Id(receiverName).Dot("AddChildAtButKeepParent").Call(jen.Id("addThis"), jen.Id("ind").Op("+").Lit(1)),
					jen.Return(),
				),
			),
			jen.Id(receiverName).Dot("AddChild").Call(jen.Id("addThis")),
		).
		Line()

	// .UpdateExpressions() (int, ErrorGroup)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("UpdateExpressions").
		Params().
		Params(jen.Int(), jen.Qual(vitPackage, "ErrorGroup")).
		BlockFunc(func(g *jen.Group) {
			// initialize 'sum' and 'errs' variables
			g.Var().Id("sum").Int()
			g.Var().Id("errs").Qual(vitPackage, "ErrorGroup")
			g.Line()
			// now handle changes for all necessary properties
			addMultiple(g, mapProperties(comp.Properties, func(prop vit.PropertyDefinition, propID string) jen.Code {
				if prop.HasTag(typeTag) && !prop.HasTag(onChangeTag) {
					// We will not handle changes for properties with custom types.
					// Except for ones that also have the 'onchange' tag set. In this case the custom type must implement the 'ShouldEvaluate', 'Update' and 'GetExpression' methods from the Value interface.
					return nil
				}
				var changeHandler jen.Code
				if handlerName, ok := prop.Tags[onChangeTag]; ok {
					changeHandler = jen.Id(receiverName).Dot(handlerName).Call(jen.Id(receiverName).Dot(propID))
				}
				return jen.If(jen.List(jen.Id("changed"), jen.Id("err")).Op(":=").Id(receiverName).Dot(propID).Dot("Update").Call(jen.Id(receiverName)).Op(";").Id("changed").Op("||").Id("err").Op("!=").Nil()).Block(
					jen.Id("sum").Op("++"),
					jen.If(jen.Id("err").Op("!=").Nil()).Block(
						jen.Id("errs").Dot("Add").Call(jen.Qual(vitPackage, "NewPropertyError").Call(
							jen.Lit(compName),
							jen.Lit(propID),
							jen.Id(receiverName).Dot("id"),
							jen.Id("err"),
						)),
					),
					changeHandler,
				)
			}))
			g.Line()
			g.Comment("this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables")
			g.Id("n").Op(",").Id("err").Op(":=").Id(receiverName).Dot("UpdatePropertiesInContext").Call(jen.Id(receiverName)) // n, err := receiver.UpdatePropertiesInContext(receiver) // just approximate code, names will vary
			g.Id("sum").Op("+=").Id("n")                                                                                      // sum += n
			g.Id("errs").Dot("AddGroup").Call(jen.Id("err"))                                                                  // errs.AddGroup(err)
			g.Id("n").Op(",").Id("err").Op("=").Id(receiverName).Dot(comp.BaseName).Dot("UpdateExpressions").Call()           // n, err = receiver.BaseComponent.UpdateExpressions()
			g.Id("sum").Op("+=").Id("n")                                                                                      // sum += n
			g.Id("errs").Dot("AddGroup").Call(jen.Id("err"))                                                                  // errs.AddGroup(err)
			g.Return(jen.Id("sum"), jen.Id("errs"))                                                                           // return sum, errs
		}).
		Line()

	// .As(*Component) (bool)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("As").
		Params(jen.Id("target").Op("*").Qual(vitPackage, "Component")).
		Params(jen.Bool()).
		Block(
			jen.If(jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").Parens(jen.Op("*").Id("target")).Op(".").Parens(jen.Op("*").Id(compName)).Op(";").Id("ok")).Block(
				jen.Op("*").Id("target").Op("=").Id(receiverName),
				jen.Return(jen.True()),
			),
			jen.Return(jen.Id(receiverName).Dot("Item").Dot("As").Call(jen.Id("target"))),
		).
		Line()

	// ID() string
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("ID").
		Params().
		Params(jen.String()).
		Block(jen.Return(jen.Id(receiverName).Dot("id"))).
		Line()

	// Finish() error
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("Finish").
		Params().
		Params(jen.Error()).
		Block(
			jen.Return(jen.Id(receiverName).Dot("RootC").Call().Dot("FinishInContext").Call(jen.Id(receiverName))),
		).
		Line()

	f.Add(generateStaticAttributeMethod(receiverName, compName, comp))

	return nil
}

func generateComponentEnums(compName string, comp *vit.ComponentDefinition) jen.Code {
	var code = new(jen.Group)
	for _, enum := range comp.Enumerations {
		typeName := fmt.Sprintf("%s_%s", compName, enum.Name)
		code.Type().Id(typeName).Uint().Line()
		code.Const().DefsFunc(func(g *jen.Group) {
			for _, value := range orderEnumValues(enum.Values) {
				g.Id(fmt.Sprintf("%s_%s", typeName, value.name)).Id(typeName).Op("=").Lit(value.value)
			}
		}).Line()
	}
	return code
}

func generateStaticAttributeMethod(receiverName, compName string, comp *vit.ComponentDefinition) jen.Code {
	var didGenerateSomething bool

	// staticAttribute(name string) (interface{}, bool)
	code := jen.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("staticAttribute").
		Params(jen.Id("name").String()).
		Params(jen.Interface(), jen.Bool()).
		Block(
			jen.Switch(jen.Id("name")).BlockFunc(func(g *jen.Group) {
				// list all values of embedded enums
				for _, enum := range comp.Enumerations {
					if !enum.Embedded {
						continue
					}
					didGenerateSomething = true
					for _, value := range orderEnumValues(enum.Values) {
						g.Case(jen.Lit(value.name)).Block(
							jen.Return().List(
								// the type needs to be converted to uint to be usable in an expression
								jen.Uint().Call(jen.Id(fmt.Sprintf("%s_%s_%s", compName, enum.Name, value.name))),
								jen.True(),
							),
						)
					}
				}
				g.Default().Block(
					jen.Return().List(jen.Nil(), jen.False()),
				)
			}),
		).
		Line()

	if !didGenerateSomething {
		return jen.Null()
	}
	return code
}

func generateComponentDefinition(comp *vit.ComponentDefinition) *jen.Statement {
	return nil
}

// mapProperties calls function 'f' for every property and combines their generated codes into one block.
// It skips properties that are considered to be internal.
// 'f' will be called with it's property definition and identifier of the property.
// It should the corresponding code for the property. It may return nil to indicate that no code should be added.
func mapProperties(props []vit.PropertyDefinition, f func(vit.PropertyDefinition, string) jen.Code) []jen.Code {
	var result []jen.Code
	for _, prop := range props {
		if isInternalProperty(prop) {
			continue
		}
		// NOTE: properties with multiple identifiers are currently not supported
		result = append(result, f(prop, prop.Identifier[0]))
	}
	return result
}

func prepend(first jen.Code, rest []jen.Code, tail ...jen.Code) []jen.Code {
	return append([]jen.Code{first}, append(rest, tail...)...)
}

// addMultiple adds all code to the group as separate statements
func addMultiple(g *jen.Group, code []jen.Code) {
	for _, c := range code {
		if c != nil {
			g.Add(c)
		}
	}
}

// vitTypeInfo returns statements that describe the type and constructor for a given property.
// It might return an error if the property contains incompatible tags or the type is unknown.
// For example for a property of type 'float' the returned code might look like this:
//     propType:    vit.FloatType
//     constructor: *vit.NewFloatValue("", nil),
func vitTypeInfo(comp *vit.ComponentDefinition, prop vit.PropertyDefinition) (propType *jen.Statement, constructor *jen.Statement, err error) {
	// handles gen-initializer and gen-type
	if init, ok := prop.Tags[initializerTag]; ok {
		constructor = jen.Id(init) // a custom initializer is provided
		// this also requires a custom initializer
		if typeString, ok := prop.Tags[typeTag]; ok {
			propType = jen.Id(typeString) // a custom type is provided

			if prop.HasTag(optionalTag) {
				// a custom type can't be optional at the same time
				err = fmt.Errorf("property %s cannot combine %q with %q", prop.Identifier, typeTag, optionalTag)
			}

			return
		} else {
			// only a custom type but no initializer was provided
			err = fmt.Errorf("property %s has %q tag but no %q tag", prop.Identifier, initializerTag, typeTag)
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
		constructor = jen.Op("*").Qual(vitPackage, "NewEmptyStringValue").Call()
	case "int":
		propType = jen.Qual(vitPackage, "IntValue")
		constructor = jen.Op("*").Qual(vitPackage, "NewEmptyIntValue").Call()
	case "float":
		propType = jen.Qual(vitPackage, "FloatValue")
		constructor = jen.Op("*").Qual(vitPackage, "NewEmptyFloatValue").Call()
	case "bool":
		propType = jen.Qual(vitPackage, "BoolValue")
		constructor = jen.Op("*").Qual(vitPackage, "NewEmptyBoolValue").Call()
	case "color":
		propType = jen.Qual(vitPackage, "ColorValue")
		constructor = jen.Op("*").Qual(vitPackage, "NewEmptyColorValue").Call()
	case "var":
		propType = jen.Qual(vitPackage, "AnyValue")
		constructor = jen.Op("*").Qual(vitPackage, "NewEmptyAnyValue").Call()
	case "component":
		propType = jen.Qual(vitPackage, "ComponentDefValue")
		constructor = jen.Op("*").Qual(vitPackage, "NewEmptyComponentDefValue").Call()
	default:
		if _, ok := comp.GetEnum(prop.VitType); ok {
			propType = jen.Qual(vitPackage, "IntValue")
			constructor = jen.Op("*").Qual(vitPackage, "NewEmptyIntValue").Call()
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

// generateCallbackForAddedChild checks if the component provides a callback for the event that a child has been added to the component.
// If that is the case it returns a defer statement that should be added to the top of all methods that add a child.
// If no callback is provided it returns nil.
func generateCallbackForAddedChild(comp *vit.ComponentDefinition, receiverName, parameterName string) jen.Code {
	if childrenProp, ok := getProperty(comp, "children"); ok && childrenProp.HasTag(onChangeTag) {
		// if the children property is explicitly provided and has the 'onchange' tag we will call the provided method with the added child.
		return jen.Defer().Id(receiverName).Dot(childrenProp.Tags[onChangeTag]).Call(jen.Id(parameterName))
	}
	return nil
}

// valueChange returns the code for changing a value of a property that is appropriate for the type of the property.
func valueChange(vitType string, listDimentions int) *jen.Statement {
	switch vitType {
	case "component":
		if listDimentions == 0 {
			return jen.Id("ChangeComponent").Call(jen.Id("value").Assert(jen.Index().Op("*").Qual(vitPackage, "ComponentDefinition")).Index(jen.Lit(0)))
		} else {
			return jen.Id("ChangeComponents").Call(jen.Id("value").Assert(jen.Index().Op("*").Qual(vitPackage, "ComponentDefinition")))
		}
	default:
		return jen.Id("ChangeCode").Call(jen.Id("value").Assert(jen.String()), jen.Id("position"))
	}
}

// capitalizes the first letter of a string
func firstLetterUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Creates a copy of the given map but excludes the given keys.
func copyMapExcluding(m map[string]string, exclude ...string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if !stringSliceContains(exclude, k) {
			result[k] = v
		}
	}
	return result
}

// Returns true if the string slice contains the given string.
func stringSliceContains(s []string, v string) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}

// unpointer removes the first '*' from statements.
// TODO: currently this works by just always removing the first code without checking if that's actually a '*'.
func unpointer(code *jen.Statement) *jen.Statement {
	statement := (*code)[1:]
	return &statement
}

// isReadable return true if the property can be read from outside of the component.
func isReadable(prop vit.PropertyDefinition) bool {
	return !prop.HasTag(privateTag)
}

// isWritable returns true if the property can be written to from outside of the component.
func isWritable(prop vit.PropertyDefinition) bool {
	return !prop.HasTag(privateTag) && !prop.ReadOnly
}

// isInternalProperty returns true if the property is internal and no code should be generated for it.
func isInternalProperty(prop vit.PropertyDefinition) bool {
	return prop.Identifier[0] == "children"
}

// getProperty returns the property with the given name from the component.
// The book indicated that a property was found.
func getProperty(prop *vit.ComponentDefinition, identifier string) (*vit.PropertyDefinition, bool) {
	for _, p := range prop.Properties {
		if p.Identifier[0] == identifier {
			return &p, true
		}
	}
	return nil, false
}

func orderEnumValues(values map[string]int) []enumValue {
	var list = make(enumValueList, 0, len(values))
	for k, v := range values {
		list = append(list, enumValue{k, v})
	}
	sort.Sort(list)
	return list
}

type enumValue struct {
	name  string
	value int
}

type enumValueList []enumValue

func (v enumValueList) Len() int {
	return len(v)
}

func (v enumValueList) Less(i, j int) bool {
	return v[i].value < v[j].value
}

func (v enumValueList) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
