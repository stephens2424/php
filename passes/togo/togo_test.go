package togo

import (
	"bytes"
	goast "go/ast"
	"go/build"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stephens2424/php"
)

func TestTranslation(t *testing.T) {
	testsDir := path.Join(build.Default.GOPATH, "src", "github.com/stephens2424/php/passes/togo/tests")
	phpFiles, err := filepath.Glob(testsDir + "/*.php")
	if err != nil {
		t.Fatal(err)
	}

	for _, phpFile := range phpFiles {
		phpStr, err := readFile(phpFile)
		if err != nil {
			t.Fatal("couldn't read file", phpFile, err)
		}
		parseFile(t, phpFile, phpStr)
	}
}

func readFile(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(f)
	return string(b), err
}

func parseFile(t *testing.T, phpFilename, phpStr string) {
	file, errs := php.NewParser().Parse(phpFilename, phpStr)
	if len(errs) != 0 {
		t.Errorf("found errors while parsing %s: %s", phpFilename, errs)
		return
	}

	buf := &bytes.Buffer{}

	nodes := []goast.Node{}
	for _, phpNode := range file.Nodes {
		nodes = append(nodes, ToGo(phpNode))
	}

	err := format.Node(buf, token.NewFileSet(), File(nodes...))
	if err != nil {
		t.Errorf("error while formatting %s: %s", phpFilename, err)
		return
	}

	goStr, err := readFile(phpFilename[:len(phpFilename)-3] + "go")
	if err == nil && buf.String() != goStr {
		t.Errorf("mistranlation:\n\n===php===\n\n%s\n\n===expected===\n\n%s\n\n===got===\n\n%s\n\n", phpStr, goStr, buf.String())
	}
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
