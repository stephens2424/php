package php

import (
	"testing"

	"github.com/stephens2424/php/token"
)

var src = `<?php

echo "one";

class Two {

  public $three = "four";

  public function five() {
    return function () {
      return 6;
    };
  }
}
`

func TestLines(t *testing.T) {
	p := NewParser(src)
	nodes, errs := p.Parse()
	if len(errs) > 0 {
		for _, err := range errs {
			t.Log(err)
		}
		t.FailNow()
	}

	echo := nodes[0]
	assertPosition(t, echo.Begin(), token.Position{3, 0, 7, ""})
}

func assertPosition(t *testing.T, actual, expected token.Position) {
	if actual != expected {
		t.Errorf("Found %+v, expected %+v\n", actual, expected)
	}
}
