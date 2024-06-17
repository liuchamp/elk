package imp

import (
	"fmt"
	"github.com/masseelch/elk/pkg/utils/write"
	"go/ast"
	"go/token"
	"path"
	"path/filepath"
	"strings"

	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/internal/consts"
	"github.com/masseelch/elk/pkg/utils"
	"golang.org/x/tools/go/ast/astutil"
)

const impPkgName = consts.ImpPkgName

// ImpOuter  pr 前缀
func ImpOuter(g *gen.Graph, pr string) error {

	for _, n := range g.Nodes {
		if err := genImp(g, pr, n); err != nil {
			return err
		}
	}
	return nil
}

func genImp(g *gen.Graph, pr string, n *gen.Type) error {

	// 创建文件集
	fset := token.NewFileSet()

	// 创建包名
	file := &ast.File{
		Name: ast.NewIdent(impPkgName),
	}
	// 创建导入语句
	pkgNameEntroot := filepath.Base(n.Config.Package)
	astutil.AddNamedImport(fset, file, pkgNameEntroot, n.Config.Package)
	pkgNameEntBoPath := path.Join(n.Config.Package, consts.MainRepoPath, consts.BoPkgName)
	astutil.AddNamedImport(fset, file, consts.BoPkgName, pkgNameEntBoPath)
	pkgNameEntDTOPath := path.Join(n.Config.Package, consts.MainRepoPath, consts.DtoPkgName)
	astutil.AddNamedImport(fset, file, consts.DtoPkgName, pkgNameEntDTOPath)

	entityName := getEntityName(n)
	pkgNameEntEntityPath := path.Join(n.Config.Package, entityName)
	astutil.AddNamedImport(fset, file, entityName, pkgNameEntEntityPath)

	astutil.AddNamedImport(fset, file, "log", "github.com/go-kratos/kratos/v2/log")
	// imp 结构体
	file.Decls = append(file.Decls, createStruct(n))
	// 添加方法
	file.Decls = append(file.Decls, createMethods(n)...)
	// 添加构造函数
	file.Decls = append(file.Decls, createConstructor(n))
	// 打印生成的代码
	f := filepath.Join(pr, impPkgName, fmt.Sprintf("%s_imp.go", strings.ToLower(n.Name)))
	return write.WireGoFile(f, fset, file)
}

const innerImpSuffix = "RepoImp"

func getImpStructName(n *gen.Type) string {
	return utils.SnakeToCamel(n.Name) + innerImpSuffix
}

// createStruct 创建 RepoImp 结构体
func createStruct(n *gen.Type) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(getImpStructName(n)),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("c")},
								Type:  &ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("ent"), Sel: ast.NewIdent("Client")}},
							},
							{
								Names: []*ast.Ident{ast.NewIdent("log")},
								Type:  &ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("log"), Sel: ast.NewIdent("Helper")}},
							},
						},
					},
				},
			},
		},
	}
}

// createConstructor 创建构造函数
func createConstructor(n *gen.Type) *ast.FuncDecl {
	params := &ast.FieldList{
		List: []*ast.Field{
			{Names: []*ast.Ident{ast.NewIdent("c")}, Type: &ast.StarExpr{X: ast.NewIdent("ent.Client")}},
			{Names: []*ast.Ident{ast.NewIdent("logger")}, Type: ast.NewIdent("log.Logger")},
		},
	}
	results := &ast.FieldList{
		List: []*ast.Field{
			{Type: &ast.SelectorExpr{X: ast.NewIdent(consts.DefsvPkgName), Sel: ast.NewIdent(utils.ToCamelCase(n.Name) + consts.DefRepoSuffix)}},
		},
	}
	// 参数赋值过程
	var elts []ast.Expr
	//c
	elts = append(elts, &ast.KeyValueExpr{
		Key:   ast.NewIdent("c"),
		Value: ast.NewIdent("c"),
	})
	// log
	elts = append(elts, &ast.KeyValueExpr{
		Key: ast.NewIdent("log"),
		Value: &ast.CallExpr{
			Args: []ast.Expr{&ast.CallExpr{
				Fun: &ast.SelectorExpr{Sel: ast.NewIdent("With"), X: ast.NewIdent("log")},
				Args: []ast.Expr{
					ast.NewIdent("logger"),
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"module\"",
					},
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("\"imp/%s\"", utils.SnakeToCamel(n.Name)),
					},
				},
			}},
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("log"),
				Sel: ast.NewIdent("NewHelper"),
			},
		},
	})

	body := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: ast.NewIdent(getImpStructName(n)),
						Elts: elts,
					},
				}},
			},
		},
	}
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("New%s%s", utils.ToCamelCase(n.Name), consts.DefRepoSuffix)),
		Type: &ast.FuncType{
			Params:  params,
			Results: results,
		},
		Body: body,
	}
}

// createMethods 创建所有方法
func createMethods(n *gen.Type) []ast.Decl {
	var methods []ast.Decl

	methods = append(methods, saveImp(n))
	methods = append(methods, updateImp(n))
	methods = append(methods, patchImp(n))
	methods = append(methods, getImp(n))
	methods = append(methods, listImp(n))

	return methods
}

// createMethod 创建方法
func createMethod(receiverName, receiverType, methodName string, params, results *ast.FieldList, body *ast.BlockStmt) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(receiverName)},
					Type:  &ast.StarExpr{X: ast.NewIdent(receiverType)},
				},
			},
		},
		Name: ast.NewIdent(methodName),
		Type: &ast.FuncType{
			Params:  params,
			Results: results,
		},
		Body: body,
	}
}
