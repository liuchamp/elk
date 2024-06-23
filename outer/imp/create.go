package imp

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/liuchamp/elk/internal/consts"
	"github.com/liuchamp/elk/pkg/utils"
	"go/ast"
	"go/token"
)

func saveImp(n *gen.Type) *ast.FuncDecl {
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
						Sel: ast.NewIdent(utils.ToCamelCase(n.Name) + consts.CreateBoSuffix)},
				},
			},
		},
	}
	createResults := &ast.FieldList{
		List: []*ast.Field{
			{Type: &ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent(consts.DtoPkgName), Sel: ast.NewIdent(utils.ToCamelCase(n.Name) + consts.DtoSuffix)}}},
			{Type: ast.NewIdent("error")},
		},
	}

	const opHandle = "createOP"
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
					Sel: ast.NewIdent("Create"),
				},
			},
		},
	})

	// 设置数据
	for _, field := range n.Fields {
		bodyStmt = append(bodyStmt, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent(opHandle),
					Sel: ast.NewIdent(fmt.Sprintf("Set%s", utils.ToCamelCase(field.Name))),
				},
				Args: []ast.Expr{&ast.SelectorExpr{
					X:   ast.NewIdent("req"),
					Sel: ast.NewIdent(utils.ToCamelCase(field.Name)),
				}},
			},
		})
	}
	// 生成操作结果
	bodyStmt = append(bodyStmt, &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("m"), ast.NewIdent("err")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent(opHandle),
					Sel: ast.NewIdent("Save"),
				},
				Args: []ast.Expr{ast.NewIdent("ctx")},
			},
		},
	})
	// 设置返回
	bodyStmt = append(bodyStmt,
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("err"),
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("nil"),
							ast.NewIdent("err"),
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("m"),
				Op: token.EQL,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("nil"),
							ast.NewIdent("nil"),
						},
					},
				},
			},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent(consts.DtoPkgName),
						Sel: ast.NewIdent("New" + utils.ToCamelCase(n.Name) + consts.DtoSuffix),
					},
					Args: []ast.Expr{
						ast.NewIdent("m"),
					},
				},
				ast.NewIdent("nil"),
			},
		})
	createBody := &ast.BlockStmt{
		List: bodyStmt,
	}

	return createMethod(mpsv, getImpStructName(n), consts.DefCreateFuncName, createParams, createResults, createBody)

}
