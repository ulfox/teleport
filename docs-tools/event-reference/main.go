package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sync"
)

// Top-level directories for source code where EmitAuditEvent is called
var pkg = []string{
	"api",
	"e",
	"integration",
	"lib",
}

type AuditEventEmitCollection struct {
	Calls []*ast.CallExpr
	mu    *sync.Mutex
}

func (a *AuditEventEmitCollection) addEmitCall(c *ast.CallExpr) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Calls = append(a.Calls, c)
}

// TODO: Use astutil.Apply

func main() {
	fs := token.NewFileSet()
	for _, p := range pkg {
		a, err := parser.ParseDir(fs, p, nil, 0)
		if err != nil {
			// TODO: Replace with proper logger call
			fmt.Fprintf(os.Stderr, "error parsing Go source files: %v", err)
			os.Exit(1)
		}

		for _, c := range a {
			for _, f := range c.Files {
				for _, d := range f.Decls {
					// TODO: Call the astutil.Apply function apply function (https://pkg.go.dev/golang.org/x/tools/go/ast/astutil#Apply)
					// find EmitAuditEvent calls and do something with them
				}
			}
		}
	}
}
