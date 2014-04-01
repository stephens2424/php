package php

import (
	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
)

func (p *parser) parseArrayLookup(e ast.Expression) ast.Expression {
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

func (p *parser) parseArrayDeclaration() ast.Expression {
	pairs := make([]ast.ArrayPair, 0)
	p.expect(token.OpenParen)
ArrayLoop:
	for {
		var key, val ast.Expression
		switch p.peek().typ {
		case token.CloseParen:
			break ArrayLoop
		default:
			val = p.parseNextExpression()
		}
		switch p.peek().typ {
		case token.Comma:
			p.expect(token.Comma)
		case token.CloseParen:
			pairs = append(pairs, ast.ArrayPair{Key: key, Value: val})
			break ArrayLoop
		case token.ArrayKeyOperator:
			p.expect(token.ArrayKeyOperator)
			key = val
			val = p.parseNextExpression()
			if p.peek().typ == token.CloseParen {
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
	p.expect(token.CloseParen)
	return &ast.ArrayExpression{Pairs: pairs}
}
