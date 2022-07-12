package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/omniskop/vitrum/controls"
	"github.com/omniskop/vitrum/gui"
	_ "github.com/omniskop/vitrum/vit/std"
)

func main() {
	app := gui.NewApplication()
	app.SetLogger(log.New(os.Stdout, "app: ", 0))
	app.AddImportPath("sources")

	_, err := app.NewWindow("sources/appTest.vit")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
