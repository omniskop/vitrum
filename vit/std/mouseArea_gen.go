// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
)

func newFileContextForMouseArea(globalCtx *vit.GlobalContext) (*vit.FileContext, error) {
	return vit.NewFileContext(globalCtx), nil
}

type MouseArea_MouseButtons uint

const (
	MouseArea_MouseButtons_noButton     MouseArea_MouseButtons = 0
	MouseArea_MouseButtons_leftButton   MouseArea_MouseButtons = 1
	MouseArea_MouseButtons_rightButton  MouseArea_MouseButtons = 2
	MouseArea_MouseButtons_middleButton MouseArea_MouseButtons = 4
	MouseArea_MouseButtons_allButtons   MouseArea_MouseButtons = 134217727
)

type MouseArea struct {
	*Item
	id string

	acceptedButtons vit.IntValue
	containsMouse   vit.BoolValue
	enabled         vit.BoolValue
	hoverEnabled    vit.BoolValue
	mouseX          vit.FloatValue
	mouseY          vit.FloatValue
	pressed         vit.BoolValue
	pressedButtons  vit.IntValue

	onClicked vit.EventAttribute[MouseEvent]
}

// newMouseAreaInGlobal creates an appropriate file context for the component and then returns a new MouseArea instance.
// The returned error will only be set if a library import that is required by the component fails.
func newMouseAreaInGlobal(id string, globalCtx *vit.GlobalContext) (*MouseArea, error) {
	fileCtx, err := newFileContextForMouseArea(globalCtx)
	if err != nil {
		return nil, err
	}
	return NewMouseArea(id, fileCtx), nil
}
func NewMouseArea(id string, context *vit.FileContext) *MouseArea {
	m := &MouseArea{
		Item:            NewItem("", context),
		id:              id,
		acceptedButtons: *vit.NewIntValueFromCode(vit.Code{FileCtx: context, Code: "MouseButtons.leftButton", Position: nil}),
		containsMouse:   *vit.NewEmptyBoolValue(),
		enabled:         *vit.NewBoolValueFromCode(vit.Code{FileCtx: context, Code: "true", Position: nil}),
		hoverEnabled:    *vit.NewEmptyBoolValue(),
		mouseX:          *vit.NewEmptyFloatValue(),
		mouseY:          *vit.NewEmptyFloatValue(),
		pressed:         *vit.NewEmptyBoolValue(),
		pressedButtons:  *vit.NewEmptyIntValue(),
		onClicked:       *vit.NewEventAttribute[MouseEvent](),
	}
	// property assignments on embedded components
	// register listeners for when a property changes
	m.enabled.AddDependent(vit.FuncDep(m.enableDisable))
	// register event listeners
	// register enumerations
	m.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "MouseButtons",
		Position: nil,
		Values:   map[string]int{"noButton": 0, "leftButton": 1, "rightButton": 2, "middleButton": 4, "allButtons": 134217727},
	})
	// add child components

	context.RegisterComponent("", m)

	return m
}

func (m *MouseArea) String() string {
	return fmt.Sprintf("MouseArea(%s)", m.id)
}

func (m *MouseArea) Property(key string) (vit.Value, bool) {
	switch key {
	case "acceptedButtons":
		return &m.acceptedButtons, true
	case "containsMouse":
		return &m.containsMouse, true
	case "enabled":
		return &m.enabled, true
	case "hoverEnabled":
		return &m.hoverEnabled, true
	case "mouseX":
		return &m.mouseX, true
	case "mouseY":
		return &m.mouseY, true
	case "pressed":
		return &m.pressed, true
	case "pressedButtons":
		return &m.pressedButtons, true
	default:
		return m.Item.Property(key)
	}
}

func (m *MouseArea) MustProperty(key string) vit.Value {
	v, ok := m.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (m *MouseArea) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "acceptedButtons":
		err = m.acceptedButtons.SetValue(value)
	case "containsMouse":
		err = m.containsMouse.SetValue(value)
	case "enabled":
		err = m.enabled.SetValue(value)
	case "hoverEnabled":
		err = m.hoverEnabled.SetValue(value)
	case "mouseX":
		err = m.mouseX.SetValue(value)
	case "mouseY":
		err = m.mouseY.SetValue(value)
	case "pressed":
		err = m.pressed.SetValue(value)
	case "pressedButtons":
		err = m.pressedButtons.SetValue(value)
	default:
		return m.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("MouseArea", key, m.id, err)
	}
	return nil
}

func (m *MouseArea) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "acceptedButtons":
		m.acceptedButtons.SetCode(code)
	case "containsMouse":
		m.containsMouse.SetCode(code)
	case "enabled":
		m.enabled.SetCode(code)
	case "hoverEnabled":
		m.hoverEnabled.SetCode(code)
	case "mouseX":
		m.mouseX.SetCode(code)
	case "mouseY":
		m.mouseY.SetCode(code)
	case "pressed":
		m.pressed.SetCode(code)
	case "pressedButtons":
		m.pressedButtons.SetCode(code)
	default:
		return m.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (m *MouseArea) Event(name string) (vit.Listenable, bool) {
	switch name {
	case "onClicked":
		return &m.onClicked, true
	default:
		return m.Item.Event(name)
	}
}

func (m *MouseArea) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "acceptedButtons":
		return &m.acceptedButtons, true
	case "containsMouse":
		return &m.containsMouse, true
	case "enabled":
		return &m.enabled, true
	case "hoverEnabled":
		return &m.hoverEnabled, true
	case "mouseX":
		return &m.mouseX, true
	case "mouseY":
		return &m.mouseY, true
	case "pressed":
		return &m.pressed, true
	case "pressedButtons":
		return &m.pressedButtons, true
	case "onClicked":
		return &m.onClicked, true
	default:
		return m.Item.ResolveVariable(key)
	}
}

func (m *MouseArea) AddChild(child vit.Component) {
	child.SetParent(m)
	m.AddChildButKeepParent(child)
}

func (m *MouseArea) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	var targetType vit.Component = afterThis

	for ind, child := range m.Children() {
		if child.As(&targetType) {
			addThis.SetParent(m)
			m.AddChildAtButKeepParent(addThis, ind+1)
			return
		}
	}
	m.AddChild(addThis)
}

func (m *MouseArea) UpdateExpressions() (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	// properties
	if changed, err := m.acceptedButtons.Update(m); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("MouseArea", "acceptedButtons", m.id, err))
		}
	}
	if changed, err := m.containsMouse.Update(m); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("MouseArea", "containsMouse", m.id, err))
		}
	}
	if changed, err := m.enabled.Update(m); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("MouseArea", "enabled", m.id, err))
		}
	}
	if changed, err := m.hoverEnabled.Update(m); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("MouseArea", "hoverEnabled", m.id, err))
		}
	}
	if changed, err := m.mouseX.Update(m); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("MouseArea", "mouseX", m.id, err))
		}
	}
	if changed, err := m.mouseY.Update(m); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("MouseArea", "mouseY", m.id, err))
		}
	}
	if changed, err := m.pressed.Update(m); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("MouseArea", "pressed", m.id, err))
		}
	}
	if changed, err := m.pressedButtons.Update(m); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("MouseArea", "pressedButtons", m.id, err))
		}
	}

	// methods

	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	n, err := m.UpdatePropertiesInContext(m)
	sum += n
	errs.AddGroup(err)
	n, err = m.Item.UpdateExpressions()
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (m *MouseArea) As(target *vit.Component) bool {
	if _, ok := (*target).(*MouseArea); ok {
		*target = m
		return true
	}
	return m.Item.As(target)
}

func (m *MouseArea) ID() string {
	return m.id
}

func (m *MouseArea) Finish() error {
	return m.RootC().FinishInContext(m)
}

func (m *MouseArea) staticAttribute(name string) (interface{}, bool) {
	switch name {
	case "noButton":
		return uint(MouseArea_MouseButtons_noButton), true
	case "leftButton":
		return uint(MouseArea_MouseButtons_leftButton), true
	case "rightButton":
		return uint(MouseArea_MouseButtons_rightButton), true
	case "middleButton":
		return uint(MouseArea_MouseButtons_middleButton), true
	case "allButtons":
		return uint(MouseArea_MouseButtons_allButtons), true
	default:
		return nil, false
	}
}
