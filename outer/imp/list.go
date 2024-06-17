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

	//const opHandle = "query"
	var bodyStmt []ast.Stmt

	// 生成opHandle
	bodyStmt = append(bodyStmt, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.Ident{
				Name: "panic",
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"unimplemented\"",
				},
			},
		},
	})
	//bodyStmt = append(bodyStmt, &ast.AssignStmt{
	//	Lhs: []ast.Expr{
	//		ast.NewIdent(opHandle),
	//	},
	//	Tok: token.DEFINE,
	//	Rhs: []ast.Expr{
	//		&ast.CallExpr{
	//			Fun: &ast.SelectorExpr{
	//				X:   &ast.SelectorExpr{X: ast.NewIdent(mpsv + ".c"), Sel: ast.NewIdent(utils.ToCamelCase(n.Name))},
	//				Sel: ast.NewIdent("Query"),
	//			},
	//		},
	//	},
	//})

	//bodyStmt = append(bodyStmt, fieldQueryStmt(n.Name, n.ID))
	//// 设置数据
	//for _, field := range n.Fields {
	//	bodyStmt = append(bodyStmt, fieldQueryStmt(n.Name, field))
	//}

	createBody := &ast.BlockStmt{
		List: bodyStmt,
	}
	return createMethod(mpsv, getImpStructName(n), consts.DefListFuncName, createParams, createResults, createBody)
}
