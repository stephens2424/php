package printer

import (
	"bytes"
	"testing"

	"github.com/stephens2424/php"
)

type Test struct {
	Before, After string
}

func TestPrinter(t *testing.T) {
	t.SkipNow() // disabling this test before some ast restructuring

	for _, test := range tests {
		p := php.NewParser()
		file, err := p.Parse("test.php", test.Before)
		if err != nil {
			t.Error("parsing error:", err)
			continue
		}

		if len(file.Nodes) < 1 {
			continue
		}

		buf := &bytes.Buffer{}
		pr := NewPrinter(buf)
		pr.PrintNode(file.Nodes[0])

		if buf.String() != test.After {
			t.Fatalf("formatted text did not match\nFormatted\n\n%s\n\nExpected\n\n%s\n", buf.String(), test.After)
		}
	}
}

var tests = []Test{
	{
		Before: ``,
		After:  ``,
	},
	{
		Before: `<?php $var = "x"; `,
		After: `<?php
$var = "x";
`,
	},
}
