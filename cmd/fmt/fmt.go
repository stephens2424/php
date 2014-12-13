package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/stephens2424/php"
	"github.com/stephens2424/php/ast/printer"
)

func main() {
	flag.Parse()
	for _, arg := range flag.Args() {

		fmt.Println(arg)
		fmt.Println()

		src, err := ioutil.ReadFile(arg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		p := printer.NewPrinter(os.Stdout)
		nodes, errs := php.NewParser(string(src)).Parse()
		if len(errs) != 0 {
			log.Fatal(errs)
		}
		for _, node := range nodes {
			p.PrintNode(node)
		}
	}
}
