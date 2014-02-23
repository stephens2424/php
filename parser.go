package php

import (
	"fmt"

	"stephensearles.com/php/ast"
)

type parser struct {
	lexer *lexer

	previous   []Item
	idx        int
	current    Item
	errors     []error
	parenLevel int
	errorMap   map[int]bool

	debug     bool
	MaxErrors int
}

func NewParser(input string) *parser {
	return newParser(input)
}

func newParser(input string) *parser {
	p := &parser{
		idx:       -1,
		MaxErrors: 10,
		lexer:     newLexer(input),
		errorMap:  make(map[int]bool),
	}
	return p
}

func (p *parser) Parse() []ast.Node {
	return p.parse()
}

func (p *parser) parse() []ast.Node {
	defer func() {
		if len(p.errors) > 0 {
			for _, err := range p.errors {
				fmt.Println(err)
			}
		}
		if r := recover(); r != nil {
			if p.debug {
				panic(r)
			} else {
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
		return ast.Echo(ast.Literal{Type: ast.String})
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

func (p *parser) peek() (i Item) {
	p.next()
	i = p.current
	p.backup()
	return
}

func (p *parser) expect(i ItemType) {
	p.next()
	if p.current.typ != i {
		p.expected(i)
	}
}

func (p *parser) expected(i ItemType) {
	p.errorf("Found %s, expected %s", p.current, i)
}

func (p *parser) errorf(str string, args ...interface{}) {
	if _, ok := p.errorMap[p.current.pos.Line]; ok {
		return
	}
	errString := fmt.Sprintf(str, args...)
	p.errorMap[p.current.pos.Line] = true
	p.errors = append(p.errors, fmt.Errorf("%s: %s", p.errorPrefix(), errString))
	if len(p.errors) > p.MaxErrors {
		panic("too many errors")
	}
}

func (p *parser) errorPrefix() string {
	return fmt.Sprintf("%s %d", p.lexer.file, p.current.pos.Line)
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

func (p *parser) parseNextExpression() ast.Expression {
	p.next()
	return p.parseExpression()
}

func newUnaryOperation(operator Item, expr ast.Expression) ast.OperatorExpression {
	t := ast.Numeric
	if operator.val == "!" {
		t = ast.Boolean
	}
	return ast.OperatorExpression{
		Type:     t,
		Operand1: expr,
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

func (p *parser) parseFunctionCall() *ast.FunctionCallExpression {
	expr := &ast.FunctionCallExpression{}
	if p.current.typ != itemNonVariableIdentifier {
		p.expected(itemNonVariableIdentifier)
	}
	expr.FunctionName = p.current.val
	expr.Arguments = make([]ast.Expression, 0)
	p.expect(itemOpenParen)
	first := true
	p.next()
	for {
		if p.current.typ == itemCloseParen {
			break
		}
		if !first {
			p.backup()
			p.expect(itemArgumentSeparator)
			p.next()
		} else {
			first = false
		}
		arg := p.parseExpression()
		if arg == nil {
			break
		}
		expr.Arguments = append(expr.Arguments, arg)
		p.next()
	}
	return expr
}

func (p *parser) parseStmt() ast.Statement {
	switch p.current.typ {
	case itemBlockBegin:
		p.backup()
		return p.parseBlock()
	case itemGlobal:
		p.next()
		ident := ast.GlobalIdentifier{ast.NewIdentifier(p.current.val)}
		p.expect(itemStatementEnd)
		return ident
	case itemIdentifier:
		switch p.peek().typ {
		case itemUnaryOperator:
			expr := ast.ExpressionStmt{p.parseExpression()}
			p.expect(itemStatementEnd)
			return expr
		default:
			n := ast.AssignmentStmt{}
			n.Assignee = p.parseIdentifier().(ast.Assignable)
			p.expect(itemAssignmentOperator)
			n.Operator = p.current.val
			p.next()
			n.Value = p.parseExpression()
			p.expect(itemStatementEnd)
			return n
		}
	case itemUnaryOperator:
		expr := ast.ExpressionStmt{p.parseExpression()}
		p.expect(itemStatementEnd)
		return expr
	case itemFunction:
		return p.parseFunctionStmt()
	case itemEcho:
		p.next()
		expr := p.parseExpression()
		p.expect(itemStatementEnd)
		return ast.Echo(expr)
	case itemIf:
		return p.parseIf()
	case itemWhile:
		return p.parseWhile()
	case itemDo:
		return p.parseDo()
	case itemFor:
		return p.parseFor()
	case itemForeach:
		return p.parseForeach()
	case itemNonVariableIdentifier:
		stmt := p.parseFunctionCall()
		p.expect(itemStatementEnd)
		return stmt
	case itemClass:
		return p.parseClass()
	case itemReturn:
		p.next()
		stmt := ast.ReturnStmt{}
		if p.current.typ != itemStatementEnd {
			stmt.Expression = p.parseExpression()
			p.expect(itemStatementEnd)
		}
		return stmt
	case itemThrow:
		stmt := ast.ThrowStmt{Expression: p.parseNextExpression()}
		p.expect(itemStatementEnd)
		return stmt
	case itemExit:
		stmt := ast.ExitStmt{}
		if p.peek().typ == itemOpenParen {
			p.expect(itemOpenParen)
			stmt.Expression = p.parseNextExpression()
			p.expect(itemCloseParen)
		}
		p.expect(itemStatementEnd)
		return stmt
	case itemInclude:
		stmt := ast.IncludeStmt{Expression: p.parseNextExpression()}
		p.expect(itemStatementEnd)
		return stmt
	case itemIgnoreErrorOperator:
		// Ignore this operator
		return p.parseStmt()
	default:
		p.errorf("Found %s, expected html or php begin", p.current)
		return nil
	}
}

func (p *parser) parseWhile() ast.Statement {
	p.expect(itemOpenParen)
	term := p.parseNextExpression()
	p.expect(itemCloseParen)
	block := p.parseBlock()
	return &ast.WhileStmt{
		Termination: term,
		LoopBlock:   block,
	}
}

func (p *parser) parseForeach() ast.Statement {
	stmt := &ast.ForeachStmt{}
	p.expect(itemOpenParen)
	stmt.Source = p.parseNextExpression()
	p.expect(itemAsOperator)
	p.expect(itemIdentifier)
	stmt.Value = ast.NewIdentifier(p.current.val)
	if p.peek().typ == itemArrayKeyOperator {
		stmt.Key = stmt.Value
		p.expect(itemArrayKeyOperator)
		p.expect(itemIdentifier)
		stmt.Value = ast.NewIdentifier(p.current.val)
	}
	p.expect(itemCloseParen)
	stmt.LoopBlock = p.parseBlock()
	return stmt
}

func (p *parser) parseFor() ast.Statement {
	stmt := &ast.ForStmt{}
	p.expect(itemOpenParen)
	stmt.Initialization = p.parseNextExpression()
	p.expect(itemStatementEnd)
	stmt.Termination = p.parseNextExpression()
	p.expect(itemStatementEnd)
	stmt.Iteration = p.parseNextExpression()
	p.expect(itemCloseParen)
	stmt.LoopBlock = p.parseBlock()
	return stmt
}

func (p *parser) parseDo() ast.Statement {
	block := p.parseBlock()
	p.expect(itemWhile)
	p.expect(itemOpenParen)
	term := p.parseNextExpression()
	p.expect(itemCloseParen)
	p.expect(itemStatementEnd)
	return &ast.DoWhileStmt{
		Termination: term,
		LoopBlock:   block,
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
	return p.parseClassFields(ast.Class{Name: name})
}

func (p *parser) parseClassFields(c ast.Class) ast.Class {
	c.Methods = make([]ast.Method, 0)
	c.Properties = make([]ast.Property, 0)
	p.next()
	var vis ast.Visibility
	for p.current.typ != itemBlockEnd {
		switch p.current.typ {
		case itemPrivate:
			vis = ast.Private
		case itemProtected:
			vis = ast.Protected
		case itemPublic:
			vis = ast.Public
		default:
			vis = ast.Public
			p.backup()
		}
		p.next()
		switch p.current.typ {
		case itemFunction:
			c.Methods = append(c.Methods, ast.Method{
				Visibility:   vis,
				FunctionStmt: p.parseFunctionStmt(),
			})
		case itemIdentifier:
			c.Properties = append(c.Properties, ast.Property{
				Visibility: vis,
				Name:       p.current.val,
			})
		default:
		}
		p.next()
	}
	return c
}
