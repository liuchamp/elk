package merge

import (
	"go/ast"
	"go/token"

	log "github.com/sirupsen/logrus"
)

func HeaderMerge(new, old, target *ast.File) error {
	nImp := getImports(new)
	lImp := getImports(old)
	var mergedImp []ast.Spec

	for _, spec := range lImp {
		mergedImp = append(mergedImp, spec)
	}

	for s, l := range nImp {
		n, ok := lImp[s]
		if !ok {
			mergedImp = append(mergedImp, l)
			continue
		}
		//存在
		// 并且都没有设置name
		if n.Name == nil && l.Name == nil {
			continue
		}
		// 并且都有设置name
		if n.Name != nil && l.Name != nil {
			if n.Name.Name == n.Name.Name {
				continue
			}
			mergedImp = append(mergedImp, l)
			log.Warnf("import confict %s , name: %s --> %s", s, n.Name.Name, l.Name.Name)
			continue
		}
		// 其他情况， 保留原有的导入代码
	}
	mgimp := ast.GenDecl{
		Doc:   new.Doc,
		Tok:   token.IMPORT,
		Specs: mergedImp,
	}
	target.Decls = append(target.Decls, &mgimp)
	return nil
}

func getImports(f *ast.File) map[string]*ast.ImportSpec {
	restMap := make(map[string]*ast.ImportSpec)
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.IMPORT {
			continue
		}

		for _, spec := range genDecl.Specs {
			imp, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			if imp.Path != nil {
				restMap[imp.Path.Value] = imp
			} else {
				log.Printf("Cat not parser import without path")
			}
		}
	}

	return restMap
}
