package ast

import (
	"fmt"
	"strings"

	"github.com/stephens2424/php/token"
)

// Node encapsulates every AST node.
type Node interface {
	String() string
	Children() []Node
	Tokens() token.Stream
}

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

func (i Identifier) Tokens() token.Stream {
	return token.NewList(token.NewItem(token.Identifier, i.Value))
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

func (v Variable) Tokens() token.Stream {
	l := token.NewList()
	l.Push(token.NewItem(token.VariableOperator, "$"))
	l.PushStream(v.Name.Tokens())
	return l
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

func (g GlobalDeclaration) Tokens() token.Stream {
	l := token.NewList(token.Keyword(token.Global))
	for i, id := range g.Identifiers {
		l.PushStream(id.Tokens())
		if i+1 < len(g.Identifiers) {
			l.PushKeyword(token.Comma)
		}
	}
	l.PushKeyword(token.StatementEnd)
	return l
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
	return &Variable{Name: Identifier{Value: name}, Type: AnyType}
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

func (e EmptyStatement) String() string       { return "" }
func (e EmptyStatement) Children() []Node     { return nil }
func (e EmptyStatement) Tokens() token.Stream { return token.NewList() }

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

func (b BinaryExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(b.Antecedent.Tokens())
	l.Push(token.Item{Typ: token.TokenMap[b.Operator], Val: b.Operator})
	l.PushStream(b.Subsequent.Tokens())
	return l
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

func (t TernaryExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(t.Condition.Tokens())
	l.PushKeyword(token.TernaryOperator1)
	l.PushStream(t.True.Tokens())
	l.PushKeyword(token.TernaryOperator2)
	l.PushStream(t.False.Tokens())
	return l
}

// UnaryExpression is an expression that applies an operator to only one operand. The
// operator may precede or follow the operand.
type UnaryExpression struct {
	Operand   Expression
	Operator  string
	Preceding bool
}

func (u UnaryExpression) Children() []Node {
	return nil
}

func (u UnaryExpression) String() string {
	if u.Preceding {
		return u.Operator + u.Operand.String()
	}
	return u.Operand.String() + u.Operator
}

func (u UnaryExpression) EvaluatesTo() Type {
	return Unknown
}

func (u UnaryExpression) Tokens() token.Stream {
	op := token.Item{Val: u.Operator, Typ: token.UnaryOperator}
	l := token.NewList()
	if u.Preceding {
		l.Push(op)
	}
	l.PushStream(u.Operand.Tokens())
	if !u.Preceding {
		l.Push(op)
	}
	return l
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

func (e EchoStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Echo)
	for i, expr := range e.Expressions {
		if i > 0 {
			l.PushKeyword(token.Comma)
			l.PushStream(expr.Tokens())
		}
	}
	l.PushKeyword(token.StatementEnd)
	return l
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

func (r ReturnStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Return)
	if r.Expression != nil {
		l.PushStream(r.Expression.Tokens())
	}
	l.PushKeyword(token.StatementEnd)
	return l
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

func (b BreakStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Break)
	if b.Expression != nil {
		l.PushStream(b.Expression.Tokens())
	}
	l.PushKeyword(token.StatementEnd)
	return l
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

func (b ContinueStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Continue)
	if b.Expression != nil {
		l.PushStream(b.Expression.Tokens())
	}
	l.PushKeyword(token.StatementEnd)
	return l
}

type ThrowStmt struct {
	Expression
}

func (b ThrowStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Throw)
	if b.Expression != nil {
		l.PushStream(b.Expression.Tokens())
	}
	l.PushKeyword(token.StatementEnd)
	return l
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

func (e Include) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Include)
	for i, expr := range e.Expressions {
		if i > 0 {
			l.PushKeyword(token.Comma)
			l.PushStream(expr.Tokens())
		}
	}
	l.PushKeyword(token.StatementEnd)
	return l
}

type ExitStmt struct {
	Expression Expression
}

func (b ExitStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Exit)
	if b.Expression != nil {
		l.PushStream(b.Expression.Tokens())
	}
	l.PushKeyword(token.StatementEnd)
	return l
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

func (b NewExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.NewOperator)
	l.PushStream(b.Class.Tokens())
	l.PushKeyword(token.OpenParen)
	for i, arg := range b.Arguments {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(arg.Tokens())
	}
	l.PushKeyword(token.CloseParen)
	return l
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

func (a AssignmentExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(a.Assignee.Tokens())
	l.Push(token.Item{Typ: token.AssignmentOperator, Val: a.Operator})
	l.PushStream(a.Value.Tokens())
	return l
}

type Assignable interface {
	Node
	AssignableType() Type
}

type FunctionCallStmt struct {
	FunctionCallExpression
}

func (f FunctionCallStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(f.FunctionCallExpression.Tokens())
	l.PushKeyword(token.StatementEnd)
	return l
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

func (f FunctionCallExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(f.FunctionName.Tokens())
	l.PushKeyword(token.OpenParen)
	for i, arg := range f.Arguments {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(arg.Tokens())
	}
	l.PushKeyword(token.CloseParen)
	return l
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

func (b Block) Tokens() token.Stream {
	l := token.NewList()
	for _, s := range b.Statements {
		l.PushStream(s.Tokens())
	}
	return l
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

func (f FunctionStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(f.FunctionDefinition.Tokens())
	l.PushStream(f.Body.Tokens())
	return l
}

type AnonymousFunction struct {
	ClosureVariables []FunctionArgument
	Arguments        []FunctionArgument
	Body             *Block
}

func (a AnonymousFunction) EvaluatesTo() Type {
	return Function
}

func (a AnonymousFunction) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Function)
	l.PushKeyword(token.OpenParen)
	for i, arg := range a.Arguments {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(arg.Tokens())
	}
	l.PushKeyword(token.CloseParen)
	if len(a.ClosureVariables) > 0 {
		l.PushKeyword(token.OpenParen)
		for i, arg := range a.ClosureVariables {
			if i > 0 {
				l.PushKeyword(token.Comma)
			}
			l.PushStream(arg.Tokens())
		}
		l.PushKeyword(token.CloseParen)
	}
	l.PushStream(a.Body.Tokens())
	return l
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
	Arguments []FunctionArgument
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

func (fd FunctionDefinition) Tokens() token.Stream {
	l := token.NewList(token.Keyword(token.Function))
	l.Push(token.Item{Typ: token.Identifier, Val: fd.Name})
	l.PushKeyword(token.OpenParen)
	for i, arg := range fd.Arguments {
		l.PushStream(arg.Tokens())
		if i+1 < len(fd.Arguments) {
			l.PushKeyword(token.Comma)
		}
	}
	l.Push(token.Keyword(token.CloseParen))
	return l
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

func (fa FunctionArgument) Tokens() token.Stream {
	l := token.NewList()
	if fa.TypeHint != "" {
		l.Push(token.Item{Typ: token.Identifier, Val: fa.TypeHint})
	}
	l.PushStream(fa.Variable.Tokens())
	if fa.Default != nil {
		l.PushKeyword(token.AssignmentOperator)
		l.PushStream(fa.Default.Tokens())
	}
	return l
}

type Class struct {
	Name       string
	Extends    string
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

func (c Class) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Class)
	l.Push(token.Item{Typ: token.Identifier, Val: c.Name})
	if c.Extends != "" {
		l.PushKeyword(token.Extends)
		l.Push(token.Item{Typ: token.Identifier, Val: c.Extends})
	}
	for i, imp := range c.Implements {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.Push(token.Item{Typ: token.Identifier, Val: imp})
	}
	l.PushKeyword(token.BlockBegin)
	for _, c := range c.Constants {
		l.PushStream(c.Tokens())
	}
	for _, p := range c.Properties {
		l.PushStream(p.Tokens())
	}
	for _, m := range c.Methods {
		l.PushStream(m.Tokens())
	}
	l.PushKeyword(token.BlockEnd)
	return l
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

func (i Interface) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Interface)
	l.Push(token.Item{Typ: token.Identifier, Val: i.Name})

	for i, imp := range i.Inherits {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.Push(token.Item{Typ: token.Identifier, Val: imp})
	}

	l.PushKeyword(token.BlockBegin)
	for _, c := range i.Constants {
		l.PushStream(c.Tokens())
	}

	for _, m := range i.Methods {
		l.PushStream(m.Tokens())
	}

	l.PushKeyword(token.BlockEnd)
	return l
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

func (p Property) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(p.Visibility.Token())
	l.Push(token.Item{Typ: token.Identifier, Val: p.Name})
	if p.Initialization != nil {
		l.PushKeyword(token.AssignmentOperator)
		l.PushStream(p.Initialization.Tokens())
	}
	return l
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

func (p PropertyExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(p.Receiver.Tokens())
	l.PushKeyword(token.ObjectOperator)
	l.PushStream(p.Name.Tokens())
	return l
}

type ClassExpression struct {
	Receiver   Expression
	Expression Expression
	Type       Type
}

func NewClassExpression(r string, e Expression) *ClassExpression {
	return &ClassExpression{
		Receiver:   Identifier{Value: r},
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

func (c ClassExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(c.Receiver.Tokens())
	l.PushKeyword(token.ScopeResolutionOperator)
	l.PushStream(c.Expression.Tokens())
	return l
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

func (m Method) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(m.Visibility.Token())
	l.PushStream(m.FunctionStmt.Tokens())
	return l
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

func (m MethodCallExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(m.Receiver.Tokens())
	l.PushKeyword(token.ObjectOperator)
	l.PushStream(m.FunctionCallExpression.Tokens())
	return l
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

func (i IfStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.If)
	l.PushKeyword(token.OpenParen)
	l.PushStream(i.Condition.Tokens())
	l.PushKeyword(token.CloseParen)
	l.PushStream(i.TrueBranch.Tokens())
	if i.FalseBranch != nil {
		l.PushKeyword(token.Else)
		l.PushStream(i.FalseBranch.Tokens())
	}
	return l
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

func (s SwitchStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Switch)
	l.PushKeyword(token.BlockBegin)
	for _, c := range s.Cases {
		l.PushStream(c.Tokens())
	}
	if s.DefaultCase != nil {
		l.PushKeyword(token.Default)
		l.PushKeyword(token.TernaryOperator2)
		l.PushStream(s.DefaultCase.Tokens())
	}
	l.PushKeyword(token.BlockEnd)
	return l
}

type SwitchCase struct {
	Expression Expression
	Block      Block
}

func (s SwitchCase) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Case)
	l.PushStream(s.Expression.Tokens())
	l.PushKeyword(token.TernaryOperator2)
	l.PushStream(s.Block.Tokens())
	return l
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

func (f ForStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.For)
	l.PushKeyword(token.OpenParen)
	for i, e := range f.Initialization {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(e.Tokens())
	}
	for i, e := range f.Termination {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(e.Tokens())
	}
	for i, e := range f.Iteration {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(e.Tokens())
	}
	l.PushKeyword(token.CloseParen)
	l.PushStream(f.LoopBlock.Tokens())
	return l
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

func (w WhileStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.While)
	l.PushKeyword(token.OpenParen)
	l.PushStream(w.Termination.Tokens())
	l.PushKeyword(token.CloseParen)
	l.PushStream(w.LoopBlock.Tokens())
	return l
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

func (w DoWhileStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Do)
	l.PushStream(w.LoopBlock.Tokens())
	l.PushKeyword(token.While)
	l.PushKeyword(token.OpenParen)
	l.PushStream(w.Termination.Tokens())
	l.PushKeyword(token.CloseParen)
	l.PushKeyword(token.StatementEnd)
	return l
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

func (t TryStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Try)
	l.PushStream(t.TryBlock.Tokens())
	for _, c := range t.CatchStmts {
		l.PushStream(c.Tokens())
	}
	if t.FinallyBlock != nil {
		l.PushKeyword(token.Finally)
		l.PushStream(t.FinallyBlock.Tokens())
	}
	return l
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

func (c CatchStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Catch)
	l.PushKeyword(token.OpenParen)
	l.Push(token.Item{Typ: token.Identifier, Val: c.CatchType})
	l.PushStream(c.CatchVar.Tokens())
	l.PushKeyword(token.CloseParen)
	l.PushStream(c.CatchBlock.Tokens())
	return l
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

func (l Literal) Tokens() token.Stream {
	switch l.Type {
	case String:
		return token.NewList(token.Item{Typ: token.StringLiteral, Val: l.Value})
	case Integer, Float:
		return token.NewList(token.Item{Typ: token.NumberLiteral, Val: l.Value})
	case Boolean:
		return token.NewList(token.Item{Typ: token.BooleanLiteral, Val: l.Value})
	case Null:
		return token.NewList(token.Item{Typ: token.Null, Val: "null"})
	}
	panic("invalid literal type")
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

func (f ForeachStmt) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Foreach)
	l.PushKeyword(token.OpenParen)
	l.PushStream(f.Source.Tokens())
	l.PushKeyword(token.AsOperator)
	if f.Key != nil {
		l.PushStream(f.Key.Tokens())
		l.PushKeyword(token.ArrayKeyOperator)
	}
	l.PushStream(f.Value.Tokens())
	l.PushKeyword(token.CloseParen)
	l.PushStream(f.LoopBlock.Tokens())
	return l
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

func (a ArrayExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushKeyword(token.Array)
	l.PushKeyword(token.OpenParen)
	for i, pair := range a.Pairs {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(pair.Tokens())
	}
	return l
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

func (p ArrayPair) Tokens() token.Stream {
	l := token.NewList()
	if p.Key != nil {
		l.PushStream(p.Key.Tokens())
		l.PushKeyword(token.ArrayKeyOperator)
	}
	l.PushStream(p.Value.Tokens())
	return l
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

func (a ArrayLookupExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(a.Array.Tokens())
	l.PushKeyword(token.ArrayLookupOperatorLeft)
	if a.Index != nil {
		l.PushStream(a.Index.Tokens())
	}
	l.PushKeyword(token.ArrayLookupOperatorRight)
	return l
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

func (a ArrayAppendExpression) Tokens() token.Stream {
	l := token.NewList()
	l.PushStream(a.Array.Tokens())
	l.PushKeyword(token.ArrayLookupOperatorLeft)
	l.PushKeyword(token.ArrayLookupOperatorRight)
	return l
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

func (s ShellCommand) Tokens() token.Stream {
	return token.NewList(token.Item{Typ: token.ShellCommand, Val: s.Command})
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

func (l ListStatement) Tokens() token.Stream {
	s := token.NewList()
	s.PushKeyword(token.List)
	for i, a := range l.Assignees {
		if i > 0 {
			s.PushKeyword(token.Comma)
		}
		s.PushStream(a.Tokens())
	}
	s.Push(token.Item{Typ: token.AssignmentOperator, Val: l.Operator})
	s.PushStream(l.Value.Tokens())
	return s
}

type StaticVariableDeclaration struct {
	Declarations []Expression
}

type DeclareBlock struct {
	Statements   *Block
	Declarations []string
}
