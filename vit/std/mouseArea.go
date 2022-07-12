package std

import "fmt"

type MouseEvent struct {
	X, Y    int
	Buttons MouseArea_MouseButtons
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
		fmt.Println("fired")
		m.onClicked.Fire(&e)
	}
}
