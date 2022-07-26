package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/omniskop/vitrum/controls"
	"github.com/omniskop/vitrum/gui"
	_ "github.com/omniskop/vitrum/vit/std"
)

func main() {
	app := gui.NewApplication()
	app.SetLogger(log.New(os.Stdout, "app: ", 0))
	app.AddImportPath("sources")

	window, err := app.NewWindow("sources/appTest.vit")
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

	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
