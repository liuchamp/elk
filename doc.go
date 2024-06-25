package elk

import (
	"encoding/json"
	"entgo.io/ent/entc/gen"
	"errors"
	"io"
	"os"
)
import "github.com/getkin/kin-openapi/openapi3"

type DocConfig struct {
	Title       string
	Description string
	Version     string
}

// GenerateSpec enables the OpenAPI-Spec generator. Data will be written to given filename.
func GenerateSpec(out string, hooks ...Hook) ExtensionOption {
	return func(ex *Extension) error {
		if out == "" {
			return errors.New("spec filename cannot be empty")
		}
		ex.hooks = append(ex.hooks, ex.SpecGenerator(out))
		ex.specHooks = append(ex.specHooks, hooks...)
		return nil
	}
}

// SpecGenerator 生成数据
func (e *Extension) SpecGenerator(out string) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			// Let ent create all the files.
			if err := next.Generate(g); err != nil {
				return err
			}
			// Start the Generator chain.
			var chain Generator = generate(g)
			// Add user hooks to chain.
			for i := len(e.specHooks) - 1; i >= 0; i-- {
				chain = e.specHooks[i](chain)
			}
			// Create a fresh spec.
			s := initSpec()
			// Run the generators.
			if err := chain.Generate(s); err != nil {
				return err
			}
			// Dump the spec.
			b, err := json.MarshalIndent(s, "", "  ")
			if err != nil {
				return err
			}
			return os.WriteFile(out, b, 0664)
		})
	}
}

// initSpec returns an empty spec ready to receive data.
func initSpec() *openapi3.T {
	return &openapi3.T{
		Extensions: nil,
		OpenAPI:    "",
		Components: &openapi3.Components{

			Schemas:       openapi3.Schemas{},
			Parameters:    openapi3.ParametersMap{},
			Headers:       openapi3.Headers{},
			RequestBodies: openapi3.RequestBodies{},
			Responses:     openapi3.ResponseBodies{},
		},
		Info: &openapi3.Info{
			Title:          "Ent Schema API",
			Description:    "This is an auto generated API description made out of an Ent schema definition",
			Version:        "0.0.0",
			Extensions:     nil,
			TermsOfService: "",
			Contact:        nil,
			License:        nil,
		},
		Paths:        nil,
		Security:     nil,
		Servers:      nil,
		Tags:         nil,
		ExternalDocs: nil,
	}
}

// generate is the default Generator to fill a given spec.
func generate(g *gen.Graph) GenerateFunc {
	return func(t *openapi3.T) error {
		// Add all views to the schemas.
		if err := viewSchemas(g, t); err != nil {
			return err
		}
		// Add all error responses.
		//errResponses(s)
		// Create the paths.
		if err := paths(g, t); err != nil {
			return err
		}
		return nil
	}
}

// SpecTitle sets the title of the Info block.
func SpecTitle(v string) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *openapi3.T) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Info.Title = v
			return nil
		})
	}
}

// SpecDescription sets the title of the Info block.
func SpecDescription(v string) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *openapi3.T) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Info.Description = v
			return nil
		})
	}
}

// SpecVersion sets the version of the Info block.
func SpecVersion(v string) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *openapi3.T) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Info.Version = v
			return nil
		})
	}
}

// TODO: Rest of Info block ...

// SpecSecuritySchemes sets the security schemes of the Components block.
func SpecSecuritySchemes(schemes openapi3.SecuritySchemes) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *openapi3.T) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Components.SecuritySchemes = schemes
			return nil
		})
	}
}

// SpecSecurity sets the global security Spec.
func SpecSecurity(sec openapi3.SecurityRequirements) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *openapi3.T) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Security = sec
			return nil
		})
	}
}

// SpecDump dumps the current specs content to the given io.Writer.
func SpecDump(out io.Writer) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *openapi3.T) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			j, err := json.MarshalIndent(spec, "", "  ")
			if err != nil {
				return err
			}
			_, err = out.Write(j)
			return err
		})
	}
}
