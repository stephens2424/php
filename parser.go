package php

import (
	"fmt"

	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
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

	Debug       bool
	PrintTokens bool
	MaxErrors   int
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
			errors = append([]error{fmt.Errorf("%s", r)}, p.errors...)
			if p.Debug {
				for _, err := range p.errors {
					fmt.Println(err)
				}
				panic(r)
			}
		}
	}()
	// expecting either token.HTML or token.PHPBegin
	nodes = make([]ast.Node, 0, 1)
TokenLoop:
	for {
		p.next()
		switch p.current.typ {
		case token.EOF:
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
	case token.HTML:
		return ast.Echo(ast.Literal{Type: ast.String})
	case token.PHPBegin:
		return nil
	case token.PHPEnd:
		return nil
	}
	return p.parseStmt()
}

func (p *parser) next() {
	p.idx += 1
	if len(p.previous) <= p.idx {
		p.current = p.lexer.nextItem()
		if p.PrintTokens {
			fmt.Println(p.current)
		}
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

func (p *parser) expectCurrent(i ...token.Token) {
	for _, typ := range i {
		if p.current.typ == typ {
			return
		}
	}
	p.expected(i...)
}

func (p *parser) expectAndNext(i ...token.Token) {
	defer p.next()
	for _, typ := range i {
		if p.current.typ == typ {
			return
		}
	}
	p.expected(i...)
}

func (p *parser) expect(i ...token.Token) {
	p.next()
	p.expectCurrent(i...)
}

func (p *parser) expected(i ...token.Token) {
	p.errorf("Found %s, expected %s", p.current, i)
}

func (p *parser) errorf(str string, args ...interface{}) {
	if p.errorCount > p.MaxErrors {
		panic("too many errors")
	}
	if _, ok := p.errorMap[p.current.pos.Line]; ok {
		return
	}
	errString := fmt.Sprintf(str, args...)
	p.errorCount += 1
	p.errors = append(p.errors, fmt.Errorf("%s: %s", p.errorPrefix(), errString))
	p.errorMap[p.current.pos.Line] = true
}

func (p *parser) errorPrefix() string {
	return fmt.Sprintf("%s %d", p.lexer.file, p.current.pos.Line)
}

func (p *parser) parseNextExpression() ast.Expression {
	p.next()
	return p.parseExpression()
}

func (p *parser) parseStmt() ast.Statement {
	switch p.current.typ {
	case token.BlockBegin:
		p.backup()
		return p.parseBlock()
	case token.Global:
		p.next()
		g := &ast.GlobalDeclaration{
			Identifiers: make([]*ast.Variable, 0, 1),
		}
		for p.current.typ == token.VariableOperator {
			variable, ok := p.parseVariable().(*ast.Variable)
			if !ok {
				p.errorf("global declarations must be of standard variables")
				break
			}
			g.Identifiers = append(g.Identifiers, variable)
			if p.peek().typ != token.Comma {
				break
			}
			p.expect(token.Comma)
			p.next()
		}
		p.expectStmtEnd()
		return g
	case token.Namespace:
		p.expect(token.Identifier)
		p.expectStmtEnd()
		// We are ignoring this for now
		return nil
	case token.Use:
		p.expect(token.Identifier)
		if p.peek().typ == token.AsOperator {
			p.expect(token.AsOperator)
			p.expect(token.Identifier)
		}
		p.expectStmtEnd()
		// We are ignoring this for now
		return nil
	case token.VariableOperator, token.UnaryOperator:
		expr := ast.ExpressionStmt{p.parseExpression()}
		p.expectStmtEnd()
		return expr
	case token.Function:
		return p.parseFunctionStmt()
	case token.PHPEnd:
		if p.peek().typ == token.EOF {
			return nil
		}
		p.expect(token.HTML)
		expr := ast.Echo(&ast.Literal{Type: ast.String})
		p.next()
		if p.current.typ != token.EOF {
			p.expectCurrent(token.PHPBegin)
		}
		return expr
	case token.Echo:
		exprs := []ast.Expression{
			p.parseNextExpression(),
		}
		for p.peek().typ == token.Comma {
			p.expect(token.Comma)
			exprs = append(exprs, p.parseNextExpression())
		}
		p.expectStmtEnd()
		return ast.Echo(exprs...)
	case token.If:
		return p.parseIf()
	case token.While:
		return p.parseWhile()
	case token.Do:
		return p.parseDo()
	case token.For:
		return p.parseFor()
	case token.Foreach:
		return p.parseForeach()
	case token.Switch:
		return p.parseSwitch()
	case token.Abstract:
		fallthrough
	case token.Class:
		return p.parseClass()
	case token.Interface:
		return p.parseInterface()
	case token.Return:
		p.next()
		stmt := ast.ReturnStmt{}
		if p.current.typ != token.StatementEnd {
			stmt.Expression = p.parseExpression()
			p.expectStmtEnd()
		}
		return stmt
	case token.Break:
		p.next()
		stmt := ast.BreakStmt{}
		if p.current.typ != token.StatementEnd {
			stmt.Expression = p.parseExpression()
			p.expectStmtEnd()
		}
		return stmt
	case token.Continue:
		p.next()
		stmt := ast.ContinueStmt{}
		if p.current.typ != token.StatementEnd {
			stmt.Expression = p.parseExpression()
			p.expectStmtEnd()
		}
		return stmt
	case token.Throw:
		stmt := ast.ThrowStmt{Expression: p.parseNextExpression()}
		p.expectStmtEnd()
		return stmt
	case token.Exit:
		stmt := ast.ExitStmt{}
		if p.peek().typ == token.OpenParen {
			p.expect(token.OpenParen)
			if p.peek().typ != token.CloseParen {
				stmt.Expression = p.parseNextExpression()
			}
			p.expect(token.CloseParen)
		}
		p.expectStmtEnd()
		return stmt
	case token.Try:
		stmt := &ast.TryStmt{}
		stmt.TryBlock = p.parseBlock()
		for p.expect(token.Catch); p.current.typ == token.Catch; p.next() {
			caught := &ast.CatchStmt{}
			p.expect(token.OpenParen)
			p.expect(token.Identifier)
			caught.CatchType = p.current.val
			p.expect(token.VariableOperator)
			p.expect(token.Identifier)
			caught.CatchVar = ast.NewVariable(p.current.val)
			p.expect(token.CloseParen)
			caught.CatchBlock = p.parseBlock()
			stmt.CatchStmts = append(stmt.CatchStmts, caught)
		}
		p.backup()
		return stmt
	case token.IgnoreErrorOperator:
		// Ignore this operator
		p.next()
		return p.parseStmt()
	default:
		expr := p.parseExpression()
		if expr != nil {
			p.expectStmtEnd()
			return ast.ExpressionStmt{expr}
		}
		p.errorf("Found %s, statement or expression", p.current)
		return nil
	}
}

func (p *parser) expectStmtEnd() {
	if p.peek().typ != token.PHPEnd {
		p.expect(token.StatementEnd)
	}
}

func (p *parser) parseBlock() *ast.Block {
	p.expect(token.BlockBegin)
	b := p.parseStatementsUntil(token.BlockEnd)
	p.expectCurrent(token.BlockEnd)
	return b
}

func (p *parser) parseStatementsUntil(endTokens ...token.Token) *ast.Block {
	block := &ast.Block{}
	breakTypes := map[token.Token]bool{}
	for _, typ := range endTokens {
		breakTypes[typ] = true
	}
	for {
		p.next()
		if _, ok := breakTypes[p.current.typ]; ok {
			break
		}
		stmt := p.parseStmt()
		if stmt == nil {
			return block
		}
		block.Statements = append(block.Statements, stmt)
	}
	return block
}
