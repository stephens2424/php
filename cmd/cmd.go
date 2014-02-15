package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"stephensearles.com/php"
	"stephensearles.com/php/passes/typechecking"
)

func main() {
	flag.Parse()
	fBytes, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	parser := php.NewParser(string(fBytes))
	nodes := parser.Parse()
	walker := typecheck.Walker{}
	for _, node := range nodes {
		walker.Walk(node)
	}
}
