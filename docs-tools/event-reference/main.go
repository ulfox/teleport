package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/tools/go/ast/astutil"
)

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
	s := token.NewFileSet()
	filepath.Walk(path.Join("..", ".."), func(pth string, i fs.FileInfo, err error) error {
		if strings.HasSuffix(i.Name(), ".go") {
			f, err := parser.ParseFile(s, pth, nil, 0)
			for _, d := range f.Decls {
				astutil.Apply(d, func(c *astutil.Cursor) bool {
					if i, ok := c.Node().(*ast.Ident); ok && i.Name == "EmitAuditEvent" {
						// TODO: Figure out what "i" indicates within the source, and how to get
						// the Metadata of the argument to EmitAuditEvent
						fmt.Printf("this is an ident: %v\n", i)
					}
					return true
				}, nil)
			}
			if err != nil {
				// TODO: Replace with proper logger call
				fmt.Fprintf(os.Stderr, "error parsing Go source files: %v", err)
				os.Exit(1)
			}
		}
		return nil
	})
}
