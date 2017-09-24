package ast

import (
	"fmt"
)

// Walker is a file walker
type Walker interface {
	Walk(n Node)
	Errorf(fmt string, params ...interface{})
}

// DefaultWalker is the default walker implementation
type DefaultWalker struct {
	Nodes  []Node
	Errors []error
}

func (d *DefaultWalker) Errorf(format string, params ...interface{}) {
	fmt.Printf(format+"\n", params...)
	d.Errors = append(d.Errors, fmt.Errorf(format, params...))
}
