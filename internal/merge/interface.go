package merge

import (
	log "github.com/sirupsen/logrus"
	"go/ast"
	"go/token"
)

func InterfaceMerge(new, old, target *ast.File) error {
	nStruct := getInterface(new)
	lStruct := getInterface(old)

	mergedImp := make(map[string]*ast.InterfaceType)

	conflictMap := make(map[string]bool)

	for s, st := range lStruct {
		sf1, ok := nStruct[s]
		if !ok {
			mergedImp[s] = st
			continue
		}
		cpSta, err := interfaceComparators(st, sf1)
		if err != nil {
			return err
		}
		if !cpSta {
			conflictMap[s] = true
			mergedImp[s+"old"] = st

		}
	}

	for s, l := range nStruct {
		_, ok := lStruct[s]
		if !ok {
			mergedImp[s] = l
			continue
		}
		_, iOk := conflictMap[s]
		if !iOk {
			continue
		}
		// 如果存在， 新的函数重命名
		mergedImp[s+"_new"] = l
	}

	for s, structType := range mergedImp {
		var ssInner []ast.Spec
		ssInner = append(ssInner, &ast.TypeSpec{Name: ast.NewIdent(s), Type: structType})
		stmp := &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: ssInner,
		}
		target.Decls = append(target.Decls, stmp)
	}

	return nil
}

func getInterface(f *ast.File) map[string]*ast.InterfaceType {
	restMap := make(map[string]*ast.InterfaceType)
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			imp, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			stcInfo, okS := imp.Type.(*ast.InterfaceType)
			if !okS {
				continue
			}
			restMap[imp.Name.Name] = stcInfo
		}
	}

	return restMap
}

func interfaceComparators(st *ast.InterfaceType, sf1 *ast.InterfaceType) (bool, error) {
	funcWapper := &ast.File{
		Name: &ast.Ident{
			Name: "ent",
		},
	}

	funcWapper.Decls = []ast.Decl{
		&ast.GenDecl{Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{Type: sf1, Name: ast.NewIdent("test")},
			},
		},
	}
	f1string, err := FormatNode(funcWapper)
	if err != nil {

		log.Error("cat not parser data", err)
		return false, err
	}
	funcWapper.Decls = []ast.Decl{
		&ast.GenDecl{Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{Type: st, Name: ast.NewIdent("test")},
			},
		},
	}
	f2string, err := FormatNode(funcWapper)
	if err != nil {

		log.Error("cat not parser data", err)
		return false, err
	}
	return f2string == f1string, nil
}
