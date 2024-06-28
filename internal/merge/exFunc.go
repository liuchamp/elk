package merge

import (
	"go/ast"
)

func ExFuncMerge(new, old, target *ast.File) error {
	nFunc := getExFunc(new)
	lFunc := getExFunc(old)

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

func getExFunc(f *ast.File) map[string]*ast.FuncDecl {
	restMap := make(map[string]*ast.FuncDecl)
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		print(funcDecl)

		if funcDecl.Recv == nil {
			continue
		}
		Sft := funcDecl.Recv.List[0].Names[0].Name
		restMap[Sft+funcDecl.Name.Name] = funcDecl
	}

	return restMap
}
