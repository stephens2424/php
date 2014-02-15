package ast

type Node interface{}

type Identifier struct {
	Name string
	Type Type
}

func (i Identifier) EvaluatesTo() Type {
	return i.Type
}

func NewIdentifier(name string) Identifier {
	return Identifier{name, AnyType}
}

type Statement interface{}
type Expression interface {
	EvaluatesTo() Type
}

const AnyType = String | Integer | Float | Boolean | Null | Resource | Array | Object

type OperatorExpression struct {
	Operand1 Expression
	Operand2 Expression
	Operand3 Expression
	Type     Type
}

func (o OperatorExpression) EvaluatesTo() Type {
	return o.Type
}

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

func (f FunctionCallExpression) EvaluatesTo() Type {
	return String | Integer | Float | Boolean | Null | Resource | Array | Object
}

type Block struct {
	Statements []Statement
	Scope      Scope
}

type FunctionStmt struct {
	FunctionDefinition
	Body Block
}

type FunctionDefinition struct {
	Name      string
	Arguments []FunctionArgument
}

type FunctionArgument struct {
	TypeHint   string
	Identifier Identifier
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
	*FunctionStmt
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
	TrueBlock  Statement
	FalseBlock Statement
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

func (l Literal) EvaluatesTo() Type {
	return l.Type
}

type ForeachStmt struct {
	Source    Expression
	Key       *Identifier
	Value     Identifier
	LoopBlock Block
}
