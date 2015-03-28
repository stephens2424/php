package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/stephens2424/php/token"
)

// Node encapsulates every AST node.
type Node interface {
	String() string
	Children() []Node
}

type Declaration interface {
	Node
	Declares() DeclarationType
}

type DeclarationType int

const (
	VariableDeclaration DeclarationType = iota
	FunctionDeclaration
	ClassDeclaration
)

type Format struct{}

// An Identifier is a raw string that can be used to identify
// a variable, function, class, constant, property, etc.
type Identifier struct {
	Parent Node
	Value  string
}

func (i Identifier) EvaluatesTo() Type {
	return String
}

func (i Identifier) String() string {
	return i.Value
}

func (i Identifier) Children() []Node {
	return nil
}

type Variable struct {

	// Name is the identifier for the variable, which may be
	// a dynamic expression.
	Name Expression
	Type Type
}

func (v Variable) String() string {
	return "$" + v.Name.String()
}

func (v Variable) Children() []Node {
	return []Node{v.Name}
}

type GlobalDeclaration struct {
	Identifiers []*Variable
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
func (i Variable) AssignableType() Type {
	return i.Type
}

// EvaluatesTo returns the known type of the variable.
func (i Variable) EvaluatesTo() Type {
	return i.Type
}

// NewVariable intializes a variable node with its name being a simple
// identifier and its type set to AnyType. The name argument should not
// include the $ operator.
func NewVariable(name string) *Variable {
	return &Variable{Name: &Identifier{Value: name}, Type: AnyType}
}

// A statement is an executable piece of code. It may be as simple as
// a function call or a variable assignment. It also includes things like
// "if".
type Statement interface {
	Node
}

// EmptyStatement represents a statement that does nothing.
type EmptyStatement struct {
}

func (e EmptyStatement) String() string        { return "" }
func (e EmptyStatement) Children() []Node      { return nil }
func (e EmptyStatement) Print(f Format) string { return ";" }

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
type BinaryExpression struct {
	Antecedent Expression
	Subsequent Expression
	Type       Type
	Operator   string
}

func (b BinaryExpression) Children() []Node {
	return []Node{b.Antecedent, b.Subsequent}
}

func (b BinaryExpression) String() string {
	return b.Operator
}

func (b BinaryExpression) EvaluatesTo() Type {
	return b.Type
}

type TernaryExpression struct {
	Condition, True, False Expression
	Type                   Type
}

func (t TernaryExpression) Children() []Node {
	return []Node{t.Condition, t.True, t.False}
}

func (t TernaryExpression) String() string {
	return "?:"
}

func (t TernaryExpression) EvaluatesTo() Type {
	return t.Type
}

// UnaryExpression is an expression that applies an operator to only one operand. The
// operator may precede or follow the operand.
type UnaryExpression struct {
	Operand   Expression
	Operator  string
	Preceding bool
}

func (u UnaryExpression) Children() []Node {
	return []Node{u.Operand}
}

func (u UnaryExpression) String() string {
	if u.Preceding {
		return u.Operator + " (preceding)"
	}
	return u.Operator
}

func (u UnaryExpression) EvaluatesTo() Type {
	return Unknown
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
	Expressions []Expression
}

func (i Include) Children() []Node {
	n := make([]Node, len(i.Expressions))
	for idx, expr := range i.Expressions {
		n[idx] = expr
	}
	return n
}

func (i Include) String() string {
	return "include"
}

func (i Include) EvaluatesTo() Type {
	return AnyType
}

type ExitStmt struct {
	Expression Expression
}

func (e ExitStmt) Children() []Node {
	return nil
}

func (e ExitStmt) String() string {
	return "exit"
}

type NewExpression struct {
	Class     Expression
	Arguments []Expression
}

func (n NewExpression) EvaluatesTo() Type {
	return Object
}

func (n NewExpression) String() string {
	return "new"
}

func (c NewExpression) Children() []Node {
	n := make([]Node, len(c.Arguments)+1)
	n[0] = c.Class
	for i, arg := range c.Arguments {
		n[i+1] = arg
	}
	return n
}

type AssignmentExpression struct {
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

type Assignable interface {
	Node
	AssignableType() Type
}

type FunctionCallStmt struct {
	FunctionCallExpression
}
type FunctionCallExpression struct {
	FunctionName Expression
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

type AnonymousFunction struct {
	ClosureVariables []*FunctionArgument
	Arguments        []*FunctionArgument
	Body             *Block
}

func (a AnonymousFunction) EvaluatesTo() Type {
	return Function
}
func (a AnonymousFunction) Children() []Node {
	n := []Node{}
	for _, c := range a.ClosureVariables {
		n = append(n, c)
	}
	for _, a := range a.Arguments {
		n = append(n, a)
	}
	n = append(n, a.Body)
	return n
}

func (a AnonymousFunction) String() string {
	return "anonymous function"
}

type FunctionDefinition struct {
	Name      string
	Arguments []*FunctionArgument
}

func (fd FunctionDefinition) Children() []Node {
	n := make([]Node, len(fd.Arguments))
	for i, arg := range fd.Arguments {
		n[i] = arg
	}
	return n
}

func (fd FunctionDefinition) String() string {
	return fmt.Sprintf("function %s( %s )", fd.Name, fd.Arguments)
}

type FunctionArgument struct {
	TypeHint string
	Default  Expression
	Variable *Variable
}

func (fa FunctionArgument) String() string {
	return fmt.Sprintf("Arg: %s", fa.TypeHint)
}

func (fa FunctionArgument) Children() []Node {
	n := []Node{
		fa.Variable,
	}
	if fa.Default != nil {
		n = append(n, fa.Default)
	}
	return n
}

type Class struct {
	Name       string
	Extends    string
	Implements []string
	Methods    []*Method
	Properties []*Property
	Constants  []*Constant
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
	*Variable
	Value interface{}
}

type ConstantExpression struct {
	*Variable
}

type Interface struct {
	Name      string
	Inherits  []string
	Methods   []Method
	Constants []Constant
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

func (p Property) Children() []Node {
	return []Node{p.Initialization}
}

type PropertyExpression struct {
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

type ClassExpression struct {
	Receiver   Expression
	Expression Expression
	Type       Type
}

func NewClassExpression(r string, e Expression) *ClassExpression {
	return &ClassExpression{
		Receiver:   &Identifier{Value: r},
		Expression: e,
	}
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

func (c ClassExpression) AssignableType() Type {
	return c.Type
}

type Method struct {
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

func (v Visibility) Token() token.Token {
	switch v {
	case Private:
		return token.Private
	case Protected:
		return token.Protected
	case Public:
		return token.Public
	}
	panic("invalid visibility value")
}

type IfStmt struct {
	Branches  []IfBranch
	ElseBlock Statement
}

type IfBranch struct {
	Condition Expression
	Block     Statement
}

func (i IfBranch) String() string {
	return i.Condition.String()
}

func (i IfBranch) Children() []Node {
	return []Node{i.Block}
}

func (i IfStmt) String() string {
	return "if"
}

func (i IfStmt) Children() []Node {
	n := make([]Node, 0, 3)
	for _, branch := range i.Branches {
		n = append(n, branch)
	}
	if i.ElseBlock != nil {
		n = append(n, i.ElseBlock)
	}
	return n
}

type SwitchStmt struct {
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
	Initialization []Expression
	Termination    []Expression
	Iteration      []Expression
	LoopBlock      Statement
}

func (f ForStmt) String() string {
	return "for"
}

func (f ForStmt) Children() []Node {
	nodes := []Node{}
	for _, stmt := range f.Initialization {
		nodes = append(nodes, stmt)
	}
	for _, stmt := range f.Termination {
		nodes = append(nodes, stmt)
	}
	for _, stmt := range f.Iteration {
		nodes = append(nodes, stmt)
	}
	return nodes
}

type WhileStmt struct {
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
	CatchBlock *Block
	CatchType  string
	CatchVar   *Variable
}

func (c CatchStmt) String() string {
	return fmt.Sprintf("catch %s %s", c.CatchType, c.CatchVar)
}

func (c CatchStmt) Children() []Node {
	return []Node{c.CatchBlock}
}

type Literal struct {
	Type  Type
	Value string
}

func (l Literal) String() string {
	return fmt.Sprintf("Literal-%s: %s", l.Type, l.Value)
}

func (l Literal) EvaluatesTo() Type {
	return l.Type
}

func (l Literal) Children() []Node {
	return nil
}

type ForeachStmt struct {
	Source    Expression
	Key       *Variable
	Value     *Variable
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
	Key   Expression
	Value Expression
}

func (p ArrayPair) Children() []Node {
	if p.Key != nil {
		return []Node{p.Key, p.Value}
	}
	return []Node{p.Value}
}
func (p ArrayPair) String() string {
	return fmt.Sprintf("%s => %s", p.Key, p.Value)
}

func (a ArrayExpression) EvaluatesTo() Type {
	return Array
}

func (a ArrayExpression) AssignableType() Type {
	return AnyType
}

type ArrayLookupExpression struct {
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
	Array Expression
}

func (a ArrayAppendExpression) EvaluatesTo() Type {
	return AnyType
}

func (a ArrayAppendExpression) AssignableType() Type {
	return AnyType
}

func (a ArrayAppendExpression) Children() []Node {
	return nil
}

func (a ArrayAppendExpression) String() string {
	return a.Array.String() + "[]"
}

type ShellCommand struct {
	Command string
}

func (s ShellCommand) String() string {
	return fmt.Sprintf("`%s`", s.Command)
}

func (s ShellCommand) EvaluatesTo() Type {
	return String
}
func (s ShellCommand) Children() []Node {
	return nil
}

type ListStatement struct {
	Assignees []Assignable
	Value     Expression
	Operator  string
}

func (l ListStatement) EvaluatesTo() Type {
	return Array
}

func (l ListStatement) String() string {
	return fmt.Sprintf("list(%s)", l.Assignees)
}

func (l ListStatement) Children() []Node {
	return []Node{l.Value}
}

type StaticVariableDeclaration struct {
	Declarations []Expression
}

func (s StaticVariableDeclaration) Children() []Node {
	return nil
}
func (s StaticVariableDeclaration) String() string {
	buf := bytes.NewBufferString("static ")
	for i, d := range s.Declarations {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(d.String())
	}
	return buf.String()
}

type DeclareBlock struct {
	Statements   *Block
	Declarations []string
}

func (d DeclareBlock) Children() []Node {
	return d.Statements.Children()
}

func (d DeclareBlock) String() string {
	return "declare{}"
}

type File struct {
	Name      string
	Namespace Namespace
	Nodes     []Node
}

type FileSet struct {
	Files      map[string]File
	Namespaces map[string]Namespace
}

type Namespace struct {
	Name  string
	Decls []Declaration
}
