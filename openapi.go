package elk

import (
	"fmt"
	"net/http"
	"strings"

	"entgo.io/ent/entc/gen"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
)

var (
	_int32    = &openapi3.Types{openapi3.TypeInteger}
	_int64    = &openapi3.Types{openapi3.TypeInteger}
	_float    = &openapi3.Types{openapi3.TypeNumber}
	_double   = &openapi3.Types{openapi3.TypeNumber}
	_string   = &openapi3.Types{openapi3.TypeString}
	_bool     = &openapi3.Types{openapi3.TypeBoolean}
	_dateTime = &openapi3.Types{openapi3.TypeString}
	oasTypes  = map[string]*openapi3.Types{
		"bool":            _bool,
		"time.Time":       _dateTime,
		"time.Duration":   _int64,
		"enum":            _string,
		"string":          _string,
		"uuid.UUID":       _string,
		"int":             _int32,
		"int8":            _int32,
		"int16":           _int32,
		"int32":           _int32,
		"uint":            _int32,
		"uint8":           _int32,
		"uint16":          _int32,
		"uint32":          _int32,
		"int64":           _int64,
		"uint64":          _int64,
		"float32":         _float,
		"float64":         _double,
		"json.RawMessage": _string,
		"[16]byte":        _string,
		"[]byte":          _string,
	}
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
		fs := openapi3.Schemas{}
		// We can already add the schema fields.
		for _, f := range v.Fields {
			sf, err := newField(f)
			if err != nil {
				return err
			}
			fs[f.Name] = sf
		}
		s.Components.Schemas[n] = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:       &openapi3.Types{"object"},
				Properties: fs,
			},
		}
	}
	// Loop over the views again and this time fill the edges.
	//for n, v := range vs {
	//	es := make(spec.Edges, len(v.Edges))
	//	for _, e := range v.Edges {
	//		es[e.Edge.Name] = spec.Edge{
	//			Ref:    s.Components.Schemas[e.Name],
	//			Unique: e.Unique,
	//		}
	//	}
	//	s.Components.Schemas[n].Edges = es
	//}
	return nil
}

// newField constructs a spec.Field out of a gen.Field.
func newField(f *gen.Field) (*openapi3.SchemaRef, error) {
	t, err := oasType(f)
	if err != nil {
		return nil, err
	}
	//e, err := exampleValue(f)
	//if err != nil {
	//	return nil, err
	//}
	return &openapi3.SchemaRef{

		Value: &openapi3.Schema{
			Type: t,
			//Example: e,
			//Required: !f.Optional,
		},
	}, nil
}

// paths adds all views to the specs schemas.
func paths(g *gen.Graph, s *openapi3.T) error {
	for _, n := range g.Nodes {
		// Add schema operations.
		ops, err := nodeOperations(n)
		if err != nil {
			return err
		}
		// root for all operations on this node.
		root := "/" + strcase.KebabCase(n.Name)
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
		if contains(ops, opPatch) {
			path(s, root+"/{id}").Patch, err = patchOp(s, n)
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
		//es, err := filterEdges(n)
		//if err != nil {
		//	return err
		//}
		//for _, e := range es {
		//	p := path(s, root+"/{id}/"+strcase.KebabCase(e.Name))
		//	if e.Unique {
		//		p.Get, err = readEdgeOp(s, n, e)
		//		if err != nil {
		//			return err
		//		}
		//	} else {
		//		p.Get, err = listEdgeOp(s, n, e)
		//		if err != nil {
		//			return err
		//		}
		//	}
		//}
	}
	return nil
}

func requestBody(n *gen.Type, method string) (*openapi3.RequestBodyRef, error) {

	data := openapi3.Schemas{}

	var rq []string
	for _, field := range n.Fields {
		tps, ok := oasTypes[field.Type.String()]
		if !ok {
			return nil, fmt.Errorf("cat not support type %s", field.Type.String())
		}
		v := &openapi3.Schema{Type: tps}
		if method == opCreate || method == opUpdate {
			rq = append(rq, field.Name)
		}
		data[n.Name] = &openapi3.SchemaRef{
			Value: v,
		}
	}
	content := openapi3.Content{
		"application/json": &openapi3.MediaType{
			Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{
				Type:       &openapi3.Types{openapi3.TypeObject},
				Properties: data,
				Required:   rq,
			}},
		},
	}
	return &openapi3.RequestBodyRef{Value: openapi3.NewRequestBody().WithRequired(true).WithContent(content).WithDescription(fmt.Sprintf("%s %s request body", method, n.Name))}, nil
}

// createOp returns the spec description for a create operation on the given node.
func createOp(s *openapi3.T, n *gen.Type) (*openapi3.Operation, error) {
	//ant, err := schemaAnnotation(n)
	//if err != nil {
	//	return nil, err
	//}
	req, err := requestBody(n, opCreate)
	if err != nil {
		return nil, err
	}
	//v, err := newView(n, ant.CreateGroups)
	//if err != nil {
	//	return nil, err
	//}
	//rspName, err := v.Name()
	//if err != nil {
	//	return nil, err
	//}
	resp := openapi3.WithStatus(200, &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("create a new %s ok ", n.Name)).
			WithContent(openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/"+strcase.UpperCamelCase(n.Name), nil))),
	})

	errResp := optStatus()
	errResp = append([]openapi3.NewResponsesOption{resp}, errResp...)
	rsp := openapi3.NewResponses(errResp...)
	return &openapi3.Operation{
		Summary:     fmt.Sprintf("Create a new %s", n.Name),
		Description: fmt.Sprintf("Creates a new %s and persists it to storage.", n.Name),
		Tags:        []string{n.Name},
		OperationID: opCreate + n.Name,
		RequestBody: req,
		Responses:   rsp,
		//Security:    ant.CreateSecurity,
	}, nil
}

// readOp returns a spec.Operation for a read operation on the given node.
func readOp(s *openapi3.T, n *gen.Type) (*openapi3.Operation, error) {
	//ant, err := schemaAnnotation(n)
	//if err != nil {
	//	return nil, err
	//}
	//v, err := newView(n, ant.ReadGroups)
	//if err != nil {
	//	return nil, err
	//}
	//rspName, err := v.Name()
	//if err != nil {
	//	return nil, err
	//}
	t, err := oasType(n.ID)
	if err != nil {
		return nil, err
	}
	resp := openapi3.WithStatus(200, &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("patch a  %s ok ", n.Name)).
			WithContent(openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/"+strcase.UpperCamelCase(n.Name), nil))),
	})

	errResp := optStatus()
	errResp = append([]openapi3.NewResponsesOption{resp}, errResp...)
	rsp := openapi3.NewResponses(errResp...)
	return &openapi3.Operation{
		Summary:     fmt.Sprintf("Find a %s by ID", n.Name),
		Description: fmt.Sprintf("Finds the %s with the requested ID and returns it.", n.Name),
		Tags:        []string{n.Name},
		OperationID: opRead + n.Name,
		Parameters: openapi3.Parameters{&openapi3.ParameterRef{
			Value: openapi3.NewPathParameter("id").WithSchema(&openapi3.Schema{Type: t}).WithDescription(fmt.Sprintf("ID of the %s to update", n.Name)),
		}},
		Responses: rsp,
		//Security: ant.ReadSecurity,
	}, nil
}

func patchOp(s *openapi3.T, n *gen.Type) (*openapi3.Operation, error) {
	//ant, err := schemaAnnotation(n)
	//if err != nil {
	//	return nil, err
	//}
	req, err := requestBody(n, opPatch)
	if err != nil {
		return nil, err
	}
	//v, err := newView(n, ant.UpdateGroups)
	//if err != nil {
	//	return nil, err
	//}
	//rspName, err := v.Name()
	//if err != nil {
	//	return nil, err
	//}
	t, err := oasType(n.ID)
	if err != nil {
		return nil, err
	}
	resp := openapi3.WithStatus(200, &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("patch a  %s ok ", n.Name)).
			WithContent(openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/"+strcase.UpperCamelCase(n.Name), nil))),
	})

	errResp := optStatus()
	errResp = append([]openapi3.NewResponsesOption{resp}, errResp...)
	rsp := openapi3.NewResponses(errResp...)
	return &openapi3.Operation{
		Summary:     fmt.Sprintf("patch a %s", n.Name),
		Description: fmt.Sprintf("patch a %s and persists changes to storage.", n.Name),
		Tags:        []string{n.Name},
		OperationID: opUpdate + n.Name,
		Parameters: openapi3.Parameters{&openapi3.ParameterRef{
			Value: openapi3.NewPathParameter("id").WithSchema(&openapi3.Schema{Type: t}).WithDescription(fmt.Sprintf("ID of the %s to update", n.Name)),
		}},
		RequestBody: req,
		Responses:   rsp,
		//Security: ant.UpdateSecurity,
	}, nil
}

func updateOp(s *openapi3.T, n *gen.Type) (*openapi3.Operation, error) {
	//ant, err := schemaAnnotation(n)
	//if err != nil {
	//	return nil, err
	//}
	req, err := requestBody(n, opUpdate)
	if err != nil {
		return nil, err
	}
	//v, err := newView(n, ant.UpdateGroups)
	//if err != nil {
	//	return nil, err
	//}
	//rspName, err := v.Name()
	//if err != nil {
	//	return nil, err
	//}
	t, err := oasType(n.ID)
	if err != nil {
		return nil, err
	}
	resp := openapi3.WithStatus(200, &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("patch a  %s ok ", n.Name)).
			WithContent(openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/"+strcase.UpperCamelCase(n.Name), nil))),
	})

	errResp := optStatus()
	errResp = append([]openapi3.NewResponsesOption{resp}, errResp...)
	rsp := openapi3.NewResponses(errResp...)
	return &openapi3.Operation{
		Summary:     fmt.Sprintf("Updates a %s", n.Name),
		Description: fmt.Sprintf("Updates a %s and persists changes to storage.", n.Name),
		Tags:        []string{n.Name},
		OperationID: opUpdate + n.Name,
		Parameters: openapi3.Parameters{&openapi3.ParameterRef{
			Value: openapi3.NewPathParameter("id").WithSchema(&openapi3.Schema{Type: t}).WithDescription(fmt.Sprintf("ID of the %s to update", n.Name)),
		}},
		RequestBody: req,
		Responses:   rsp,
		//Security: ant.UpdateSecurity,
	}, nil
}

func deleteOp(s *openapi3.T, n *gen.Type) (*openapi3.Operation, error) {
	//ant, err := schemaAnnotation(n)
	//if err != nil {
	//	return nil, err
	//}
	t, err := oasType(n.ID)
	if err != nil {
		return nil, err
	}
	resp := openapi3.WithStatus(200, &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("delete a  %s ok ", n.Name)).
			WithContent(openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/response/Delete", nil))),
	})

	errResp := optStatus()
	errResp = append([]openapi3.NewResponsesOption{resp}, errResp...)
	rsp := openapi3.NewResponses(errResp...)
	return &openapi3.Operation{
		Summary:     fmt.Sprintf("Deletes a %s by ID", n.Name),
		Description: fmt.Sprintf("Deletes the %s with the requested ID.", n.Name),
		Tags:        []string{n.Name},
		OperationID: opDelete + n.Name,
		Parameters: openapi3.Parameters{&openapi3.ParameterRef{
			Value: openapi3.NewPathParameter("id").WithSchema(&openapi3.Schema{Type: t}).WithDescription(fmt.Sprintf("ID of the %s to update", n.Name)),
		}},
		Responses: rsp,
		//Security: ant.DeleteSecurity,
	}, nil
}

func listOp(s *openapi3.T, n *gen.Type) (*openapi3.Operation, error) {
	//ant, err := schemaAnnotation(n)
	//if err != nil {
	//	return nil, err
	//}
	//v, err := newView(n, ant.ListGroups)
	//if err != nil {
	//	return nil, err
	//}
	//rspName, err := v.Name()
	//if err != nil {
	//	return nil, err
	//}
	respbody := openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:  &openapi3.Types{openapi3.TypeArray},
		Items: openapi3.NewSchemaRef("#/components/schemas/"+strcase.UpperCamelCase(n.Name), nil),
	})
	resp := openapi3.WithStatus(200, &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("query list  %s ok ", n.Name)).
			WithContent(openapi3.NewContentWithJSONSchemaRef(respbody)),
	})

	errResp := optStatus()
	errResp = append([]openapi3.NewResponsesOption{resp}, errResp...)
	rsp := openapi3.NewResponses(errResp...)

	return &openapi3.Operation{
		Summary:     fmt.Sprintf("List %s", n.Name),
		Description: fmt.Sprintf("List %s.", n.Name),
		Tags:        []string{n.Name},
		OperationID: opList + n.Name,
		Parameters:  nil,
		Responses:   rsp,
		//Security: ant.ListSecurity,
	}, nil
}

//
//func readEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
//	op, err := readOp(s, e.Type)
//	if err != nil {
//		return nil, err
//	}
//	nrop, err := readOp(s, n)
//	if err != nil {
//		return nil, err
//	}
//	ant, err := edgeAnnotation(e)
//	if err != nil {
//		return nil, err
//	}
//	// Alter incorrect fields.
//	op.Summary = fmt.Sprintf("Find the attached %s", e.Type.Name)
//	op.Description = fmt.Sprintf("Find the attached %s of the %s with the given ID", e.Type.Name, n.Name)
//	op.Tags = []string{n.Name}
//	op.Parameters = nrop.Parameters
//	op.Parameters[0].Description = fmt.Sprintf("ID of the %s", n.Name)
//	op.OperationID = opRead + n.Name + strcase.UpperCamelCase(e.Name)
//	op.Responses[strconv.Itoa(http.StatusOK)].Response.Description = fmt.Sprintf(
//		"%s attached to %s with requested ID was found", e.Type.Name, n.Name,
//	)
//	op.Security = ant.Security
//	return op, nil
//}
//
//func listEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
//	op, err := listOp(s, e.Type)
//	if err != nil {
//		return nil, err
//	}
//	rop, err := readOp(s, n)
//	if err != nil {
//		return nil, err
//	}
//	ant, err := edgeAnnotation(e)
//	if err != nil {
//		return nil, err
//	}
//	// Alter incorrect fields.
//	op.Summary = fmt.Sprintf("Find the attached %s", rules.Pluralize(e.Type.Name))
//	op.Description = fmt.Sprintf("Find the attached %s of the %s with the given ID", rules.Pluralize(e.Type.Name), n.Name)
//	op.Tags = []string{n.Name}
//	op.OperationID = opList + n.Name + strcase.UpperCamelCase(e.Name)
//	op.Parameters = append(op.Parameters, rop.Parameters...)
//	op.Responses[strconv.Itoa(http.StatusOK)].Response.Description = fmt.Sprintf(
//		"%s attached to %s with requested ID was found", rules.Pluralize(e.Type.Name), n.Name,
//	)
//	op.Security = ant.Security
//	return op, nil
//}

// path returns the correct spec.Path for the given root. Creates and sets a fresh instance if non does yet exist.
func path(s *openapi3.T, root string) *openapi3.PathItem {
	if s.Paths == nil {
		s.Paths = openapi3.NewPaths()
	}
	if tdx := s.Paths.Value(root); tdx == nil {
		s.Paths.Set(root, &openapi3.PathItem{})
	}
	return s.Paths.Value(root)
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
func oasType(f *gen.Field) (*openapi3.Types, error) {
	if f.IsEnum() {
		return _string, nil
	}

	s := f.Type.String()
	if strings.Contains(s, "[]") {
		return &openapi3.Types{"array"}, nil
		//ending := strings.Replace(s, "[]", "", 1)
		//t, ok := oasTypes[ending]
		//if !ok {
		//	return nil, fmt.Errorf("no OAS-type exists for %q", s)
		//}
		//return &spec.Type{
		//	Type:   "array",
		//	Format: t.Format,
		//	Items:  t,
		//}, nil
	}
	t, ok := oasTypes[s]
	if !ok {
		return nil, fmt.Errorf("no OAS-type exists for %q", s)
	}
	return t, nil
}

func optStatus(opts ...int) []openapi3.NewResponsesOption {
	optErrMap := map[int]openapi3.NewResponsesOption{
		http.StatusBadRequest:          openapi3.WithStatus(http.StatusBadRequest, &openapi3.ResponseRef{Ref: "#/components/schemas/400"}),
		http.StatusNotFound:            openapi3.WithStatus(http.StatusNotFound, &openapi3.ResponseRef{Ref: "#/components/schemas/404"}),
		http.StatusInternalServerError: openapi3.WithStatus(http.StatusInternalServerError, &openapi3.ResponseRef{Ref: "#/components/schemas/500"}),
	}
	var rest []openapi3.NewResponsesOption
	if opts == nil {
		for _, option := range optErrMap {
			rest = append(rest, option)
		}
		return rest
	}
	for _, opt := range opts {
		if s, ok := optErrMap[opt]; ok {
			rest = append(rest, s)
		}
	}
	return rest

}
