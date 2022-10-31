package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/omniskop/vitrum/gui"
	"github.com/omniskop/vitrum/vit/vpath"

	_ "github.com/omniskop/vitrum/controls"
	_ "github.com/omniskop/vitrum/vit/std"
)

//go:embed sources
var sources embed.FS

func main() {
	app := gui.NewApplication()
	app.SetLogger(log.New(os.Stdout, "app: ", 0))
	app.AddImportPath(vpath.FS(sources, "sources"))

	window, err := app.NewWindow(vpath.FS(sources, "sources/appTest.vit"))
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		tck := time.NewTicker(time.Second)
		for range tck.C {
			window.SetVariable("globalText", time.Now().Format("15:04:05"))
		}
	}()

	window.SetVariable("globalCallback", func(v int) int {
		fmt.Println("callback here!")
		return v + 1
	})

	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
