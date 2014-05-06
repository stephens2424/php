package php

import (
	"fmt"

	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
)

// Parser handles scanning through and parsing a PHP source string into an AST. It is configurable
// to have various types of debugging features.
type Parser struct {
	Debug       bool // Debug causes the parser to print all errors to stdout and relay any panic upon internal panic recovery.
	PrintTokens bool // PrintTokens causes the parser to print all tokens received from the lexer to stdout.
	MaxErrors   int  // Indicates the number of errors to allow before triggering a panic. The default is 10.

	lexer      *lexer
	previous   []Item
	idx        int
	current    Item
	errors     []error
	parenLevel int
	errorMap   map[int]bool
	errorCount int

	instantiation bool
}

// NewParser readies a parser object for the given input string.
func NewParser(input string) *Parser {
	p := &Parser{
		idx:       -1,
		MaxErrors: 10,
		lexer:     newLexer(input),
		errorMap:  make(map[int]bool),
	}
	return p
}

// Parse consumes the input string to produce an AST that represents it.
func (p *Parser) Parse() (nodes []ast.Node, errors []error) {
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
	nodes = make([]ast.Node, 0, 1)
TokenLoop:
	for {
		p.next()
		switch p.current.typ {
		case token.EOF:
			break TokenLoop
		default:
			n := p.parseNode()
			if n != nil {
				nodes = append(nodes, n)
			}
		}
	}
	errors = p.errors
	return nodes, errors
}

func (p *Parser) parseNode() ast.Node {
	switch p.current.typ {
	case token.HTML:
		return ast.Echo(ast.Literal{Type: ast.String, Value: p.current.val})
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
		p.current = p.lexer.nextItem()
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

func (p *Parser) peek() (i Item) {
	p.next()
	i = p.current
	p.backup()
	return
}

func (p *Parser) expectCurrent(i ...token.Token) {
	for _, typ := range i {
		if p.current.typ == typ {
			return
		}
	}
	p.expected(i...)
}

func (p *Parser) expectAndNext(i ...token.Token) {
	defer p.next()
	for _, typ := range i {
		if p.current.typ == typ {
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
	nextTyp := p.peek().typ
	for _, typ := range i {
		if nextTyp == typ {
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
	if _, ok := p.errorMap[p.current.pos.Line]; ok {
		return
	}
	errString := fmt.Sprintf(str, args...)
	p.errorCount += 1
	p.errors = append(p.errors, fmt.Errorf("%s: %s", p.errorPrefix(), errString))
	p.errorMap[p.current.pos.Line] = true
}

func (p *Parser) errorPrefix() string {
	return fmt.Sprintf("%s %d", p.lexer.file, p.current.pos.Line)
}

func (p *Parser) parseNextExpression() ast.Expression {
	p.next()
	return p.parseExpression()
}
