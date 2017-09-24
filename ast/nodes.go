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

// Statement describes an executable piece of code. It may be as simple as
// a function call or a variable assignment. It also includes things like
// "if".
type Statement interface {
	Node

	// Declares returns the DeclarationType of the statement
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

// Format is a format
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

// Variable is a variable
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

// GlobalDeclaration is a global declaration
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

// Dynamic is a dynamic expression
type Dynamic Expr

// Static checks wether input is a identifier and returns it
func Static(d Dynamic) *Identifier {
	switch d := d.(type) {
	case Identifier:
		return &d
	case *Identifier:
		return d
	}

	return nil
}

// Expr describes a snippet of code that evaluates to a single
// value when run and does not represent a program instruction.
type Expr interface {
	Statement
	EvaluatesTo() Type
}

// BinaryExpr is an expression that applies an operator to one, two, or three
// operands. The operator determines how many operands it should contain.
type BinaryExpr struct {
	Antecedent Expr
	Subsequent Expr
	Type       Type
	Operator   string
}

func (b BinaryExpr) Children() []Node {
	return []Node{b.Antecedent, b.Subsequent}
}

func (b BinaryExpr) String() string {
	return b.Operator
}

func (b BinaryExpr) EvaluatesTo() Type {
	return b.Type
}

func (b BinaryExpr) Declares() DeclarationType { return NoDeclaration }

// TernaryCallExpr is a ternary call expression
type TernaryCallExpr struct {
	Condition, True, False Expr
	Type                   Type
}

func (t TernaryCallExpr) Children() []Node {
	return []Node{t.Condition, t.True, t.False}
}

func (t TernaryCallExpr) String() string {
	return "?:"
}

func (t TernaryCallExpr) EvaluatesTo() Type {
	return t.Type
}

func (t TernaryCallExpr) Declares() DeclarationType { return NoDeclaration }

// UnaryCallExpr is an expression that applies an operator to only one operand. The
// operator may precede or follow the operand.
type UnaryCallExpr struct {
	Operand   Expr
	Operator  string
	Preceding bool
}

func (u UnaryCallExpr) Children() []Node {
	return []Node{u.Operand}
}

func (u UnaryCallExpr) String() string {
	if u.Preceding {
		return u.Operator + " (preceding)"
	}
	return u.Operator
}

func (u UnaryCallExpr) EvaluatesTo() Type {
	return Unknown
}

func (u UnaryCallExpr) Declares() DeclarationType { return NoDeclaration }

// ExprStmt is a expression statement
type ExprStmt struct {
	Expr
}

func (e ExprStmt) String() string {
	if e.Expr != nil {
		return e.Expr.String()
	}
	return ""
}

func (e ExprStmt) Children() []Node {
	if e.Expr != nil {
		return []Node{e.Expr}
	}
	return nil
}

func (e ExprStmt) Declares() DeclarationType { return NoDeclaration }

// Echo returns a new echo statement.
func Echo(exprs ...Expr) EchoStmt {
	return EchoStmt{Expressions: exprs}
}

// EchoStmt represents an echo statement. It may be either a literal statement
// or it may be from data outside PHP-mode, such as "here" in: <? not here ?> here <? not here ?>
type EchoStmt struct {
	Expressions []Expr
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
	Expr
}

func (r ReturnStmt) String() string {
	return fmt.Sprintf("return")
}

func (r ReturnStmt) Children() []Node {
	if r.Expr == nil {
		return nil
	}
	return []Node{r.Expr}
}

func (r ReturnStmt) Declares() DeclarationType { return NoDeclaration }

// BreakStmt is a break statement
type BreakStmt struct {
	Expr
}

func (b BreakStmt) Children() []Node {
	if b.Expr != nil {
		return b.Expr.Children()
	}
	return nil
}

func (b BreakStmt) String() string {
	return "break"
}

// ContinueStmt is a continue statement
type ContinueStmt struct {
	Expr
}

func (c ContinueStmt) String() string {
	return "continue"
}

func (c ContinueStmt) Children() []Node {
	if c.Expr != nil {
		return c.Expr.Children()
	}
	return nil
}

// ThrowStmt is a throw statment
type ThrowStmt struct {
	Expr
}

func (t ThrowStmt) Declares() DeclarationType { return NoDeclaration }

// IncludeStmt is a include statment
type IncludeStmt struct {
	Include
}

// Include is a include statement
type Include struct {
	Expressions []Expr
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

// ExitStmt is an exit statment
type ExitStmt struct {
	Expr Expr
}

func (e ExitStmt) Children() []Node {
	return nil
}

func (e ExitStmt) String() string {
	return "exit"
}

func (e ExitStmt) Declares() DeclarationType { return NoDeclaration }

// NewCallExpr is a `new call` expression
type NewCallExpr struct {
	Class     Dynamic
	Arguments []Expr
}

func (n NewCallExpr) EvaluatesTo() Type {
	if static := Static(n.Class); static != nil {
		return ObjectType{static.Value}
	}
	return Object
}

func (n NewCallExpr) String() string {
	return "new"
}

func (n NewCallExpr) Children() []Node {
	nodes := make([]Node, len(n.Arguments)+1)
	nodes[0] = n.Class
	for i, arg := range n.Arguments {
		nodes[i+1] = arg
	}
	return nodes
}

func (n NewCallExpr) Declares() DeclarationType { return NoDeclaration }

// AssignmentExpr is an assighment expression
type AssignmentExpr struct {
	Assignee Assignable
	Value    Expr
	Operator string
}

func (a AssignmentExpr) String() string {
	return a.Operator
}

func (a AssignmentExpr) Children() []Node {
	return []Node{
		a.Assignee,
		a.Value,
	}
}

func (a AssignmentExpr) EvaluatesTo() Type {
	return a.Value.EvaluatesTo()
}

func (AssignmentExpr) Declares() DeclarationType { return NoDeclaration }

// Assignable is an assignable
type Assignable interface {
	Dynamic
	AssignableType() Type
}

// FunctionCallStmt is a function call statement
type FunctionCallStmt struct {
	FunctionCallExpr
}

// FunctionCallExpr is a function call expression
type FunctionCallExpr struct {
	FunctionName Dynamic
	Arguments    []Expr
}

func (f FunctionCallExpr) EvaluatesTo() Type {
	return Unknown
}

func (f FunctionCallExpr) String() string {
	return fmt.Sprintf("%s()", f.FunctionName)
}

func (f FunctionCallExpr) Children() []Node {
	n := make([]Node, len(f.Arguments))
	for i, a := range f.Arguments {
		n[i] = a
	}
	return n
}

func (f FunctionCallExpr) Declares() DeclarationType { return NoDeclaration }

// Block is a block
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

func (Block) Declares() DeclarationType { return NoDeclaration }

// FunctionStmt is a function statment
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

// AnonymousFunction is an anonymous function
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

// FunctionDefinition is a function defintion
type FunctionDefinition struct {
	Name      string
	Arguments []*FunctionArgument
	Type      string
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

// FunctionArgument is a function argument
type FunctionArgument struct {
	TypeHint string
	Default  Expr
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

// Class is a class
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

// Constant is a constant
type Constant struct {
	Name  string
	Value interface{}
}

func (c Constant) Children() []Node { return nil }
func (c Constant) String() string   { return c.Name }

// ConstantExpr is a constant expression
type ConstantExpr struct {
	*Variable
}

func (c Constant) Declares() DeclarationType { return ConstantDeclaration }

func (c Constant) EvaluatesTo() Type { return Unknown }

// Interface is an interface
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

// Property is a property
type Property struct {
	Name           string
	Visibility     Visibility
	Type           Type
	Initialization Expr
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

// PropertyCallExpr is a property call expression
type PropertyCallExpr struct {
	Receiver Dynamic
	Name     Dynamic
	Type     Type
}

func (p PropertyCallExpr) String() string {
	return fmt.Sprintf("%s->%s", p.Receiver, p.Name)
}

func (p PropertyCallExpr) AssignableType() Type {
	return p.Type
}

func (p PropertyCallExpr) EvaluatesTo() Type {
	return Unknown
}

func (p PropertyCallExpr) Children() []Node {
	return []Node{
		p.Receiver,
		p.Name,
	}
}

func (p PropertyCallExpr) Declares() DeclarationType { return NoDeclaration }

// ClassExpr is a class expression
type ClassExpr struct {
	Receiver Dynamic
	Expr     Dynamic
	Type     Type
}

// NewClassExpression returns a new ClassExpression
func NewClassExpression(r string, e Expr) *ClassExpr {
	return &ClassExpr{
		Receiver: &Identifier{Value: r},
		Expr:     e,
	}
}

func (c ClassExpr) EvaluatesTo() Type {
	return Unknown
}

func (c ClassExpr) String() string {
	return fmt.Sprintf("%s::", c.Receiver)
}

func (c ClassExpr) Children() []Node {
	return []Node{c.Receiver, c.Expr}
}

func (c ClassExpr) AssignableType() Type {
	return c.Type
}

func (c ClassExpr) Declares() DeclarationType { return NoDeclaration }

// Method is a method
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

// MethodCallExpr is a method call expression
type MethodCallExpr struct {
	Receiver Dynamic
	*FunctionCallExpr
}

func (m MethodCallExpr) Children() []Node {
	return []Node{
		m.Receiver,
		m.FunctionCallExpr,
	}
}

func (m MethodCallExpr) String() string {
	return fmt.Sprintf("%s->", m.Receiver)
}

// Visibility is a visibility
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

// IfStmt is an if statment
type IfStmt struct {
	Branches  []IfBranch
	ElseBlock Statement
}

// IfBranch is an if branch
type IfBranch struct {
	Condition Expr
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

// SwitchStmt is a switch statment
type SwitchStmt struct {
	Expr        Expr
	Cases       []*SwitchCase
	DefaultCase *Block
}

func (s SwitchStmt) String() string {
	return "switch"
}

func (s SwitchStmt) Children() []Node {
	n := []Node{
		s.Expr,
	}
	for _, c := range s.Cases {
		n = append(n, c)
	}
	if s.DefaultCase != nil {
		n = append(n, s.DefaultCase)
	}
	return n
}

func (SwitchStmt) Declares() DeclarationType { return NoDeclaration }

// SwitchCase is a switch case
type SwitchCase struct {
	Expr  Expr
	Block Block
}

func (s SwitchCase) String() string {
	return "case"
}

func (s SwitchCase) Children() []Node {
	return []Node{
		s.Expr,
		s.Block,
	}
}

// ForStmt is a for statment
type ForStmt struct {
	Initialization []Expr
	Termination    []Expr
	Iteration      []Expr
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

func (ForStmt) Declares() DeclarationType { return NoDeclaration }

// WhileStmt is a while statment
type WhileStmt struct {
	Termination Expr
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

func (WhileStmt) Declares() DeclarationType { return NoDeclaration }

// DoWhileStmt is a do while statement
type DoWhileStmt struct {
	Termination Expr
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

func (DoWhileStmt) Declares() DeclarationType { return NoDeclaration }

// TryStmt is a try statment
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

func (TryStmt) Declares() DeclarationType { return NoDeclaration }

// CatchStmt is a catch statment
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

// Literal is a literal
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

func (Literal) Declares() DeclarationType { return NoDeclaration }

// ForeachStmt is a for each statment
type ForeachStmt struct {
	Source    Expr
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

func (ForeachStmt) Declares() DeclarationType { return NoDeclaration }

// ArrayExpr is an array expression
type ArrayExpr struct {
	ArrayType
	Pairs []ArrayPair
}

func (a ArrayExpr) String() string {
	return "array"
}

func (a ArrayExpr) Children() []Node {
	n := make([]Node, len(a.Pairs))
	for i, p := range a.Pairs {
		n[i] = p
	}
	return n
}

func (ArrayExpr) Declares() DeclarationType { return NoDeclaration }

// ArrayPair is an array pair
type ArrayPair struct {
	Key   Expr
	Value Expr
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

func (a ArrayExpr) EvaluatesTo() Type {
	return Array
}

func (a ArrayExpr) AssignableType() Type {
	return Unknown
}

// ArrayLookupExpr is an array lookup expression
type ArrayLookupExpr struct {
	Array Dynamic
	Index Expr
}

func (ArrayLookupExpr) Declares() DeclarationType { return NoDeclaration }

func (a ArrayLookupExpr) String() string {
	return fmt.Sprintf("%s[", a.Array)
}

func (a ArrayLookupExpr) Children() []Node {
	return []Node{a.Index}
}

func (a ArrayLookupExpr) EvaluatesTo() Type {
	return Unknown
}

func (a ArrayLookupExpr) AssignableType() Type {
	return Unknown
}

// ArrayAppendExpr is an array append expression
type ArrayAppendExpr struct {
	Array Dynamic
}

func (a ArrayAppendExpr) EvaluatesTo() Type {
	return Unknown
}

func (a ArrayAppendExpr) AssignableType() Type {
	return Unknown
}

func (a ArrayAppendExpr) Children() []Node {
	return nil
}

func (a ArrayAppendExpr) String() string {
	return a.Array.String() + "[]"
}

func (ArrayAppendExpr) Declares() DeclarationType { return NoDeclaration }

// ShellCommand is a shell command
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

func (ShellCommand) Declares() DeclarationType { return NoDeclaration }

// ListStatement is a list statement
type ListStatement struct {
	Assignees []Assignable
	Value     Expr
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

func (ListStatement) Declares() DeclarationType { return NoDeclaration }

// StaticVariableDeclaration is a static variable declaration
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

// DeclareBlock is a declare block
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

func (DeclareBlock) Declares() DeclarationType { return NoDeclaration }
