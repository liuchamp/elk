package imp

import (
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/internal/consts"
	"github.com/masseelch/elk/pkg/utils"
	"go/ast"
	"go/token"
)

func listImp(n *gen.Type) *ast.FuncDecl {
	const mpsv = "imc"
	// Create
	createParams := &ast.FieldList{
		List: []*ast.Field{
			{Names: []*ast.Ident{ast.NewIdent("ctx")}, Type: ast.NewIdent("context.Context")},
			{
				Names: []*ast.Ident{ast.NewIdent("req")},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(consts.BoPkgName),
						Sel: ast.NewIdent(utils.ToCamelCase(n.Name) + consts.ListBoSuffix)},
				},
			},
		},
	}
	createResults := &ast.FieldList{
		List: []*ast.Field{
			{ // 数组值
				Type: &ast.ArrayType{
					Elt: &ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent(consts.DtoPkgName), Sel: ast.NewIdent(utils.ToCamelCase(n.Name) + consts.DtoSuffix)}},
				},
			},
			{Type: ast.NewIdent("int")}, // 总值
			{Type: ast.NewIdent("error")},
		},
	}

	const opHandle = "query"
	var bodyStmt []ast.Stmt

	// 生成opHandle
	bodyStmt = append(bodyStmt, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent(opHandle),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.SelectorExpr{X: ast.NewIdent(mpsv + ".c"), Sel: ast.NewIdent(utils.ToCamelCase(n.Name))},
					Sel: ast.NewIdent("Query"),
				},
			},
		},
	})

	// id 过滤
	entityName := getEntityName(n)
	bodyStmt = append(bodyStmt, fieldQueryStmt(entityName, n.ID))
	//// 设置数据
	for _, field := range n.Fields {
		bodyStmt = append(bodyStmt, fieldQueryStmt(entityName, field))
	}

	const countHandle = "total"
	// 设置total 统计
	bodyStmt = append(bodyStmt,
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: countHandle,
				},
				&ast.Ident{
					Name: "err",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: opHandle,
						},
						Sel: &ast.Ident{
							Name: "Count",
						},
					},
					Args: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
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
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: mpsv,
									},
									Sel: &ast.Ident{
										Name: "log",
									},
								},
								Sel: &ast.Ident{
									Name: "Errorf",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"cat not get query total %v\"",
								},
								&ast.Ident{
									Name: "err",
								},
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: "nil",
							},
							&ast.BasicLit{
								Kind:  token.INT,
								Value: "0",
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		})

	// 设置偏移
	bodyStmt = append(bodyStmt,
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: opHandle,
					},
					Sel: &ast.Ident{
						Name: "Limit",
					},
				},
				Args: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "req",
							},
							Sel: &ast.Ident{
								Name: "GetLimit",
							},
						},
					},
				},
			},
		},
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: opHandle,
					},
					Sel: &ast.Ident{
						Name: "Offset",
					},
				},
				Args: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "req",
							},
							Sel: &ast.Ident{
								Name: "GetOffset",
							},
						},
					},
				},
			},
		})

	// 获取查询结果
	const restHandle = "ms"
	bodyStmt = append(bodyStmt,
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: restHandle,
				},
				&ast.Ident{
					Name: "err",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "query",
						},
						Sel: &ast.Ident{
							Name: "All",
						},
					},
					Args: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
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
								Name: "nil",
							},
							&ast.BasicLit{
								Kind:  token.INT,
								Value: "0",
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: restHandle,
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
								Name: "nil",
							},
							&ast.BasicLit{
								Kind:  token.INT,
								Value: "0",
							},
							&ast.Ident{
								Name: "nil",
							},
						},
					},
				},
			},
		},
	)
	// 将数据转化为dto
	const restDtoHandle = "rest"
	bodyStmt = append(bodyStmt,
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: restDtoHandle,
							},
						},
						Type: &ast.ArrayType{
							Elt: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent(consts.DtoPkgName),
									Sel: ast.NewIdent(utils.ToCamelCase(n.Name) + consts.DtoSuffix),
								},
							},
						},
					},
				},
			},
		},
		&ast.RangeStmt{
			Key: &ast.Ident{
				Name: "_",
			},
			Value: &ast.Ident{
				Name: "m",
			},
			Tok: token.DEFINE,
			X: &ast.Ident{
				Name: restHandle,
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: restDtoHandle,
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.Ident{
									Name: "append",
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: restDtoHandle,
									},
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent(consts.DtoPkgName),
											Sel: ast.NewIdent("New" + utils.ToCamelCase(n.Name) + consts.DtoSuffix),
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "m",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		})

	// 设置返回值
	bodyStmt = append(bodyStmt,
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.Ident{
					Name: restDtoHandle,
				},
				&ast.Ident{
					Name: countHandle,
				},
				&ast.Ident{
					Name: "nil",
				},
			},
		},
	)
	createBody := &ast.BlockStmt{
		List: bodyStmt,
	}
	return createMethod(mpsv, getImpStructName(n), consts.DefListFuncName, createParams, createResults, createBody)
}
