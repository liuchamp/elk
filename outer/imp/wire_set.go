package imp

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/masseelch/elk/pkg/utils/write"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"path/filepath"
)

// 生成依赖注册代码

func wireGen(g *gen.Graph, pr string) error {
	// 创建文件集
	fset := token.NewFileSet()

	// 创建包名
	file := &ast.File{
		Name: ast.NewIdent(impPkgName),
	}

	astutil.AddNamedImport(fset, file, "wire", "github.com/google/wire")

	var args []ast.Expr
	for _, n := range g.Nodes {
		args = append(args, ast.NewIdent(GetImpNewFuncName(n)))
	}
	file.Decls = append(file.Decls, &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					ast.NewIdent("ProviderSet"),
				},
				Values: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("wire"),
							Sel: ast.NewIdent("NewSet"),
						},
						Args: args,
					},
				},
			},
		},
	})

	// 打印生成的代码
	f := filepath.Join(pr, impPkgName, fmt.Sprintf("imp.go"))
	return write.WireGoFile(f, fset, file)
}
