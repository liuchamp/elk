package merge

import (
	"fmt"
	"go/ast"
)

type Hook func(new, old, target *ast.File) error

func NewCompare(srcOld any, srcNew any) (*ast.File, error) {
	var oldNodes, newNodes *ast.File
	var err error
	// 如果是[]byte , 那么直接解析， 如果是string, 当文件处理
	if fileName, ok := srcOld.(string); ok {
		_, oldNodes, err = parseFile(fileName)
		if err != nil {
			return nil, err
		}
	} else if fileBs, ok2 := srcOld.([]byte); ok2 {
		_, oldNodes, err = parseBytes(fileBs)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("cat not get file source %V", srcOld)
	}

	if fileName, ok := srcNew.(string); ok {
		_, newNodes, err = parseFile(fileName)
		if err != nil {
			return nil, err
		}
	} else if fileBs, ok2 := srcNew.([]byte); ok2 {
		_, newNodes, err = parseBytes(fileBs)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("cat not get file source %V", srcNew)
	}
	target := ast.File{
		Doc:      newNodes.Doc,
		Name:     oldNodes.Name,
		Comments: newNodes.Comments,
	}
	err = coinsOptChain(newNodes, oldNodes, &target, HeaderMerge, VarMerge, ConstsMerge, StructMerge, InterfaceMerge, FuncMerge, ExFuncMerge)
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func coinsOptChain(new, old, target *ast.File, pssors ...Hook) error {

	for _, pssor := range pssors {
		err := pssor(new, old, target)
		if err != nil {
			return err
		}
	}
	return nil
}
