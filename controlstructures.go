package php

import "stephensearles.com/php/ast"

func (p *parser) parseIf() *ast.IfStmt {
	p.expect(itemOpenParen)
	n := &ast.IfStmt{}
	p.next()
	n.Condition = p.parseExpression()
	p.expect(itemCloseParen)
	p.next()
	n.TrueBranch = p.parseStmt()
	p.next()
	switch p.current.typ {
	case itemElseIf:
		n.FalseBranch = p.parseIf()
	case itemElse:
		p.next()
		n.FalseBranch = p.parseStmt()
	default:
		n.FalseBranch = ast.Block{}
		p.backup()
	}
	return n
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
	if p.peek().typ == itemAmpersandOperator {
		p.expect(itemAmpersandOperator)
	}
	p.expect(itemVariableOperator)
	p.next()
	first := ast.NewVariable("$" + p.current.val)
	if p.peek().typ == itemArrayKeyOperator {
		stmt.Key = first
		p.expect(itemArrayKeyOperator)
		if p.peek().typ == itemAmpersandOperator {
			p.expect(itemAmpersandOperator)
		}
		p.expect(itemVariableOperator)
		p.next()
		stmt.Value = ast.NewVariable("$" + p.current.val)
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
