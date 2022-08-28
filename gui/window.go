package gui

import (
	"fmt"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/parse"
	"github.com/omniskop/vitrum/vit/std"
	"github.com/tdewolff/canvas"
	gioRenderer "github.com/tdewolff/canvas/renderers/gio"
)

type componentHandler struct {
	mouse            map[*std.MouseArea]bool
	key              map[*std.KeyArea]bool
	focusedComponent vit.FocusableComponent
	logger           *log.Logger
}

func newComponentHandler(log *log.Logger) *componentHandler {
	return &componentHandler{
		mouse:  make(map[*std.MouseArea]bool),
		key:    make(map[*std.KeyArea]bool),
		logger: log,
	}
}

func (h *componentHandler) RegisterComponent(id string, comp vit.Component) {
	switch comp := comp.(type) {
	case *std.MouseArea:
		h.mouse[comp] = true
	case *std.KeyArea:
		h.key[comp] = true
	}
}

func (h *componentHandler) UnregisterComponent(id string, comp vit.Component) {
	switch comp := comp.(type) {
	case *std.MouseArea:
		if _, ok := h.mouse[comp]; ok {
			delete(h.mouse, comp)
		}
	case *std.KeyArea:
		if _, ok := h.key[comp]; ok {
			delete(h.key, comp)
		}
	}
}

func (h *componentHandler) resetFocus() {
	if h.focusedComponent != nil {
		h.focusedComponent.Blur()
		h.focusedComponent = nil
	}
}

func (h *componentHandler) RequestFocus(comp vit.FocusableComponent) {
	if h.focusedComponent != nil {
		h.focusedComponent.Blur()
	}
	h.focusedComponent = comp
	comp.Focus()
}

func (h *componentHandler) TriggerMouseEvent(e pointer.Event, metric unit.Metric) {
	mouseEvent := std.MouseEvent{
		X: int(e.Position.X / metric.PxPerDp),
		Y: int(e.Position.Y / metric.PxPerDp),
	}
	if e.Buttons&pointer.ButtonPrimary > 0 {
		h.resetFocus()
		mouseEvent.Buttons |= std.MouseArea_MouseButtons_leftButton
	}
	if e.Buttons&pointer.ButtonSecondary > 0 {
		mouseEvent.Buttons |= std.MouseArea_MouseButtons_rightButton
	}
	if e.Buttons&pointer.ButtonTertiary > 0 {
		mouseEvent.Buttons |= std.MouseArea_MouseButtons_middleButton
	}
	for ma := range h.mouse {
		ma.TriggerEvent(mouseEvent)
	}
}

func (h *componentHandler) TriggerKeyEvent(e std.KeyEvent) {
	for ka := range h.key {
		ka.TriggerEvent(e)
	}
}

func (h *componentHandler) Logger() *log.Logger {
	return h.logger
}

type Window struct {
	manager       *parse.Manager
	handler       *componentHandler
	mainComponent vit.Component
	gioWindow     *app.Window
	logger        *log.Logger
}

func NewWindow(source string, log *log.Logger) (*Window, error) {
	w := &Window{
		manager: parse.NewManager(),
		handler: newComponentHandler(log),
		logger:  log,
	}
	err := w.manager.SetSource(source)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Window) AddImportPath(filePath string) error {
	return w.manager.AddImportPath(filePath)
}

func (w *Window) updateExpressions() vit.ErrorGroup {
evaluateExpressions:
	n, errs := w.mainComponent.UpdateExpressions(nil)
	if errs.Failed() {
		return errs
	}
	if n > 0 {
		goto evaluateExpressions
	}
	return vit.ErrorGroup{}
}

func (w *Window) prepare() error {
	err := w.manager.Run(w.handler)
	if err != nil {
		return err
	}
	w.mainComponent = w.manager.MainComponent()

	w.gioWindow = app.NewWindow(func(m unit.Metric, cfg *app.Config) {
		if v, ok := w.mainComponent.Property("title"); ok {
			title := v.GetValue()
			if titleString, ok := title.(string); ok {
				cfg.Title = titleString
			}
		}
		if v, ok := w.mainComponent.Property("width"); ok {
			width := v.GetValue()
			if widthFloat, ok := width.(float64); ok {
				cfg.Size.X = m.Dp(unit.Dp(widthFloat))
			}
		}
		if v, ok := w.mainComponent.Property("height"); ok {
			height := v.GetValue()
			if heightFloat, ok := height.(float64); ok {
				cfg.Size.Y = m.Dp(unit.Dp(heightFloat))
			}
		}
		if v, ok := w.mainComponent.Property("maxWidth"); ok {
			width := v.GetValue()
			if widthFloat, ok := width.(float64); ok {
				cfg.MaxSize.X = m.Dp(unit.Dp(widthFloat))
			}
		}
		if v, ok := w.mainComponent.Property("maxHeight"); ok {
			height := v.GetValue()
			if heightFloat, ok := height.(float64); ok {
				cfg.MaxSize.Y = m.Dp(unit.Dp(heightFloat))
			}
		}
		if v, ok := w.mainComponent.Property("minWidth"); ok {
			width := v.GetValue()
			if widthFloat, ok := width.(float64); ok {
				cfg.MinSize.X = m.Dp(unit.Dp(widthFloat))
			}
		}
		if v, ok := w.mainComponent.Property("minHeight"); ok {
			height := v.GetValue()
			if heightFloat, ok := height.(float64); ok {
				cfg.MinSize.Y = m.Dp(unit.Dp(heightFloat))
			}
		}
	})

	return nil
}

func (w *Window) run() error {
	if w.mainComponent == nil {
		return fmt.Errorf("window: main component is not set")
	}
	if w.gioWindow == nil {
		return fmt.Errorf("window: no underlying gio window set")
	}

	var ops op.Ops
	for {
		e := <-w.gioWindow.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			// handle user interaction
			var keysOfInterest = allSpecialKeys
			for _, ev := range e.Queue.Events(w) {
				switch event := ev.(type) {
				case pointer.Event:
					w.handler.TriggerMouseEvent(event, gtx.Metric)
				case key.Event:
					code, ok := keyCodeMapping[event.Name]
					if ok {
						// special character
						w.handler.TriggerKeyEvent(std.KeyEvent{
							Pressed: event.State == key.Press,
							Code:    string(code),
						})
					} else {
						// regular letter
						// here we are only interested in key releases as we have received the press as an edit event already
						if event.State == key.Release {
							r, _ := utf8.DecodeRuneInString(event.Name)
							w.handler.TriggerKeyEvent(std.KeyEvent{
								Pressed: false,
								Letter:  r,
							})
						}
					}
				case key.EditEvent:
					// a key was pressed
					r, _ := utf8.DecodeRuneInString(event.Text)
					// we wan't to be informed about changes to this specific key
					keysOfInterest = append(keysOfInterest, string(unicode.ToUpper(r)))
					w.handler.TriggerKeyEvent(std.KeyEvent{
						Pressed: true,
						Letter:  r,
					})
				default:
					// fmt.Printf("unknown event: %T %v\n", ev, ev)
				}
			}

			// render new frame
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// calculate window bounds
				physicalPixelWidth := float64(gtx.Constraints.Max.X) - float64(gtx.Constraints.Min.X)
				physicalPixelHeight := float64(gtx.Constraints.Max.Y) - float64(gtx.Constraints.Min.Y)
				// bounds of the window before scaling to the actual physical pixel size
				virtualWindowBounds := vit.NewRect(0, 0, physicalPixelWidth/float64(gtx.Metric.PxPerDp), physicalPixelHeight/float64(gtx.Metric.PxPerDp))

				// update dimensions of main component to reflect window size
				w.mainComponent.SetProperty("width", virtualWindowBounds.Width())
				w.mainComponent.SetProperty("height", virtualWindowBounds.Height())
				// update all expressions
				errs := w.updateExpressions()
				if errs.Failed() {
					w.logger.Println(fmt.Errorf("window update:"))
					w.logger.Println(parse.FormatError(errs))
				}

				// We give NewContain the virtual size of the window and it will calculate the necessary scaling factor
				// to fit the physical pixel size of the window.
				c := gioRenderer.NewContain(gtx, virtualWindowBounds.Width(), virtualWindowBounds.Height())
				ctx := canvas.NewContext(c)
				ctx.SetCoordSystem(canvas.CartesianIV) // move origin of the context to the top left corner
				err := w.mainComponent.Draw(
					vit.DrawingContext{ctx},
					virtualWindowBounds,
				)
				if err != nil {
					w.logger.Println(fmt.Errorf("window draw: %v", err))
				}
				return c.Dimensions()
			})

			// register input operations for the next frame
			pointer.InputOp{
				Tag:   w,
				Types: pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Scroll,
			}.Add(gtx.Ops)
			key.FocusOp{
				Tag: w,
			}.Add(gtx.Ops)
			key.InputOp{ // this also enables EditEvents
				Tag:  w,
				Hint: key.HintAny,
				Keys: key.Set(strings.Join(keysOfInterest, "|")),
			}.Add(gtx.Ops)

			e.Frame(gtx.Ops)
		}
	}
}

func (w *Window) SetVariable(name string, value interface{}) error {
	err := w.manager.SetVariable(name, value)
	if err != nil {
		return err
	}
	if w.gioWindow != nil {
		w.gioWindow.Invalidate()
	}
	return nil
}
