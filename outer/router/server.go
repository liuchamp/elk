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
const serverPkgName = "server"

func wireGen(g *gen.Graph, pr string) error {
	// 创建文件集
	fset := token.NewFileSet()

	// 创建包名
	file := &ast.File{
		Name: ast.NewIdent(serverPkgName),
	}

	astutil.AddNamedImport(fset, file, "log", "github.com/go-kratos/kratos/v2/log")
	astutil.AddNamedImport(fset, file, consts.DefsvPkgName, consts.GetDefPackageName(g.Package))

	// struct定义

	file.Decls = append(file.Decls, defineServer(g))

	// 构造函数添加

	// 打印生成的代码
	f := filepath.Join(pr, serverPkgName, fmt.Sprintf("server.go"))
	return write.WireGoFile(f, fset, file)
}

func defineCons(g *gen.Graph) *ast.FuncDecl {
	var params []*ast.Field
	var elts []ast.Expr
	// log
	params = append(params, &ast.Field{
		Names: []*ast.Ident{
			&ast.Ident{
				Name: "logger",
			},
		},
		Type: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "log",
			},
			Sel: &ast.Ident{
				Name: "Logger",
			},
		},
	})
	elts = append(elts, &ast.KeyValueExpr{
		Key: &ast.Ident{
			Name: "log",
		},
		Value: &ast.Ident{
			Name: "logInfo",
		},
	})

	for _, n := range g.Nodes {
		params = append(params, &ast.Field{
			Names: []*ast.Ident{
				&ast.Ident{
					Name: utils.SnakeToCamel(n.Name),
				},
			},
			Type: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: consts.DefsvPkgName,
				},
				Sel: &ast.Ident{
					Name: utils.ToCamelCase(n.Name) + consts.DefRepoSuffix,
				},
			},
		})
	}
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "NewHttpServer",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: params,
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: "httpServer",
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "logInfo",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "log",
								},
								Sel: &ast.Ident{
									Name: "NewHelper",
								},
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "log",
										},
										Sel: &ast.Ident{
											Name: "With",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "logger",
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"module\"",
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"imp/pet\"",
										},
									},
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "srv",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.Ident{
									Name: "httpServer",
								},
								Elts: elts,
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{
							Name: "srv",
						},
					},
				},
			},
		},
	}
}
func defineServer(g *gen.Graph) *ast.GenDecl {

	var feilds []*ast.Field
	// log
	feilds = append(feilds, &ast.Field{
		Names: []*ast.Ident{
			ast.NewIdent("log"),
		},
		Type: &ast.StarExpr{
			X: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "log",
				},
				Sel: &ast.Ident{
					Name: "Helper",
				},
			},
		},
	})

	// 各个def 添加
	for _, node := range g.Nodes {
		feilds = append(feilds, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(node.Name),
			},
			Type: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: consts.DefsvPkgName,
				},
				Sel: &ast.Ident{
					Name: node.Name + consts.DefRepoSuffix,
				},
			},
		})
	}

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "httpServer",
				},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: feilds,
					},
				},
			},
		},
	}
}
