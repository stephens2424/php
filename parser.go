package php

import (
	"fmt"
	"stephensearles.com/php/ast"
)

type parser struct {
	lexer *lexer

	current item
	errors  []error
}

func newParser(input string) *parser {
	p := &parser{
		lexer: newLexer(input),
	}
	return p
}

func (p *parser) parse() ([]ast.Node, error) {
	// expecting either itemHTML or itemPHPBegin
	for i := range p.lexer.items {
		if i.typ == itemHTML {
			continue
		}
		if i.typ != itemPHPBegin {
			return nil, fmt.Errorf("Found %s, expected html or php begin", i)
		}
		if i.typ == itemEOF {
			break
		}
		p.parsePHP()
	}
	return nil, nil
}

func (p *parser) parsePHP() ([]ast.Node, error) {
	nodes := make([]ast.Node, 0, 1)
	for i := range p.lexer.items {
		if i.typ == itemPHPEnd {
			break
		}
		if i.typ == itemIf {
			subtree := p.parseIf()
			nodes = append(nodes, subtree)
		}
	}
	return nodes, nil
}

func (p *parser) next() {
	p.current = p.lexer.nextItem()
}

func (p *parser) expect(i itemType) {
	p.next()
	if p.current.typ != i {
		p.expected(i)
	}
}

func (p *parser) expected(i itemType) {
}

func (p *parser) errorf(str string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Errorf(str, args...))
}

func (p *parser) parseIf() ast.IfStmt {
	p.expect(itemOpenParen)
	n := ast.IfStmt{}
	n.Condition = p.parseExpression()
	p.expect(itemCloseParen)
	n.TrueBlock = p.parseStmt()
	n.FalseBlock = ast.Block{}
	p.next()
	if p.current.typ == itemElse {
		p.next()
		n.FalseBlock = p.parseStmt()
	}
	return n
}

func (p *parser) parseExpression() ast.Expression {
	return nil
}

func (p *parser) parseStmt() ast.Statement {
	return nil
}
