package main

import (
	"fmt"

	"github.com/omniskop/vitrum/vit/parse"
	"github.com/omniskop/vitrum/vit/script"
)

// //go: generate go run generate.go -i Item.txt -o Item.go

func main() {
	// err := parse.Open("test.vit")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	script.Setup()

	mngr := parse.NewManager()

	mngr.AddImportPath("sources")
	err := mngr.SetSource("sources/test.vit")
	// err := mngr.SetSource("firefly/LaunchWindow.vit")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = mngr.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("======================================")

	component := mngr.MainComponent()
	_ = component

	// fmt.Println(component.Children()[0].MustProperty("horizontalAlignment"))
	// c, _ := component.Children()[0].Property("color")
	// fmt.Println(c.(color.Color).RGBA())
}
