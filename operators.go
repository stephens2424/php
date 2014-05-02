package php

import (
	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
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
	case token.UnaryOperator, token.BitwiseNotOperator:
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

func newUnaryOperation(operator Item, expr ast.Expression) ast.OperatorExpression {
	t := ast.Numeric
	if operator.val == "!" {
		t = ast.Boolean
	}
	return ast.OperatorExpression{
		Type:     t,
		Operand1: expr,
		Operator: operator.val,
	}
}

func (p *Parser) newBinaryOperation(operator Item, expr1, expr2 ast.Expression) ast.Expression {
	t := ast.Numeric
	switch operator.typ {
	case token.AssignmentOperator:
		return p.parseAssignmentOperation(expr1, expr2, operator)
	case token.ComparisonOperator, token.AndOperator, token.OrOperator, token.WrittenAndOperator, token.WrittenOrOperator, token.WrittenXorOperator:
		t = ast.Boolean
	case token.ConcatenationOperator:
		t = ast.String
	case token.AmpersandOperator, token.BitwiseXorOperator, token.BitwiseOrOperator, token.BitwiseShiftOperator:
		t = ast.AnyType
	}
	return ast.OperatorExpression{
		Type:     t,
		Operand1: expr1,
		Operand2: expr2,
		Operator: operator.val,
	}
}

func (p *Parser) parseBinaryOperation(lhs ast.Expression, operator Item, originalParenLevel int) ast.Expression {
	p.next()
	rhs := p.parseOperand()
	currentPrecedence := operatorPrecedence[operator.typ]
	for {
		nextOperator := p.peek()
		nextPrecedence, ok := operatorPrecedence[nextOperator.typ]
		if !ok || nextPrecedence < currentPrecedence {
			break
		}
		rhs = p.parseOperation(originalParenLevel, rhs)
	}
	return p.newBinaryOperation(operator, lhs, rhs)
}

func (p *Parser) parseTernaryOperation(lhs ast.Expression) ast.Expression {
	var truthy ast.Expression
	if p.peek().typ == token.TernaryOperator2 {
		truthy = lhs
	} else {
		truthy = p.parseNextExpression()
	}
	p.expect(token.TernaryOperator2)
	falsy := p.parseNextExpression()
	return &ast.OperatorExpression{
		Operand1: lhs,
		Operand2: truthy,
		Operand3: falsy,
		Type:     truthy.EvaluatesTo() | falsy.EvaluatesTo(),
		Operator: "?:",
	}
}

func (p *Parser) parseUnaryExpressionRight(operand ast.Expression, operator Item) ast.Expression {
	return newUnaryOperation(operator, operand)
}

func (p *Parser) parseUnaryExpressionLeft(operand ast.Expression, operator Item) ast.Expression {
	return newUnaryOperation(operator, operand)
}
