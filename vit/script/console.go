package script

import (
	"fmt"
	"strings"
)

func stringify(args []interface{}) []string {
	var out []string
	for _, arg := range args {
		// TODO: It would probably a good idea to improve the conversion to string here.
		//       To make it more robust for different types.
		out = append(out, fmt.Sprintf("%v", arg))
	}
	return out
}

// The console object that is globally defined in JavaScript.
// This currently only implements a subset of all methods.
var globalConsole = map[string]interface{}{
	"log": func(args ...interface{}) {
		fmt.Println(fmt.Sprintf("JS log: %s", strings.Join(stringify(args), " ")))
	},
	"info": func(args ...interface{}) {
		fmt.Println(fmt.Sprintf("JS info: %s", strings.Join(stringify(args), " ")))
	},
	"warn": func(args ...interface{}) {
		fmt.Println(fmt.Sprintf("JS warn: %s", strings.Join(stringify(args), " ")))
	},
	"error": func(args ...interface{}) {
		fmt.Println(fmt.Sprintf("JS error: %s", strings.Join(stringify(args), " ")))
	},
}
