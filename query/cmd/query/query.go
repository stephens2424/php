// Command query implements a simple CLI for querying a PHP AST.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/parser"
	"github.com/stephens2424/php/query"
)

func main() {
	recursive := flag.Bool("r", false, "Recursive")
	flag.Parse()
	selector := strings.Join(flag.Args(), " ")
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	g := newGatherer(*recursive)
	if err := filepath.Walk(dir, g.walkFile); err != nil {
		panic(err)
	}

	selected, _ := query.Select(g.nodes).Select(selector)

	/*
		for _, sel := range selected {
			pos := sel.Node.Begin()
			fmt.Println(pos)
		}
	*/

	fmt.Println(len(selected), "found")
}

func newGatherer(recursive bool) gatherer {
	return gatherer{recursive, make([]ast.Node, 0)}
}

type gatherer struct {
	recursive bool
	nodes     []ast.Node
}

func (g *gatherer) walkFile(path string, info os.FileInfo, err error) error {
	// TODO handle error
	_ = err

	if info.IsDir() || !strings.HasSuffix(path, ".php") {
		return nil
	}

	if info.IsDir() && !g.recursive {
		return filepath.SkipDir
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	src, err := ioutil.ReadAll(f)

	if err != nil {
		return err
	}

	p := parser.NewParser()
	file, err := p.Parse("test.php", string(src))
	if err != nil {
		return err
	}

	g.nodes = append(g.nodes, file.Nodes...)
	return nil
}
