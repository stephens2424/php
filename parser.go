package php

import (
	"fmt"

	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/lexer"
	"github.com/stephens2424/php/token"
)

// Parser handles scanning through and parsing a PHP source string into an AST. It is configurable
// to have various types of debugging features.
type Parser struct {
	Debug       bool // Debug causes the parser to print all errors to stdout and relay any panic upon internal panic recovery.
	PrintTokens bool // PrintTokens causes the parser to print all tokens received from the lexer to stdout.
	MaxErrors   int  // Indicates the number of errors to allow before triggering a panic. The default is 10.
	FileSet     ast.FileSet

	lexer      token.Stream
	previous   []token.Item
	idx        int
	current    token.Item
	errors     []error
	parenLevel int
	errorMap   map[int]bool
	errorCount int

	namespace ast.Namespace
	scope     ast.Scope

	instantiation bool
}

// NewParser readies a parser
func NewParser() *Parser {
	p := &Parser{
		idx:       -1,
		MaxErrors: 10,
		errorMap:  make(map[int]bool),
	}
	return p
}

// Parse consumes the input string to produce an AST that represents it.
func (p *Parser) Parse(filepath, input string) (*ast.File, []error) {
	var errors []error
	file := &ast.File{}
	p.lexer = lexer.NewLexer(input)
	defer func() {
		if r := recover(); r != nil {
			errors = append([]error{fmt.Errorf("%s", r)}, p.errors...)
			if p.Debug {
				for _, err := range p.errors {
					fmt.Println(err)
				}
				panic(r)
			}
		}
	}()
	// expecting either token.HTML or token.PHPBegin
	file.Nodes = make([]ast.Node, 0, 1)
TokenLoop:
	for {
		p.next()
		switch p.current.Typ {
		case token.EOF:
			break TokenLoop
		default:
			n := p.parseNode()
			if n != nil {
				file.Nodes = append(file.Nodes, n)
			}
		}
	}
	errors = p.errors
	return file, errors
}

func (p *Parser) parseNode() ast.Node {
	switch p.current.Typ {
	case token.HTML:
		return ast.Echo(ast.Literal{Type: ast.String, Value: p.current.Val})
	case token.PHPBegin:
		return nil
	case token.PHPEnd:
		return nil
	}
	return p.parseStmt()
}

func (p *Parser) next() {
	p.idx += 1
	if len(p.previous) <= p.idx {
		p.current = p.lexer.Next()
		if p.PrintTokens {
			fmt.Println(p.current)
		}
		p.previous = append(p.previous, p.current)
	} else {
		p.current = p.previous[p.idx]
	}
}

func (p *Parser) backup() {
	p.idx -= 1
	p.current = p.previous[p.idx]
}

func (p *Parser) peek() (i token.Item) {
	p.next()
	i = p.current
	p.backup()
	return
}

func (p *Parser) expectCurrent(i ...token.Token) {
	for _, Typ := range i {
		if p.current.Typ == Typ {
			return
		}
	}
	p.expected(i...)
}

func (p *Parser) expectAndNext(i ...token.Token) {
	defer p.next()
	for _, Typ := range i {
		if p.current.Typ == Typ {
			return
		}
	}
	p.expected(i...)
}

func (p *Parser) expect(i ...token.Token) {
	p.next()
	p.expectCurrent(i...)
}

func (p *Parser) expected(i ...token.Token) {
	p.errorf("Found %s, expected %s", p.current, i)
}

func (p *Parser) accept(i ...token.Token) bool {
	nextTyp := p.peek().Typ
	for _, Typ := range i {
		if nextTyp == Typ {
			p.next()
			return true
		}
	}
	return false
}

func (p *Parser) errorf(str string, args ...interface{}) {
	if p.errorCount > p.MaxErrors {
		panic("too many errors")
	}
	if _, ok := p.errorMap[p.current.Begin.Line]; ok {
		return
	}
	errString := fmt.Sprintf(str, args...)
	p.errorCount += 1
	p.errors = append(p.errors, fmt.Errorf("%s: %s", p.errorPrefix(), errString))
	p.errorMap[p.current.Begin.Line] = true
}

func (p *Parser) errorPrefix() string {
	return fmt.Sprintf("%d", p.current.Begin.Line)
}

func (p *Parser) parseNextExpression() ast.Expression {
	p.next()
	return p.parseExpression()
}
