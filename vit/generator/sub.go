package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/omniskop/vitrum/vit"
)

// This file contains functions that generate code describing specific data structures.

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
	return jen.Id(fmt.Sprintf("%#v", comp))
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
