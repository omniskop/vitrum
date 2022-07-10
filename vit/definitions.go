package vit

import (
	"fmt"
	"strings"
)

// ComponentDefinition contains everything about a component that is defined in a vit file
type ComponentDefinition struct {
	Pos          PositionRange
	BaseName     string                 // name of the instantiated base component
	ID           string                 // custom id of the component
	Properties   []PropertyDefinition   // all explicitly defined or declared properties
	Children     []*ComponentDefinition // child components
	Enumerations []Enumeration          // all explicitly defined enumerations
	Events       []EventDefinition
}

// IdentifierIsKnown returns true of the given identifier has already been defined
func (d *ComponentDefinition) IdentifierIsKnown(identifier []string) bool {
	if len(identifier) == 0 {
		return false
	}
	for _, p := range d.Properties {
		if slicesEqual(p.Identifier, identifier) {
			return true
		}
	}
	for _, e := range d.Enumerations {
		if e.Name == identifier[0] {
			return true
		}
	}

	return false
}

// String returns a human readable string representation of the component definition
func (d *ComponentDefinition) String() string {
	if d == nil {
		return "<nil>"
	}
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s { ", d.BaseName))

	for _, p := range d.Properties {
		s.WriteString(fmt.Sprintf("%s, ", p.String()))
	}

	if len(d.Children) == 0 {

	} else if len(d.Children) == 1 {
		s.WriteString(fmt.Sprintf("%d child", len(d.Children)))
	} else {
		s.WriteString(fmt.Sprintf("%d children", len(d.Children)))
	}

	// TODO: add other fields

	s.WriteString(" }")

	return s.String()
}

func (d *ComponentDefinition) GetEnum(name string) (*Enumeration, bool) {
	for _, e := range d.Enumerations {
		if e.Name == name {
			return &e, true
		}
	}
	return nil, false
}

// PropertyDefinition contains everything about a defined or declared PropertyDefinition
type PropertyDefinition struct {
	Pos        PositionRange          // position of the property declaration
	ValuePos   *PositionRange         // Position of the property value
	Identifier []string               // Identifier of this property. This will usually be only one value but can contain multiple parts for example with 'Anchors.fill'
	Expression string                 // Expression string that defines the property. Can be empty.
	Components []*ComponentDefinition // Only set if this properties type is a component or list of components.

	VitType        string            // data type of the property in vit terms, not go
	ListDimensions int               // If this property is a list this indicated the Number of dimensions it has. 0 means that it's not a list at all.
	ReadOnly       bool              // Readonly properties are statically defined on the component itself and cannot be changed directly. They will however be recalculated if one of the expressions dependencies should change.
	Static         bool              // Static properties are defined on the component itself. They will only be evaluated once when the component is loaded and are constant from that point on.
	StaticValue    interface{}       // The evaluated value of a static property.
	Tags           map[string]string // optional tags of the property
}

// String returns a human readable string representation of the property
func (p PropertyDefinition) String() string {
	ident := strings.Join(p.Identifier, ".")
	if len(p.Components) != 0 {
		return fmt.Sprintf("%s (%s): %v", ident, p.VitType, p.Components)
	}

	return fmt.Sprintf("%s (%s): %s", ident, p.VitType, p.Expression)
}

func (p PropertyDefinition) Position() PositionRange {
	return p.Pos
}

// HasTag returns true if the property has *one or more* of the given tags set and false otherwise.
func (p PropertyDefinition) HasTag(tag string, moreTags ...string) bool {
	// Singular name has been chosen to reflect the most common use case and to hopefully make it clearer that it will return true even if only one tag is set.
	if _, ok := p.Tags[tag]; ok {
		return true
	}
	for _, tag := range moreTags {
		if _, ok := p.Tags[tag]; ok {
			return true
		}
	}
	return false
}

func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, e := range a {
		if e != b[i] {
			return false
		}
	}
	return true
}
