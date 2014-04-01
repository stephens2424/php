package php

import (
	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
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

var operatorPrecedence = map[token.Token]int{
	token.ArrayLookupOperatorLeft: 19,
	token.UnaryOperator:           18,
	token.BitwiseNotOperator:      18,
	token.CastOperator:            18,
	token.InstanceofOperator:      17,
	token.NegationOperator:        16,
	token.MultOperator:            15,
	token.AdditionOperator:        14,
	token.SubtractionOperator:     14,
	token.ConcatenationOperator:   14,

	token.BitwiseShiftOperator: 13,
	token.ComparisonOperator:   12,
	token.EqualityOperator:     11,

	token.AmpersandOperator:  10,
	token.BitwiseXorOperator: 9,
	token.BitwiseOrOperator:  8,
	token.AndOperator:        7,
	token.OrOperator:         6,
	token.TernaryOperator1:   5,
	token.TernaryOperator2:   5,
	token.AssignmentOperator: 4,
	token.WrittenAndOperator: 3,
	token.WrittenXorOperator: 2,
	token.WrittenOrOperator:  1,
}

func (p *parser) parseExpression() (expr ast.Expression) {
	// consume expression
	originalParenLev := p.parenLevel
	switch p.current.typ {
	case token.IgnoreErrorOperator:
		p.next()
		return p.parseExpression()
	case token.Function:
		return p.parseAnonymousFunction()
	case token.NewOperator:
		expr = p.parseInstantiation()
		expr = p.parseOperation(originalParenLev, expr)
		return

	case token.List:
		l := &ast.ListStatement{
			Assignees: make([]*ast.Variable, 0),
		}
		p.expect(token.OpenParen)
		for {
			p.expect(token.VariableOperator)
			p.expect(token.Identifier)
			l.Assignees = append(l.Assignees, ast.NewVariable(p.current.val))
			if p.peek().typ != token.Comma {
				break
			}
			p.expect(token.Comma)
		}
		p.expect(token.CloseParen)
		p.expect(token.AssignmentOperator)
		l.Operator = p.current.val
		l.Value = p.parseNextExpression()
		p.expectStmtEnd()
		return l

	case token.UnaryOperator, token.NegationOperator, token.AmpersandOperator, token.CastOperator, token.SubtractionOperator, token.BitwiseNotOperator:
		op := p.current
		return p.parseUnaryExpressionRight(p.parseNextExpression(), op)
	case token.VariableOperator:
		fallthrough
	case token.Array:
		fallthrough
	case token.Identifier, token.StringLiteral, token.NumberLiteral, token.BooleanLiteral, token.Null, token.Self, token.Static, token.Parent, token.ShellCommand:
		expr = p.parseOperation(originalParenLev, p.expressionize())
	case token.Include:
		inc := ast.Include{Expressions: make([]ast.Expression, 0)}
		for {
			inc.Expressions = append(inc.Expressions, p.parseNextExpression())
			if p.peek().typ != token.Comma {
				break
			}
			p.expect(token.Comma)
		}
		expr = inc
	case token.OpenParen:
		p.parenLevel += 1
		p.next()
		expr = p.parseExpression()
		p.expect(token.CloseParen)
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
	case token.IgnoreErrorOperator:
		return p.parseOperation(originalParenLevel, lhs)
	case token.UnaryOperator, token.BitwiseNotOperator:
		expr = p.parseUnaryExpressionLeft(lhs, p.current)
	case token.AdditionOperator, token.SubtractionOperator, token.ConcatenationOperator, token.ComparisonOperator, token.MultOperator, token.AndOperator, token.OrOperator, token.AmpersandOperator, token.BitwiseXorOperator, token.BitwiseOrOperator, token.BitwiseShiftOperator, token.WrittenAndOperator, token.WrittenXorOperator, token.WrittenOrOperator, token.InstanceofOperator:
		expr = p.parseBinaryOperation(lhs, p.current, originalParenLevel)
	case token.TernaryOperator1:
		expr = p.parseTernaryOperation(lhs)
	case token.CloseParen:
		if p.parenLevel <= originalParenLevel {
			p.backup()
			return lhs
		}
		p.parenLevel -= 1
		return p.parseOperation(originalParenLevel, lhs)
	case token.ScopeResolutionOperator:
		p.next()
		expr = &ast.ClassExpression{Receiver: lhs, Expression: p.expressionize()}
	case token.ArrayLookupOperatorLeft:
		expr = p.parseArrayLookup(lhs)
	case token.ObjectOperator:
		expr = p.parseObjectLookup(lhs)
	case token.AssignmentOperator:
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
	case token.UnaryOperator, token.NegationOperator, token.CastOperator, token.SubtractionOperator, token.AmpersandOperator, token.BitwiseNotOperator:
		op := p.current
		p.next()
		return p.parseUnaryExpressionRight(p.expressionize(), op)
	case token.OpenParen:
		// Only parse open parentheses as a front matter to expression terms
		// so we don't get any dynamic function calls here.
		return p.parseExpression()
	}

	for {
		switch p.current.typ {
		case token.ShellCommand:
			return &ast.ShellCommand{Command: p.current.val}
		case token.StringLiteral, token.BooleanLiteral, token.NumberLiteral, token.Null:
			return p.parseLiteral()
		case token.UnaryOperator:
			expr = newUnaryOperation(p.current, expr)
			p.next()
		case token.Array:
			expr = p.parseArrayDeclaration()
			p.next()
		case token.VariableOperator:
			expr = p.parseVariable()
			p.next()
			// Array lookup with curly braces is a special case that is only supported by PHP in
			// simple contexts.
			if p.current.typ == token.BlockBegin {
				expr = p.parseArrayLookup(expr)
				p.next()
			}
		case token.Identifier:
			if p.peek().typ == token.OpenParen {
				// Function calls are okay here because we know they came with
				// a non-dynamic identifier.
				expr = p.parseFunctionCall(ast.Identifier{Value: p.current.val})
				p.next()
				continue
			}
			fallthrough
		case token.Self, token.Static, token.Parent:
			if p.peek().typ == token.ScopeResolutionOperator {
				r := p.current.val
				p.expect(token.ScopeResolutionOperator)
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
	case token.StringLiteral:
		return &ast.Literal{Type: ast.String}
	case token.BooleanLiteral:
		return &ast.Literal{Type: ast.Boolean}
	case token.NumberLiteral:
		return &ast.Literal{Type: ast.Float}
	case token.Null:
		return &ast.Literal{Type: ast.Null}
	}
	p.errorf("Unknown literal type")
	return nil
}

func (p *parser) parseVariable() ast.Expression {
	p.expectCurrent(token.VariableOperator)
	switch p.next(); {
	case isKeyword(p.current.typ):
		// keywords are all valid variable names
		fallthrough
	case p.current.typ == token.Identifier:
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
	p.expect(token.OpenParen)
	if p.peek().typ != token.CloseParen {
		f.Arguments = append(f.Arguments, p.parseFunctionArgument())
	}

Loop:
	for {
		switch p.peek().typ {
		case token.Comma:
			p.expect(token.Comma)
			f.Arguments = append(f.Arguments, p.parseFunctionArgument())
		case token.CloseParen:
			break Loop
		default:
			p.errorf("unexpected argument separator:", p.current)
			return f
		}
	}
	p.expect(token.CloseParen)

	// Closure variables
	if p.peek().typ == token.Use {
		p.expect(token.Use)
		p.expect(token.OpenParen)
		f.ClosureVariables = append(f.ClosureVariables, p.parseFunctionArgument())
	ClosureLoop:
		for {
			switch p.peek().typ {
			case token.Comma:
				p.expect(token.Comma)
				f.ClosureVariables = append(f.ClosureVariables, p.parseFunctionArgument())
			case token.CloseParen:
				break ClosureLoop
			default:
				p.errorf("unexpected argument separator:", p.current)
				return f
			}
		}
		p.expect(token.CloseParen)
	}

	f.Body = p.parseBlock()
	return f
}
