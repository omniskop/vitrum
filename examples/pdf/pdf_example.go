package main

import (
	"fmt"
	"log"
	"os"

	"github.com/omniskop/vitrum/pdf"
	"github.com/omniskop/vitrum/vit/parse"
	"github.com/omniskop/vitrum/vit/vpath"

	_ "github.com/omniskop/vitrum/controls"
	_ "github.com/omniskop/vitrum/vit/std"
)

func main() {
	doc, err := pdf.NewDocument(vpath.Local("examples/pdf/document.vit"))
	if err != nil {
		fmt.Println(err)
		return
	}
	doc.SetLogger(log.New(os.Stdout, "app: ", 0))

	f, err := os.Create("output.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = doc.Render(f)
	if err != nil {
		fmt.Println(parse.FormatError(err))
		return
	}

	fmt.Println("saved to output.pdf")
}
