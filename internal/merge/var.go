package merge

import (
	"go/ast"
	"go/token"

	log "github.com/sirupsen/logrus"
)

func VarMerge(new, old, target *ast.File) error {
	nVars := getVars(new, token.VAR)
	lVars := getVars(old, token.VAR)

	var mergedImp []ast.Spec

	for _, spec := range lVars {
		//l, ok := nVars[s]
		mergedImp = append(mergedImp, spec)
	}
	for s, l := range nVars {
		_, ok := lVars[s]
		if !ok {
			mergedImp = append(mergedImp, l)
			continue
		}
		// 如果存在， 保持原有值
	}
	mgimp := ast.GenDecl{
		Doc:   new.Doc,
		Tok:   token.VAR,
		Specs: mergedImp,
	}
	target.Decls = append(target.Decls, &mgimp)

	return nil
}

func getVars(f *ast.File, t token.Token) map[string]*ast.ValueSpec {
	restMap := make(map[string]*ast.ValueSpec)
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != t {
			continue
		}

		for _, spec := range genDecl.Specs {
			imp, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if imp.Names != nil {
				restMap[imp.Names[0].Name] = imp
			} else {
				log.Printf("Cat not parser import without path")
			}
		}
	}

	return restMap
}
