package bo

import (
	"entgo.io/ent/entc/gen"
	Tyfield2 "entgo.io/ent/schema/field"
	"fmt"
	"github.com/fatih/structtag"
	"github.com/masseelch/elk/internal/consts"
	"github.com/masseelch/elk/pkg/utils"
	"github.com/masseelch/elk/pkg/utils/write"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"path"
	"path/filepath"
	"strings"
)

// 生成request param 数据

const boPkgName = consts.BoPkgName

const (
	CreateSuffix = consts.CreateBoSuffix
	UpdateSuffix = consts.UpdateBoSuffix
	PatchSuffix  = consts.PatchBoSuffix
	GetSuffix    = consts.GetBoSuffix
	ListSuffix   = consts.ListBoSuffix
)

// BoOuter  pr 前缀
func BoOuter(g *gen.Graph, pr string) error {

	for _, n := range g.Nodes {
		if err := boGen(g, pr, n); err != nil {
			return err
		}
	}
	return nil
}

func boGen(g *gen.Graph, pr string, n *gen.Type) error {

	// 创建文件集
	fset := token.NewFileSet()

	// 创建包名
	file := &ast.File{
		Name: ast.NewIdent(boPkgName),
	}
	// 创建导入语句
	pkgNameEntroot := filepath.Base(n.Config.Package)
	astutil.AddNamedImport(fset, file, pkgNameEntroot, n.Config.Package)

	// 创建结构体声明
	var createFs []*ast.Field
	var updateFs []*ast.Field
	var patchFs []*ast.Field
	var listFs []*ast.Field
	var getFs []*ast.Field

	// id 处理， patch, update， 需要填写
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
	updateFs = append(updateFs, &idfInfo)
	patchFs = append(patchFs, &idfInfo)

	gtFs := queryFieldDuty(n.ID, n, fset, file)
	listFs = append(listFs, gtFs...)
	getFs = append(getFs, gtFs...)

	for _, field := range n.Fields {
		notPatchFs := doFieldDuty(field, n, fset, file, false)
		createFs = append(createFs, notPatchFs)
		updateFs = append(updateFs, notPatchFs)
		patchFs = append(patchFs, doFieldDuty(field, n, fset, file, true))

		gets := queryFieldDuty(field, n, fset, file)
		listFs = append(listFs, gets...)
		getFs = append(getFs, gets...)
	}

	// 页查询额外参数
	listFs = append(listFs, QueryParam()...)

	file.Decls = append(file.Decls, boStruct(createFs, CreateSuffix, n.Name))
	file.Decls = append(file.Decls, boStruct(updateFs, UpdateSuffix, n.Name))
	file.Decls = append(file.Decls, boStruct(patchFs, PatchSuffix, n.Name))
	file.Decls = append(file.Decls, boStruct(listFs, ListSuffix, n.Name))
	file.Decls = append(file.Decls, boStruct(getFs, GetSuffix, n.Name))

	// 打印生成的代码
	f := filepath.Join(pr, boPkgName, fmt.Sprintf("%s_bo.go", strings.ToLower(n.Name)))
	return write.WireGoFile(f, fset, file)
}

func boStruct(fs []*ast.Field, suffix string, name string) *ast.GenDecl {
	updateType := &ast.StructType{
		Fields: &ast.FieldList{
			List: fs,
		},
	}
	updateDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(fmt.Sprintf("%s"+suffix, name)),
				Type: updateType,
			},
		},
	}
	return updateDecl
}
func doFieldDuty(field *gen.Field, n *gen.Type, fset *token.FileSet, file *ast.File, isPatch bool) *ast.Field {
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
	if isPatch {
		f.Type = &ast.StarExpr{
			X: ast.NewIdent(field.Type.String()),
		}
	}
	f.Tag = &ast.BasicLit{
		Kind:  token.STRING,
		Value: tag(field),
	}
	return f
}

func queryFieldDuty(field *gen.Field, n *gen.Type, fset *token.FileSet, file *ast.File) []*ast.Field {
	fn := utils.ToCamelCase(field.Name)
	f := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(fn)},
		Type: &ast.StarExpr{
			X: ast.NewIdent(field.Type.String()),
		},
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
	f.Tag = &ast.BasicLit{
		Kind:  token.STRING,
		Value: tagQuery(field),
	}
	return []*ast.Field{f}
}

func tag(n *gen.Field) string {

	tgs := structtag.Tags{}
	_ = tgs.Set(&structtag.Tag{
		Key:     "json",
		Name:    n.Name,
		Options: []string{"omitempty"},
	})

	return tag2String(tgs)
}

func tagQuery(n *gen.Field) string {
	tgs := structtag.Tags{}
	_ = tgs.Set(&structtag.Tag{
		Key:     "json",
		Name:    n.Name,
		Options: []string{"omitempty"},
	})
	if n.Annotations == nil {
		return fmt.Sprintf("`%s`", tgs.String())
	}

	//@todo 参数验证规则
	//for s, a := range n.Annotations {
	//	if s=="validate" {
	//
	//		_ = tgs.Set(&structtag.Tag{
	//			Key:     "validate",
	//			Name:    ,
	//			Options: []string{"omitempty"},
	//		})
	//	}
	//}

	return tag2String(tgs)
}
func tag2String(tgs structtag.Tags) string {

	return fmt.Sprintf("`%s`", tgs.String())
}

func QueryParam() []*ast.Field {

	pageSizeTags := structtag.Tags{}
	_ = pageSizeTags.Set(&structtag.Tag{
		Key:     "json",
		Name:    "page_size",
		Options: []string{"omitempty"},
	})
	pageTags := structtag.Tags{}
	_ = pageTags.Set(&structtag.Tag{
		Key:     "json",
		Name:    "page",
		Options: []string{"omitempty"},
	})
	pageSizeParam := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent("PageSize")},
		Type:  ast.NewIdent("int"),
		Tag: &ast.BasicLit{
			Kind:  token.STRING,
			Value: tag2String(pageSizeTags),
		},
	}
	pageParam := &ast.Field{
		Names: []*ast.Ident{ast.NewIdent("Page")},
		Type:  ast.NewIdent("int"),
		Tag: &ast.BasicLit{
			Kind:  token.STRING,
			Value: tag2String(pageTags),
		},
	}

	return []*ast.Field{pageParam, pageSizeParam}
}