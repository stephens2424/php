package ast

// superGlobalScope represents the scope containing superglobals such as $_GET
type SuperGlobalScope struct {
	Identifiers []Variable
}

// globalScope represents the global scope on which functions and classes are
// defined. This is always within a namespace, but in many cases that may just
// be the default global namespace ("\")
type GlobalScope struct {
	*Namespace
	*Scope
}

// scope represents a particular local scope (such as within a function).
type Scope struct {
	Identifiers      map[string][]*Variable
	DynamicVariables []*Variable
	EnclosingScope   *Scope
	GlobalScope      *GlobalScope
	SuperGlobalScope *SuperGlobalScope
}

func (s *Scope) Variable(v *Variable) {
	switch i := v.Name.(type) {
	case Identifier:
		s.Identifiers[i.Value] = append(s.Identifiers[i.Value], v)
	default:
		s.DynamicVariables = append(s.DynamicVariables, v)
	}
}

type File struct {
	Name      string
	Namespace Namespace
	Nodes     []Node
}

type FileSet struct {
	Files           map[string]*File
	Namespaces      map[string]*Namespace
	GlobalNamespace *Namespace
	*Scope
}

func NewFileSet() *FileSet {
	return &FileSet{
		Files:           make(map[string]*File),
		Namespaces:      make(map[string]*Namespace),
		GlobalNamespace: NewNamespace("/"),
		Scope:           NewScope(nil, &GlobalScope{}, &SuperGlobalScope{}),
	}
}

func (f *FileSet) Namespace(name string) *Namespace {
	_, ok := f.Namespaces[name]
	if !ok {
		f.Namespaces[name] = NewNamespace(name)
	}
	return f.Namespaces[name]
}

type Namespace struct {
	Name                 string
	ClassesAndInterfaces map[string]Statement
	Constants            map[string][]*Variable
	Functions            map[string]*FunctionStmt
}

func NewNamespace(name string) *Namespace {
	return &Namespace{
		Name:                 name,
		ClassesAndInterfaces: map[string]Statement{},
		Constants:            map[string][]*Variable{},
		Functions:            map[string]*FunctionStmt{},
	}
}

type Classer interface {
	Node
	ClassName() string
}

func (c Class) ClassName() string     { return c.Name }
func (i Interface) ClassName() string { return i.Name }

func NewScope(parent *Scope, global *GlobalScope, superGlobal *SuperGlobalScope) *Scope {
	return &Scope{
		Identifiers:      map[string][]*Variable{},
		EnclosingScope:   parent,
		GlobalScope:      global,
		SuperGlobalScope: superGlobal,
	}
}
