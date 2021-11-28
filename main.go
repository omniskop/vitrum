package main

import (
	"fmt"
	"os"

	"github.com/omniskop/vitrum/vit/parse"
)

// //go: generate go run generate.go -i Item.txt -o Item.go

func main() {
	// err := parse.Open("test.vit")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	file, err := os.Open("test.vit")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	lexer := parse.NewLexer(file, "test.vit")

	vitDocument, err := parse.Parse(parse.NewTokenBuffer(lexer.Lex))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(vitDocument)
}
