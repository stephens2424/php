package php

import (
	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/token"
)

func (p *Parser) parseArrayLookup(e ast.Expression) ast.Expression {
	p.expectCurrent(token.ArrayLookupOperatorLeft, token.BlockBegin)
	switch typ := p.peek().typ; typ {
	case token.ArrayLookupOperatorRight, token.BlockBegin:
		p.expect(token.ArrayLookupOperatorRight, token.BlockEnd)
		return ast.ArrayAppendExpression{Array: e}
	}
	p.next()
	expr := &ast.ArrayLookupExpression{
		Array: e,
		Index: p.parseExpression(),
	}
	p.expect(token.ArrayLookupOperatorRight, token.BlockEnd)
	return expr
}

func (p *Parser) parseArrayDeclaration() ast.Expression {
	var endType token.Token
	pairs := make([]ast.ArrayPair, 0)
	p.expectCurrent(token.Array, token.ArrayLookupOperatorLeft)
	switch p.current.typ {
	case token.Array:
		p.expect(token.OpenParen)
		endType = token.CloseParen
	case token.ArrayLookupOperatorLeft:
		endType = token.ArrayLookupOperatorRight
	}
ArrayLoop:
	for {
		var key, val ast.Expression
		switch p.peek().typ {
		case endType:
			break ArrayLoop
		default:
			val = p.parseNextExpression()
		}
		switch p.peek().typ {
		case token.Comma:
			p.expect(token.Comma)
		case endType:
			pairs = append(pairs, ast.ArrayPair{Key: key, Value: val})
			break ArrayLoop
		case token.ArrayKeyOperator:
			p.expect(token.ArrayKeyOperator)
			key = val
			val = p.parseNextExpression()
			if p.peek().typ == endType {
				pairs = append(pairs, ast.ArrayPair{Key: key, Value: val})
				break ArrayLoop
			}
			p.expect(token.Comma)
		default:
			p.errorf("expected => or ,")
			return nil
		}
		pairs = append(pairs, ast.ArrayPair{Key: key, Value: val})
	}
	p.expect(endType)
	return &ast.ArrayExpression{Pairs: pairs}
}

func (p *Parser) parseList() ast.Expression {
	l := &ast.ListStatement{
		Assignees: make([]ast.Assignable, 0),
	}
	p.expect(token.OpenParen)
	for {
		if p.accept(token.Comma) {
			continue
		}
		if p.peek().typ == token.CloseParen {
			break
		}
		p.next()
		op, ok := p.parseOperand().(ast.Assignable)
		if ok {
			l.Assignees = append(l.Assignees, op)
		} else {
			p.errorf("%v list element is not assignable", op)
		}
		if p.peek().typ != token.Comma {
			break
		}
		p.expect(token.Comma)
	}
	p.expect(token.CloseParen)
	p.expect(token.AssignmentOperator)
	l.Operator = p.current.val
	l.Value = p.parseNextExpression()
	return l

}
