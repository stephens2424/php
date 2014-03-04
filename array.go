package php

import "stephensearles.com/php/ast"

func (p *parser) parseArrayLookup(e ast.Expression) ast.Expression {
	p.expect(itemArrayLookupOperatorLeft)
	if p.peek().typ == itemArrayLookupOperatorRight {
		p.expect(itemArrayLookupOperatorRight)
		return ast.ArrayAppendExpression{Array: e}
	}
	p.next()
	expr := &ast.ArrayLookupExpression{
		Array: e,
		Index: p.parseExpression(),
	}
	p.expect(itemArrayLookupOperatorRight)
	switch p.peek().typ {
	case itemArrayLookupOperatorLeft:
		return p.parseArrayLookup(expr)
	case itemObjectOperator:
		return p.parseObjectLookup(expr)
	}
	return expr
}

func (p *parser) parseArrayDeclaration() ast.Expression {
	pairs := make([]ast.ArrayPair, 0)
	p.expect(itemOpenParen)
ArrayLoop:
	for {
		var key, val ast.Expression
		switch p.peek().typ {
		case itemCloseParen:
			break ArrayLoop
		default:
			val = p.parseNextExpression()
		}
		switch p.peek().typ {
		case itemComma:
			p.expect(itemComma)
		case itemCloseParen:
			pairs = append(pairs, ast.ArrayPair{Key: key, Value: val})
			break ArrayLoop
		case itemArrayKeyOperator:
			p.expect(itemArrayKeyOperator)
			key = val
			val = p.parseNextExpression()
			if p.peek().typ == itemCloseParen {
				pairs = append(pairs, ast.ArrayPair{Key: key, Value: val})
				break ArrayLoop
			}
			p.expect(itemComma)
		default:
			p.errorf("expected => or ,")
			return nil
		}
		pairs = append(pairs, ast.ArrayPair{Key: key, Value: val})
	}
	p.expect(itemCloseParen)
	return &ast.ArrayExpression{Pairs: pairs}
}
