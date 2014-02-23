package ast

import (
	"fmt"
)

// Node encapsulates every AST node.
type Node interface {
	Position() Position
}

type Position int

type BaseNode struct {
	pos int
}

func (b BaseNode) Position() Position {
	return Position(b.pos)
}

// An Identifier is specifically a variable in code.
type Identifier struct {
	BaseNode
	Name string
	Type Type
}

// EvaluatesTo returns the known type of the variable.
func (i Identifier) EvaluatesTo() Type {
	return i.Type
}

// NewIdentifier intializes an identifier node with its type set to AnyType.
func NewIdentifier(name string) Identifier {
	return Identifier{Name: name, Type: AnyType}
}

// A statement is an executable piece of code. It may be as simple as
// a function call or a variable assignment. It also includes things like
// "if".
type Statement interface {
	Node
}

// An Expression is a snippet of code that evaluates to a single value when run
// and does not represent a program instruction.
type Expression interface {
	Node
	EvaluatesTo() Type
}

// AnyType is a bitmask of all the valid types.
const AnyType = String | Integer | Float | Boolean | Null | Resource | Array | Object

// OperatorExpression is an expression that applies an operator to one, two, or three
// operands. The operator determines how many operands it should contain.
type OperatorExpression struct {
	BaseNode
	Operand1 Expression
	Operand2 Expression
	Operand3 Expression
	Type     Type
	Operator string
}

func (o OperatorExpression) String() string {
	if o.Operand2 == nil {
		return fmt.Sprintf("(%s%v~%v)", o.Operand1, o.Operator, o.Type)
	}
	if o.Operand3 == nil {
		return fmt.Sprintf("(%s %v %s~%v)", o.Operand1, o.Operator, o.Operand2, o.Type)
	}
	return fmt.Sprintf("(%s ? %s : %s~%v)", o.Operand1, o.Operand2, o.Operand3, o.Type)
}

func (o OperatorExpression) EvaluatesTo() Type {
	return o.Type
}

// Echo returns a new echo statement.
func Echo(expr Expression) EchoStmt {
	return EchoStmt{Expression: expr}
}

// Echo represents an echo statement. It may be either a literal statement
// or it may be from data outside PHP-mode, such as "here" in: <? not here ?> here <? not here ?>
type EchoStmt struct {
	BaseNode
	Expression Expression
}

// ReturnStmt represents a function return.
type ReturnStmt struct {
	Expression
}

type ThrowStmt struct {
	Expression
}

type NewExpression struct {
	Expression
}

// AssignmentStmt represents an assignment.
type AssignmentStmt struct {
	BaseNode
	Assignee Identifier
	Value    Expression
	Operator string
}

type FunctionCallStmt struct {
	FunctionCallExpression
}

type FunctionCallExpression struct {
	BaseNode
	FunctionName string
	Arguments    []Expression
}

func (f FunctionCallExpression) EvaluatesTo() Type {
	return String | Integer | Float | Boolean | Null | Resource | Array | Object
}

type Block struct {
	BaseNode
	Statements []Statement
	Scope      Scope
}

type FunctionStmt struct {
	BaseNode
	FunctionDefinition
	Body Block
}

type FunctionDefinition struct {
	BaseNode
	Name      string
	Arguments []FunctionArgument
}

type FunctionArgument struct {
	BaseNode
	TypeHint   string
	Identifier Identifier
}

type Class struct {
	BaseNode
	Name       string
	Extends    *Class
	Implements []*Interface
	Methods    []Method
	Properties []Property
}

type Constant struct {
	BaseNode
	Identifier
	Value interface{}
}

type ConstantExpression struct {
	BaseNode
	Identifier
}

type Interface struct {
	BaseNode
	Methods []FunctionDefinition
}

type Property struct {
	BaseNode
	Name       string
	Visibility Visibility
}

type PropertyExpression struct {
	BaseNode
	Receiver Identifier
	Name     string
}

func (p PropertyExpression) EvaluatesTo() Type {
	return AnyType
}

type Method struct {
	BaseNode
	*FunctionStmt
	Visibility Visibility
}

type MethodCallExpression struct {
	Receiver Identifier
	*FunctionCallExpression
}

type Visibility int

const (
	Private Visibility = iota
	Protected
	Public
)

type IfStmt struct {
	BaseNode
	Condition  Expression
	TrueBlock  Statement
	FalseBlock Statement
}

type ForStmt struct {
	BaseNode
	Initialization Expression
	Termination    Expression
	Iteration      Expression
	LoopBlock      Block
}

type WhileStmt struct {
	BaseNode
	Termination Expression
	LoopBlock   Block
}

type DoWhileStmt struct {
	BaseNode
	Termination Expression
	LoopBlock   Statement
}

type TryStmt struct {
	BaseNode
	TryBlock     *Block
	FinallyBlock *Block
	CatchStmts   []*CatchStmt
}

type CatchStmt struct {
	BaseNode
	CatchBlock Block
	CatchType  Type
}

type Literal struct {
	BaseNode
	Type Type
}

func (l Literal) EvaluatesTo() Type {
	return l.Type
}

type ForeachStmt struct {
	BaseNode
	Source    Expression
	Key       *Identifier
	Value     Identifier
	LoopBlock Block
}

type ArrayExpression struct {
	BaseNode
	ArrayType
	Pairs []ArrayPair
}

type ArrayPair struct {
	Key   Expression
	Value Expression
}

func (a ArrayExpression) EvaluatesTo() Type {
	return Array
}

type ArrayLookupExpression struct {
	BaseNode
	Array Identifier
	Index Expression
}

func (a ArrayLookupExpression) EvaluatesTo() Type {
	return AnyType
}
