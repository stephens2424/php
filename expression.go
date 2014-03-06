package php

import "stephensearles.com/php/ast"

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
	itemArrayLookupOperatorLeft: 19,
	itemUnaryOperator:           18,
	itemCastOperator:            18,
	itemInstanceofOperator:      17,
	itemNegationOperator:        16,
	itemMultOperator:            15,
	itemAdditionOperator:        14,
	itemSubtractionOperator:     14,
	itemConcatenationOperator:   14,

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
	case itemIgnoreErrorOperator:
		p.next()
		return p.parseExpression()
	case itemFunction:
		return p.parseAnonymousFunction()
	case itemNewOperator:
		return p.parseInstantiation()
	case itemVariableOperator:
		if p.peek().typ == itemAssignmentOperator {
			assignee := p.parseIdentifier().(ast.Assignable)
			p.next()
			return ast.AssignmentExpression{
				Assignee: assignee,
				Operator: p.current.val,
				Value:    p.parseNextExpression(),
			}
		}
		fallthrough
	case itemArray:
		fallthrough
	case itemUnaryOperator, itemNegationOperator, itemAmpersandOperator, itemCastOperator, itemSubtractionOperator:
		fallthrough
	case itemIdentifier, itemStringLiteral, itemNumberLiteral, itemBooleanLiteral, itemInclude, itemNull, itemSelf, itemStatic, itemParent:
		expr = p.parseOperation(originalParenLev, p.expressionize())
	case itemOpenParen:
		p.parenLevel += 1
		p.next()
		expr = p.parseExpression()
		p.expect(itemCloseParen)
		p.parenLevel -= 1
		expr = p.parseOperation(originalParenLev, expr)
	default:
		p.errorf("Expected expression. Found %s", p.current)
		return
	}
	if p.parenLevel != originalParenLev {
		p.errorf("unbalanced parens: %d prev: %d", p.parenLevel, originalParenLev)
		return
	}
	return
}

func (p *parser) parseOperation(originalParenLevel int, lhs ast.Expression) (expr ast.Expression) {
	p.next()
	switch p.current.typ {
	case itemIgnoreErrorOperator:
		return p.parseOperation(originalParenLevel, lhs)
	case itemUnaryOperator:
		expr = p.parseUnaryExpressionLeft(lhs, p.current)
	case itemAdditionOperator, itemSubtractionOperator, itemConcatenationOperator, itemComparisonOperator, itemMultOperator, itemAndOperator, itemOrOperator, itemAmpersandOperator, itemBitwiseXorOperator, itemBitwiseOrOperator, itemBitwiseShiftOperator, itemWrittenAndOperator, itemWrittenXorOperator, itemWrittenOrOperator, itemInstanceofOperator:
		expr = p.parseBinaryOperation(lhs, p.current, originalParenLevel)
	case itemTernaryOperator1:
		expr = p.parseTernaryOperation(lhs)
	case itemCloseParen:
		if p.parenLevel <= originalParenLevel {
			p.backup()
			return lhs
		}
		p.parenLevel -= 1
		return p.parseOperation(originalParenLevel, lhs)
	case itemAssignmentOperator:
		assignee := lhs.(ast.Assignable)
		expr = ast.AssignmentExpression{
			Assignee: assignee,
			Operator: p.current.val,
			Value:    p.parseNextExpression(),
		}
	default:
		p.backup()
		return lhs
	}
	return p.parseOperation(originalParenLevel, expr)
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

func newBinaryOperation(operator Item, expr1, expr2 ast.Expression) ast.OperatorExpression {
	t := ast.Numeric
	switch operator.typ {
	case itemComparisonOperator, itemAndOperator, itemOrOperator, itemWrittenAndOperator, itemWrittenOrOperator, itemWrittenXorOperator:
		t = ast.Boolean
	case itemConcatenationOperator:
		t = ast.String
	case itemAmpersandOperator, itemBitwiseXorOperator, itemBitwiseOrOperator, itemBitwiseShiftOperator:
		t = ast.AnyType
	}
	return ast.OperatorExpression{
		Type:     t,
		Operand1: expr1,
		Operand2: expr2,
		Operator: operator.val,
	}
}

func (p *parser) parseBinaryOperation(lhs ast.Expression, operator Item, originalParenLevel int) ast.Expression {
	p.next()
	rhs := p.expressionize()
	for {
		nextOperator := p.peek()
		nextOperatorPrecedence, ok := operatorPrecedence[nextOperator.typ]
		if ok && nextOperatorPrecedence > operatorPrecedence[operator.typ] {
			rhs = p.parseOperation(originalParenLevel, rhs)
		} else {
			break
		}
	}
	return newBinaryOperation(operator, lhs, rhs)
}

func (p *parser) parseTernaryOperation(lhs ast.Expression) ast.Expression {
	truthy := p.parseNextExpression()
	p.expect(itemTernaryOperator2)
	falsy := p.parseNextExpression()
	return &ast.OperatorExpression{
		Operand1: lhs,
		Operand2: truthy,
		Operand3: falsy,
		Type:     truthy.EvaluatesTo() | falsy.EvaluatesTo(),
		Operator: "?:",
	}
}

func (p *parser) parseUnaryExpressionRight(operand ast.Expression, operator Item) ast.Expression {
	return newUnaryOperation(operator, operand)
}

func (p *parser) parseUnaryExpressionLeft(operand ast.Expression, operator Item) ast.Expression {
	return newUnaryOperation(operator, operand)
}

// expressionize takes the current token and returns it as the simplest
// expression for that token. That means an expression with no operators
// except for the object operator.
func (p *parser) expressionize() ast.Expression {
	switch p.current.typ {
	case itemIgnoreErrorOperator:
		p.next()
		return p.expressionize()
	case itemUnaryOperator, itemNegationOperator, itemCastOperator, itemSubtractionOperator, itemAmpersandOperator:
		op := p.current
		p.next()
		return p.parseUnaryExpressionRight(p.expressionize(), op)
	case itemArray:
		return p.parseArrayDeclaration()
	case itemVariableOperator:
		return p.parseIdentifier()
	case itemStringLiteral, itemBooleanLiteral, itemNumberLiteral, itemNull:
		return p.parseLiteral()
	case itemIdentifier, itemSelf, itemStatic, itemParent:
		if p.peek().typ == itemOpenParen {
			var expr ast.Expression
			expr = p.parseFunctionCall()
			for p.peek().typ == itemObjectOperator {
				expr = p.parseObjectLookup(expr)
			}
			return expr
		}
		if p.peek().typ == itemScopeResolutionOperator {
			r := p.current.val
			p.expect(itemScopeResolutionOperator)
			return &ast.ClassExpression{
				Receiver:   r,
				Expression: p.parseNextExpression(),
			}
		}
		return ast.ConstantExpression{
			Variable: ast.NewVariable(p.current.val),
		}
	case itemOpenParen:
		return p.parseExpression()
	case itemInclude:
		return ast.Include{Expression: p.parseNextExpression()}
	}
	// error?
	return nil
}

func (p *parser) parseLiteral() *ast.Literal {
	switch p.current.typ {
	case itemStringLiteral:
		return &ast.Literal{Type: ast.String}
	case itemBooleanLiteral:
		return &ast.Literal{Type: ast.Boolean}
	case itemNumberLiteral:
		return &ast.Literal{Type: ast.Float}
	case itemNull:
		return &ast.Literal{Type: ast.Null}
	}
	p.errorf("Unknown literal type")
	return nil
}

func (p *parser) parseVariable() ast.Expression {
	p.expectCurrent(itemVariableOperator)
	switch p.next(); {
	case isKeyword(p.current.typ):
		// keywords are all valid variable names
		fallthrough
	case p.current.typ == itemIdentifier:
		expr := ast.NewVariable("$" + p.current.val)
		return expr
	default:
		return p.parseExpression()
	}
}

func (p *parser) parseIdentifier() (expr ast.Expression) {
	expr = p.parseVariable()
	switch pk := p.peek(); pk.typ {
	case itemScopeResolutionOperator:
		r := "$" + p.current.val
		p.expect(itemScopeResolutionOperator)
		p.next()
		return &ast.ClassExpression{
			Receiver:   r,
			Expression: p.expressionize(),
		}
	case itemObjectOperator:
		for p.peek().typ == itemObjectOperator {
			expr = p.parseObjectLookup(expr)
		}
	case itemArrayLookupOperatorLeft:
		return p.parseArrayLookup(expr)
	case itemOpenParen:
		var expr ast.Expression
		expr = p.parseFunctionArguments(&ast.FunctionCallExpression{
			FunctionName: expr,
		})
		if p.peek().typ == itemObjectOperator {
			expr = p.parseObjectLookup(expr)
		}
	}
	return expr
}

func (p *parser) parseAnonymousFunction() ast.Expression {
	f := &ast.AnonymousFunction{}
	f.Arguments = make([]ast.FunctionArgument, 0)
	p.expect(itemOpenParen)
	if p.peek().typ == itemCloseParen {
		p.expect(itemCloseParen)
		return f
	}
	f.Arguments = append(f.Arguments, p.parseFunctionArgument())
	for {
		switch p.peek().typ {
		case itemComma:
			p.expect(itemComma)
			f.Arguments = append(f.Arguments, p.parseFunctionArgument())
		case itemCloseParen:
			p.expect(itemCloseParen)
			return f
		default:
			p.errorf("unexpected argument separator:", p.current)
			return f
		}
	}
	f.Body = p.parseBlock()
	return f
}
