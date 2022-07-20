package gui

import (
	"fmt"
	"io"
	"log"
	"sync"

	"gioui.org/app"
	"github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/parse"
)

type Application struct {
	importPaths []string
	windowsMux  sync.RWMutex
	windows     []*Window
	log         *log.Logger
}

func NewApplication() *Application {
	return &Application{
		log: log.New(io.Discard, "", 0),
	}
}

func (a *Application) AddImportPath(path string) {
	a.importPaths = append(a.importPaths, path)
}

func (a *Application) SetLogger(log *log.Logger) {
	a.log = log
}

func (a *Application) NewWindow(path string) (*Window, error) {
	w, err := NewWindow(path, a.log)
	if err != nil {
		return nil, err
	}
	var errs vit.ErrorGroup
	for _, path := range a.importPaths {
		err := w.AddImportPath(path)
		if err != nil {
			errs.Add(err)
		}
	}
	a.AddWindow(w)
	if errs.Failed() {
		return w, errs
	}
	return w, nil
}

func (a *Application) AddWindow(w *Window) {
	a.windowsMux.Lock()
	a.windows = append(a.windows, w)
	a.windowsMux.Unlock()
}

func (a *Application) RemoveWindow(w *Window) {
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

func (a *Application) runWindow(w *Window) {
	// TODO: figure out how we should handle errors here
	err := w.prepare()
	if err != nil {
		str := parse.FormatError(err)
		a.log.Println(fmt.Errorf("window preparation:"))
		a.log.Println(str)
		return
	}
	err = w.run()
	if err != nil {
		a.log.Println(fmt.Errorf("window failed: %v", err))
	}
	a.RemoveWindow(w)
}
