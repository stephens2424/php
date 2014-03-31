package php

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestFiles(t *testing.T) {
	files, err := ioutil.ReadDir("testfiles")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		filename := file.Name()
		if strings.HasPrefix(filename, ".") || strings.HasSuffix(filename, ".php") {
			continue
		}
		src, err := ioutil.ReadFile(path.Join("testfiles", file.Name()))
		if err != nil {
			t.Error(err)
			continue
		}

		p := NewParser(string(src))
		_, errs := p.Parse()
		for _, err := range errs {
			t.Error(err)
		}
	}
}
