package gui

import (
	"fmt"
	"io"
	"log"
	"sync"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/omniskop/vitrum/vit"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/gio"
)

type Application struct {
	windowsMux sync.RWMutex
	windows    []*AppWindow
	log        *log.Logger
}

func NewApplication() *Application {
	return &Application{
		log: log.New(io.Discard, "", 0),
	}
}

func (a *Application) SetLogger(log *log.Logger) {
	a.log = log
}

func (a *Application) AddWindow(w *AppWindow) {
	a.windowsMux.Lock()
	a.windows = append(a.windows, w)
	a.windowsMux.Unlock()
}

func (a *Application) RemoveWindow(w *AppWindow) {
	a.windowsMux.Lock()
	for i, w2 := range a.windows {
		if w2 == w {
			a.windows = append(a.windows[:i], a.windows[i+1:]...)
			break
		}
	}
	a.windowsMux.Unlock()
}

func (a *Application) Run() error {
	a.windowsMux.RLock()
	for _, w := range a.windows {
		go a.runWindow(w)
	}
	a.windowsMux.RUnlock()
	app.Main()
	return nil
}

func (a *Application) runWindow(w *AppWindow) {
	err := w.run(a.log)
	if err != nil {
		a.log.Println(fmt.Errorf("window failed: %v", err))
	}
	a.RemoveWindow(w)
}

type AppWindow struct {
	mainComponent vit.Component
	window        *app.Window
	canvas        *canvas.Canvas
}

func NewAppWindow(comp vit.Component) *AppWindow {
	w := app.NewWindow(func(m unit.Metric, cfg *app.Config) {
		if v, ok := comp.Property("title"); ok {
			title := v.GetValue()
			if titleString, ok := title.(string); ok {
				cfg.Title = titleString
			}
		}
		if v, ok := comp.Property("width"); ok {
			width := v.GetValue()
			if widthFloat, ok := width.(float64); ok {
				cfg.Size.X = m.Dp(unit.Dp(widthFloat))
			}
		}
		if v, ok := comp.Property("height"); ok {
			height := v.GetValue()
			if heightFloat, ok := height.(float64); ok {
				cfg.Size.Y = m.Dp(unit.Dp(heightFloat))
			}
		}
		if v, ok := comp.Property("maxWidth"); ok {
			width := v.GetValue()
			if widthFloat, ok := width.(float64); ok {
				cfg.MaxSize.X = m.Dp(unit.Dp(widthFloat))
			}
		}
		if v, ok := comp.Property("maxHeight"); ok {
			height := v.GetValue()
			if heightFloat, ok := height.(float64); ok {
				cfg.MaxSize.Y = m.Dp(unit.Dp(heightFloat))
			}
		}
		if v, ok := comp.Property("minWidth"); ok {
			width := v.GetValue()
			if widthFloat, ok := width.(float64); ok {
				cfg.MinSize.X = m.Dp(unit.Dp(widthFloat))
			}
		}
		if v, ok := comp.Property("minHeight"); ok {
			height := v.GetValue()
			if heightFloat, ok := height.(float64); ok {
				cfg.MinSize.Y = m.Dp(unit.Dp(heightFloat))
			}
		}
	})
	return &AppWindow{
		mainComponent: comp,
		window:        w,
	}
}

func (w *AppWindow) updateExpressions() vit.ErrorGroup {
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

func (w *AppWindow) run(log *log.Logger) error {
	var ops op.Ops
	for {
		e := <-w.window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				windowBounds := vit.Rect{float64(gtx.Constraints.Min.X), float64(gtx.Constraints.Min.Y), float64(gtx.Constraints.Max.X), float64(gtx.Constraints.Max.Y)}

				//

				w.mainComponent.SetProperty("width", windowBounds.Width())
				w.mainComponent.SetProperty("height", windowBounds.Height())
				err := w.updateExpressions()
				if err.Failed() {
					log.Println(fmt.Errorf("window update: %v", err))
				}

				//

				c := gio.New(gtx, windowBounds.Width(), windowBounds.Height())
				ctx := canvas.NewContext(c)
				// move origin of the context to the top left corner
				ctx.SetCoordSystem(canvas.CartesianIV)
				w.mainComponent.Draw(
					vit.DrawingContext{ctx},
					windowBounds,
				)
				return c.Dimensions()
			})

			e.Frame(gtx.Ops)
		}
	}
}
