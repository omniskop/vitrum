package parse

import "github.com/omniskop/vitrum/vit"

// unitType describes the type of a code unit that has been parsed
type unitType int

const (
	unitTypeNil          unitType = iota // no valid unit was parsed
	unitTypeEOF                          // the file has ended
	unitTypeComponent                    // a component definition has been parsed
	unitTypeComponentEnd                 // a component has ended
	unitTypeProperty                     // a property has been parsed
	unitTypeEnum                         // an enum has been parsed
)

// String returns the name of a unitType as a string
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

type unit struct {
	position vit.PositionRange
	kind     unitType
	value    interface{}
}

func nilUnit() unit {
	return unit{position: vit.PositionRange{}, kind: unitTypeNil, value: nil}
}

func eofUnit(position vit.Position) unit {
	return unit{vit.NewRangeFromPosition(position), unitTypeEOF, nil}
}

func componentUnit(position vit.PositionRange, comp *vit.ComponentDefinition) unit {
	return unit{position, unitTypeComponent, comp}
}

func componentEndUnit(position vit.Position) unit {
	return unit{vit.NewRangeFromPosition(position), unitTypeComponentEnd, nil}
}

func propertyUnit(position vit.PositionRange, prop vit.PropertyDefinition) unit {
	return unit{position, unitTypeProperty, prop}
}

func enumUnit(position vit.PositionRange, enum vit.Enumeration) unit {
	return unit{position, unitTypeEnum, enum}
}
