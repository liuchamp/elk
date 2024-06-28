package merge

import (
	"go/ast"
	"go/token"
)

func ConstsMerge(new, old, target *ast.File) error {
	nVars := getVars(new, token.CONST)
	lVars := getVars(old, token.CONST)

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
		Tok:   token.CONST,
		Specs: mergedImp,
	}
	target.Decls = append(target.Decls, &mgimp)

	return nil
}
