package php

import (
	"stephensearles.com/php/ast"
)

/*

Valid Expression Patterns
Expr [Binary Op] Expr
[Unary Op] Expr
Expr [Unary Op]
Expr [Tertiary Op 1] Expr [Tertiary Op 2] Expr
Identifier
Literal
Function Call

Parentesis always triggers sub-expression

non-associative clone new clone and new
left  [ array()
right ++ -- ~ (int) (float) (string) (array) (object) (bool) @  types and increment/decrement
non-associative instanceof  types
right ! logical
left  * / % arithmetic
left  + - . arithmetic and string
left  << >> bitwise
non-associative < <= > >= comparison
non-associative == != === !== <>  comparison
left  & bitwise and references
left  ^ bitwise
left  | bitwise
left  &&  logical
left  ||  logical
left  ? : ternary
right = += -= *= /= .= %= &= |= ^= <<= >>= => assignment
left  and logical
left  xor logical
left  or  logical
left  , many uses

*/

var operatorPrecedence = map[ItemType]int{
	itemArrayAccessOperator:   19,
	itemUnaryOperator:         18,
	itemCastOperator:          18,
	itemInstanceofOperator:    17,
	itemNegationOperator:      16,
	itemMultOperator:          15,
	itemAdditionOperator:      14,
	itemSubtractionOperator:   14,
	itemConcatenationOperator: 14,

	itemBitwiseShiftOperator: 13,
	itemComparisonOperator:   12,
	itemEqualityOperator:     11,

	itemAmpersandOperator:  10,
	itemBitwiseXorOperator: 9,
	itemBitwiseOrOperator:  8,
	itemAndOperator:        7,
	itemOrOperator:         6,
	itemTernaryOperator1:   5,
	itemTernaryOperator2:   5,
	itemAssignmentOperator: 4,
	itemWrittenAndOperator: 3,
	itemWrittenXorOperator: 2,
	itemWrittenOrOperator:  1,
}

func (p *parser) parseExpression() (expr ast.Expression) {
	// consume expression
	originalParenLev := p.parenLevel
	switch p.current.typ {
	case itemUnaryOperator:
		op := p.current
		expr = p.parseUnaryExpressionRight(p.parseNextExpression(), op)
	case itemIdentifier, itemNonVariableIdentifier, itemStringLiteral, itemNumberLiteral, itemBooleanLiteral:
		expr = p.parseOperation(p.expressionize())
	case itemOpenParen:
		p.parenLevel += 1
		p.next()
		expr = p.parseExpression()
	default:
		p.errorf("Expected expression")
		return nil
	}
	// if the last item was a close paren, that close wasn't part of the expression
	// reset the parenLevel to not include it and backup the parser
	if p.current.typ == itemCloseParen {
		p.parenLevel += 1
		p.backup()
	}
	if p.parenLevel != originalParenLev {
		p.errorf("unbalanced parens: %d", p.parenLevel)
		return nil
	}
	return
}

func (p *parser) parseOperation(lhs ast.Expression) (expr ast.Expression) {
	p.next()
	switch p.current.typ {
	case itemUnaryOperator:
		expr = p.parseUnaryExpressionLeft(lhs, p.current)
	case itemAdditionOperator, itemSubtractionOperator, itemConcatenationOperator, itemComparisonOperator, itemMultOperator, itemAndOperator, itemOrOperator, itemAmpersandOperator, itemBitwiseXorOperator, itemBitwiseOrOperator, itemBitwiseShiftOperator, itemWrittenAndOperator, itemWrittenXorOperator, itemWrittenOrOperator:
		expr = p.parseBinaryOperation(lhs, p.current)
	case itemCloseParen:
		p.parenLevel -= 1
		return p.parseOperation(lhs)
	default:
		p.backup()
		return lhs
	}
	return p.parseOperation(expr)
}

func (p *parser) parseBinaryOperation(lhs ast.Expression, operator Item) ast.Expression {
	p.next()
	rhs := p.expressionize()
	for {
		p.next()
		nextOperator := p.current
		p.backup()
		nextOperatorPrecedence, ok := operatorPrecedence[nextOperator.typ]
		if ok && nextOperatorPrecedence > operatorPrecedence[operator.typ] {
			rhs = p.parseOperation(rhs)
		} else {
			break
		}
	}
	return newBinaryOperation(operator, lhs, rhs)
}

func (p *parser) parseUnaryExpressionRight(operand ast.Expression, operator Item) ast.Expression {
	return newUnaryOperation(operator, operand)
}

func (p *parser) parseUnaryExpressionLeft(operand ast.Expression, operator Item) ast.Expression {
	return newUnaryOperation(operator, operand)
}

func (p *parser) expressionize() ast.Expression {
	switch p.current.typ {
	case itemIdentifier:
		return ast.NewIdentifier(p.current.val)
	case itemStringLiteral:
		return ast.Literal{ast.String}
	case itemBooleanLiteral:
		return ast.Literal{ast.Boolean}
	case itemNumberLiteral:
		return ast.Literal{ast.Float}
	case itemNonVariableIdentifier:
		return p.parseFunctionCall()
	case itemOpenParen:
		return p.parseExpression()
	}
	// error?
	return nil
}
