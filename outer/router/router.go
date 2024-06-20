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
	var rouerImp []*ast.FuncDecl

	for _, n := range g.Nodes {
		operationIds = append(operationIds, optIds(n)...)
		iInterfaceTypes = append(iInterfaceTypes, interfaceGen(n)...)
		register = append(register, registerGen(n)...)

		rouerImp = append(rouerImp, routerImpGen(n)...)
	}
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
		},
	)
	for _, decl := range rouerImp {
		file.Decls = append(file.Decls, decl)
	}

	f := filepath.Join(pr, routerPkgName, fmt.Sprintf("router.go"))
	return write.WireGoFile(f, fset, file)
}

func routerImpGen(n *gen.Type) []*ast.FuncDecl {

	getRouterFunc := routerGetImp(n)
	//patchRouterFunc := routerPatchImp(n)
	//createRouterFunc := routerCreateImp(n)
	//listRouterFunc := routerListImp(n)

	return []*ast.FuncDecl{getRouterFunc}
}

func routerGetImp(n *gen.Type) *ast.FuncDecl {

	// query 查询
	//&ast.IfStmt{
	//	Init: &ast.AssignStmt{
	//		Lhs: []ast.Expr{
	//			&ast.Ident{
	//				Name: "err",
	//			},
	//		},
	//		Tok: token.DEFINE,
	//		Rhs: []ast.Expr{
	//			&ast.CallExpr{
	//				Fun: &ast.SelectorExpr{
	//					X: &ast.Ident{
	//						Name: "ctx",
	//					},
	//					Sel: &ast.Ident{
	//						Name: "BindQuery",
	//					},
	//				},
	//				Args: []ast.Expr{
	//					&ast.UnaryExpr{
	//						Op: token.AND,
	//						X: &ast.Ident{
	//							Name: "in",
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//	Cond: &ast.BinaryExpr{
	//		X: &ast.Ident{
	//			Name: "err",
	//		},
	//		Op: token.NEQ,
	//		Y: &ast.Ident{
	//			Name: "nil",
	//		},
	//	},
	//	Body: &ast.BlockStmt{
	//		List: []ast.Stmt{
	//			&ast.ReturnStmt{
	//				Results: []ast.Expr{
	//					&ast.Ident{
	//						Name: "err",
	//					},
	//				},
	//			},
	//		},
	//	},
	//},

	bodyStmlx := []ast.Stmt{
		// in 定义
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "in",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: consts.BoPkgName,
							},
							Sel: &ast.Ident{
								Name: fmt.Sprintf("%s%s", utils.ToCamelCase(n.Name), consts.GetBoSuffix),
							},
						},
					},
				},
			},
		},

		// id 绑定
		&ast.IfStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "err",
					},
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "ctx",
							},
							Sel: &ast.Ident{
								Name: "BindVars",
							},
						},
						Args: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.Ident{
									Name: "in",
								},
							},
						},
					},
				},
			},
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: "err",
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		// 设置操作id
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "http",
					},
					Sel: &ast.Ident{
						Name: "SetOperation",
					},
				},
				Args: []ast.Expr{
					&ast.Ident{
						Name: "ctx",
					},
					&ast.Ident{
						Name: fmt.Sprintf("Operation%s%s", utils.ToCamelCase(n.Name), consts.DefGetFuncName),
					},
				},
			},
		},
		// 回调server 声明 h
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "h",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "ctx",
						},
						Sel: &ast.Ident{
							Name: "Middleware",
						},
					},
					Args: []ast.Expr{
						&ast.FuncLit{
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
											Type: &ast.InterfaceType{
												Methods: &ast.FieldList{},
											},
										},
									},
								},
								Results: &ast.FieldList{
									List: []*ast.Field{
										&ast.Field{
											Type: &ast.InterfaceType{
												Methods: &ast.FieldList{},
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
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "srv",
													},
													Sel: &ast.Ident{
														Name: fmt.Sprintf("%s%s", consts.DefGetFuncName, utils.ToCamelCase(n.Name)),
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "ctx",
													},
													&ast.TypeAssertExpr{
														X: &ast.Ident{
															Name: "req",
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
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// 执行回调h
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "out",
				},
				&ast.Ident{
					Name: "err",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: "h",
					},
					Args: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
						},
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.Ident{
								Name: "in",
							},
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
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		// 返回值判断
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "reply",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.TypeAssertExpr{
					X: &ast.Ident{
						Name: "out",
					},
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
			},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "ctx",
						},
						Sel: &ast.Ident{
							Name: "Result",
						},
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.INT,
							Value: "200",
						},
						&ast.Ident{
							Name: "reply",
						},
					},
				},
			},
		},
	}

	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: fmt.Sprintf("_%s_%s_0_HTTP_Handler", utils.ToCamelCase(n.Name), consts.DefGetFuncName),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "srv",
							},
						},
						Type: &ast.Ident{
							Name: "HTTPServer",
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
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
												Name: "http",
											},
											Sel: &ast.Ident{
												Name: "Context",
											},
										},
									},
								},
							},
							Results: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Type: &ast.Ident{
											Name: "error",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.FuncLit{
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
													Name: "http",
												},
												Sel: &ast.Ident{
													Name: "Context",
												},
											},
										},
									},
								},
								Results: &ast.FieldList{
									List: []*ast.Field{
										&ast.Field{
											Type: &ast.Ident{
												Name: "error",
											},
										},
									},
								},
							},
							Body: &ast.BlockStmt{
								List: bodyStmlx,
							},
						},
					},
				},
			},
		},
	}
}
func registerGen(n *gen.Type) []ast.Stmt {
	genfun := func(fn string, path string) ast.Stmt {

		return &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "r",
					},
					Sel: &ast.Ident{
						Name: fn,
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: path, //fmt.Sprintf("\"/%s/{id}\"", utils.SnakeToCamel(n.Name)),
					},
					&ast.CallExpr{
						Fun: &ast.Ident{
							Name: fmt.Sprintf("_%s_%s_0_HTTP_Handler", utils.ToCamelCase(n.Name), fn),
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "srv",
							},
						},
					},
				},
			},
		}
	}
	getFunc := genfun(consts.DefGetFuncName, fmt.Sprintf("\"/%ss/{id}\"", utils.SnakeToCamel(n.Name)))
	listFunc := genfun(consts.DefListFuncName, fmt.Sprintf("\"/%ss\"", utils.SnakeToCamel(n.Name)))
	patchFunc := genfun(consts.DefPatchFuncName, fmt.Sprintf("\"/%ss/{id}\"", utils.SnakeToCamel(n.Name)))
	createFunc := genfun(consts.DefCreateFuncName, fmt.Sprintf("\"/%ss\"", utils.SnakeToCamel(n.Name)))

	return []ast.Stmt{createFunc, patchFunc, getFunc, listFunc}
}
func interfaceGen(n *gen.Type) []*ast.Field {

	iifc := func(fn, bosfix string) *ast.Field {
		return &ast.Field{
			Names: []*ast.Ident{
				&ast.Ident{
					Name: fmt.Sprintf("%s%s", utils.ToCamelCase(n.Name), fn),
				},
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
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
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: consts.BoPkgName,
									},
									Sel: &ast.Ident{
										Name: fmt.Sprintf("%s%s", utils.ToCamelCase(n.Name), bosfix),
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
										Name: fmt.Sprintf("%s%s", utils.ToCamelCase(n.Name), consts.DtoSuffix),
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
		}
	}
	iis := []*ast.Field{
		iifc(consts.DefGetFuncName, consts.GetBoSuffix),
		iifc(consts.DefPatchFuncName, consts.PatchBoSuffix),
		iifc(consts.DefCreateFuncName, consts.CreateBoSuffix),
		{
			Names: []*ast.Ident{
				&ast.Ident{
					Name: fmt.Sprintf("%s%s", utils.ToCamelCase(n.Name), consts.DefListFuncName),
				},
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
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
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: consts.BoPkgName,
									},
									Sel: &ast.Ident{
										Name: fmt.Sprintf("%s%s", utils.ToCamelCase(n.Name), consts.ListBoSuffix),
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
										Name: fmt.Sprintf("%s%s", utils.ToCamelCase(n.Name), consts.DtoSuffix),
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
		},
	}
	return iis
}

func optIds(n *gen.Type) []ast.Spec {
	getOptName, getOptID := genOptID(n)
	pathOptName, patchOptIDe := patchOptID(n)
	listOptName, listOptIDe := listOptID(n)
	createOptName, createOptIDe := createOptID(n)

	opst := []ast.Spec{
		&ast.ValueSpec{
			Names: []*ast.Ident{ast.NewIdent(getOptName)},
			Values: []ast.Expr{
				&ast.BasicLit{Kind: token.STRING, Value: getOptID},
			},
		},
		&ast.ValueSpec{
			Names: []*ast.Ident{ast.NewIdent(pathOptName)},
			Values: []ast.Expr{
				&ast.BasicLit{Kind: token.STRING, Value: patchOptIDe},
			},
		},
		&ast.ValueSpec{
			Names: []*ast.Ident{ast.NewIdent(listOptName)},
			Values: []ast.Expr{
				&ast.BasicLit{Kind: token.STRING, Value: listOptIDe},
			},
		},
		&ast.ValueSpec{
			Names: []*ast.Ident{ast.NewIdent(createOptName)},
			Values: []ast.Expr{
				&ast.BasicLit{Kind: token.STRING, Value: createOptIDe},
			},
		},
	}
	return opst
}

func genOptID(n *gen.Type) (string, string) {
	return fmt.Sprintf("Operation%s%s", utils.ToCamelCase(n.Name), consts.DefGetFuncName),
		fmt.Sprintf("/%s/%s", utils.SnakeToCamel(n.Name), utils.SnakeToCamel(consts.DefGetFuncName))
}
func listOptID(n *gen.Type) (string, string) {
	return fmt.Sprintf("Operation%s%s", utils.ToCamelCase(n.Name), consts.DefListFuncName),
		fmt.Sprintf("/%s/%s", utils.SnakeToCamel(n.Name), utils.SnakeToCamel(consts.DefListFuncName))
}
func patchOptID(n *gen.Type) (string, string) {
	return fmt.Sprintf("Operation%s%s", utils.ToCamelCase(n.Name), consts.DefPatchFuncName),
		fmt.Sprintf("/%s/%s", utils.SnakeToCamel(n.Name), utils.SnakeToCamel(consts.DefPatchFuncName))
}
func createOptID(n *gen.Type) (string, string) {
	return fmt.Sprintf("Operation%s%s", utils.ToCamelCase(n.Name), consts.DefCreateFuncName),
		fmt.Sprintf("/%s/%s", utils.SnakeToCamel(n.Name), utils.SnakeToCamel(consts.DefCreateFuncName))
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
