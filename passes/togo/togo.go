package togo

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"go/format"
	"go/token"
	"io"
	"strings"

	"github.com/stephens2424/php"
	phpast "github.com/stephens2424/php/ast"
	"golang.org/x/tools/imports"
)

type Togo struct {
	currentScope phpast.Scope
}

func TranspileFile(goFilename, phpFilename, phpStr string, gosrc io.Writer) error {
	file, errs := php.NewParser().Parse(phpFilename, phpStr)
	if len(errs) != 0 {
		return fmt.Errorf("found errors while parsing %s: %s", phpFilename, errs)
	}

	tg := Togo{}

	nodes := []goast.Node{}
	for _, phpNode := range file.Nodes {
		nodes = append(nodes, tg.ToGoStmt(phpNode.(phpast.Statement)))
	}

	buf := &bytes.Buffer{}

	err := format.Node(buf, token.NewFileSet(), File(phpFilename[:len(phpFilename)-4], nodes...))
	if err != nil {
		return fmt.Errorf("error while formatting %s: %s", phpFilename, err)
	}

	imported, err := imports.Process(goFilename, buf.Bytes(), &imports.Options{AllErrors: true, Comments: true, TabIndent: true, TabWidth: 8})
	if err != nil {
		return fmt.Errorf("error while getting imports for %s: %s", phpFilename, err)
	}

	_, err = gosrc.Write(imported)
	return err
}

func File(name string, nodes ...goast.Node) *goast.File {
	f := &goast.File{
		Name: goast.NewIdent("translated"),
	}

	stmts := []goast.Stmt{}

	name = strings.Replace(name, `.`, "", -1)
	name = strings.Replace(name, `/`, "_", -1)
	name = strings.Title(name)

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
			Name: &goast.Ident{Name: name},
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{{
						Names: []*goast.Ident{goast.NewIdent("ctx")},
						Type:  goast.NewIdent("phpctx.PHPContext"),
					}},
				},
			},
			Body: &goast.BlockStmt{
				List: stmts,
			},
		},
	}
	return f
}
