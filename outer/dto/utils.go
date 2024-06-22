package dto

import (
	"go/ast"
	"go/token"
	"path/filepath"

	"github.com/masseelch/elk/internal/consts"
	"github.com/masseelch/elk/pkg/utils/write"
)

func genUtilFile(pr string) error {

	// 创建文件集
	fset := token.NewFileSet()

	// 创建包名
	file := &ast.File{
		Name: ast.NewIdent(dtoPkgName),
	}

	// page result
	file.Decls = append(file.Decls, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: consts.PageResultNameStr,
				},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "Total",
									},
								},
								Type: &ast.Ident{
									Name: "int",
								},
								Tag: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "`json:\"total\"`",
								},
							},
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "Data",
									},
								},
								Type: &ast.Ident{
									Name: "any",
								},
								Tag: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "`json:\"data\"`",
								},
							},
						},
					},
				},
			},
		},
	})
	// 打印生成的代码
	f := filepath.Join(pr, dtoPkgName, "utils.go")
	return write.WireGoFile(f, fset, file)
}
