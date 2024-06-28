package merge

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

// parseFile parses a Go source file and returns its AST.
func parseFile(filename string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}
	return fset, node, nil
}

// parseBytes parses a Go source bytes and returns its AST.
func parseBytes(bs []byte) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", string(bs), parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}
	return fset, node, nil
}

// FormatNode formats an AST node and returns the formatted code as a string.
func FormatNode(node *ast.File) (string, error) {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	// 使用 go/format 包来格式化代码
	if err := format.Node(&buf, fset, node); err != nil {
		panic(err)
	}
	return buf.String(), nil
}
