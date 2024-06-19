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

// 生成依赖注册代码

func serviceGen(g *gen.Graph, pr string) error {
	for _, node := range g.Nodes {
		err := domain(g, node, pr)
		if err != nil {
			return err
		}
	}
	return nil
}

func serviceHeader(g *gen.Graph) (*token.FileSet, *ast.File) {
	// 创建文件集
	fset := token.NewFileSet()

	// 创建包名
	file := &ast.File{
		Name: ast.NewIdent(serverPkgName),
	}

	astutil.AddNamedImport(fset, file, "log", "github.com/go-kratos/kratos/v2/log")
	astutil.AddNamedImport(fset, file, consts.DefsvPkgName, consts.GetDefPackageName(g.Package))
	astutil.AddNamedImport(fset, file, consts.BoPkgName, consts.GetBoPackageName(g.Package))
	astutil.AddNamedImport(fset, file, consts.DtoPkgName, consts.GetDtoPackageName(g.Package))

	return fset, file
}
func domain(g *gen.Graph, n *gen.Type, pr string) error {
	fset, file := serviceHeader(g)

	file.Decls = append(file.Decls, domainGet(n))
	// 打印生成的代码
	f := filepath.Join(pr, routerPkgName, serverPkgName, fmt.Sprintf("%s.go", utils.SnakeToCamel(n.Name)))
	return write.WireGoFile(f, fset, file)
}

func domainGet(n *gen.Type) *ast.FuncDecl {

	// body体
	body := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "m",
				},
				&ast.Ident{
					Name: "err",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "h",
							},
							Sel: &ast.Ident{
								Name: utils.SnakeToCamel(n.Name),
							},
						},
						Sel: &ast.Ident{
							Name: consts.DefGetFuncName,
						},
					},
					Args: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
						},
						&ast.Ident{
							Name: "req",
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: "err",
				},
				Op: token.EQL,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: "m",
							},
							&ast.Ident{
								Name: "nil",
							},
						},
					},
				},
			},
		},
	}

	// 异常处理
	body = append(body, &ast.SwitchStmt{
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "ent",
								},
								Sel: &ast.Ident{
									Name: "IsNotFound",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "err",
								},
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "ent",
								},
								Sel: &ast.Ident{
									Name: "IsNotSingular",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "err",
								},
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "ent",
								},
								Sel: &ast.Ident{
									Name: "IsValidationError",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "err",
								},
							},
						},
					},
				},
				&ast.CaseClause{},
			},
		},
	})

	// 返回值
	body = append(body, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.Ident{
				Name: "nil",
			},
			&ast.Ident{
				Name: "nil",
			},
		},
	})
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Names: []*ast.Ident{
						&ast.Ident{
							Name: "h",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: "httpServer",
						},
					},
				},
			},
		},
		Name: ast.NewIdent(GetName(n)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ctx",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Context",
							},
						},
					},
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "req",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: consts.BoPkgName,
								},
								Sel: &ast.Ident{
									Name: utils.ToCamelCase(n.Name) + consts.GetBoSuffix,
								},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: consts.DtoPkgName,
								},
								Sel: &ast.Ident{
									Name: utils.ToCamelCase(n.Name) + consts.DtoSuffix,
								},
							},
						},
					},
					&ast.Field{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}
