package elk

import (
	"entgo.io/ent/entc/gen"
)
import "github.com/getkin/kin-openapi/openapi3"

type DocConfig struct {
	Title       string
	Description string
	Version     string
}

func GenerateSpec(opts ...HandlerOption) ExtensionOption {
	return func(ex *Extension) error {
		for _, opt := range opts {
			if err := opt(ex); err != nil {
				return err
			}
		}
		ex.hooks = append(ex.hooks, DocGenerator(ex.docCfg))
		return nil
	}
}

func DocGenerator(c DocConfig) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			// 收集model 信息，
			// dto

			// bo

			// paths
			return nil
		})
	}
}

func getTitle(c DocConfig) *openapi3.T {
	openapi := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:       c.Title,
			Description: c.Description,
			Version:     c.Version,
		},
	}

	return openapi
}
