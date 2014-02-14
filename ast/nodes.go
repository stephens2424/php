package ast

type Identifier struct {
	Name string
}

type Statement interface{}

type EchoStmt Expression

type AssignmentStmt struct {
	Assignee Identifier
	Value    Expression
}

type FunctionCallStmt struct {
	FunctionCallExpression
}

type FunctionCallExpression struct {
	FunctionName string
	Arguments    []Expression
}

type Expression interface{}

type Block struct {
	Statments []Statement
	Scope     Scope
}

type FunctionStmt struct {
	FunctionDefinition
	Body Block
}

type FunctionDefinition struct {
	Name      string
	Arguments []Identifier
}

type Class struct {
	Name       string
	Extends    *Class
	Implements []*Interface
	Methods    []Method
}

type Constant struct {
	Identifier
	Value interface{}
}

type Interface struct {
	Methods []FunctionDefinition
}

type Method struct {
	FunctionStmt
	Visibility Visibility
}

type Visibility int

const (
	Private Visibility = iota
	Protected
	Public
)

type IfStmt struct {
	Condition  Expression
	TrueBlock  Block
	FalseBlock Block
}

type ForStmt struct {
	Initialization Expression
	Termination    Expression
	Iteration      Expression
	LoopBlock      Block
}

type WhileStmt struct {
	Termination Expression
	LoopBlock   Block
}

type DoWhileStmt struct {
	Termination Expression
	LoopBlock   Block
}

type TryStmt struct {
	TryBlock     *Block
	FinallyBlock *Block
	CatchStmts   []*CatchStmt
}

type CatchStmt struct {
	CatchBlock Block
	CatchType  Type
}

type Literal struct {
	Type Type
}

type ForeachStmt struct {
	Source    Expression
	Key       *Identifier
	Value     Identifier
	LoopBlock Block
}
