package printing

import (
	"fmt"
	"io"
	"os"
	"strings"

	"stephensearles.com/php/ast"
)

type Walker struct {
	tabLevel int
	ast.DefaultWalker
	W io.Writer
}

func NewWalker(w io.Writer) *Walker {
	return &Walker{W: os.Stdout}
}

func (w *Walker) Walk(node ast.Node) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			fmt.Printf("%T - %+v\n", node, node)
			panic(r)
		}
	}()
	fmt.Fprintf(w.W, "%s(%T)%s\n", strings.Repeat("\t", w.tabLevel), node, node.String())
	switch children := node.Children(); children {
	case nil:
	default:
		w.tabLevel += 1
		for _, child := range children {
			if child != nil {
				w.Walk(child)
			}
		}
		w.tabLevel -= 1
	}
}
