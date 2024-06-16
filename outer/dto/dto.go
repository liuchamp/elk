package dto

import (
	"fmt"
	"github.com/masseelch/elk/pkg/utils/write"
	"go/ast"
	"go/token"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"entgo.io/ent/entc/gen"
	Tyfield2 "entgo.io/ent/schema/field"
	"github.com/fatih/structtag"
	"github.com/masseelch/elk/pkg/utils"
)

// DtoOuter  pr 前缀
func DtoOuter(g *gen.Graph, pr string) error {

	for _, n := range g.Nodes {
		if err := genDto(g, pr, n); err != nil {
			return err
		}
	}
	return nil
}

const po = "po"
const dtoPkgName = "dto"

// genDto 生成标准dto
// 全部的dto 都生成到dto 目录， 正常输出路径 n.Config.Package/../biz/repo
func genDto(g *gen.Graph, pr string, n *gen.Type) error {
	// 创建文件集
	fset := token.NewFileSet()

	// 创建包名
	file := &ast.File{
		Name: ast.NewIdent(dtoPkgName),
	}
	// 创建导入语句
	pkgNameEntroot := filepath.Base(n.Config.Package)
	astutil.AddNamedImport(fset, file, pkgNameEntroot, n.Config.Package)

	// 创建结构体声明
	var dtoFields []*ast.Field
	var elts []ast.Expr

	idfInfo := ast.Field{
		Names: []*ast.Ident{ast.NewIdent(strings.ToUpper(n.ID.Name))},
		Type:  ast.NewIdent(n.ID.Type.String()),
		Tag: &ast.BasicLit{
			Kind:  token.STRING,
			Value: tag(n.ID),
		},
	}
	if n.ID.Type.PkgPath != "" {
		if n.ID.Type.PkgName != "" {
			astutil.AddNamedImport(fset, file, n.ID.Type.PkgName, n.ID.Type.PkgPath)
		} else {
			astutil.AddImport(fset, file, n.ID.Type.PkgPath)
		}
	}
	elts = append(elts, &ast.KeyValueExpr{
		Key:   ast.NewIdent(idfInfo.Names[0].Name),
		Value: ast.NewIdent(po + "." + idfInfo.Names[0].Name),
	})
	dtoFields = append(dtoFields, &idfInfo)

	for _, field := range n.Fields {
		fn, f := funcFieldDuty(field, n, fset, file)
		elts = append(elts, &ast.KeyValueExpr{
			Key:   ast.NewIdent(fn),
			Value: ast.NewIdent(po + "." + fn),
		})
		dtoFields = append(dtoFields, f)
	}
	//{{- range $e := $v.Edges }}
	//{{ $e.StructField }} {{ if $e.Unique }}*{{ end }}{{ $e.Name }}{{ if not $e.Unique }}s{{ end }}
	//{{- with tagLookup $e.StructTag "json" }} `json:"{{ . }}"{{ end }}`
	//{{- end }}
	//for _, edge := range n.Edges {
	//
	//}

	dtoStructType := &ast.StructType{
		Fields: &ast.FieldList{
			List: dtoFields,
		},
	}
	structDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(fmt.Sprintf("%sDTO", n.Name)),
				Type: dtoStructType,
			},
		},
	}
	file.Decls = append(file.Decls, structDecl)

	// 创建新po to dto

	// 创建参数列表
	params := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent(po)},
			Type:  &ast.StarExpr{X: ast.NewIdent("ent." + n.Name)},
		},
	}
	// 创建函数体
	body := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X: &ast.CompositeLit{
							Type: ast.NewIdent(fmt.Sprintf("%sDTO", n.Name)),
							Elts: elts,
						},
					},
				},
			},
		},
	}

	// 创建函数声明
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("New%sDTO", n.Name)),
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: params},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{X: ast.NewIdent(fmt.Sprintf("%sDTO", n.Name))},
					},
				},
			},
		},
		Body: body,
	}

	file.Decls = append(file.Decls, funcDecl)
	// 打印生成的代码
	f := filepath.Join(pr, dtoPkgName, fmt.Sprintf("%s_dto.go", strings.ToLower(n.Name)))
	return write.WireGoFile(f, fset, file)
}

func funcFieldDuty(field *gen.Field, n *gen.Type, fset *token.FileSet, file *ast.File) (string, *ast.Field) {
	fn := utils.ToCamelCase(field.Name)
	f := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(fn)},
		Type:  ast.NewIdent(field.Type.String()),
	}
	if field.Type.Type == Tyfield2.TypeEnum { // 屏蔽路径枚举
		pkgName := strings.SplitN(field.Type.String(), ".", 2)[0]
		pkgPath := path.Join(n.Config.Package, pkgName)
		astutil.AddNamedImport(fset, file, pkgName, pkgPath)
	} else if n.ID.Type.PkgPath != "" {
		if n.ID.Type.PkgName != "" {
			astutil.AddNamedImport(fset, file, n.ID.Type.PkgName, n.ID.Type.PkgPath)
		} else {
			astutil.AddImport(fset, file, n.ID.Type.PkgPath)
		}
	}
	if field.NillableValue() {
		f.Type = &ast.StarExpr{
			X: ast.NewIdent(field.Type.String()),
		}
	}
	f.Tag = &ast.BasicLit{
		Kind:  token.STRING,
		Value: tag(field),
	}
	return fn, f
}

func tag(n *gen.Field) string {

	tgs := structtag.Tags{}
	_ = tgs.Set(&structtag.Tag{
		Key:     "json",
		Name:    n.Name,
		Options: []string{"omitempty"},
	})

	return fmt.Sprintf("`%s`", tgs.String())
}

func GetDtoPkg(pr string) string {
	return path.Join(pr, dtoPkgName)
}
