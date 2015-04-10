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

// A statement is an executable piece of code. It may be as simple as
// a function call or a variable assignment. It also includes things like
// "if".
type Statement interface {
	Node
	Declares() DeclarationType
}

// DeclarationType identifies the type of declarative statements.
type DeclarationType int

const (
	// NoDeclaration identifies a statement which declares nothing in the local namespace.
	NoDeclaration DeclarationType = iota

	// ConstantDeclaration identifies a statement which declares a constant in the local namespace.
	ConstantDeclaration

	// FunctionDeclaration identfies a statement which declares a function in the local namespace.
	FunctionDeclaration

	// ClassDeclaration identfies a statement which declares a class in the local namespace.
	ClassDeclaration

	// InterfaceDeclaration identfies a statement which declares a interface in the local namespace.
	InterfaceDeclaration
)

type Format struct{}

// An Identifier is a raw string that can be used to identify
// a variable, function, class, constant, property, etc.
type Identifier struct {
	Parent Node
	Value  string
}

func (i Identifier) Declares() DeclarationType { return NoDeclaration }

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
	Name Dynamic
	Type Type
}

// NewVariable intializes a variable node with its name being a simple
// identifier and its type set to Unknown. The name argument should not
// include the $ operator.
func NewVariable(name string) *Variable {
	return &Variable{Name: &Identifier{Value: name}, Type: Unknown}
}

func (v Variable) String() string {
	return "$" + v.Name.String()
}

func (v Variable) Children() []Node {
	return []Node{v.Name}
}

func (v Variable) AssignableType() Type {
	return v.Type
}

// EvaluatesTo returns the known type of the variable.
func (v Variable) EvaluatesTo() Type {
	return v.Type
}

func (v Variable) Declares() DeclarationType { return NoDeclaration }

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

func (g GlobalDeclaration) Declares() DeclarationType { return NoDeclaration }

// EmptyStatement represents a statement that does nothing.
type EmptyStatement struct{}

func (e EmptyStatement) String() string            { return "" }
func (e EmptyStatement) Children() []Node          { return nil }
func (e EmptyStatement) Print(f Format) string     { return ";" }
func (e EmptyStatement) Declares() DeclarationType { return NoDeclaration }

type Dynamic Expression

func Static(d Dynamic) *Identifier {
	switch d := d.(type) {
	case Identifier:
		return &d
	case *Identifier:
		return d
	}

	return nil
}

// An Expression is a snippet of code that evaluates to a single value when run
// and does not represent a program instruction.
type Expression interface {
	Statement
	EvaluatesTo() Type
}

// BinaryExpression is an expression that applies an operator to one, two, or three
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

func (b BinaryExpression) Declares() DeclarationType { return NoDeclaration }

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

func (t TernaryExpression) Declares() DeclarationType { return NoDeclaration }

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

func (u UnaryExpression) Declares() DeclarationType { return NoDeclaration }

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
		return []Node{e.Expression}
	}
	return nil
}

func (e ExpressionStmt) Declares() DeclarationType { return NoDeclaration }

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

func (e EchoStmt) Declares() DeclarationType { return NoDeclaration }

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

func (r ReturnStmt) Declares() DeclarationType { return NoDeclaration }

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

func (t ThrowStmt) Declares() DeclarationType { return NoDeclaration }

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
	return Unknown
}

func (i Include) Declares() DeclarationType { return NoDeclaration }

type ExitStmt struct {
	Expression Expression
}

func (e ExitStmt) Children() []Node {
	return nil
}

func (e ExitStmt) String() string {
	return "exit"
}

func (e ExitStmt) Declares() DeclarationType { return NoDeclaration }

type NewExpression struct {
	Class     Dynamic
	Arguments []Expression
}

func (n NewExpression) EvaluatesTo() Type {
	if static := Static(n.Class); static != nil {
		return ObjectType{static.Value}
	}
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

func (n NewExpression) Declares() DeclarationType { return NoDeclaration }

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

func (a AssignmentExpression) Declares() DeclarationType { return NoDeclaration }

type Assignable interface {
	Dynamic
	AssignableType() Type
}

type FunctionCallStmt struct {
	FunctionCallExpression
}

type FunctionCallExpression struct {
	FunctionName Dynamic
	Arguments    []Expression
}

func (f FunctionCallExpression) EvaluatesTo() Type {
	return Unknown
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

func (f FunctionCallExpression) Declares() DeclarationType { return NoDeclaration }

type Block struct {
	Statements []Statement
	Scope      *Scope
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

func (_ Block) Declares() DeclarationType { return NoDeclaration }

type FunctionStmt struct {
	*FunctionDefinition
	Body *Block
}

func (f FunctionStmt) Declares() DeclarationType { return FunctionDeclaration }

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

func (a AnonymousFunction) Declares() DeclarationType { return FunctionDeclaration }

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

func (c Class) Declares() DeclarationType { return ClassDeclaration }

type Constant struct {
	Name  string
	Value interface{}
}

func (c Constant) Children() []Node { return nil }
func (c Constant) String() string   { return c.Name }

type ConstantExpression struct {
	*Variable
}

func (c Constant) Declares() DeclarationType { return ConstantDeclaration }

func (c Constant) EvaluatesTo() Type { return Unknown }

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

func (i Interface) Declares() DeclarationType { return InterfaceDeclaration }

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
	Receiver Dynamic
	Name     Dynamic
	Type     Type
}

func (p PropertyExpression) String() string {
	return fmt.Sprintf("%s->%s", p.Receiver, p.Name)
}

func (p PropertyExpression) AssignableType() Type {
	return p.Type
}

func (p PropertyExpression) EvaluatesTo() Type {
	return Unknown
}

func (p PropertyExpression) Children() []Node {
	return []Node{
		p.Receiver,
	}
}

func (p PropertyExpression) Declares() DeclarationType { return NoDeclaration }

type ClassExpression struct {
	Receiver   Dynamic
	Expression Dynamic
	Type       Type
}

func NewClassExpression(r string, e Expression) *ClassExpression {
	return &ClassExpression{
		Receiver:   &Identifier{Value: r},
		Expression: e,
	}
}

func (c ClassExpression) EvaluatesTo() Type {
	return Unknown
}

func (c ClassExpression) String() string {
	return fmt.Sprintf("%s::", c.Receiver)
}

func (c ClassExpression) Children() []Node {
	return []Node{c.Receiver, c.Expression}
}

func (c ClassExpression) AssignableType() Type {
	return c.Type
}

func (c ClassExpression) Declares() DeclarationType { return NoDeclaration }

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
	Receiver Dynamic
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

func (i IfStmt) Declares() DeclarationType { return NoDeclaration }

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

func (_ SwitchStmt) Declares() DeclarationType { return NoDeclaration }

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

func (_ ForStmt) Declares() DeclarationType { return NoDeclaration }

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

func (_ WhileStmt) Declares() DeclarationType { return NoDeclaration }

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

func (_ DoWhileStmt) Declares() DeclarationType { return NoDeclaration }

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

func (_ TryStmt) Declares() DeclarationType { return NoDeclaration }

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

func (_ Literal) Declares() DeclarationType { return NoDeclaration }

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

func (_ ForeachStmt) Declares() DeclarationType { return NoDeclaration }

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

func (_ ArrayExpression) Declares() DeclarationType { return NoDeclaration }

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
	return Unknown
}

type ArrayLookupExpression struct {
	Array Dynamic
	Index Expression
}

func (_ ArrayLookupExpression) Declares() DeclarationType { return NoDeclaration }

func (a ArrayLookupExpression) String() string {
	return fmt.Sprintf("%s[", a.Array)
}

func (a ArrayLookupExpression) Children() []Node {
	return []Node{a.Index}
}

func (a ArrayLookupExpression) EvaluatesTo() Type {
	return Unknown
}

func (a ArrayLookupExpression) AssignableType() Type {
	return Unknown
}

type ArrayAppendExpression struct {
	Array Dynamic
}

func (a ArrayAppendExpression) EvaluatesTo() Type {
	return Unknown
}

func (a ArrayAppendExpression) AssignableType() Type {
	return Unknown
}

func (a ArrayAppendExpression) Children() []Node {
	return nil
}

func (a ArrayAppendExpression) String() string {
	return a.Array.String() + "[]"
}

func (_ ArrayAppendExpression) Declares() DeclarationType { return NoDeclaration }

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

func (_ ShellCommand) Declares() DeclarationType { return NoDeclaration }

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

func (_ ListStatement) Declares() DeclarationType { return NoDeclaration }

type StaticVariableDeclaration struct {
	Declarations []Dynamic
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

func (s StaticVariableDeclaration) Declares() DeclarationType { return NoDeclaration }

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

func (p DeclareBlock) Declares() DeclarationType { return NoDeclaration }
