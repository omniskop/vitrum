package main

import (
	"fmt"
	"log"
	"os"

	"github.com/omniskop/vitrum/gui"
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
