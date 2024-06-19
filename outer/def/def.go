package def

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/masseelch/elk/internal/consts"
	"github.com/masseelch/elk/pkg/utils/write"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"path/filepath"
	"strings"
)

const defPkgName = consts.DefsvPkgName

// DefOuter  pr 前缀
func DefOuter(g *gen.Graph, pr string) error {
	for _, n := range g.Nodes {
		if err := DefGen(g, pr, n); err != nil {
			return err
		}
	}
	return nil
}

func DefGen(g *gen.Graph, pr string, n *gen.Type) error {

	// 创建文件集
	fset := token.NewFileSet()

	// 创建包名
	file := &ast.File{
		Name: ast.NewIdent(defPkgName),
	}
	// 创建导入语句
	pkgNameEntroot := filepath.Base(n.Config.Package)
	astutil.AddNamedImport(fset, file, pkgNameEntroot, n.Config.Package)
	astutil.AddNamedImport(fset, file, consts.BoPkgName, consts.GetBoPackageName(g.Package))
	astutil.AddNamedImport(fset, file, consts.DtoPkgName, consts.GetDtoPackageName(g.Package))

	var methods []*ast.Field
	methods = append(methods, SaveDef(n))
	methods = append(methods, UpdateDef(n))
	methods = append(methods, PatchDef(n))
	methods = append(methods, GetDef(n))
	methods = append(methods, ListDef(n))

	// 创建接口声明
	interfaceType := &ast.InterfaceType{
		Methods: &ast.FieldList{List: methods},
	}
	interfaceDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(fmt.Sprintf("%s%s", n.Name, consts.DefRepoSuffix)),
				Type: interfaceType,
				Comment: &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: fmt.Sprintf("// %s%s is a define for opt data", n.Name, consts.DefRepoSuffix)},
					},
				},
			},
		},
	}
	file.Decls = append(file.Decls, interfaceDecl)
	// 打印生成的代码
	f := filepath.Join(pr, defPkgName, fmt.Sprintf("%s_def.go", strings.ToLower(n.Name)))
	return write.WireGoFile(f, fset, file)
}

const (
	CreateFuncName = consts.DefCreateFuncName
	UpdateFuncName = consts.DefUpdateFuncName
	PatchFuncName  = consts.DefPatchFuncName
	GetFuncName    = consts.DefGetFuncName
	ListFuncName   = consts.DefListFuncName
)

func SaveDef(n *gen.Type) *ast.Field {
	return FuncDef(n, CreateFuncName, consts.CreateBoSuffix)
}

func UpdateDef(n *gen.Type) *ast.Field {
	return FuncDef(n, UpdateFuncName, consts.UpdateBoSuffix)
}
func PatchDef(n *gen.Type) *ast.Field {
	return FuncDef(n, PatchFuncName, consts.PatchBoSuffix)
}
func GetDef(n *gen.Type) *ast.Field {
	return FuncDef(n, GetFuncName, consts.GetBoSuffix)
}
func ListDef(n *gen.Type) *ast.Field {
	fnName := ListFuncName
	BoSuffix := consts.ListBoSuffix
	method := ast.Field{
		Names: []*ast.Ident{ast.NewIdent(fnName)},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{Names: []*ast.Ident{ast.NewIdent("ctx")}, Type: ast.NewIdent("context.Context")},
					{Names: []*ast.Ident{ast.NewIdent("req")}, Type: &ast.StarExpr{X: ast.NewIdent(consts.BoPkgName + "." + fmt.Sprintf("%s"+BoSuffix, n.Name))}},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.ArrayType{Elt: &ast.StarExpr{X: ast.NewIdent(consts.DtoPkgName + "." + n.Name + consts.DtoSuffix)}}},
					//{Type: &ast.StarExpr{X: ast.NewIdent(consts.DtoPkgName + "." + n.Name + consts.DtoSuffix)}},
					{Type: ast.NewIdent("int")}, // total
					{Type: ast.NewIdent("error")},
				},
			},
		},
	}
	return &method
}

func FuncDef(n *gen.Type, fnName string, BoSuffix string) *ast.Field {
	method := ast.Field{
		Names: []*ast.Ident{ast.NewIdent(fnName)},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{Names: []*ast.Ident{ast.NewIdent("ctx")}, Type: ast.NewIdent("context.Context")},
					{Names: []*ast.Ident{ast.NewIdent("req")}, Type: &ast.StarExpr{X: ast.NewIdent(consts.BoPkgName + "." + fmt.Sprintf("%s"+BoSuffix, n.Name))}},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.StarExpr{X: ast.NewIdent(consts.DtoPkgName + "." + n.Name + consts.DtoSuffix)}},
					{Type: ast.NewIdent("error")},
				},
			},
		},
	}
	return &method
}
