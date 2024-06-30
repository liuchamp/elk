package entool

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/stoewer/go-strcase"
	"strings"
)

var ignoreFields = []string{
	"created_at",
	"updated_at",
}

func checkIgnoreSetField(n string) bool {
	for _, field := range ignoreFields {
		if field == n {
			return true
		}
	}
	return false
}

var psFieldSuffixs = []string{"_id", "_ip"}

func SetNameGen(field *gen.Field) string {
	if checkIgnoreSetField(field.Name) {
		return ""
	}
	fn := strcase.UpperCamelCase(field.Name)
	for _, suffix := range psFieldSuffixs {
		if strings.HasSuffix(field.Name, suffix) {
			fn = fmt.Sprintf("%s%s", strcase.UpperCamelCase(strings.TrimSuffix(field.Name, suffix)), strings.ToUpper(strings.TrimPrefix(suffix, "_")))
			break
		}
	}
	return fn
}
