package generator

import (
	"fmt"
	"reflect"

	"github.com/dave/jennifer/jen"
	"github.com/omniskop/vitrum/vit"
)

// This file contains functions that generate code describing specific data structures. (Sub-Generators)

// ============================== Type Definition for Enumerations =================================

// generateComponentEnums generates a new type and constants for the given enum
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
		// String method
		code.Func().Params(jen.Id("enum").Id(typeName)).Id("String").Params().Params(jen.String()).Block(
			jen.Switch(jen.Id("enum").BlockFunc(func(g *jen.Group) {
				// We'll keep track of all set values to make sure we don't double them in the switch statement.
				// Having two entries with the same value is permitted in an enum.
				setValues := make(map[int]string)
				for _, value := range orderEnumValues(enum.Values) {
					if setKey, ok := setValues[value.value]; ok {
						g.Comment(fmt.Sprintf("key %s is omitted as it's the same as %s", value.name, setKey))
						continue // this value was set already
					}
					setValues[value.value] = value.name
					g.Case(jen.Id(fmt.Sprintf("%s_%s", typeName, value.name))).Block(
						jen.Return(jen.Lit(value.name)),
					)
				}
				g.Default().Block(
					jen.Return(jen.Lit(fmt.Sprintf("<unknown%s>", enum.Name))),
				)
			})),
		).Line()

	}
	return code
}

// ======================================= Position Range ==========================================

// generatePositionRange returns code that recreates the given PositionRange.
// If hideSourceFiles is true, it will just return nil.
func generatePositionRange(pos vit.PositionRange) *jen.Statement {
	if hideSourceFiles {
		return jen.Nil()
	}
	return jen.Op("&").Qual(vitPackage, "PositionRange").Values(jen.Dict{
		jen.Id("FilePath"):    jen.Lit(pos.FilePath),
		jen.Id("StartLine"):   jen.Lit(pos.StartLine),
		jen.Id("StartColumn"): jen.Lit(pos.StartColumn),
		jen.Id("EndLine"):     jen.Lit(pos.EndLine),
		jen.Id("EndColumn"):   jen.Lit(pos.EndColumn),
	})
}

// ==================================== Component Definition =======================================

func generateComponentDefinition(comp *vit.ComponentDefinition) *jen.Statement {
	return generateFromValue(reflect.ValueOf(comp))
}

// ============================================= Map ===============================================

type validMapTypes interface {
	string | int | int64
}

func generateMap[keyType, valueType validMapTypes](m map[keyType]valueType) *jen.Statement {
	return jen.Map(jen.Id(fmt.Sprintf("%T", *new(keyType)))).Id(fmt.Sprintf("%T", *new(valueType))).ValuesFunc(func(g *jen.Group) {
		for k, v := range m {
			g.Lit(k).Op(":").Lit(v)
		}
	})
}

// ========================================= Enumeration ===========================================

func generateEnumeration(enum vit.Enumeration) *jen.Statement {
	values := orderEnumValues(enum.Values)
	valueMap := jen.Map(jen.String()).Int().ValuesFunc(func(g *jen.Group) {
		for _, v := range values {
			g.Lit(v.name).Op(":").Lit(v.value)
		}
	})
	return jen.Qual(vitPackage, "Enumeration").Values(jen.Dict{
		jen.Id("Name"):     jen.Lit(enum.Name),
		jen.Id("Embedded"): jen.Lit(enum.Embedded),
		jen.Id("Values"):   valueMap,
		jen.Id("Position"): generatePositionRange(*enum.Position),
	})
}

// =================================== Generic Type Generation =====================================

// generateFromValue generates code that recreates the given value.
// Attention: This function must not be called with recursive values. It will loop endlessly otherwise.
// It was inspired by fmt.printValue which does a similar thing. (https://cs.opensource.google/go/go/+/master:src/fmt/print.go)
func generateFromValue(value reflect.Value) *jen.Statement {
	switch f := value; value.Kind() {
	case reflect.Invalid:
		panic("invalid value")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String,
		reflect.Bool:
		return jen.Lit(value.Interface())
	case reflect.Map:
		return jen.Map(jen.Id(value.Type().Key().String())).Id(value.Type().Elem().String()).ValuesFunc(func(g *jen.Group) {
			sorted := sortMap(f)
			for _, item := range sorted {
				g.Add(jen.Id(item.key)).Op(":").Add(generateFromValue(item.value))
			}
		})
	case reflect.Struct:
		if f.IsZero() {
			return jen.Id(f.Type().String()).Values()
		}
		return jen.Id(f.Type().String()).ValuesFunc(func(g *jen.Group) {
			for i := 0; i < f.NumField(); i++ {
				fieldValue := f.Field(i)
				if fieldValue.IsZero() {
					continue // skip fields whose value is zero
				}
				if name := f.Type().Field(i).Name; name != "" {
					g.Id(name).Op(":").Add(generateFromValue(fieldValue))
				} else {
					g.Add(generateFromValue(fieldValue))
				}
			}
		})
	case reflect.Interface:
		value := f.Elem()
		if !value.IsValid() {
			return jen.Nil()
		} else {
			return jen.Add(generateFromValue(value))
		}
	case reflect.Array, reflect.Slice:
		return jen.Index().Id(f.Type().Elem().String()).ValuesFunc(func(g *jen.Group) {
			for i := 0; i < f.Len(); i++ {
				g.Add(generateFromValue(f.Index(i)))
			}
		})
	case reflect.Pointer:
		if f.Pointer() == 0 {
			jen.Nil()
		}
		return jen.Op("&").Add(generateFromValue(f.Elem()))
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return jen.Nil()
	default:
		panic(fmt.Sprintf("unsupported type: %s", value.Type()))
	}
}
