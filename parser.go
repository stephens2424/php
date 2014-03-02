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
	errorCount int

	Debug     bool
	MaxErrors int
}

func NewParser(input string) *parser {
	p := &parser{
		idx:       -1,
		MaxErrors: 10,
		lexer:     newLexer(input),
		errorMap:  make(map[int]bool),
	}
	return p
}

func (p *parser) Parse() ([]ast.Node, []error) {
	defer func() {
		if r := recover(); r != nil {
			if p.Debug {
				for _, err := range p.errors {
					fmt.Println(err)
				}
				panic(r)
			}
			p.errors = append([]error{fmt.Errorf("%s", r)}, p.errors...)
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
	return nodes, p.errors
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

func (p *parser) expectCurrent(i ItemType) {
	if p.current.typ != i {
		p.expected(i)
	}
}

func (p *parser) expectAndNext(i ItemType) {
	if p.current.typ != i {
		p.expected(i)
	}
	p.next()
}

func (p *parser) expect(i ItemType) {
	p.next()
	p.expectCurrent(i)
}

func (p *parser) expected(i ItemType) {
	p.errorf("Found %s, expected %s", p.current, i)
}

func (p *parser) errorf(str string, args ...interface{}) {
	p.errorCount += 1
	if p.errorCount > p.MaxErrors {
		panic("too many errors")
	}
	if _, ok := p.errorMap[p.current.pos.Line]; ok {
		return
	}
	errString := fmt.Sprintf(str, args...)
	p.errorMap[p.current.pos.Line] = true
	p.errors = append(p.errors, fmt.Errorf("%s: %s", p.errorPrefix(), errString))
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
	n.TrueBranch = p.parseStmt()
	p.next()
	if p.current.typ == itemElse {
		p.next()
		n.FalseBranch = p.parseStmt()
	} else {
		n.FalseBranch = ast.Block{}
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

func (p *parser) parseFunctionCall() *ast.FunctionCallExpression {
	expr := &ast.FunctionCallExpression{}
	if p.current.typ != itemNonVariableIdentifier {
		p.expected(itemNonVariableIdentifier)
	}
	expr.FunctionName = p.current.val
	expr.Arguments = make([]ast.Expression, 0)
	p.expect(itemOpenParen)
	if p.peek().typ == itemCloseParen {
		p.expect(itemCloseParen)
		return expr
	}
	expr.Arguments = append(expr.Arguments, p.parseNextExpression())
	for p.peek().typ != itemCloseParen {
		p.expect(itemComma)
		arg := p.parseNextExpression()
		if arg == nil {
			break
		}
		expr.Arguments = append(expr.Arguments, arg)
	}
	p.expect(itemCloseParen)
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
		p.expectStmtEnd()
		return ident
	case itemIdentifier:
		ident := p.parseIdentifier()
		switch p.peek().typ {
		case itemUnaryOperator:
			expr := ast.ExpressionStmt{p.parseOperation(p.parenLevel, ident)}
			p.expectStmtEnd()
			return expr
		case itemAssignmentOperator, itemArrayLookupOperatorLeft:
			n := ast.AssignmentStmt{}
			n.Assignee = ident.(ast.Assignable)
			p.expect(itemAssignmentOperator)
			n.Operator = p.current.val
			p.next()
			n.Value = p.parseExpression()
			p.expectStmtEnd()
			return n
		default:
			expr := ast.ExpressionStmt{ident}
			p.expectStmtEnd()
			return expr
		}
	case itemUnaryOperator:
		expr := ast.ExpressionStmt{p.parseExpression()}
		p.expectStmtEnd()
		return expr
	case itemFunction:
		return p.parseFunctionStmt()
	case itemEcho:
		p.next()
		expr := p.parseExpression()
		p.expectStmtEnd()
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
	case itemSwitch:
		return p.parseSwitch()
	case itemAbstract:
		fallthrough
	case itemClass:
		return p.parseClass()
	case itemInterface:
		return p.parseInterface()
	case itemReturn:
		p.next()
		stmt := ast.ReturnStmt{}
		if p.current.typ != itemStatementEnd {
			stmt.Expression = p.parseExpression()
			p.expectStmtEnd()
		}
		return stmt
	case itemBreak:
		p.next()
		stmt := ast.BreakStmt{}
		if p.current.typ != itemStatementEnd {
			stmt.Expression = p.parseExpression()
			p.expectStmtEnd()
		}
		return stmt
	case itemContinue:
		p.next()
		stmt := ast.ContinueStmt{}
		if p.current.typ != itemStatementEnd {
			stmt.Expression = p.parseExpression()
			p.expectStmtEnd()
		}
		return stmt
	case itemThrow:
		stmt := ast.ThrowStmt{Expression: p.parseNextExpression()}
		p.expectStmtEnd()
		return stmt
	case itemExit:
		stmt := ast.ExitStmt{}
		if p.peek().typ == itemOpenParen {
			p.expect(itemOpenParen)
			stmt.Expression = p.parseNextExpression()
			p.expect(itemCloseParen)
		}
		p.expectStmtEnd()
		return stmt
	case itemTry:
		stmt := &ast.TryStmt{}
		stmt.TryBlock = p.parseBlock()
		p.expect(itemCatch)
		for p.current.typ == itemCatch {
			caught := &ast.CatchStmt{}
			p.expect(itemOpenParen)
			p.expect(itemNonVariableIdentifier)
			caught.CatchType = p.current.val
			p.expect(itemIdentifier)
			caught.CatchVar = ast.NewIdentifier(p.current.val)
			p.expect(itemCloseParen)
			caught.CatchBlock = p.parseBlock()
			stmt.CatchStmts = append(stmt.CatchStmts, caught)
			p.next()
		}
		return stmt
	case itemIgnoreErrorOperator:
		// Ignore this operator
		p.next()
		return p.parseStmt()
	default:
		expr := p.parseExpression()
		if expr != nil {
			p.expectStmtEnd()
			return ast.ExpressionStmt{expr}
		}
		p.errorf("Found %s, expected html or php begin", p.current)
		return nil
	}
}

func (p *parser) expectStmtEnd() {
	if p.peek().typ != itemPHPEnd {
		p.expect(itemStatementEnd)
	}
}

func (p *parser) parseWhile() ast.Statement {
	p.expect(itemOpenParen)
	term := p.parseNextExpression()
	p.expect(itemCloseParen)
	p.next()
	block := p.parseStmt()
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
	first := ast.NewIdentifier(p.current.val)
	if p.peek().typ == itemArrayKeyOperator {
		stmt.Key = first
		p.expect(itemArrayKeyOperator)
		p.expect(itemIdentifier)
		stmt.Value = ast.NewIdentifier(p.current.val)
	} else {
		stmt.Value = first
	}
	p.expect(itemCloseParen)
	p.next()
	stmt.LoopBlock = p.parseStmt()
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
	p.next()
	stmt.LoopBlock = p.parseStmt()
	return stmt
}

func (p *parser) parseDo() ast.Statement {
	block := p.parseBlock()
	p.expect(itemWhile)
	p.expect(itemOpenParen)
	term := p.parseNextExpression()
	p.expect(itemCloseParen)
	p.expectStmtEnd()
	return &ast.DoWhileStmt{
		Termination: term,
		LoopBlock:   block,
	}
}

func (p *parser) parseSwitch() ast.Statement {
	stmt := ast.SwitchStmt{}
	p.expect(itemOpenParen)
	stmt.Expression = p.parseExpression()
	p.expectCurrent(itemCloseParen)
	p.expect(itemBlockBegin)
	p.next()
	for {
		switch p.current.typ {
		case itemCase:
			expr := p.parseNextExpression()
			p.expect(itemTernaryOperator2)
			p.next()
			stmt.Cases = append(stmt.Cases, &ast.SwitchCase{
				Expression: expr,
				Block:      *(p.parseSwitchBlock()),
			})
		case itemDefault:
			p.expect(itemTernaryOperator2)
			p.next()
			stmt.DefaultCase = p.parseSwitchBlock()
		case itemBlockEnd:
			return stmt
		default:
			p.errorf("Unexpected item in switch statement:", p.current)
			return nil
		}
	}
}

func (p *parser) parseSwitchBlock() *ast.Block {
	needBlockEnd := false
	if p.current.typ == itemBlockBegin {
		needBlockEnd = true
		p.next()
	}
	block := &ast.Block{
		Statements: make([]ast.Statement, 0),
	}
stmtLoop:
	for {
		switch p.current.typ {
		case itemBlockEnd:
			if needBlockEnd {
				needBlockEnd = false
				p.next()
			}
			fallthrough
		case itemCase, itemDefault:
			break stmtLoop
		default:
			stmt := p.parseStmt()
			if stmt == nil {
				p.errorf("Invalid statement in switch block", p.current)
				break stmtLoop
			}
			block.Statements = append(block.Statements, stmt)
			p.next()
		}
	}
	if needBlockEnd {
		p.errorf("switch case needs block end")
	}
	return block
}

func (p *parser) parseFunctionStmt() *ast.FunctionStmt {
	stmt := &ast.FunctionStmt{}
	stmt.FunctionDefinition = p.parseFunctionDefinition()
	stmt.Body = p.parseBlock()
	return stmt
}

func (p *parser) parseFunctionDefinition() *ast.FunctionDefinition {
	def := &ast.FunctionDefinition{}
	p.expect(itemNonVariableIdentifier)
	def.Name = p.current.val
	def.Arguments = make([]ast.FunctionArgument, 0)
	p.expect(itemOpenParen)
	if p.peek().typ == itemCloseParen {
		p.expect(itemCloseParen)
		return def
	}
	def.Arguments = append(def.Arguments, p.parseFunctionArgument())
	for {
		switch p.peek().typ {
		case itemComma:
			p.expect(itemComma)
			def.Arguments = append(def.Arguments, p.parseFunctionArgument())
		case itemCloseParen:
			p.expect(itemCloseParen)
			return def
		default:
			p.errorf("unexpected argument separator:", p.current)
			return def
		}
	}
}

func (p *parser) parseFunctionArgument() ast.FunctionArgument {
	arg := ast.FunctionArgument{}
	switch p.peek().typ {
	case itemNonVariableIdentifier, itemArray:
		p.next()
		arg.TypeHint = p.current.val
	}
	p.expect(itemIdentifier)
	arg.Identifier = ast.NewIdentifier(p.current.val)
	if p.peek().typ == itemAssignmentOperator {
		p.expect(itemAssignmentOperator)
		p.next()
		arg.Default = p.parseLiteral()
	}
	return arg
}

func (p *parser) parseBlock() *ast.Block {
	block := &ast.Block{}
	p.expect(itemBlockBegin)
	for p.peek().typ != itemBlockEnd {
		p.next()
		block.Statements = append(block.Statements, p.parseStmt())
	}
	p.next()
	return block
}

func (p *parser) parseClass() ast.Class {
	if p.current.typ == itemAbstract {
		p.expect(itemClass)
	}
	p.expect(itemNonVariableIdentifier)
	name := p.current.val
	if p.peek().typ == itemExtends {
		p.expect(itemExtends)
		p.expect(itemNonVariableIdentifier)
	}
	if p.peek().typ == itemImplements {
		p.expect(itemImplements)
		p.expect(itemNonVariableIdentifier)
		for p.peek().typ == itemComma {
			p.expect(itemComma)
			p.expect(itemNonVariableIdentifier)
		}
	}
	p.expect(itemBlockBegin)
	return p.parseClassFields(ast.Class{Name: name})
}

func (p *parser) parseVisibility() (vis ast.Visibility, found bool) {
	switch p.peek().typ {
	case itemPrivate:
		vis = ast.Private
	case itemPublic:
		vis = ast.Public
	case itemProtected:
		vis = ast.Protected
	default:
		return ast.Public, false
	}
	p.next()
	return vis, true
}

func (p *parser) parseAbstract() bool {
	if p.peek().typ == itemAbstract {
		p.next()
		return true
	}
	return false
}

func (p *parser) parseClassFields(c ast.Class) ast.Class {
	c.Methods = make([]ast.Method, 0)
	c.Properties = make([]ast.Property, 0)
	for p.peek().typ != itemBlockEnd {
		vis, foundVis := p.parseVisibility()
		abstract := p.parseAbstract()
		if foundVis == false {
			vis, foundVis = p.parseVisibility()
		}
		if p.peek().typ == itemStatic {
			p.next()
		}
		p.next()
		switch p.current.typ {
		case itemFunction:
			if abstract {
				f := p.parseFunctionDefinition()
				m := ast.Method{
					Visibility:   vis,
					FunctionStmt: &ast.FunctionStmt{FunctionDefinition: f},
				}
				c.Methods = append(c.Methods, m)
				p.expect(itemStatementEnd)
			} else {
				c.Methods = append(c.Methods, ast.Method{
					Visibility:   vis,
					FunctionStmt: p.parseFunctionStmt(),
				})
			}
		case itemIdentifier:
			prop := ast.Property{
				Visibility: vis,
				Name:       p.current.val,
			}
			if p.peek().typ == itemAssignmentOperator {
				p.expect(itemAssignmentOperator)
				prop.Initialization = p.parseNextExpression()
			}
			c.Properties = append(c.Properties, prop)
			p.expect(itemStatementEnd)
		case itemConst:
			constant := ast.Constant{}
			p.expect(itemNonVariableIdentifier)
			constant.Identifier = ast.NewIdentifier(p.current.val)
			if p.peek().typ == itemAssignmentOperator {
				p.expect(itemAssignmentOperator)
				constant.Value = p.parseNextExpression()
			}
			c.Constants = append(c.Constants, constant)
			p.expect(itemStatementEnd)
		default:
			p.errorf("unexpected class member %v", p.current)
		}
	}
	p.expect(itemBlockEnd)
	return c
}

func (p *parser) parseInterface() *ast.Interface {
	i := &ast.Interface{
		Inherits: make([]string, 0),
	}
	p.expect(itemNonVariableIdentifier)
	i.Name = p.current.val
	if p.peek().typ == itemExtends {
		p.expect(itemExtends)
		for {
			p.expect(itemNonVariableIdentifier)
			i.Inherits = append(i.Inherits, p.current.val)
			if p.peek().typ != itemComma {
				break
			}
			p.next()
		}
	}
	p.expect(itemBlockBegin)
	for p.peek().typ != itemBlockEnd {
		vis, _ := p.parseVisibility()
		if p.peek().typ == itemStatic {
			p.next()
		}
		p.next()
		switch p.current.typ {
		case itemFunction:
			f := p.parseFunctionDefinition()
			m := ast.Method{
				Visibility:   vis,
				FunctionStmt: &ast.FunctionStmt{FunctionDefinition: f},
			}
			i.Methods = append(i.Methods, m)
			p.expect(itemStatementEnd)
		default:
			p.errorf("unexpected interface member %v", p.current)
		}
		p.next()
	}
	return i
}
