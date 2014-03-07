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

func (p *parser) Parse() (nodes []ast.Node, errors []error) {
	defer func() {
		if r := recover(); r != nil {
			if p.Debug {
				for _, err := range p.errors {
					fmt.Println(err)
				}
				panic(r)
			}
			errors = append([]error{fmt.Errorf("%s", r)}, p.errors...)
		}
	}()
	// expecting either itemHTML or itemPHPBegin
	nodes = make([]ast.Node, 0, 1)
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
	errors = p.errors
	return nodes, errors
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

func (p *parser) parseNextExpression() ast.Expression {
	p.next()
	return p.parseExpression()
}

func (p *parser) parseFunctionCall() *ast.FunctionCallExpression {
	expr := &ast.FunctionCallExpression{}
	if p.current.typ != itemIdentifier {
		p.expected(itemIdentifier)
	}
	expr.FunctionName = ast.FunctionNameExpression{Name: p.current.val}
	return p.parseFunctionArguments(expr)
}

func (p *parser) parseFunctionArguments(expr *ast.FunctionCallExpression) *ast.FunctionCallExpression {
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
		g := &ast.GlobalDeclaration{
			Identifiers: make([]*ast.Variable, 0, 1),
		}
		for p.current.typ == itemVariableOperator {
			variable, ok := p.parseVariable().(*ast.Variable)
			if !ok {
				p.errorf("global declarations must be of standard variables")
				break
			}
			g.Identifiers = append(g.Identifiers, variable)
			if p.peek().typ != itemComma {
				break
			}
			p.expect(itemComma)
			p.next()
		}
		p.expectStmtEnd()
		return g
	case itemNamespace:
		p.expect(itemIdentifier)
		p.expectStmtEnd()
		// We are ignoring this for now
		return nil
	case itemUse:
		p.expect(itemIdentifier)
		if p.peek().typ == itemAsOperator {
			p.expect(itemAsOperator)
			p.expect(itemIdentifier)
		}
		p.expectStmtEnd()
		// We are ignoring this for now
		return nil
	case itemVariableOperator:
		ident := p.parseIdentifier()
		switch p.peek().typ {
		case itemUnaryOperator:
			expr := ast.ExpressionStmt{p.parseOperation(p.parenLevel, ident)}
			p.expectStmtEnd()
			return expr
		case itemAssignmentOperator, itemArrayLookupOperatorLeft:
			n := ast.AssignmentStmt{}
			var ok bool
			n.Assignee, ok = ident.(ast.Assignable)
			if !ok {
				p.errorf("%s is not assignable", n.Assignee)
			}
			p.expect(itemAssignmentOperator)
			n.Operator = p.current.val
			p.next()
			n.Value = p.parseExpression()
			p.expectStmtEnd()
			return n
		case itemOpenParen:
			var expr ast.Expression
			expr = p.parseFunctionArguments(&ast.FunctionCallExpression{
				FunctionName: ident,
			})
			if p.peek().typ == itemObjectOperator {
				expr = p.parseObjectLookup(expr)
			}
			p.expectStmtEnd()
			return expr
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
		exprs := []ast.Expression{
			p.parseNextExpression(),
		}
		for p.peek().typ == itemComma {
			p.expect(itemComma)
			exprs = append(exprs, p.parseNextExpression())
		}
		p.expectStmtEnd()
		return ast.Echo(exprs...)
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
			if p.peek().typ != itemCloseParen {
				stmt.Expression = p.parseNextExpression()
			}
			p.expect(itemCloseParen)
		}
		p.expectStmtEnd()
		return stmt
	case itemTry:
		stmt := &ast.TryStmt{}
		stmt.TryBlock = p.parseBlock()
		for p.expect(itemCatch); p.current.typ == itemCatch; p.next() {
			caught := &ast.CatchStmt{}
			p.expect(itemOpenParen)
			p.expect(itemIdentifier)
			caught.CatchType = p.current.val
			p.expect(itemVariableOperator)
			p.expect(itemIdentifier)
			caught.CatchVar = ast.NewVariable(p.current.val)
			p.expect(itemCloseParen)
			caught.CatchBlock = p.parseBlock()
			stmt.CatchStmts = append(stmt.CatchStmts, caught)
		}
		p.backup()
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
func (p *parser) parseFunctionStmt() *ast.FunctionStmt {
	stmt := &ast.FunctionStmt{}
	stmt.FunctionDefinition = p.parseFunctionDefinition()
	stmt.Body = p.parseBlock()
	return stmt
}

func (p *parser) parseFunctionDefinition() *ast.FunctionDefinition {
	def := &ast.FunctionDefinition{}
	p.expect(itemIdentifier)
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
	case itemIdentifier, itemArray:
		p.next()
		arg.TypeHint = p.current.val
	}
	if p.peek().typ == itemAmpersandOperator {
		p.next()
	}
	p.expect(itemVariableOperator)
	p.next()
	arg.Variable = ast.NewVariable(p.current.val)
	if p.peek().typ == itemAssignmentOperator {
		p.expect(itemAssignmentOperator)
		p.next()
		arg.Default = p.parseExpression()
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
