package ast

// superGlobalScope represents the scope containing superglobals such as $_GET
type SuperGlobalScope struct {
	Identifiers map[string]*Variable
}

func NewSuperGlobalScope() *SuperGlobalScope {
	return &SuperGlobalScope{map[string]*Variable{}}
}

// globalScope represents the global scope on which functions and classes are
// defined. This is always within a namespace, but in many cases that may just
// be the default global namespace ("\")
type GlobalScope struct {
	*Namespace
	*Scope
}

func NewGlobalScope(ns *Namespace) *GlobalScope {
	return &GlobalScope{ns, nil}
}

// scope represents a particular local scope (such as within a function).
type Scope struct {
	Identifiers      map[string]VariableGroup
	DynamicVariables []*Variable
	EnclosingScope   *Scope
	GlobalScope      *GlobalScope
	SuperGlobalScope *SuperGlobalScope
}

type VariableGroup struct {
	References []*Variable
	Type       Type
}

func (s *Scope) Variable(v *Variable) {
	static := Static(v.Name)
	if static == nil {
		s.DynamicVariables = append(s.DynamicVariables, v)
		return
	}
	vg, ok := s.Identifiers[static.Value]
	if !ok {
		vg = VariableGroup{Type: Unknown}
	}
	vg.Type = vg.Type.Union(v.EvaluatesTo())
	for _, ref := range vg.References {
		ref.Type = vg.Type
	}
	vg.References = append(vg.References, v)
	s.Identifiers[static.Value] = vg
}

type File struct {
	Name      string
	Namespace *Namespace
	Nodes     []Node
}

type FileSet struct {
	Files           map[string]*File
	Namespaces      map[string]*Namespace
	GlobalNamespace *Namespace
	*Scope
}

func NewFileSet() *FileSet {
	ns := NewNamespace("/")
	gscope := NewGlobalScope(ns)
	scope := NewScope(nil, gscope, &SuperGlobalScope{})
	gscope.Scope = scope
	return &FileSet{
		Files:           make(map[string]*File),
		Namespaces:      make(map[string]*Namespace),
		GlobalNamespace: ns,
		Scope:           scope,
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
		Identifiers:      map[string]VariableGroup{},
		EnclosingScope:   parent,
		GlobalScope:      global,
		SuperGlobalScope: superGlobal,
	}
}
