package php

import (
	"strings"

	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/lexer"
	"github.com/stephens2424/php/token"
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

	switch p.current.Typ {
	case token.IgnoreErrorOperator:
		return p.parseNextExpression()
	case token.List:
		expr = p.parseList()
	case token.AmpersandOperator, token.SubtractionOperator:
		expr := p.parseUnaryExpressionRight(p.parseOperand(), p.current)
		return p.parseOperation(originalParenLev, expr)
	case token.UnaryOperator,
		token.NegationOperator,
		token.CastOperator,
		token.BitwiseNotOperator,
		token.ArrayLookupOperatorLeft,
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
		token.Include,
		token.Exit,
		token.ShellCommand:
		expr = p.parseOperation(originalParenLev, p.parseOperand())
	case token.OpenParen:
		// check for a cast operator that happens to have had spaces in it, and was thus lexed incorrectly
		if op := p.checkForCast(); op != nil {
			expr = p.parseUnaryExpressionRight(p.parseNextExpression(), *op)
			break
		}
		p.parenLevel++
		p.next()
		expr = p.parseExpression()
		p.expect(token.CloseParen)
		p.parenLevel--
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

func (p *Parser) checkForCast() *token.Item {
	if t := p.peek(); p.isCastType(t.Val) {
		p.next()
		if p.accept(token.CloseParen) {
			t.Val = "(" + t.Val + ")"
			t.Typ = token.CastOperator
			return &t
		}
		p.backup()
	}
	return nil

}

func (p *Parser) isCastType(s string) bool {
	switch strings.ToLower(s) {
	case "int":
	case "integer":
	case "float":
	case "bool":
	case "boolean":
	case "string":
	case "null":
	case "array":
	case "object":
	default:
		return false
	}
	return true
}

func (p *Parser) parseOperation(originalParenLevel int, lhs ast.Expression) (expr ast.Expression) {
	p.next()
	switch operationTypeForToken(p.current.Typ) {
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
		p.parenLevel--
		expr = lhs
	case subexpressionBeginOperation:
		// Check if we have a paren directly after a literal
		if _, ok := lhs.(*ast.Literal); ok {
			// If we do, we might be in a particular case of a function call with NULL as the function name. Let callers handle this
			p.backup()
			return lhs
		}
		p.parenLevel++
		expr = p.parseNextExpression()
	default:
		p.backup()
		return lhs
	}

	return p.parseOperation(originalParenLevel, expr)
}

func (p *Parser) parseAssignmentOperation(lhs, rhs ast.Expression, operator token.Item) (expr ast.Expression) {
	assignee, ok := lhs.(ast.Assignable)
	if !ok {
		p.errorf("%s is not assignable", lhs)
	}
	expr = ast.AssignmentExpression{
		Assignee: assignee,
		Operator: operator.Val,
		Value:    rhs,
	}
	return expr
}

// parseOperand takes the current token and returns it as the simplest
// expression for that token. That means an expression with no operators
// except for the object operator.
func (p *Parser) parseOperand() (expr ast.Expression) {

	// These cases must come first and not repeat
	switch p.current.Typ {
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
		return p.parseUnaryExpressionRight(p.parseExpression(), op)
	case token.OpenParen:
		if op := p.checkForCast(); op != nil {
			p.next()
			return p.parseUnaryExpressionRight(p.parseExpression(), *op)
		}
	case token.Include:
		return p.parseInclude()
	case token.Function:
		return p.parseAnonymousFunction()
	case token.NewOperator:
		return p.parseInstantiation()
	case token.ArrayLookupOperatorLeft:
		return p.parseArrayDeclaration()
	}

	switch p.current.Typ {
	case token.ShellCommand:
		return &ast.ShellCommand{Command: p.current.Val}
	case
		token.StringLiteral,
		token.BooleanLiteral,
		token.NumberLiteral,
		token.Null:
		return p.parseLiteral()
	case token.UnaryOperator:
		expr = p.parseUnaryExpressionLeft(expr, p.current)
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
	case token.ArrayLookupOperatorLeft, token.BlockBegin:
		expr = p.parseArrayLookup(expr)
		p.next()
	case token.Identifier, token.Exit:
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
		switch p.current.Typ {
		case token.UnaryOperator:
			expr = p.parseUnaryExpressionRight(expr, p.current)
			return
		case token.ObjectOperator:
			expr = p.parseObjectLookup(expr)
			p.next()
		case token.ArrayLookupOperatorLeft, token.BlockBegin:
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
}

func (p *Parser) parseLiteral() ast.Expression {
	switch p.current.Typ {
	case token.StringLiteral:
		return &ast.Literal{Type: ast.String, Value: p.current.Val}
	case token.BooleanLiteral:
		return &ast.Literal{Type: ast.Boolean, Value: p.current.Val}
	case token.NumberLiteral:
		return &ast.Literal{Type: ast.Float, Value: p.current.Val}
	case token.Null:
		if p.peek().Typ == token.OpenParen {
			expr := p.parseIdentifier()
			p.backup()
			return expr
		}
		return &ast.Literal{Type: ast.Null, Value: p.current.Val}
	}
	p.errorf("Unknown literal type")
	return nil
}

func (p *Parser) parseVariable() ast.Expression {
	var expr *ast.Variable
	p.expectCurrent(token.VariableOperator)
	switch p.next(); {
	case lexer.IsKeyword(p.current.Typ, p.current.Val):
		// keywords are all valid variable names
		fallthrough
	case p.current.Typ == token.Identifier:
		expr = ast.NewVariable(p.current.Val)
	case p.current.Typ == token.BlockBegin:
		expr = &ast.Variable{Name: p.parseNextExpression()}
		p.expect(token.BlockEnd)
	case p.current.Typ == token.VariableOperator:
		expr = &ast.Variable{Name: p.parseVariable()}
	default:
		p.errorf("unexpected variable operand %s", p.current)
		return nil
	}

	p.scope.Variable(expr)
	return expr
}

func (p *Parser) parseInclude() ast.Expression {
	inc := ast.Include{Expressions: make([]ast.Expression, 0)}
	for {
		inc.Expressions = append(inc.Expressions, p.parseNextExpression())
		if p.peek().Typ != token.Comma {
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
	switch typ := p.peek().Typ; {
	case typ == token.OpenParen && !p.instantiation:
		// Function calls are okay here because we know they came with
		// a non-dynamic identifier.
		expr = p.parseFunctionCall(&ast.Identifier{Value: p.current.Val})
		p.next()
	case typ == token.ScopeResolutionOperator:
		classIdent := p.current.Val
		p.next() // get onto ::, then we get to the next expr
		p.next()
		expr = ast.NewClassExpression(classIdent, p.parseOperand())
		p.next()
	case p.instantiation:
		defer p.next()
		return &ast.Identifier{Value: p.current.Val}
	default:
		name := p.current.Val
		v := ast.NewVariable(p.current.Val)
		expr = ast.ConstantExpression{
			Variable: v,
		}
		p.namespace.Constants[name] = append(p.namespace.Constants[name], v)
		p.next()
	}
	return expr
}

// parseScopeResolutionFromKeyword specifically parses self::, static::, and parent::
func (p *Parser) parseScopeResolutionFromKeyword() ast.Expression {
	if p.peek().Typ == token.ScopeResolutionOperator {
		r := p.current.Val
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
	switch p.current.Typ {
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
