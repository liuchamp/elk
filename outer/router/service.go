package router

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/internal/consts"
	"github.com/masseelch/elk/pkg/utils"
	"github.com/masseelch/elk/pkg/utils/write"
	"golang.org/x/tools/go/ast/astutil"
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
		Name: ast.NewIdent(routerPkgName),
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
	file.Decls = append(file.Decls, domainList(n))
	file.Decls = append(file.Decls, domainCreate(n))
	file.Decls = append(file.Decls, domainPatch(n))
	// 打印生成的代码
	f := filepath.Join(pr, routerPkgName, fmt.Sprintf("%s.go", utils.SnakeToCamel(n.Name)))
	return write.WireGoFile(f, fset, file)
}

func domainPatch(n *gen.Type) *ast.FuncDecl {

	// body体
	body := bodyHis(n, consts.DefPatchFuncName)

	// 异常处理
	body = append(body, bodyErrDuty(n)...)

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
		Recv: funcHead(n),
		Name: funcName(n, consts.DefPatchFuncName),
		Type: funcParam(n, consts.PatchBoSuffix),
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}
func domainCreate(n *gen.Type) *ast.FuncDecl {

	// body体
	body := bodyHis(n, consts.DefCreateFuncName)

	// 异常处理
	body = append(body, bodyErrDuty(n)...)

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
		Recv: funcHead(n),
		Name: funcName(n, consts.DefCreateFuncName),
		Type: funcParam(n, consts.CreateBoSuffix),
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}
func domainList(n *gen.Type) *ast.FuncDecl {
	// body体
	body := bodyHis(n, consts.DefListFuncName)

	// 异常处理
	body = append(body, bodyErrDuty(n)...)

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
		Recv: funcHead(n),
		Name: funcName(n, consts.DefListFuncName),
		Type: funcParam(n, consts.ListBoSuffix),
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}
func domainGet(n *gen.Type) *ast.FuncDecl {

	// body体
	body := bodyHis(n, consts.DefGetFuncName)
	// 异常处理
	body = append(body, bodyErrDuty(n)...)

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
		Recv: funcHead(n),
		Name: funcName(n, consts.DefGetFuncName),
		Type: funcParam(n, consts.GetBoSuffix),
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

func bodyHis(n *gen.Type, fn string) []ast.Stmt {

	lts := []ast.Expr{
		&ast.Ident{
			Name: "m",
		},
		&ast.Ident{
			Name: "err",
		},
	}
	rest := []ast.Stmt{
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
	}
	if fn == consts.DefListFuncName {
		lts = []ast.Expr{
			&ast.Ident{
				Name: "m",
			},
			&ast.Ident{
				Name: "total",
			},
			&ast.Ident{
				Name: "err",
			},
		}

		rest = []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X: &ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "dto",
								},
								Sel: &ast.Ident{
									Name: consts.PageResultNameStr,
								},
							},
							Elts: []ast.Expr{
								&ast.KeyValueExpr{
									Key: &ast.Ident{
										Name: "Total",
									},
									Value: &ast.Ident{
										Name: "total",
									},
								},
								&ast.KeyValueExpr{
									Key: &ast.Ident{
										Name: "Data",
									},
									Value: &ast.Ident{
										Name: "m",
									},
								},
							},
						},
					},
					&ast.Ident{
						Name: "nil",
					},
				},
			},
		}
	}
	body := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: lts,
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
							Name: fn,
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
				List: rest,
			},
		},
	}

	return body
}
func funcHead(n *gen.Type) *ast.FieldList {
	return &ast.FieldList{
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
	}
}

func funcName(n *gen.Type, fn string) *ast.Ident {
	return ast.NewIdent(fmt.Sprintf("%s%s", fn, utils.ToCamelCase(n.Name)))
}

func funcParam(n *gen.Type, fn string) *ast.FuncType {

	rest := &ast.StarExpr{
		X: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: consts.DtoPkgName,
			},
			Sel: &ast.Ident{
				Name: utils.ToCamelCase(n.Name) + consts.DtoSuffix,
			},
		},
	}
	if fn == consts.ListBoSuffix {
		rest = &ast.StarExpr{
			X: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: consts.DtoPkgName,
				},
				Sel: &ast.Ident{
					Name: consts.PageResultNameStr,
				},
			},
		}
	}
	return &ast.FuncType{
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
								Name: utils.ToCamelCase(n.Name) + fn,
							},
						},
					},
				},
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Type: rest,
				},
				&ast.Field{
					Type: &ast.Ident{
						Name: "error",
					},
				},
			},
		},
	}
}

func bodyErrDuty(n *gen.Type) []ast.Stmt {

	return []ast.Stmt{
		&ast.SwitchStmt{
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
		},
	}
}
