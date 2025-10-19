package main

import (
	"fmt"
	"sort/service/parser"
)

func main() {
	config, err := parser.ParseArgs()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config)
}
