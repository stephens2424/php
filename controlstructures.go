package php

import (
	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
)

func (p *Parser) parseIf() *ast.IfStmt {
	p.expect(token.OpenParen)
	n := &ast.IfStmt{}
	p.next()
	n.Condition = p.parseExpression()
	p.expect(token.CloseParen)

	if p.peek().typ == token.TernaryOperator2 {
		p.expect(token.TernaryOperator2)
		n.TrueBranch = p.parseStatementsUntil(token.EndIf, token.ElseIf, token.Else)
	} else {
		p.next()
		n.TrueBranch = p.parseStmt()
		p.next()
	}

	switch p.current.typ {
	case token.ElseIf:
		n.FalseBranch = p.parseIf()
	case token.Else:
		p.next()
		if p.current.typ == token.TernaryOperator2 {
			n.FalseBranch = p.parseStatementsUntil(token.EndIf)
		} else {
			n.FalseBranch = p.parseStmt()
		}
	default:
		n.FalseBranch = ast.Block{}
		p.backup()
	}
	return n
}

func (p *Parser) parseWhile() ast.Statement {
	p.expect(token.OpenParen)
	term := p.parseNextExpression()
	p.expect(token.CloseParen)
	p.next()
	block := p.parseStmt()
	return &ast.WhileStmt{
		Termination: term,
		LoopBlock:   block,
	}
}

func (p *Parser) parseForeach() ast.Statement {
	stmt := &ast.ForeachStmt{}
	p.expect(token.OpenParen)
	stmt.Source = p.parseNextExpression()
	p.expect(token.AsOperator)
	if p.peek().typ == token.AmpersandOperator {
		p.expect(token.AmpersandOperator)
	}
	p.expect(token.VariableOperator)
	p.next()
	first := ast.NewVariable(p.current.val)
	if p.peek().typ == token.ArrayKeyOperator {
		stmt.Key = first
		p.expect(token.ArrayKeyOperator)
		if p.peek().typ == token.AmpersandOperator {
			p.expect(token.AmpersandOperator)
		}
		p.expect(token.VariableOperator)
		p.next()
		stmt.Value = ast.NewVariable(p.current.val)
	} else {
		stmt.Value = first
	}
	p.expect(token.CloseParen)
	p.next()
	stmt.LoopBlock = p.parseStmt()
	return stmt
}

func (p *Parser) parseFor() ast.Statement {
	stmt := &ast.ForStmt{}
	p.expect(token.OpenParen)
	stmt.Initialization = p.parseNextExpression()
	p.expect(token.StatementEnd)
	stmt.Termination = p.parseNextExpression()
	p.expect(token.StatementEnd)
	stmt.Iteration = p.parseNextExpression()
	p.expect(token.CloseParen)
	p.next()
	stmt.LoopBlock = p.parseStmt()
	return stmt
}

func (p *Parser) parseDo() ast.Statement {
	block := p.parseBlock()
	p.expect(token.While)
	p.expect(token.OpenParen)
	term := p.parseNextExpression()
	p.expect(token.CloseParen)
	p.expectStmtEnd()
	return &ast.DoWhileStmt{
		Termination: term,
		LoopBlock:   block,
	}
}

func (p *Parser) parseSwitch() ast.Statement {
	stmt := ast.SwitchStmt{}
	p.expect(token.OpenParen)
	stmt.Expression = p.parseExpression()
	p.expectCurrent(token.CloseParen)
	p.expect(token.BlockBegin)
	p.next()
	for {
		switch p.current.typ {
		case token.Case:
			expr := p.parseNextExpression()
			p.expect(token.TernaryOperator2)
			p.next()
			stmt.Cases = append(stmt.Cases, &ast.SwitchCase{
				Expression: expr,
				Block:      *(p.parseSwitchBlock()),
			})
		case token.Default:
			p.expect(token.TernaryOperator2)
			p.next()
			stmt.DefaultCase = p.parseSwitchBlock()
		case token.BlockEnd:
			return stmt
		default:
			p.errorf("Unexpected token. in switch statement:", p.current)
			return nil
		}
	}
}

func (p *Parser) parseSwitchBlock() *ast.Block {
	needBlockEnd := false
	if p.current.typ == token.BlockBegin {
		needBlockEnd = true
		p.next()
	}
	block := &ast.Block{
		Statements: make([]ast.Statement, 0),
	}
stmtLoop:
	for {
		switch p.current.typ {
		case token.BlockEnd:
			if needBlockEnd {
				needBlockEnd = false
				p.next()
			}
			fallthrough
		case token.Case, token.Default:
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
