package parser

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestFiles(t *testing.T) {
	files, err := ioutil.ReadDir("../testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		filename := file.Name()
		if strings.HasPrefix(filename, ".") || !strings.HasSuffix(filename, ".php") {
			continue
		}
		if filename != "badfile.php" {
			continue
		}
		src, err := ioutil.ReadFile(path.Join("../testdata", filename))
		if err != nil {
			t.Error(err)
			continue
		}

		p := NewParser()
		p.PrintTokens = true
		if _, err = p.Parse("test.php", string(src)); err != nil {
			t.Error(filename, err)
		}
	}
}
