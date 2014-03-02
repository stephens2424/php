package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"stephensearles.com/php"
	//"stephensearles.com/php/passes/typechecking"
	"stephensearles.com/php/passes/printing"
)

func main() {
	flag.Parse()
	fmt.Println(flag.Arg(0))
	fBytes, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	walker := printing.Walker{}
	parser := php.NewParser(string(fBytes))
	nodes, errs := parser.Parse()
	if len(errs) != 0 {
		if len(nodes) != 0 && nodes[0] != nil {
			walker.Walk(nodes[0])
		}
		for _, err := range errs {
			fmt.Println(err)
		}
	}
	/*
		  nodes := parser.Parse()
				walker := typecheck.Walker{}
				for _, node := range nodes {
					walker.Walk(node)
				}
	*/
}
