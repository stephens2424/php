package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stephens2424/php/passes/format"
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
		f := format.NewFormatter(os.Stdout)
		err = f.Format(string(src))
		if err != nil {
			fmt.Println(err)
		}
	}
}
