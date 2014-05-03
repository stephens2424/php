package php

import (
	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
)

func (p *Parser) parseBlock() *ast.Block {
	p.expect(token.BlockBegin)
	b := p.parseStatementsUntil(token.BlockEnd)
	p.expectCurrent(token.BlockEnd)
	return b
}

func (p *Parser) parseStatementsUntil(endTokens ...token.Token) *ast.Block {
	block := &ast.Block{}
	breakTypes := map[token.Token]bool{}
	for _, typ := range endTokens {
		breakTypes[typ] = true
	}
	for {
		p.next()
		if _, ok := breakTypes[p.current.typ]; ok {
			break
		}
		stmt := p.parseStmt()
		if stmt == nil {
			return block
		}
		block.Statements = append(block.Statements, stmt)
	}
	return block
}

func (p *Parser) parseExpressionsUntil(separator token.Token, endTokens ...token.Token) []ast.Expression {
	exprs := make([]ast.Expression, 0, 1)
	breakTypes := map[token.Token]bool{}
	for _, typ := range endTokens {
		breakTypes[typ] = true
	}
	p.next()
	first := true
	for {
		if _, ok := breakTypes[p.current.typ]; ok {
			break
		} else if first {
			first = false
		} else {
			p.expectCurrent(separator)
			p.next()
		}
		expr := p.parseExpression()
		if expr == nil {
			return exprs
		}
		exprs = append(exprs, expr)
		p.next()
	}
	return exprs
}
