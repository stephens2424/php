package togo

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"go/format"
	"go/token"
	"testing"

	"github.com/stephens2424/php"
)

func TestIf(t *testing.T) {
	phpNodes, errs := php.NewParser(`<?php
	if ("hello" == "world") {
	}`).Parse()
	if len(errs) != 0 {
		t.Fatalf("found errors during parsing:", errs)
	}
	buf := &bytes.Buffer{}

	nodes := []goast.Node{}
	for _, phpNode := range phpNodes {
		nodes = append(nodes, ToGo(phpNode))
	}

	err := format.Node(buf, token.NewFileSet(), File(nodes...))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(buf.String())
}

func File(nodes ...goast.Node) *goast.File {
	f := &goast.File{
		Name: goast.NewIdent("translated"),
	}

	stmts := []goast.Stmt{}

	for _, n := range nodes {
		switch g := n.(type) {
		case goast.Stmt:
			stmts = append(stmts, g)
		case goast.Expr:
			stmts = append(stmts, &goast.ExprStmt{g})
		}

	}

	f.Decls = []goast.Decl{
		&goast.FuncDecl{
			Name: &goast.Ident{Name: "main"},
			Type: &goast.FuncType{
				Params: &goast.FieldList{},
			},
			Body: &goast.BlockStmt{
				List: stmts,
			},
		},
	}
	return f
}
