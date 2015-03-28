package php

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestFiles(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		filename := file.Name()
		if strings.HasPrefix(filename, ".") || !strings.HasSuffix(filename, ".php") {
			continue
		}
		src, err := ioutil.ReadFile(path.Join("testdata", filename))
		if err != nil {
			t.Error(err)
			continue
		}

		p := NewParser()
		_, errs := p.Parse("test.php", string(src))
		for _, err := range errs {
			t.Error(filename, err)
		}
	}
}
