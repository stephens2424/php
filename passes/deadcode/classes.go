package deadcode

import "github.com/stephens2424/php/ast"

// DeadClasses returns a list of dead classes
func DeadClasses(fs *ast.FileSet, entryPoints []string) []ast.Node {
	knownClasses := AllTheClasses(fs)

	for _, filename := range entryPoints {
		f, ok := fs.Files[filename]
		if !ok {
			continue
		}

		EliminateClasses(f.Nodes, knownClasses)
	}

	nodes := []ast.Node{}
	for _, f := range knownClasses {
		nodes = append(nodes, f)
	}

	return nodes
}

// EliminateClasses eliminates all dead classes
func EliminateClasses(nodes []ast.Node, knownClasses map[string]ast.Node) {
	for _, node := range nodes {
		switch node := node.(type) {
		case ast.NewCallExpr:
			if static := ast.Static(node.Class); static != nil {
				delete(knownClasses, static.Value)
			}
		case *ast.NewCallExpr:
			if static := ast.Static(node.Class); static != nil {
				delete(knownClasses, static.Value)
			}
		case ast.ClassExpr:
			if static := ast.Static(node.Receiver); static != nil {
				delete(knownClasses, static.Value)
			}
		case *ast.ClassExpr:
			if static := ast.Static(node.Receiver); static != nil {
				delete(knownClasses, static.Value)
			}
		}
		EliminateClasses(node.Children(), knownClasses)
	}
}

// AllTheClasses returns a list of all classes
func AllTheClasses(fs *ast.FileSet) map[string]ast.Node {
	namedClasses := map[string]ast.Node{}
	for _, n := range fs.GlobalNamespace.ClassesAndInterfaces {
		if class, ok := n.(*ast.Class); ok {
			namedClasses[class.Name] = class
		}
	}

	for _, ns := range fs.Namespaces {
		for _, class := range ns.ClassesAndInterfaces {
			if class, ok := class.(*ast.Class); ok {
				namedClasses[class.Name] = class
			}
		}
	}
	return namedClasses
}
