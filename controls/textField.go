package controls

import (
	"unicode/utf8"

	vit "github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/std"
)

func (f *TextField) wasCompleted(*struct{}) {
	keyAreaComponent, _ := f.Context().GetComponentByID("keyArea")
	keyArea := keyAreaComponent.(*std.KeyArea)
	event, _ := keyArea.Event("onKeyDown")
	event.(*vit.EventAttribute[std.KeyEvent]).AddListener(vit.ListenerCB(f.keyPressed))

	mouseAreaComponent, _ := f.Context().GetComponentByID("mouseArea")
	mouseArea := mouseAreaComponent.(*std.MouseArea)
	event, _ = mouseArea.Event("onClicked")
	event.(*vit.EventAttribute[std.MouseEvent]).AddListener(vit.ListenerCB(f.clicked))
}

func (f *TextField) keyPressed(event *std.KeyEvent) {
	text := f.text.String()

	if event.Letter == 0 {
		switch event.Code {
		case "Backspace":
			if len(text) > 0 {
				_, size := utf8.DecodeLastRuneInString(text)
				text = text[:len(text)-size]
			}
		}
	} else {
		text += string(event.Letter)
	}

	f.text.SetValue(text)
}

func (f *TextField) clicked(event *std.MouseEvent) {
	f.Context().Global.Environment.RequestFocus(f)
}

func (f *TextField) Focus() {
	f.focused.SetValue(true)
}

func (f *TextField) Blur() {
	f.focused.SetValue(false)
}
