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
	itemBitwiseNotOperator:      18,
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
		expr = p.parseInstantiation()
		expr = p.parseOperation(originalParenLev, expr)
		return

	case itemList:
		l := &ast.ListStatement{
			Assignees: make([]*ast.Variable, 0),
		}
		p.expect(itemOpenParen)
		for {
			p.expect(itemVariableOperator)
			p.expect(itemIdentifier)
			l.Assignees = append(l.Assignees, ast.NewVariable(p.current.val))
			if p.peek().typ != itemComma {
				break
			}
			p.expect(itemComma)
		}
		p.expect(itemCloseParen)
		p.expect(itemAssignmentOperator)
		l.Operator = p.current.val
		l.Value = p.parseNextExpression()
		p.expectStmtEnd()
		return l

	case itemUnaryOperator, itemNegationOperator, itemAmpersandOperator, itemCastOperator, itemSubtractionOperator, itemBitwiseNotOperator:
		op := p.current
		return p.parseUnaryExpressionRight(p.parseNextExpression(), op)
	case itemVariableOperator:
		fallthrough
	case itemArray:
		fallthrough
	case itemIdentifier, itemStringLiteral, itemNumberLiteral, itemBooleanLiteral, itemNull, itemSelf, itemStatic, itemParent, itemShellCommand:
		expr = p.parseOperation(originalParenLev, p.expressionize())
	case itemInclude:
		inc := ast.Include{Expressions: make([]ast.Expression, 0)}
		for {
			inc.Expressions = append(inc.Expressions, p.parseNextExpression())
			if p.peek().typ != itemComma {
				break
			}
			p.expect(itemComma)
		}
		expr = inc
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
	case itemUnaryOperator, itemBitwiseNotOperator:
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
	case itemScopeResolutionOperator:
		p.next()
		expr = &ast.ClassExpression{Receiver: lhs, Expression: p.expressionize()}
	case itemArrayLookupOperatorLeft:
		expr = p.parseArrayLookup(lhs)
	case itemObjectOperator:
		expr = p.parseObjectLookup(lhs)
	case itemAssignmentOperator:
		assignee, ok := lhs.(ast.Assignable)
		if !ok {
			p.errorf("%s is not assignable", lhs)
		}
		expr = ast.AssignmentExpression{
			Assignee: assignee,
			Operator: p.current.val,
			Value:    p.parseNextExpression(),
		}
		return expr
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
func (p *parser) expressionize() (expr ast.Expression) {

	// These cases must come first and not repeat
	switch p.current.typ {
	case itemUnaryOperator, itemNegationOperator, itemCastOperator, itemSubtractionOperator, itemAmpersandOperator, itemBitwiseNotOperator:
		op := p.current
		p.next()
		return p.parseUnaryExpressionRight(p.expressionize(), op)
	case itemOpenParen:
		// Only parse open parentheses as a front matter to expression terms
		// so we don't get any dynamic function calls here.
		return p.parseExpression()
	}

	for {
		switch p.current.typ {
		case itemShellCommand:
			return &ast.ShellCommand{Command: p.current.val}
		case itemStringLiteral, itemBooleanLiteral, itemNumberLiteral, itemNull:
			return p.parseLiteral()
		case itemUnaryOperator:
			expr = newUnaryOperation(p.current, expr)
			p.next()
		case itemArray:
			expr = p.parseArrayDeclaration()
			p.next()
		case itemVariableOperator:
			expr = p.parseVariable()
			p.next()
			// Array lookup with curly braces is a special case that is only supported by PHP in
			// simple contexts.
			if p.current.typ == itemBlockBegin {
				expr = p.parseArrayLookup(expr)
				p.next()
			}
		case itemIdentifier:
			if p.peek().typ == itemOpenParen {
				// Function calls are okay here because we know they came with
				// a non-dynamic identifier.
				expr = p.parseFunctionCall(ast.Identifier{Value: p.current.val})
				p.next()
				continue
			}
			fallthrough
		case itemSelf, itemStatic, itemParent:
			if p.peek().typ == itemScopeResolutionOperator {
				r := p.current.val
				p.expect(itemScopeResolutionOperator)
				expr = ast.NewClassExpression(r, p.parseNextExpression())
				return
			}
			expr = ast.ConstantExpression{
				Variable: ast.NewVariable(p.current.val),
			}
			p.next()
		default:
			p.backup()
			return
		}
	}
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
		expr := ast.NewVariable(p.current.val)
		return expr
	default:
		return p.parseExpression()
	}
}

func (p *parser) parseAnonymousFunction() ast.Expression {
	f := &ast.AnonymousFunction{}
	f.Arguments = make([]ast.FunctionArgument, 0)
	f.ClosureVariables = make([]ast.FunctionArgument, 0)
	p.expect(itemOpenParen)
	if p.peek().typ != itemCloseParen {
		f.Arguments = append(f.Arguments, p.parseFunctionArgument())
	}

Loop:
	for {
		switch p.peek().typ {
		case itemComma:
			p.expect(itemComma)
			f.Arguments = append(f.Arguments, p.parseFunctionArgument())
		case itemCloseParen:
			break Loop
		default:
			p.errorf("unexpected argument separator:", p.current)
			return f
		}
	}
	p.expect(itemCloseParen)

	// Closure variables
	if p.peek().typ == itemUse {
		p.expect(itemUse)
		p.expect(itemOpenParen)
		f.ClosureVariables = append(f.ClosureVariables, p.parseFunctionArgument())
	ClosureLoop:
		for {
			switch p.peek().typ {
			case itemComma:
				p.expect(itemComma)
				f.ClosureVariables = append(f.ClosureVariables, p.parseFunctionArgument())
			case itemCloseParen:
				break ClosureLoop
			default:
				p.errorf("unexpected argument separator:", p.current)
				return f
			}
		}
		p.expect(itemCloseParen)
	}

	f.Body = p.parseBlock()
	return f
}
