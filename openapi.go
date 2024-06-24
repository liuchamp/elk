package elk

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/liuchamp/elk/spec"
	"github.com/stoewer/go-strcase"
	"net/http"
	"strconv"
	"strings"
)

type (
	// Generator is the interface that wraps the Generate method.
	Generator interface {
		// Generate edits the given OpenAPI spec.
		Generate(t *openapi3.T) error
	}
	// The GenerateFunc type is an adapter to allow the use of ordinary
	// function as Generator. If f is a function with the appropriate signature,
	// GenerateFunc(f) is a Generator that calls f.
	GenerateFunc func(*openapi3.T) error
	// Hook defines the "spec generate middleware".
	Hook func(Generator) Generator
)

// Generate calls f(s).
func (f GenerateFunc) Generate(s *openapi3.T) error {
	return f(s)
}

// viewSchemas adds all views to the specs schemas.
func viewSchemas(g *gen.Graph, s *openapi3.T) error {
	vs, err := newViews(g)
	if err != nil {
		return err
	}
	// Create a schema for every view.
	for n, v := range vs {
		fs := make(spec.Fields, len(v.Fields))
		// We can already add the schema fields.
		for _, f := range v.Fields {
			sf, err := newField(f)
			if err != nil {
				return err
			}
			fs[f.Name] = sf
		}
		s.Components.Schemas[n] = &spec.Schema{
			Name:   n,
			Fields: fs,
		}
	}
	// Loop over the views again and this time fill the edges.
	for n, v := range vs {
		es := make(spec.Edges, len(v.Edges))
		for _, e := range v.Edges {
			es[e.Edge.Name] = spec.Edge{
				Ref:    s.Components.Schemas[e.Name],
				Unique: e.Unique,
			}
		}
		s.Components.Schemas[n].Edges = es
	}
	return nil
}

// newField constructs a spec.Field out of a gen.Field.
func newField(f *gen.Field) (*spec.Field, error) {
	t, err := oasType(f)
	if err != nil {
		return nil, err
	}
	e, err := exampleValue(f)
	if err != nil {
		return nil, err
	}
	return &spec.Field{
		Unique:   true,
		Required: !f.Optional,
		Type:     *t,
		Example:  e,
	}, nil
}

// paths adds all views to the specs schemas.
func paths(g *gen.Graph, s *spec.Spec) error {
	for _, n := range g.Nodes {
		// Add schema operations.
		ops, err := nodeOperations(n)
		if err != nil {
			return err
		}
		// root for all operations on this node.
		root := "/" + rules.Pluralize(strcase.KebabCase(n.Name))
		// Create operation.
		if contains(ops, opCreate) {
			path(s, root).Post, err = createOp(s, n)
			if err != nil {
				return err
			}

		}
		// Read operation.
		if contains(ops, opRead) {
			path(s, root+"/{id}").Get, err = readOp(s, n)
			if err != nil {
				return err
			}
		}
		// Update operation.
		if contains(ops, opUpdate) {
			path(s, root+"/{id}").Patch, err = updateOp(s, n)
			if err != nil {
				return err
			}
		}
		// Delete operation.
		if contains(ops, opDelete) {
			path(s, root+"/{id}").Delete, err = deleteOp(s, n)
			if err != nil {
				return err
			}
		}
		// List operation.
		if contains(ops, opList) {
			path(s, root).Get, err = listOp(s, n)
			if err != nil {
				return err
			}
		}
		// Sub-Resource operations.
		es, err := filterEdges(n)
		if err != nil {
			return err
		}
		for _, e := range es {
			p := path(s, root+"/{id}/"+strcase.KebabCase(e.Name))
			if e.Unique {
				p.Get, err = readEdgeOp(s, n, e)
				if err != nil {
					return err
				}
			} else {
				p.Get, err = listEdgeOp(s, n, e)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// createOp returns the spec description for a create operation on the given node.
func createOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
	ant, err := schemaAnnotation(n)
	if err != nil {
		return nil, err
	}
	req, err := requestBody(n, opCreate)
	if err != nil {
		return nil, err
	}
	v, err := newView(n, ant.CreateGroups)
	if err != nil {
		return nil, err
	}
	rspName, err := v.Name()
	if err != nil {
		return nil, err
	}
	return &spec.Operation{
		Summary:     fmt.Sprintf("Create a new %s", n.Name),
		Description: fmt.Sprintf("Creates a new %s and persists it to storage.", n.Name),
		Tags:        []string{n.Name},
		OperationID: opCreate + n.Name,
		RequestBody: req,
		Responses: map[string]*spec.OperationResponse{
			strconv.Itoa(http.StatusOK): {
				Response: spec.Response{
					Description: fmt.Sprintf("%s created", n.Name),
					Headers:     nil, // TODO
					Content: &spec.Content{
						spec.JSON: spec.MediaTypeObject{
							Unique: true,
							Ref:    s.Components.Schemas[rspName],
						},
					},
				},
			},
			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
		},
		Security: ant.CreateSecurity,
	}, nil
}

// readOp returns a spec.Operation for a read operation on the given node.
func readOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
	ant, err := schemaAnnotation(n)
	if err != nil {
		return nil, err
	}
	v, err := newView(n, ant.ReadGroups)
	if err != nil {
		return nil, err
	}
	rspName, err := v.Name()
	if err != nil {
		return nil, err
	}
	t, err := oasType(n.ID)
	if err != nil {
		return nil, err
	}
	return &spec.Operation{
		Summary:     fmt.Sprintf("Find a %s by ID", n.Name),
		Description: fmt.Sprintf("Finds the %s with the requested ID and returns it.", n.Name),
		Tags:        []string{n.Name},
		OperationID: opRead + n.Name,
		Parameters: []*spec.Parameter{{
			Name:        "id",
			In:          spec.InPath,
			Description: fmt.Sprintf("ID of the %s", n.Name),
			Required:    true,
			Schema:      *t,
		}},
		Responses: map[string]*spec.OperationResponse{
			strconv.Itoa(http.StatusOK): {
				Response: spec.Response{
					Description: fmt.Sprintf("%s with requested ID was found", n.Name),
					Headers:     nil, // TODO
					Content: &spec.Content{
						spec.JSON: spec.MediaTypeObject{
							Unique: true,
							Ref:    s.Components.Schemas[rspName],
						},
					},
				},
			},
			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
		},
		Security: ant.ReadSecurity,
	}, nil
}

func updateOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
	ant, err := schemaAnnotation(n)
	if err != nil {
		return nil, err
	}
	req, err := requestBody(n, opUpdate)
	if err != nil {
		return nil, err
	}
	v, err := newView(n, ant.UpdateGroups)
	if err != nil {
		return nil, err
	}
	rspName, err := v.Name()
	if err != nil {
		return nil, err
	}
	t, err := oasType(n.ID)
	if err != nil {
		return nil, err
	}
	return &spec.Operation{
		Summary:     fmt.Sprintf("Updates a %s", n.Name),
		Description: fmt.Sprintf("Updates a %s and persists changes to storage.", n.Name),
		Tags:        []string{n.Name},
		OperationID: opUpdate + n.Name,
		Parameters: []*spec.Parameter{{
			Name:        "id",
			In:          spec.InPath,
			Description: fmt.Sprintf("ID of the %s to update", n.Name),
			Required:    true,
			Schema:      *t,
		}},
		RequestBody: req,
		Responses: map[string]*spec.OperationResponse{
			strconv.Itoa(http.StatusOK): {
				Response: spec.Response{
					Description: fmt.Sprintf("%s updated", n.Name),
					Headers:     nil, // TODO
					Content: &spec.Content{
						spec.JSON: spec.MediaTypeObject{
							Unique: true,
							Ref:    s.Components.Schemas[rspName],
						},
					},
				},
			},
			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
		},
		Security: ant.UpdateSecurity,
	}, nil
}

func deleteOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
	ant, err := schemaAnnotation(n)
	if err != nil {
		return nil, err
	}
	t, err := oasType(n.ID)
	if err != nil {
		return nil, err
	}
	return &spec.Operation{
		Summary:     fmt.Sprintf("Deletes a %s by ID", n.Name),
		Description: fmt.Sprintf("Deletes the %s with the requested ID.", n.Name),
		Tags:        []string{n.Name},
		OperationID: opDelete + n.Name,
		Parameters: []*spec.Parameter{{
			Name:        "id",
			In:          spec.InPath,
			Description: fmt.Sprintf("ID of the %s to delete", n.Name),
			Required:    true,
			Schema:      *t,
		}},
		Responses: map[string]*spec.OperationResponse{
			strconv.Itoa(http.StatusNoContent): {
				Response: spec.Response{
					Description: fmt.Sprintf("%s with requested ID was deleted", n.Name),
					Headers:     nil, // TODO
				},
			},
			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
		},
		Security: ant.DeleteSecurity,
	}, nil
}

func listOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
	ant, err := schemaAnnotation(n)
	if err != nil {
		return nil, err
	}
	v, err := newView(n, ant.ListGroups)
	if err != nil {
		return nil, err
	}
	rspName, err := v.Name()
	if err != nil {
		return nil, err
	}
	return &spec.Operation{
		Summary:     fmt.Sprintf("List %s", rules.Pluralize(n.Name)),
		Description: fmt.Sprintf("List %s.", rules.Pluralize(n.Name)),
		Tags:        []string{n.Name},
		OperationID: opList + n.Name,
		Parameters: []*spec.Parameter{{
			Name:        "page",
			In:          spec.InQuery,
			Description: "what page to render",
			Schema:      *_int32,
		}, {
			Name:        "itemsPerPage",
			In:          spec.InQuery,
			Description: "item count to render per page",
			Schema:      *_int32,
		}},
		Responses: map[string]*spec.OperationResponse{
			strconv.Itoa(http.StatusOK): {
				Response: spec.Response{
					Description: fmt.Sprintf("result %s list", n.Name),
					Headers:     nil, // TODO
					Content: &spec.Content{
						spec.JSON: spec.MediaTypeObject{
							Ref: s.Components.Schemas[rspName],
						},
					},
				},
			},
			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
		},
		Security: ant.ListSecurity,
	}, nil
}

func readEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
	op, err := readOp(s, e.Type)
	if err != nil {
		return nil, err
	}
	nrop, err := readOp(s, n)
	if err != nil {
		return nil, err
	}
	ant, err := edgeAnnotation(e)
	if err != nil {
		return nil, err
	}
	// Alter incorrect fields.
	op.Summary = fmt.Sprintf("Find the attached %s", e.Type.Name)
	op.Description = fmt.Sprintf("Find the attached %s of the %s with the given ID", e.Type.Name, n.Name)
	op.Tags = []string{n.Name}
	op.Parameters = nrop.Parameters
	op.Parameters[0].Description = fmt.Sprintf("ID of the %s", n.Name)
	op.OperationID = opRead + n.Name + strcase.UpperCamelCase(e.Name)
	op.Responses[strconv.Itoa(http.StatusOK)].Response.Description = fmt.Sprintf(
		"%s attached to %s with requested ID was found", e.Type.Name, n.Name,
	)
	op.Security = ant.Security
	return op, nil
}

func listEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
	op, err := listOp(s, e.Type)
	if err != nil {
		return nil, err
	}
	rop, err := readOp(s, n)
	if err != nil {
		return nil, err
	}
	ant, err := edgeAnnotation(e)
	if err != nil {
		return nil, err
	}
	// Alter incorrect fields.
	op.Summary = fmt.Sprintf("Find the attached %s", rules.Pluralize(e.Type.Name))
	op.Description = fmt.Sprintf("Find the attached %s of the %s with the given ID", rules.Pluralize(e.Type.Name), n.Name)
	op.Tags = []string{n.Name}
	op.OperationID = opList + n.Name + strcase.UpperCamelCase(e.Name)
	op.Parameters = append(op.Parameters, rop.Parameters...)
	op.Responses[strconv.Itoa(http.StatusOK)].Response.Description = fmt.Sprintf(
		"%s attached to %s with requested ID was found", rules.Pluralize(e.Type.Name), n.Name,
	)
	op.Security = ant.Security
	return op, nil
}

// path returns the correct spec.Path for the given root. Creates and sets a fresh instance if non does yet exist.
func path(s *spec.Spec, root string) *spec.Path {
	if s.Paths == nil {
		s.Paths = make(map[string]*spec.Path)
	}
	if _, ok := s.Paths[root]; !ok {
		s.Paths[root] = new(spec.Path)
	}
	return s.Paths[root]
}

// schemaAnnotation returns the SchemaAnnotation of this node.
func schemaAnnotation(n *gen.Type) (*SchemaAnnotation, error) {
	ant := &SchemaAnnotation{}
	if n.Annotations != nil && n.Annotations[ant.Name()] != nil {
		if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
			return nil, err
		}
	}
	return ant, nil
}

// edgeAnnotation returns the Annotation of this edge.
func edgeAnnotation(e *gen.Edge) (*Annotation, error) {
	ant := &Annotation{}
	if e.Annotations != nil && e.Annotations[ant.Name()] != nil {
		if err := ant.Decode(e.Annotations[ant.Name()]); err != nil {
			return nil, err
		}
	}
	return ant, nil
}

// oasType returns the spec.Type to use for the given field.
func oasType(f *gen.Field) (*spec.Type, error) {
	if f.IsEnum() {
		return _string, nil
	}

	s := f.Type.String()
	if strings.Contains(s, "[]") {
		ending := strings.Replace(s, "[]", "", 1)
		t, ok := oasTypes[ending]
		if !ok {
			return nil, fmt.Errorf("no OAS-type exists for %q", s)
		}
		return &spec.Type{
			Type:   "array",
			Format: t.Format,
			Items:  t,
		}, nil
	}
	t, ok := oasTypes[s]
	if !ok {
		return nil, fmt.Errorf("no OAS-type exists for %q", s)
	}
	return t, nil
}
