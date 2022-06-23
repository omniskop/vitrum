package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/omniskop/vitrum/vit/generator"
)

var (
	inputPath   string
	outputPath  string
	packageName string
)

func main() {
	flag.StringVar(&inputPath, "i", "", "Path to the input vit file.")
	flag.StringVar(&outputPath, "o", "", "Path of the generated output file.")
	flag.StringVar(&packageName, "p", "", "Name of the package that the generated file will belong to.")

	flag.Parse()

	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "No input path specified.")
		os.Exit(64)
	}

	if outputPath == "" {
		err := generator.GenerateFromFile(inputPath, packageName, os.Stdout)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
			return
		}
	} else {
		err := generator.GenerateFromFileAndSave(inputPath, packageName, outputPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
			return
		}
	}
}
