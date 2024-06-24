package elk

import (
	"encoding/json"
	"errors"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type (
	Config struct {
		// HandlerPolicy defines the default policy for handler generation.
		// It is used if no policy is set on a (sub-)resource.
		// Defaults to policy.Expose.
		HandlerPolicy Policy
	}
	// Extension implements entc.Extension interface for providing http handler code generation.
	Extension struct {
		entc.DefaultExtension
		repoConfig RepoConfig
		specHooks  []Hook
		hooks      []gen.Hook
		templates  []*gen.Template
		config     *Config
	}
	// ExtensionOption allows managing Extension configuration using functional arguments.
	ExtensionOption func(*Extension) error
	// HandlerOption allows managing RESTGenerator configuration using function arguments.
	HandlerOption ExtensionOption
)

// NewExtension returns a new elk extension with default values.
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{
		config:     &Config{HandlerPolicy: Expose},
		repoConfig: newRepoConfig(),
	}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	if len(ex.hooks) == 0 {
		return nil, errors.New(`no generator enabled: enable one by providing either "GenerateSpec()" or "GenerateHandlers()" to "NewExtension()"`)
	}
	return ex, nil
}

// Templates of the Extension.
func (e *Extension) Templates() []*gen.Template {
	return e.templates
}

// Hooks of the Extension.
func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

// Annotations of the Extension.
func (e *Extension) Annotations() []entc.Annotation {
	return []entc.Annotation{e.config}
}

// DefaultHandlerPolicy sets the policy.Policy to use of none is given on a (sub-)schema.
func DefaultHandlerPolicy(p Policy) ExtensionOption {
	return func(ex *Extension) error {
		if err := p.Validate(); err != nil {
			return err
		}
		ex.config.HandlerPolicy = p
		return nil
	}
}

// GenerateHandlers enables generation of http crud handlers.
func GenerateHandlers(opts ...HandlerOption) ExtensionOption {
	return func(ex *Extension) error {
		ex.hooks = append(ex.hooks, RepoGenerator(ex.repoConfig))
		for _, opt := range opts {
			if err := opt(ex); err != nil {
				return err
			}
		}
		return nil
	}
}

// Name implements entc.Annotation interface.
func (c Config) Name() string {
	return "ElkConfig"
}

// Decode from ent.
func (c *Config) Decode(o interface{}) error {
	buf, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, c)
}

var _ entc.Annotation = (*Config)(nil)
