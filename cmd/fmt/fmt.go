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
		file, err := php.NewParser().Parse("test.php", string(src))
		if err != nil {
			log.Fatal(err)
		}
		for _, node := range file.Nodes {
			p.PrintNode(node)
		}
	}
}
