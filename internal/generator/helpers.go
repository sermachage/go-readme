package generator

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
)

// oneLiner returns the first sentence of a documentation string.
func oneLiner(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if idx := strings.IndexByte(s, '.'); idx >= 0 {
		return strings.TrimSpace(s[:idx+1])
	}
	return s
}

// funcDeclString formats a *ast.FuncDecl signature (without body) as a string.
func funcDeclString(decl *ast.FuncDecl) string {
	if decl == nil {
		return ""
	}
	fset := token.NewFileSet()
	// Clone the decl without body so we only print the signature.
	sig := &ast.FuncDecl{
		Recv: decl.Recv,
		Name: decl.Name,
		Type: decl.Type,
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, sig); err != nil {
		return decl.Name.Name
	}
	// Strip leading "func " so the caller can embed it naturally.
	result := buf.String()
	result = strings.TrimPrefix(result, "func ")
	return result
}

// valueDeclString formats a *ast.GenDecl (const/var block) as a string.
func valueDeclString(decl *ast.GenDecl) string {
	if decl == nil {
		return ""
	}
	fset := token.NewFileSet()
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, decl); err != nil {
		return ""
	}
	return buf.String()
}
