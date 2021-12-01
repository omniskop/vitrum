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

	component, err := parse.DoMagic("sources/test.vit")
	if err != nil {
		fmt.Println(err)
		return
	}

	// for _, component := range components {
	// 	fmt.Printf("%+v\n", component)
	// }

	fmt.Println("======================================")

	fmt.Println(component.Property("width"))
	fmt.Println(component.Property("stuff"))
	fmt.Println(component.Children()[0].Property("width"))
	fmt.Println(component.Children()[0].Children()[0].Property("width"))

	n, err := component.UpdateExpressions()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("main evaluted %d expressions\n", n)

	fmt.Println("======================================")

	component.SetProperty("width", "100")

	n, err = component.UpdateExpressions()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("main evaluted %d expressions\n", n)
}
