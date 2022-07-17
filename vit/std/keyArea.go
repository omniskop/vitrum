package std

type KeyEvent struct {
	Pressed bool
	Letter  rune
	Code    string
}

func (a *KeyArea) enableDisable() {
	if !a.enabled.Bool() {
		// KeyArea was just disabled
	}
}

func (a *KeyArea) TriggerEvent(e KeyEvent) {
	if !a.enabled.Bool() {
		return
	}
	if e.Pressed {
		a.onKeyDown.Fire(&e)
	} else {
		a.onKeyUp.Fire(&e)
	}
}
