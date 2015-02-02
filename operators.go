package php

import (
	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/token"
)

type operationType int

const (
	nilOperation operationType = 1 << iota
	unaryOperation
	binaryOperation
	ternaryOperation
	assignmentOperation
	subexpressionBeginOperation
	subexpressionEndOperation
	ignoreErrorOperation
)

func operationTypeForToken(t token.Token) operationType {
	switch t {
	case token.IgnoreErrorOperator:
		return ignoreErrorOperation
	case token.UnaryOperator,
		token.NegationOperator,
		token.CastOperator,
		token.BitwiseNotOperator:
		return unaryOperation
	case token.AdditionOperator,
		token.SubtractionOperator,
		token.ConcatenationOperator,
		token.ComparisonOperator,
		token.MultOperator,
		token.AndOperator,
		token.OrOperator,
		token.AmpersandOperator,
		token.BitwiseXorOperator,
		token.BitwiseOrOperator,
		token.BitwiseShiftOperator,
		token.WrittenAndOperator,
		token.WrittenXorOperator,
		token.WrittenOrOperator,
		token.InstanceofOperator:
		return binaryOperation
	case token.TernaryOperator1:
		return ternaryOperation
	case token.AssignmentOperator:
		return assignmentOperation
	case token.OpenParen:
		return subexpressionBeginOperation
	case token.CloseParen:
		return subexpressionEndOperation
	}
	return nilOperation
}

func (p *Parser) newBinaryOperation(operator token.Item, expr1, expr2 ast.Expression) ast.Expression {
	t := ast.Numeric
	switch operator.Typ {
	case token.AssignmentOperator:
		return p.parseAssignmentOperation(expr1, expr2, operator)
	case token.ComparisonOperator, token.AndOperator, token.OrOperator, token.WrittenAndOperator, token.WrittenOrOperator, token.WrittenXorOperator:
		t = ast.Boolean
	case token.ConcatenationOperator:
		t = ast.String
	case token.AmpersandOperator, token.BitwiseXorOperator, token.BitwiseOrOperator, token.BitwiseShiftOperator:
		t = ast.AnyType
	}
	return ast.BinaryExpression{
		Type:       t,
		Antecedent: expr1,
		Subsequent: expr2,
		Operator:   operator.Val,
	}
}

func (p *Parser) parseBinaryOperation(lhs ast.Expression, operator token.Item, originalParenLevel int) ast.Expression {
	p.next()
	rhs := p.parseOperand()
	currentPrecedence := operatorPrecedence[operator.Typ]
	for {
		nextOperator := p.peek()
		nextPrecedence, ok := operatorPrecedence[nextOperator.Typ]
		if !ok || nextPrecedence < currentPrecedence {
			break
		}
		rhs = p.parseOperation(originalParenLevel, rhs)
	}
	return p.newBinaryOperation(operator, lhs, rhs)
}

func (p *Parser) parseTernaryOperation(lhs ast.Expression) ast.Expression {
	var truthy ast.Expression
	if p.peek().Typ == token.TernaryOperator2 {
		truthy = lhs
	} else {
		truthy = p.parseNextExpression()
	}
	p.expect(token.TernaryOperator2)
	falsy := p.parseNextExpression()
	return &ast.TernaryExpression{
		Condition: lhs,
		True:      truthy,
		False:     falsy,
		Type:      truthy.EvaluatesTo() | falsy.EvaluatesTo(),
	}
}

func (p *Parser) parseUnaryExpressionRight(operand ast.Expression, operator token.Item) ast.Expression {
	return ast.UnaryExpression{
		Operand:  operand,
		Operator: operator.Val,
	}
}

func (p *Parser) parseUnaryExpressionLeft(operand ast.Expression, operator token.Item) ast.Expression {
	return ast.UnaryExpression{
		Operand:   operand,
		Operator:  operator.Val,
		Preceding: true,
	}
}
