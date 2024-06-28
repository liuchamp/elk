package merge

import (
	log "github.com/sirupsen/logrus"
	"go/ast"
)

func FuncMerge(new, old, target *ast.File) error {
	nFunc := getFunc(new)
	lFunc := getFunc(old)

	var mergedImp []*ast.FuncDecl
	conflictMap := make(map[string]bool)
	for s, spec := range lFunc {
		sf1, ok := nFunc[s]
		if ok {
			cpSta, err := funcComparators(spec, sf1)
			if err != nil {
				return err
			}
			if !cpSta {
				conflictMap[s] = true
				spec.Name = ast.NewIdent(spec.Name.Name + "_old")
			}
		}

		mergedImp = append(mergedImp, spec)
	}
	for s, l := range nFunc {
		_, ok := lFunc[s]
		if !ok {

			mergedImp = append(mergedImp, l)
			continue
		}
		_, iOk := conflictMap[s]
		if !iOk {
			continue
		}
		// 如果存在， 新的函数重命名
		l.Name = ast.NewIdent(l.Name.Name + "_new")
		mergedImp = append(mergedImp, l)
	}

	for _, decl := range mergedImp {
		target.Decls = append(target.Decls, decl)
	}

	return nil
}

func getFunc(f *ast.File) map[string]*ast.FuncDecl {
	restMap := make(map[string]*ast.FuncDecl)
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Recv != nil {
			continue
		}
		restMap[funcDecl.Name.Name] = funcDecl
	}

	return restMap
}

func funcComparators(f1, f2 *ast.FuncDecl) (bool, error) {

	funcWapper := &ast.File{
		Name: &ast.Ident{
			Name: "ent",
		},
	}

	funcWapper.Decls = []ast.Decl{
		f1,
	}

	f1string, err := FormatNode(funcWapper)
	if err != nil {

		log.Error("cat not parser data", err)
		return false, err
	}

	funcWapper.Decls = []ast.Decl{
		f2,
	}

	f2string, err := FormatNode(funcWapper)
	if err != nil {

		log.Error("cat not parser data", err)
		return false, err
	}

	return f2string == f1string, nil
}
