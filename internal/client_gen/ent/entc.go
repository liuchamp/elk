//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"path/filepath"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
)

func main() {
	ex, err := elk.NewExtension(
		elk.GenerateSpec("openapi.json"),
		elk.GenerateHandlers(),
	)
	if err != nil {
		log.Fatalf("creating elk extension: %v", err)
	}
	err = entc.Generate("./schema", &gen.Config{}, entc.Extensions(ex))
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}

	generateClient()
}

func generateClient() {
	swagger, err := util.LoadSwagger("./openapi.json")
	if err != nil {
		log.Fatalf("Failed to load swagger %v", err)
	}

	generated, err := codegen.Generate(swagger, codegen.Configuration{
		PackageName: "sub",
		OutputOptions: codegen.OutputOptions{
			SkipFmt:   false,
			SkipPrune: false,
			// AliasTypes: true,
		},
		Generate: codegen.GenerateOptions{ChiServer: false,
			EchoServer:   false,
			Client:       true,
			Models:       true,
			EmbeddedSpec: false,
		},
	})
	if err != nil {
		log.Fatalf("Generating client failed%v", err)
	}

	dir := filepath.Join(".", "stub")
	stub := filepath.Join(".", "stub", "http.go")
	perm := os.FileMode(0777)
	err = os.MkdirAll(dir, perm)

	if err != nil {
		log.Fatalf("error creating dir: %s", err)
	}

	err = os.WriteFile(stub, []byte(generated), perm)
	if err != nil {
		panic(err)
		log.Fatalf("error writing generated code to file: %s", err)
	}
}
