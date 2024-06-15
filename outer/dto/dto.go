package dto

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
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

// genDto 生成标准dto
// 全部的dto 都生成到dto 目录， 正常输出路径 n.Config.Package/../biz/repo
func genDto(g *gen.Graph, pr string, n *gen.Type) error {
	// 创建文件集
	fset := token.NewFileSet()

	// 创建包名
	file := &ast.File{
		Name: ast.NewIdent("dto"),
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
	f, err := utils.CreateFileWithDirs(filepath.Join(pr, "dto", fmt.Sprintf("%s_dto.go", strings.ToLower(n.Name))))
	if err != nil {
		return err
	}
	defer f.Close()

	// 创建一个缓冲区来存储生成的代码
	buf := new(bytes.Buffer)

	// 打印未格式化的代码到缓冲区
	if err := printer.Fprint(buf, fset, file); err != nil {
		return err
	}

	// 获取未格式化的代码
	unformattedCode := buf.Bytes()

	// 格式化代码
	formattedCode, err := format.Source(unformattedCode)
	if err != nil {
		return err
	}

	// 将格式化的代码写入文件
	if _, err := f.Write(formattedCode); err != nil {
		return err
	}

	return nil
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
