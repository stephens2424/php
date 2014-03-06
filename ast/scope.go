package ast

// superGlobalScope represents the scope containing superglobals such as $_GET
type SuperGlobalScope struct {
	Identifiers []Variable
}

// globalScope represents the global scope on which functions and classes are
// defined. This is always within a namespace, but in many cases that may just
// be the default global namespace ("\")
type GlobalScope struct {
	Functions  []FunctionStmt
	Classes    []Class
	Interfaces []Interface
	Cconstants []Constant
	Namespace  string
	*Scope
}

// scope represents a particular local scope (such as within a function).
type Scope struct {
	Identifiers      []Variable
	EnclosingScope   *Scope
	GlobalScope      *GlobalScope
	SuperGlobalScope *SuperGlobalScope
}
