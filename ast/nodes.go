package ast

import (
	"fmt"
	"strings"
)

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

type GlobalDeclaration struct {
	BaseNode
	Identifiers []*Identifier
}

func (g GlobalDeclaration) Children() []Node {
	n := make([]Node, len(g.Identifiers))
	for i, node := range g.Identifiers {
		n[i] = node
	}
	return n
}

func (g GlobalDeclaration) String() string {
	return "global"
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
	n := make([]Node, 0, 3)
	if o.Operand1 != nil {
		n = append(n, o.Operand1)
	}
	if o.Operand2 != nil {
		n = append(n, o.Operand2)
	}
	if o.Operand3 != nil {
		n = append(n, o.Operand3)
	}
	return n
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

func (e ExpressionStmt) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

func (e ExpressionStmt) Children() []Node {
	if e.Expression != nil {
		return e.Expression.Children()
	}
	return nil
}

// Echo returns a new echo statement.
func Echo(exprs ...Expression) EchoStmt {
	return EchoStmt{Expressions: exprs}
}

// Echo represents an echo statement. It may be either a literal statement
// or it may be from data outside PHP-mode, such as "here" in: <? not here ?> here <? not here ?>
type EchoStmt struct {
	BaseNode
	Expressions []Expression
}

func (e EchoStmt) String() string {
	return "Echo"
}

func (e EchoStmt) Children() []Node {
	nodes := make([]Node, len(e.Expressions))
	for i, expr := range e.Expressions {
		nodes[i] = expr
	}
	return nodes
}

// ReturnStmt represents a function return.
type ReturnStmt struct {
	Expression
}

func (r ReturnStmt) String() string {
	return fmt.Sprintf("return")
}

func (r ReturnStmt) Children() []Node {
	if r.Expression == nil {
		return nil
	}
	return []Node{r.Expression}
}

type BreakStmt struct {
	Expression
}

func (b BreakStmt) Children() []Node {
	if b.Expression != nil {
		return b.Expression.Children()
	}
	return nil
}

func (b BreakStmt) String() string {
	return "break"
}

type ContinueStmt struct {
	Expression
}

func (c ContinueStmt) String() string {
	return "continue"
}

func (c ContinueStmt) Children() []Node {
	if c.Expression != nil {
		return c.Expression.Children()
	}
	return nil
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
	BaseNode
	Class     Expression
	Arguments []Expression
}

func (n NewExpression) EvaluatesTo() Type {
	return Object
}

type ClassIdentifier struct {
	BaseNode
	ClassName string
}

func (c ClassIdentifier) EvaluatesTo() Type {
	return String
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

func (f FunctionCallExpression) String() string {
	return fmt.Sprintf("%s()", f.FunctionName)
}

func (f FunctionCallExpression) Children() []Node {
	n := make([]Node, len(f.Arguments))
	for i, a := range f.Arguments {
		n[i] = a
	}
	return n
}

type Block struct {
	BaseNode
	Statements []Statement
	Scope      Scope
}

func (b Block) String() string {
	return "{}"
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
	n := make([]Node, 0, 2)
	if f.FunctionDefinition != nil {
		n = append(n, f.FunctionDefinition)
	}
	if f.Body != nil {
		n = append(n, f.Body)
	}
	return n
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
	Default    Expression
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
	Constants  []Constant
}

func (c Class) String() string {
	return fmt.Sprintf("class %s", c.Name)
}

func (c Class) Children() []Node {
	n := make([]Node, len(c.Methods)+len(c.Properties))
	for i, p := range c.Properties {
		n[i] = p
	}
	offset := len(c.Properties)
	for i, m := range c.Methods {
		n[i+offset] = m
	}
	return n
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

func (i Interface) String() string {
	return fmt.Sprintf("interface %s extends %s", i.Name, strings.Join(i.Inherits, ", "))
}

func (i Interface) Children() []Node {
	n := make([]Node, len(i.Methods))
	for ii, method := range i.Methods {
		n[ii] = method
	}
	return n
}

type Property struct {
	BaseNode
	Name           string
	Visibility     Visibility
	Type           Type
	Initialization Expression
}

func (p Property) String() string {
	return fmt.Sprintf("Prop: %s", p.Name)
}

func (p Property) AssignableType() Type {
	return p.Type
}

type PropertyExpression struct {
	BaseNode
	Receiver Expression
	Name     Expression
	Type     Type
}

func (p PropertyExpression) String() string {
	return fmt.Sprintf("%s->%s", p.Receiver, p.Name)
}

func (p PropertyExpression) AssignableType() Type {
	return p.Type
}

func (p PropertyExpression) EvaluatesTo() Type {
	return AnyType
}

func (p PropertyExpression) Children() []Node {
	return []Node{
		p.Receiver,
	}
}

type PropertyIdentifier struct {
	BaseNode
	Name string
}

func (p PropertyIdentifier) EvaluatesTo() Type {
	return String
}

type ClassExpression struct {
	BaseNode
	Receiver   string
	Expression Expression
}

func (c ClassExpression) EvaluatesTo() Type {
	return AnyType
}

func (c ClassExpression) String() string {
	return fmt.Sprintf("%s::", c.Receiver)
}

func (c ClassExpression) Children() []Node {
	return []Node{c.Expression}
}

type Method struct {
	BaseNode
	*FunctionStmt
	Visibility Visibility
}

func (m Method) String() string {
	return m.Name
}

func (m Method) Children() []Node {
	return m.FunctionStmt.Children()
}

type MethodCallExpression struct {
	Receiver Expression
	*FunctionCallExpression
}

func (m MethodCallExpression) Children() []Node {
	return []Node{
		m.Receiver,
		m.FunctionCallExpression,
	}
}

func (m MethodCallExpression) String() string {
	return fmt.Sprintf("%s->", m.Receiver)
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
	n := make([]Node, 0, 3)
	n = append(n, i.Condition)
	n = append(n, i.TrueBranch)
	if i.FalseBranch != nil {
		n = append(n, i.FalseBranch)
	}
	return n
}

type SwitchStmt struct {
	BaseNode
	Expression  Expression
	Cases       []*SwitchCase
	DefaultCase *Block
}

func (s SwitchStmt) String() string {
	return "switch"
}

func (s SwitchStmt) Children() []Node {
	n := []Node{
		s.Expression,
	}
	for _, c := range s.Cases {
		n = append(n, c)
	}
	if s.DefaultCase != nil {
		n = append(n, s.DefaultCase)
	}
	return n
}

type SwitchCase struct {
	BaseNode
	Expression Expression
	Block      Block
}

func (s SwitchCase) String() string {
	return "case"
}

func (s SwitchCase) Children() []Node {
	return []Node{
		s.Expression,
		s.Block,
	}
}

type ForStmt struct {
	BaseNode
	Initialization Expression
	Termination    Expression
	Iteration      Expression
	LoopBlock      Statement
}

func (f ForStmt) String() string {
	return "for"
}

func (f ForStmt) Children() []Node {
	return []Node{
		f.Initialization,
		f.Termination,
		f.Iteration,
		f.LoopBlock,
	}
}

type WhileStmt struct {
	BaseNode
	Termination Expression
	LoopBlock   Statement
}

func (w WhileStmt) String() string {
	return fmt.Sprintf("while")
}

func (w WhileStmt) Children() []Node {
	return []Node{
		w.Termination,
		w.LoopBlock,
	}
}

type DoWhileStmt struct {
	BaseNode
	Termination Expression
	LoopBlock   Statement
}

func (d DoWhileStmt) String() string {
	return fmt.Sprintf("do ... while")
}

func (d DoWhileStmt) Children() []Node {
	return []Node{
		d.LoopBlock,
		d.Termination,
	}
}

type TryStmt struct {
	BaseNode
	TryBlock     *Block
	FinallyBlock *Block
	CatchStmts   []*CatchStmt
}

func (t TryStmt) String() string {
	return "try"
}

func (t TryStmt) Children() []Node {
	n := []Node{t.TryBlock}
	for _, catch := range t.CatchStmts {
		n = append(n, catch)
	}
	if t.FinallyBlock != nil {
		n = append(n, t.FinallyBlock)
	}
	return n
}

type CatchStmt struct {
	BaseNode
	CatchBlock *Block
	CatchType  string
	CatchVar   *Identifier
}

func (c CatchStmt) String() string {
	return fmt.Sprintf("catch %s %s", c.CatchType, c.CatchVar)
}

func (c CatchStmt) Children() []Node {
	return []Node{c.CatchBlock}
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

func (f ForeachStmt) String() string {
	return "foreach"
}

func (f ForeachStmt) Children() []Node {
	n := []Node{f.Source}
	if f.Key != nil {
		n = append(n, f.Key)
	}
	n = append(n, f.Value, f.LoopBlock)
	return n
}

type ArrayExpression struct {
	BaseNode
	ArrayType
	Pairs []ArrayPair
}

func (a ArrayExpression) String() string {
	return "array"
}

func (a ArrayExpression) Children() []Node {
	n := make([]Node, len(a.Pairs))
	for i, p := range a.Pairs {
		n[i] = p
	}
	return n
}

type ArrayPair struct {
	BaseNode
	Key   Expression
	Value Expression
}

func (p ArrayPair) Children() []Node {
	if p.Key != nil {
		return []Node{p.Key, p.Value}
	}
	return []Node{p.Value}
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

func (a ArrayLookupExpression) String() string {
	return fmt.Sprintf("%s[", a.Array)
}

func (a ArrayLookupExpression) Children() []Node {
	return []Node{a.Index}
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

func (a ArrayAppendExpression) String() string {
	return a.Array.String() + "[]"
}
