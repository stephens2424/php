package printing

import (
	"fmt"
	"strings"

	"stephensearles.com/php/ast"
)

type Walker struct {
	tabLevel int
	ast.DefaultWalker
}

func (w *Walker) Walk(node ast.Node) {
	fmt.Printf("%s(%T)%s\n", strings.Repeat("\t", w.tabLevel), node, node.String())
	switch children := node.Children(); children {
	case nil:
	default:
		w.tabLevel += 1
		for _, child := range children {
			w.Walk(child)
		}
		w.tabLevel -= 1
	}
}
