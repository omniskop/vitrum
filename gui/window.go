package gui

import (
	"fmt"
	"log"

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
	mouse map[*std.MouseArea]bool
}

func newComponentHandler() componentHandler {
	return componentHandler{
		mouse: make(map[*std.MouseArea]bool),
	}
}

func (h componentHandler) RegisterComponent(comp vit.Component) {
	if ma, ok := comp.(*std.MouseArea); ok {
		h.mouse[ma] = true
	}
}

func (h componentHandler) UnregisterComponent(comp vit.Component) {
	if ma, ok := comp.(*std.MouseArea); ok {
		if _, ok := h.mouse[ma]; ok {
			delete(h.mouse, ma)
		}
	}
}

func (h componentHandler) TriggerMouseEvent(e pointer.Event, metric unit.Metric) {
	mouseEvent := std.MouseEvent{
		X: int(e.Position.X / metric.PxPerDp),
		Y: int(e.Position.Y / metric.PxPerDp),
	}
	if e.Buttons&pointer.ButtonPrimary > 0 {
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

type Window struct {
	manager       *parse.Manager
	handler       componentHandler
	mainComponent vit.Component
	gioWindow     *app.Window
}

func NewWindow(source string) (*Window, error) {
	w := &Window{
		manager: parse.NewManager(),
		handler: newComponentHandler(),
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
	n, errs := w.mainComponent.UpdateExpressions()
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

func (w *Window) run(log *log.Logger) error {
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

			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				physicalPixelWidth := float64(gtx.Constraints.Max.X) - float64(gtx.Constraints.Min.X)
				physicalPixelHeight := float64(gtx.Constraints.Max.Y) - float64(gtx.Constraints.Min.Y)
				// bounds of the window before scaling to the actual physical pixel size
				virtualWindowBounds := vit.NewRect(0, 0, physicalPixelWidth/float64(gtx.Metric.PxPerDp), physicalPixelHeight/float64(gtx.Metric.PxPerDp))

				w.mainComponent.SetProperty("width", virtualWindowBounds.Width())
				w.mainComponent.SetProperty("height", virtualWindowBounds.Height())
				err := w.updateExpressions()
				if err.Failed() {
					log.Println(fmt.Errorf("window update: %v", err))
				}

				// We give NewContain the virtual size of the window and it will calculate the necessary scaling factor
				// to fit the physical pixel size of the window.
				c := gioRenderer.NewContain(gtx, virtualWindowBounds.Width(), virtualWindowBounds.Height())
				ctx := canvas.NewContext(c)
				ctx.SetCoordSystem(canvas.CartesianIV) // move origin of the context to the top left corner
				w.mainComponent.Draw(
					vit.DrawingContext{ctx},
					virtualWindowBounds,
				)
				return c.Dimensions()
			})

			for _, ev := range e.Queue.Events(w) {
				if x, ok := ev.(pointer.Event); ok {
					w.handler.TriggerMouseEvent(x, gtx.Metric)
				}
			}

			pointer.InputOp{
				Tag:   w,
				Types: pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Scroll,
			}.Add(gtx.Ops)
			key.InputOp{
				Tag: w,
			}.Add(gtx.Ops)

			e.Frame(gtx.Ops)
		}
	}
}