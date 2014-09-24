// Package query implements CSS/jQuery-esque selectors over a PHP AST
//
// Currently, the only supported selectors are:
//
// "node", where node is the type name of an AST node (e.g. EchoStmt)
//
// "parent child", where the descendance may skip levels.
package query

import (
	"fmt"
	"strings"

	"github.com/stephens2424/php/ast"
)

type Q []Node

func Select(nodes []ast.Node) Q {
	flat := make([]Node, 0, len(nodes))
	flat = flatten(nodes, flat, nil)
	return flat
}

func (q Q) Select(s string) (Q, error) {
	sel, err := ParseSelector(s)
	if err != nil {
		return nil, err
	}
	passed := make([]Node, 0)
	for _, node := range q {
		if sel.Pass(node, true) {
			passed = append(passed, node)
		}
	}
	return passed, nil
}

func ParseSelector(s string) (*Selector, error) {
	var previous *Selector
	parts := strings.Split(s, " ")
	for _, part := range parts {
		s := &Selector{parent: previous, localRules: parseRuleSet(part)}
		previous = s
	}
	return previous, nil
}

func parseRuleSet(s string) []Rule {
	if len(s) == 0 {
		return nil
	}
	return []Rule{NodeRule{s}}
}

type Selector struct {
	localRules []Rule
	parent     *Selector
}

func (s Selector) Pass(n Node, local bool) bool {
	for _, rule := range s.localRules {
		passed := rule.Pass(n)

		// Pass and no more rules
		if passed && s.parent == nil {
			return true
		}

		// Pass, but we have parent elements and rules to check
		if passed && n.Parent != nil {
			if s.parent.Pass(*n.Parent, false) {
				return true
			}
		}

		// Didn't pass, but we have a parent element we could try
		if !passed && n.Parent != nil && !local {
			if s.Pass(*n.Parent, false) {
				return true
			}
		}

	}
	return false
}

type Rule interface {
	Pass(n Node) bool
}

type NodeRule struct {
	// Type is the node type represented as a string, without the package name. For example: ReturnStmt.
	Type string
}

func (r NodeRule) Pass(n Node) bool {
	typename := fmt.Sprintf("%T", n.Node)[3:]
	typename = typename[strings.Index(typename, ".")+1:]
	return r.Type == typename
}

func flatten(nodes []ast.Node, flat []Node, parent *Node) []Node {
	for i, node := range nodes {
		self := &Node{nodes[i], parent}
		flat = append(flat, *self)
		children := node.Children()
		if len(children) != 0 {
			flat = flatten(children, flat, self)
		}
	}
	return flat
}

type Node struct {
	Node   ast.Node
	Parent *Node
}
