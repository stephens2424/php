package ast

import "fmt"

// Node encapsulates every AST node.
type Node interface {
	String() string
	Children() []Node
	Position() Position
}

type Position int

type BaseNode struct {
	pos int
}

func (b BaseNode) Position() Position {
	return Position(b.pos)
}

func (b BaseNode) Children() []Node {
	return nil
}

func (b BaseNode) String() string {
	return ""
}

// An Identifier is specifically a variable in code.
type Identifier struct {
	BaseNode
	Name string
	Type Type
}

func (i Identifier) String() string {
	return i.Name
}

type GlobalIdentifier struct {
	*Identifier
}

func (i Identifier) AssignableType() Type {
	return i.Type
}

// EvaluatesTo returns the known type of the variable.
func (i Identifier) EvaluatesTo() Type {
	return i.Type
}

// NewIdentifier intializes an identifier node with its type set to AnyType.
func NewIdentifier(name string) *Identifier {
	return &Identifier{Name: name, Type: AnyType}
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

func (o OperatorExpression) Children() []Node {
	return []Node{
		o.Operand1,
		o.Operand2,
		o.Operand3,
	}
}

func (o OperatorExpression) String() string {
	return o.Operator
}

func (o OperatorExpression) EvaluatesTo() Type {
	return o.Type
}

type ExpressionStmt struct {
	Expression
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

func (e EchoStmt) String() string {
	return "Echo"
}

func (e EchoStmt) Children() []Node {
	return []Node{e.Expression}
}

// ReturnStmt represents a function return.
type ReturnStmt struct {
	Expression
}
type BreakStmt struct {
	Expression
}
type ContinueStmt struct {
	Expression
}

type ThrowStmt struct {
	Expression
}

type IncludeStmt struct {
	Include
}

type Include struct {
	Expression
}

func (i Include) EvaluatesTo() Type {
	return AnyType
}

type ExitStmt struct {
	Expression
}

type NewExpression struct {
	Expression
}

type AssignmentExpression struct {
	BaseNode
	Assignee Assignable
	Value    Expression
	Operator string
}

func (a AssignmentExpression) String() string {
	return a.Operator
}

func (a AssignmentExpression) Children() []Node {
	return []Node{
		a.Assignee,
		a.Value,
	}
}

func (a AssignmentExpression) EvaluatesTo() Type {
	return a.Value.EvaluatesTo()
}

// AssignmentStmt represents an assignment.
type AssignmentStmt struct {
	AssignmentExpression
}

type Assignable interface {
	Node
	AssignableType() Type
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

func (b Block) Children() []Node {
	n := make([]Node, len(b.Statements))
	for i, s := range b.Statements {
		n[i] = s
	}
	return n
}

type FunctionStmt struct {
	BaseNode
	*FunctionDefinition
	Body *Block
}

func (f FunctionStmt) String() string {
	return fmt.Sprintf("Func: %s", f.Name)
}

func (f FunctionStmt) Children() []Node {
	return []Node{
		f.FunctionDefinition,
		f.Body,
	}
}

type FunctionDefinition struct {
	BaseNode
	Name      string
	Arguments []FunctionArgument
}

func (fd FunctionDefinition) Children() []Node {
	n := make([]Node, len(fd.Arguments))
	for i, arg := range fd.Arguments {
		n[i] = arg
	}
	return n
}

type FunctionArgument struct {
	BaseNode
	TypeHint   string
	Default    *Literal
	Identifier *Identifier
}

func (fa FunctionArgument) String() string {
	return fmt.Sprintf("Arg: %s", fa.TypeHint)
}

func (fa FunctionArgument) Children() []Node {
	n := []Node{
		fa.Identifier,
	}
	if fa.Default != nil {
		n = append(n, fa.Default)
	}
	return n
}

type Class struct {
	BaseNode
	Name       string
	Extends    *Class
	Implements []string
	Methods    []Method
	Properties []Property
}

type Constant struct {
	BaseNode
	*Identifier
	Value interface{}
}

type ConstantExpression struct {
	*Identifier
}

type Interface struct {
	BaseNode
	Name     string
	Inherits []string
	Methods  []Method
}

type Property struct {
	BaseNode
	Name           string
	Visibility     Visibility
	Type           Type
	Initialization Expression
}

func (p Property) AssignableType() Type {
	return p.Type
}

type PropertyExpression struct {
	BaseNode
	Receiver Expression
	Name     string
	Type     Type
}

func (p PropertyExpression) AssignableType() Type {
	return p.Type
}

func (p PropertyExpression) EvaluatesTo() Type {
	return AnyType
}

type ClassExpression struct {
	BaseNode
	Receiver   string
	Expression Expression
}

func (c ClassExpression) EvaluatesTo() Type {
	return AnyType
}

type Method struct {
	BaseNode
	*FunctionStmt
	Visibility Visibility
}

type MethodCallExpression struct {
	Receiver Expression
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
	Condition   Expression
	TrueBranch  Statement
	FalseBranch Statement
}

func (i IfStmt) String() string {
	return "if"
}

func (i IfStmt) Children() []Node {
	return []Node{
		i.Condition,
		i.TrueBranch,
		i.FalseBranch,
	}
}

type SwitchStmt struct {
	BaseNode
	Expression  Expression
	Cases       []*SwitchCase
	DefaultCase *Block
}

type SwitchCase struct {
	BaseNode
	Expression Expression
	Block      Block
}

type ForStmt struct {
	BaseNode
	Initialization Expression
	Termination    Expression
	Iteration      Expression
	LoopBlock      Statement
}

type WhileStmt struct {
	BaseNode
	Termination Expression
	LoopBlock   Statement
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
	CatchBlock *Block
	CatchType  string
	CatchVar   *Identifier
}

type Literal struct {
	BaseNode
	Type Type
}

func (l Literal) String() string {
	return fmt.Sprintf("Literal-%s", l.Type)
}

func (l Literal) EvaluatesTo() Type {
	return l.Type
}

type ForeachStmt struct {
	BaseNode
	Source    Expression
	Key       *Identifier
	Value     *Identifier
	LoopBlock Statement
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

func (a ArrayExpression) AssignableType() Type {
	return AnyType
}

type ArrayLookupExpression struct {
	BaseNode
	Array Expression
	Index Expression
}

func (a ArrayLookupExpression) EvaluatesTo() Type {
	return AnyType
}

func (a ArrayLookupExpression) AssignableType() Type {
	return AnyType
}

type ArrayAppendExpression struct {
	BaseNode
	Array Expression
}

func (a ArrayAppendExpression) EvaluatesTo() Type {
	return AnyType
}

func (a ArrayAppendExpression) AssignableType() Type {
	return AnyType
}
