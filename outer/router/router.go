package router

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/masseelch/elk/internal/consts"
	"github.com/masseelch/elk/pkg/utils"
	"github.com/masseelch/elk/pkg/utils/write"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"path/filepath"
)

const routerPkgName = "router"

// get name
const GetRouterSuffix = "Get"

func GetName(n *gen.Type) string {

	return fmt.Sprintf("%s%s", utils.ToCamelCase(n.Name), GetRouterSuffix)
}

func routerGen(g *gen.Graph, pr string) error {

	fset, file := routerHeader(g)

	// 操作id
	var operationIds []ast.Spec
	//http server 定义
	var iInterfaceTypes []*ast.Field

	// 注册router
	var register []ast.Stmt

	register = append(register, &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{
				Name: "r",
			},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "s",
					},
					Sel: &ast.Ident{
						Name: "Route",
					},
				},
				Args: []ast.Expr{
					&ast.Ident{
						Name: "pre",
					},
				},
			},
		},
	}) // 声明
	// router 实现

	file.Decls = append(file.Decls, &ast.GenDecl{
		Tok:   token.CONST,
		Specs: operationIds,
	},
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{
						Name: "HTTPServer",
					},
					Type: &ast.InterfaceType{
						Methods: &ast.FieldList{
							List: iInterfaceTypes,
						},
					},
				},
			},
		},
		&ast.FuncDecl{
			Name: &ast.Ident{
				Name: "RegisterHTTPServer",
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "pre",
								},
							},
							Type: &ast.Ident{
								Name: "string",
							},
						},
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "s",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "http",
									},
									Sel: &ast.Ident{
										Name: "Server",
									},
								},
							},
						},
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "srv",
								},
							},
							Type: &ast.Ident{
								Name: "HTTPServer",
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: register,
			},
		})

	f := filepath.Join(pr, routerPkgName, fmt.Sprintf("router.go"))
	return write.WireGoFile(f, fset, file)
}

func routerHeader(g *gen.Graph) (*token.FileSet, *ast.File) {
	// 创建文件集
	fset := token.NewFileSet()

	// 创建包名
	file := &ast.File{
		Name: ast.NewIdent(serverPkgName),
	}

	astutil.AddNamedImport(fset, file, "context", "context")
	astutil.AddNamedImport(fset, file, "http", "github.com/go-kratos/kratos/v2/transport/http")
	astutil.AddNamedImport(fset, file, consts.BoPkgName, consts.GetBoPackageName(g.Package))
	astutil.AddNamedImport(fset, file, consts.DtoPkgName, consts.GetDtoPackageName(g.Package))

	return fset, file
}
