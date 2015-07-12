package php

import (
	"bytes"
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
	FileSet     *ast.FileSet

	lexer      token.Stream
	previous   []token.Item
	idx        int
	current    token.Item
	errors     ParseErrorList
	parenLevel int

	file      *ast.File
	namespace *ast.Namespace
	scope     *ast.Scope

	// this option exists to allow parser tests to pass while scope tests may be failing
	disableScoping bool

	instantiation bool
}

// NewParser readies a parser
func NewParser() *Parser {
	p := &Parser{
		idx:       -1,
		MaxErrors: 10,
		FileSet:   ast.NewFileSet(),
	}
	return p
}

// ParseErrorList is a list of ParseErrors.
type ParseErrorList []ParseError

// Error formats p into a string.
func (p ParseErrorList) Error() string {
	if len(p) == 0 {
		return ""
	}
	if len(p) == 1 {
		return p[0].Error()
	}
	buf := &bytes.Buffer{}
	for _, s := range p[:len(p)-1] {
		buf.WriteString(s.Error())
		buf.WriteString("\n")
	}
	buf.WriteString(p[len(p)-2].Error())
	return buf.String()
}

// ParseError represents an error found during parsing.
type ParseError struct {
	error
	Line, Column int
	File         *ast.File
}

func (p ParseError) Error() string {
	return fmt.Sprintf("%s:%d: %s", p.File.Name, p.Line, p.error)
}

func (p ParseError) String() string {
	return p.Error()
}

// Parse consumes the input string to produce an AST that represents it.
func (p *Parser) Parse(filepath, input string) (file *ast.File, err error) {
	file = &ast.File{Namespace: p.FileSet.GlobalNamespace}
	p.file = file
	p.scope = p.FileSet.Scope
	p.namespace = p.FileSet.GlobalNamespace
	p.lexer = token.Subset(lexer.NewLexer(input), token.Significant)

	p.FileSet.Files[filepath] = p.file
	defer func() {
		if r := recover(); r != nil {
			err = append(ParseErrorList{errorf(p, "%s", r)}, p.errors...)
			if p.Debug {
				for _, err := range p.errors {
					fmt.Println(err)
				}
				panic(r)
			}
		}
	}()
	// expecting either token.HTML or token.PHPBegin
	p.file.Nodes = make([]ast.Node, 0, 1)
TokenLoop:
	for {
		p.next()
		switch p.current.Typ {
		case token.EOF:
			break TokenLoop
		default:
			n := p.parseNode()
			if n != nil {
				p.file.Nodes = append(p.file.Nodes, n)
			}
		}
	}
	if p.errors != nil {
		err = p.errors
	}
	return
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
	p.idx++
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
	p.idx--
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
	if len(p.errors) > p.MaxErrors {
		panic("too many errors")
	}
	p.errors = append(p.errors, errorf(p, str, args...))
}

func errorf(p *Parser, str string, args ...interface{}) ParseError {
	e := ParseError{error: fmt.Errorf(str, args...)}
	if p != nil {
		e.File = p.file
		e.Line = 0
	}
	return e
}

func (p *Parser) errorPrefix() string {
	return fmt.Sprintf("%d", p.current.Begin.Line)
}

func (p *Parser) parseNextExpression() ast.Expression {
	p.next()
	return p.parseExpression()
}
