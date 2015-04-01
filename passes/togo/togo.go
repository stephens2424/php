package togo

import (
	phpast "github.com/stephens2424/php/ast"
)

type Togo struct {
	currentScope phpast.Scope
}
