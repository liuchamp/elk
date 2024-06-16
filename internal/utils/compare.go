package utils

import (
	"go/ast"
	"reflect"
)

// FindFuncDecl 查找给定文件中的指定函数声明
func FindFuncDecl(f *ast.File, funcName string) *ast.FuncDecl {
	for _, decl := range f.Decls {
		if fn, isFn := decl.(*ast.FuncDecl); isFn && fn.Name.Name == funcName {
			return fn
		}
	}
	return nil
}

// CompareFuncs 比较两个函数声明是否相同
func CompareFuncs(f1, f2 *ast.FuncDecl) bool {
	return reflect.DeepEqual(f1, f2)
}
