package std

import "fmt"

type MouseEvent struct {
	X, Y    int
	Buttons MouseArea_MouseButtons
}

func (e *MouseEvent) MaybeSet(input interface{}) error {
	if input == nil {
		return fmt.Errorf("value is nil")
	}
	switch input := input.(type) {
	case *MouseEvent:
		e.X = input.X
		e.Y = input.Y
		e.Buttons = input.Buttons
	case MouseEvent:
		e.X = input.X
		e.Y = input.Y
		e.Buttons = input.Buttons
	case map[string]interface{}:
		if x, ok := input["x"]; ok {
			if x, ok := x.(int); ok {
				e.X = x
			}
		}
		if y, ok := input["y"]; ok {
			if y, ok := y.(int); ok {
				e.Y = y
			}
		}
		if buttons, ok := input["buttons"]; ok {
			if buttons, ok := buttons.(int); ok {
				e.Buttons = MouseArea_MouseButtons(buttons)
			}
		}
	default:
		return fmt.Errorf("value of type %T can't be converted to MouseEvent", input)
	}
	return nil
}

func (m *MouseArea) enableDisable() {
	if !m.enabled.Bool() {
		// MouseArea was just disabled
		m.containsMouse.SetBoolValue(false)
		m.pressed.SetBoolValue(false)
		m.pressedButtons.SetIntValue(0)
	}
}

// TriggerEvent will be called by vitrum when a mouse event is received
func (m *MouseArea) TriggerEvent(e MouseEvent) {
	if !m.enabled.Bool() {
		return
	}
	if !m.Bounds().Contains(float64(e.X), float64(e.Y)) {
		m.containsMouse.SetBoolValue(false)
		return
	}
	m.containsMouse.SetBoolValue(true)
	m.mouseX.SetFloatValue(float64(e.X))
	m.mouseY.SetFloatValue(float64(e.Y))

	filtered := e.Buttons & MouseArea_MouseButtons(m.acceptedButtons.Int())
	wasPressed := m.pressed.Bool()
	m.pressed.SetBoolValue(filtered > 0)

	m.pressedButtons.SetIntValue(int(e.Buttons))

	if filtered == 0 && wasPressed {
		m.onClicked.Fire(&e)
	}
}
