package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gedview"
	"os"
)

var path *string = flag.String("path", "", "path to GEDCOM file to inspect")

func main() {
	flag.Parse()

	file, err := os.Open(*path)
	if err != nil {
		panic(err)
		return
	}

	tree, err := gedview.CreateAST(file)
	if err != nil {
		panic(err)
		return
	}

	resp, err := json.MarshalIndent(tree, "", "\t")
	if err != nil {
		panic(err)
		return
	}

	fmt.Println(string(resp))
}
