package php

import (
	"fmt"
	"stephensearles.com/php/ast"
)

type parser struct {
	lexer *lexer

	previous []item
	idx      int
	current  item
	errors   []error

	parenLevel int
	debug      bool
}

func newParser(input string) *parser {
	p := &parser{
		idx:   -1,
		lexer: newLexer(input),
	}
	return p
}

func (p *parser) parse() []ast.Node {
	defer func() {
		if !p.debug {
			if r := recover(); r != nil {
				fmt.Println(p.errors)
				fmt.Println(r)
			}
		}
	}()
	// expecting either itemHTML or itemPHPBegin
	nodes := make([]ast.Node, 0, 1)
TokenLoop:
	for {
		p.next()
		switch p.current.typ {
		case itemEOF:
			break TokenLoop
		default:
			n := p.parseNode()
			if n != nil {
				nodes = append(nodes, n)
			}
		}
	}
	return nodes
}

func (p *parser) parseNode() ast.Node {
	switch p.current.typ {
	case itemHTML:
		return ast.EchoStmt(ast.Literal{ast.String})
	case itemPHPBegin:
		return nil
	case itemPHPEnd:
		return nil
	}
	return p.parseStmt()
}

func (p *parser) next() {
	p.idx += 1
	if len(p.previous) <= p.idx {
		p.current = p.lexer.nextItem()
		p.previous = append(p.previous, p.current)
	} else {
		p.current = p.previous[p.idx]
	}
}

func (p *parser) backup() {
	p.idx -= 1
	p.current = p.previous[p.idx]
}

func (p *parser) expect(i itemType) {
	p.next()
	if p.current.typ != i {
		p.expected(i)
	}
}

func (p *parser) expected(i itemType) {
	p.errorf("Found %s, expected %s", p.current, i)
}

func (p *parser) errorf(str string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Errorf(str, args...))
	if len(p.errors) > 0 {
		panic("too many errors")
	}
}

func (p *parser) parseIf() *ast.IfStmt {
	p.expect(itemOpenParen)
	n := &ast.IfStmt{}
	p.next()
	n.Condition = p.parseExpression()
	p.expect(itemCloseParen)
	p.next()
	n.TrueBlock = p.parseStmt()
	p.next()
	if p.current.typ == itemElse {
		p.next()
		n.FalseBlock = p.parseStmt()
	} else {
		n.FalseBlock = ast.Block{}
		p.backup()
	}
	return n
}

/*

Expression = Identifier |
             Literal |
             FunctionCall |
             {Unary Operator} Expression |
             OpenParen Expression CloseParen |
             Expression {Unary Operator} |
             Expression {Binary Operator} Expression |
             Expression {Tertiary Operator 1} Expression {Tertiary Operator 2} Expression

*/
func (p *parser) parseExpression() (expr ast.Expression) {
	// check for a unary operation that begins with the operator
	if p.current.typ == itemUnaryOperator {
		return newUnaryOperation(p.current, p.parseExpression())
	}
	if p.current.typ == itemNegationOperator {
		return newUnaryOperation(p.current, p.parseExpression())
	}

	// get the left hand side of the expression (could be the end of it)
	switch p.current.typ {
	case itemStringLiteral:
		expr = ast.Literal{ast.String}
	case itemNumberLiteral:
		expr = ast.Literal{ast.Float} //need to do integers...
	case itemBooleanLiteral:
		expr = ast.Literal{ast.Boolean}
	case itemIdentifier:
		expr = ast.NewIdentifier(p.current.val)
	case itemOpenParen:
		expr = p.parseExpression()
		p.expect(itemCloseParen)
	case itemNonVariableIdentifier:
		expr = p.parseFunctionCall()
	}

	// look for an operator
	p.next()
	switch p.current.typ {
	case itemUnaryOperator:
		return newUnaryOperation(p.current, expr)
	case itemAdditionOperator:
		return newBinaryOperation(p.current, expr, p.parseExpression())
	case itemSubtractionOperator:
		return newBinaryOperation(p.current, expr, p.parseExpression())
	case itemMultOperator:
		return newBinaryOperation(p.current, expr, p.parseExpression())
	case itemComparisonOperator:
		return newBinaryOperation(p.current, expr, p.parseExpression())
	case itemConcatenationOperator:
		return ast.OperatorExpression{
			Operand1: expr,
			Operand2: p.parseExpression(),
			Type:     ast.String,
		}
	}
	p.backup()
	return expr
}

func newUnaryOperation(operator item, expr ast.Expression) ast.OperatorExpression {
	t := ast.Numeric
	if operator.val == "!" {
		t = ast.Boolean
	}
	return ast.OperatorExpression{
		Type:     t,
		Operand1: expr,
	}
}

func newBinaryOperation(operator item, expr1, expr2 ast.Expression) ast.OperatorExpression {
	t := ast.Numeric
	switch operator.typ {
	case itemComparisonOperator:
		t = ast.Boolean
	case itemConcatenationOperator:
		t = ast.String
	}
	return ast.OperatorExpression{
		Type:     t,
		Operand1: expr1,
		Operand2: expr2,
	}
}

func (p *parser) parseFunctionCall() ast.FunctionCallExpression {
	expr := ast.FunctionCallExpression{}
	if p.current.typ != itemNonVariableIdentifier {
		p.expected(itemNonVariableIdentifier)
	}
	expr.FunctionName = p.current.val
	expr.Arguments = make([]ast.Expression, 0)
	p.expect(itemOpenParen)
	first := true
	for {
		p.next()
		if p.current.typ == itemCloseParen {
			break
		}
		if !first {
			p.expect(itemArgumentSeparator)
		} else {
			first = false
		}
		expr.Arguments = append(expr.Arguments, p.parseExpression())
	}
	return expr
}

func (p *parser) parseStmt() ast.Statement {
	switch p.current.typ {
	case itemBlockBegin:
		return p.parseBlock()
	case itemIdentifier:
		n := ast.AssignmentStmt{}
		n.Assignee = ast.NewIdentifier(p.current.val)
		p.expect(itemAssignmentOperator)
		p.next()
		n.Value = p.parseExpression()
		p.expect(itemStatementEnd)
		return n
	case itemFunction:
		return p.parseFunctionStmt()
	case itemEcho:
		p.next()
		expr := p.parseExpression()
		p.expect(itemStatementEnd)
		return ast.EchoStmt(expr)
	case itemIf:
		return p.parseIf()
	case itemNonVariableIdentifier:
		stmt := p.parseExpression()
		p.expect(itemStatementEnd)
		return stmt
	case itemClass:
		return p.parseClass()
	default:
		p.errorf("Found %s, expected html or php begin", p.current)
		return nil
	}
}

func (p *parser) parseFunctionStmt() *ast.FunctionStmt {
	stmt := &ast.FunctionStmt{}
	p.expect(itemNonVariableIdentifier)
	stmt.Name = p.current.val
	stmt.Arguments = make([]ast.FunctionArgument, 0)
	p.expect(itemOpenParen)
	first := true
	for {
		p.next()
		if p.current.typ == itemCloseParen {
			break
		}
		p.backup()
		if !first {
			p.expect(itemArgumentSeparator)
		} else {
			first = false
		}
		p.next()
		arg := ast.FunctionArgument{}
		if p.current.typ == itemNonVariableIdentifier {
			arg.TypeHint = p.current.val
		} else {
			p.backup()
		}
		p.expect(itemIdentifier)
		arg.Identifier = ast.NewIdentifier(p.current.val)
		stmt.Arguments = append(stmt.Arguments, arg)
	}
	stmt.Body = p.parseBlock()
	return stmt
}

func (p *parser) parseBlock() ast.Block {
	block := ast.Block{}
	p.expect(itemBlockBegin)
	for {
		p.next()
		block.Statements = append(block.Statements, p.parseStmt())
		if p.next(); p.current.typ == itemBlockEnd {
			break
		}
		p.backup()
	}
	return block
}

func (p *parser) parseClass() ast.Class {
	p.expect(itemNonVariableIdentifier)
	name := p.current.val
	p.next()
	if p.current.typ == itemExtends {
		p.expect(itemNonVariableIdentifier)
	} else {
		p.backup()
	}
	p.expect(itemBlockBegin)
	return ast.Class{
		Name:    name,
		Methods: p.parseMethods(),
	}
}

func (p *parser) parseMethods() (methods []ast.Method) {
	methods = make([]ast.Method, 0)
	p.next()
	for p.current.typ != itemBlockEnd {
		m := ast.Method{}
		switch p.current.typ {
		case itemPrivate:
			m.Visibility = ast.Private
			p.expect(itemFunction)
		case itemProtected:
			m.Visibility = ast.Protected
			p.expect(itemFunction)
		case itemPublic:
			m.Visibility = ast.Public
			p.expect(itemFunction)
		case itemFunction:
			m.Visibility = ast.Public
		default:
			p.expected(itemFunction)
		}
		m.FunctionStmt = p.parseFunctionStmt()
		methods = append(methods, m)
		p.next()
	}
	return methods
}
