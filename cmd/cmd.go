package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"

	"stephensearles.com/php"
	//"stephensearles.com/php/passes/typechecking"
	"stephensearles.com/php/passes/printing"
)

func main() {
	astonerror := flag.Bool("astonerror", false, "Print the AST on errors")
	ast := flag.Bool("ast", false, "Print the AST")
	showErrors := flag.Bool("showerrors", true, "show errors. If this is false, astonerror will be ignored")
	debugMode := flag.Bool("debug", false, "if true, panic on finding any error")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	verbose := flag.Bool("verbose", false, "print all filenames")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var files, errors int
	for _, filename := range flag.Args() {
		if *verbose {
			fmt.Println(filename)
		}
		files += 1
		fBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			continue
		}
		walker := printing.NewWalker()
		Parser := php.NewParser(string(fBytes))
		if *debugMode {
			Parser.Debug = true
			Parser.MaxErrors = 0
			Parser.PrintTokens = true
		}
		nodes, errs := Parser.Parse()
		if *ast && len(nodes) != 0 && nodes[0] != nil {
			for _, node := range nodes {
				walker.Walk(node)
			}
		}
		if len(errs) != 0 {
			errors += 1
			if *showErrors {
				if !*verbose {
					fmt.Println(filename)
				}
				if !*ast && *astonerror && len(nodes) != 0 && nodes[0] != nil {
					for _, node := range nodes {
						walker.Walk(node)
					}
				}
				for _, err := range errs {
					fmt.Println(err)
				}
			}
		}
	}
	fmt.Printf("Compiled %d files. %d files with errors - %f%% success\n", flag.NArg(), errors, 100*(1-(float64(errors)/float64(files))))
}
