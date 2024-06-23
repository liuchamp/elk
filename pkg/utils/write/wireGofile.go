package write

import (
	"bytes"
	"github.com/liuchamp/elk/pkg/utils"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"golang.org/x/tools/imports"
)

func WireGoFile(path string, fset *token.FileSet, astfile *ast.File) error {

	f, err := utils.CreateFileWithDirs(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// 创建一个缓冲区来存储生成的代码
	buf := new(bytes.Buffer)

	// 打印未格式化的代码到缓冲区
	if err := printer.Fprint(buf, fset, astfile); err != nil {
		return err
	}

	// 获取未格式化的代码
	unformattedCode := buf.Bytes()

	// 格式化代码
	formattedCode, err := format.Source(unformattedCode)
	if err != nil {
		return err
	}

	// 格式化代码并移除未使用的包
	rmFormattedCode, err := imports.Process(path, formattedCode, nil)
	if err != nil {
		return err
	}

	// 将格式化的代码写入文件
	if _, err := f.Write(rmFormattedCode); err != nil {
		return err
	}

	return nil

}
