package php

import (
	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
)

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

	/*
	   PHP's documentation would have this operator be at 4, but it also notes:

	       Although = has a lower precedence than most other operators, PHP will
	       still allow expressions similar to the following: if (!$a = foo()), in
	       which case the return value of foo() is put into $a.

	   Thus, we put it at 17, pending further testing.
	*/
	token.AssignmentOperator: 17,
	token.WrittenAndOperator: 3,
	token.WrittenXorOperator: 2,
	token.WrittenOrOperator:  1,
}

func (p *Parser) parseExpression() (expr ast.Expression) {
	originalParenLev := p.parenLevel

	switch p.current.typ {
	case token.IgnoreErrorOperator:
		return p.parseNextExpression()
	case token.List:
		expr = p.parseList()
	case
		token.UnaryOperator,
		token.NegationOperator,
		token.AmpersandOperator,
		token.CastOperator,
		token.SubtractionOperator,
		token.BitwiseNotOperator:
		op := p.current
		expr = p.parseUnaryExpressionRight(p.parseNextExpression(), op)
	case
		token.Function,
		token.NewOperator,
		token.VariableOperator,
		token.Array,
		token.Identifier,
		token.StringLiteral,
		token.NumberLiteral,
		token.BooleanLiteral,
		token.Null,
		token.Self,
		token.Static,
		token.Parent,
		token.ShellCommand:
		expr = p.parseOperation(originalParenLev, p.parseOperand())
	case token.Include:
		expr = p.parseInclude()
	case token.OpenParen:
		p.parenLevel += 1
		p.next()
		expr = p.parseExpression()
		p.expect(token.CloseParen)
		p.parenLevel -= 1
		expr = p.parseOperation(originalParenLev, expr)
	default:
		p.errorf("Expected expression. Found %s", p.current)
	}
	if p.parenLevel != originalParenLev {
		p.errorf("unbalanced parens: %d prev: %d", p.parenLevel, originalParenLev)
		return
	}
	return
}

func (p *Parser) parseOperation(originalParenLevel int, lhs ast.Expression) (expr ast.Expression) {
	p.next()
	switch operationTypeForToken(p.current.typ) {
	case ignoreErrorOperation:
		return p.parseOperation(originalParenLevel, lhs)
	case unaryOperation:
		expr = p.parseUnaryExpressionLeft(lhs, p.current)
	case assignmentOperation, binaryOperation:
		expr = p.parseBinaryOperation(lhs, p.current, originalParenLevel)
	case ternaryOperation:
		expr = p.parseTernaryOperation(lhs)
	case subexpressionEndOperation:
		if p.parenLevel == originalParenLevel {
			p.backup()
			return lhs
		}
		p.parenLevel -= 1
		expr = lhs
	case subexpressionBeginOperation:
		// Check if we have a paren directly after a literal
		if _, ok := lhs.(*ast.Literal); ok {
			// If we do, we might be in a particular case of a function call with NULL as the function name. Let callers handle this
			p.backup()
			return lhs
		}
		p.parenLevel += 1
		expr = p.parseNextExpression()
	default:
		p.backup()
		return lhs
	}

	return p.parseOperation(originalParenLevel, expr)
}

func (p *Parser) parseAssignmentOperation(lhs, rhs ast.Expression, operator Item) (expr ast.Expression) {
	assignee, ok := lhs.(ast.Assignable)
	if !ok {
		p.errorf("%s is not assignable", lhs)
	}
	expr = ast.AssignmentExpression{
		Assignee: assignee,
		Operator: operator.val,
		Value:    rhs,
	}
	return expr
}

// parseOperand takes the current token and returns it as the simplest
// expression for that token. That means an expression with no operators
// except for the object operator.
func (p *Parser) parseOperand() (expr ast.Expression) {

	// These cases must come first and not repeat
	switch p.current.typ {
	case token.IgnoreErrorOperator:
		p.next()
		return p.parseOperand()
	case
		token.UnaryOperator,
		token.NegationOperator,
		token.CastOperator,
		token.SubtractionOperator,
		token.AmpersandOperator,
		token.BitwiseNotOperator:
		op := p.current
		p.next()
		return p.parseUnaryExpressionRight(p.parseOperand(), op)
	case token.Function:
		return p.parseAnonymousFunction()
	case token.NewOperator:
		return p.parseInstantiation()
	}

	switch p.current.typ {
	case token.ShellCommand:
		return &ast.ShellCommand{Command: p.current.val}
	case
		token.StringLiteral,
		token.BooleanLiteral,
		token.NumberLiteral,
		token.Null:
		return p.parseLiteral()
	case token.UnaryOperator:
		expr = newUnaryOperation(p.current, expr)
		p.next()
		return

	case token.Array:
		expr = p.parseArrayDeclaration()
		p.next()
	case token.VariableOperator:
		expr = p.parseVariableOperand()
	case token.ObjectOperator:
		expr = p.parseObjectLookup(expr)
		p.next()
	case token.ArrayLookupOperatorLeft:
		expr = p.parseArrayLookup(expr)
		p.next()
	case token.Identifier:
		expr = p.parseIdentifier()
	case token.Self, token.Static, token.Parent:
		expr = p.parseScopeResolutionFromKeyword()
	default:
		p.backup()
		return
	}

	return p.parseOperandComponent(expr)
}

func (p *Parser) parseOperandComponent(lhs ast.Expression) (expr ast.Expression) {
	expr = lhs
	for {
		switch p.current.typ {
		case token.UnaryOperator:
			expr = newUnaryOperation(p.current, expr)
			return
		case token.ObjectOperator:
			expr = p.parseObjectLookup(expr)
			p.next()
		case token.ArrayLookupOperatorLeft:
			expr = p.parseArrayLookup(expr)
			p.next()
		case token.OpenParen:
			p.backup()
			if p.instantiation {
				return
			}
			expr = p.parseFunctionCall(expr)
			p.next()
		default:
			p.backup()
			return
		}
	}
	return
}

func (p *Parser) parseLiteral() ast.Expression {
	switch p.current.typ {
	case token.StringLiteral:
		return &ast.Literal{Type: ast.String, Value: p.current.val}
	case token.BooleanLiteral:
		return &ast.Literal{Type: ast.Boolean, Value: p.current.val}
	case token.NumberLiteral:
		return &ast.Literal{Type: ast.Float, Value: p.current.val}
	case token.Null:
		if p.peek().typ == token.OpenParen {
			expr := p.parseIdentifier()
			p.backup()
			return expr
		}
		return &ast.Literal{Type: ast.Null, Value: p.current.val}
	}
	p.errorf("Unknown literal type")
	return nil
}

func (p *Parser) parseVariable() ast.Expression {
	p.expectCurrent(token.VariableOperator)
	switch p.next(); {
	case isKeyword(p.current.typ, p.current.val):
		// keywords are all valid variable names
		fallthrough
	case p.current.typ == token.Identifier:
		expr := ast.NewVariable(p.current.val)
		return expr
	case p.current.typ == token.BlockBegin:
		return ast.Variable{Name: p.parseExpression()}
	case p.current.typ == token.VariableOperator:
		return ast.Variable{Name: p.parseVariable()}
	default:
		p.errorf("unexpected variable operand %s", p.current)
		return nil
	}
}

func (p *Parser) parseInclude() ast.Expression {
	inc := ast.Include{Expressions: make([]ast.Expression, 0)}
	for {
		inc.Expressions = append(inc.Expressions, p.parseNextExpression())
		if p.peek().typ != token.Comma {
			break
		}
		p.expect(token.Comma)
	}
	return inc
}

func (p *Parser) parseIgnoreError() ast.Expression {
	p.next()
	return p.parseExpression()
}

func (p *Parser) parseNew(originalParenLev int) ast.Expression {
	expr := p.parseInstantiation()
	expr = p.parseOperation(originalParenLev, expr)
	return expr
}

func (p *Parser) parseIdentifier() (expr ast.Expression) {
	switch p.peek().typ {
	case token.OpenParen:
		// Function calls are okay here because we know they came with
		// a non-dynamic identifier.
		expr = p.parseFunctionCall(ast.Identifier{Value: p.current.val})
		p.next()
	case token.ScopeResolutionOperator:
		classIdent := p.current.val
		p.next() // get onto ::, then we get to the next expr
		p.next()
		expr = ast.NewClassExpression(classIdent, p.parseOperand())
		p.next()
	default:
		expr = ast.ConstantExpression{
			Variable: ast.NewVariable(p.current.val),
		}
		p.next()
	}
	return expr
}

// parseScopeResolutionFromKeyword specifically parses self::, static::, and parent::
func (p *Parser) parseScopeResolutionFromKeyword() ast.Expression {
	if p.peek().typ == token.ScopeResolutionOperator {
		r := p.current.val
		p.expect(token.ScopeResolutionOperator)
		p.next()
		expr := ast.NewClassExpression(r, p.parseOperand())
		p.next()
		return expr
	}
	// TODO Error
	p.next()
	return nil
}

func (p *Parser) parseVariableOperand() ast.Expression {
	expr := p.parseVariable()
	p.next()
	// Array lookup with curly braces is a special case that is only supported by PHP in
	// simple contexts.
	switch p.current.typ {
	case token.BlockBegin:
		expr = p.parseArrayLookup(expr)
		p.next()
	case token.ScopeResolutionOperator:
		expr = &ast.ClassExpression{Receiver: expr, Expression: p.parseNextExpression()}
		p.next()
	case token.OpenParen:
		p.backup()
		expr = p.parseFunctionCall(expr)
		p.next()
	}

	return expr
}
