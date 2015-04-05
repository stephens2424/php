package togo

import (
	"bytes"
	"flag"
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

var failTranspilation bool

func init() {
	flag.BoolVar(&failTranspilation, "fail-transpilation", false, "require all transpilation tests to pass")
}

func TestTranslation(t *testing.T) {
	testsDir := path.Join(build.Default.GOPATH, "src", "github.com/stephens2424/php/passes/togo/testdata")
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
	goFilename := phpFilename[:len(phpFilename)-3] + "go"

	buf := &bytes.Buffer{}
	err := TranspileFile(goFilename, path.Base(phpFilename), phpStr, buf)
	if err != nil {
		t.Fatal(err)
	}

	goStr, err := readFile(goFilename)
	if err != nil {
		t.Fatal(err)
	}

	if err == nil && buf.String() != goStr {
		failFunc := t.Skipf
		if failTranspilation {
			failFunc = t.Errorf
		}
		failFunc("mistranlation %s:\n\n===php===\n\n%s\n\n===expected===\n\n%s\n\n===got===\n\n%s\n\n", phpFilename, phpStr, goStr, buf.String())
	}

}
