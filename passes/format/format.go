package format

import (
	"io"
	"strings"

	"github.com/stephens2424/php"
)

type Formatter struct {
	Indent   string
	tabLevel int
	w        io.Writer
}

func NewFormatter(w io.Writer) *Formatter {
	return &Formatter{
		w:      w,
		Indent: "\t",
	}
}

func (f *Formatter) Format(src string) error {
	walker := &formatWalker{f}
	p := php.NewParser(src)
	nodes, errs := p.Parse()
	if len(errs) > 0 {
		return errs[0]
	}
	for _, node := range nodes {
		err := walker.Walk(node)
		if err != nil {
			return err
		}
	}
	walker.print("\n")
	return nil
}

func (f *Formatter) tab() string {
	return strings.Repeat(f.Indent, f.tabLevel)
}

type formatWalker struct {
	*Formatter
}
