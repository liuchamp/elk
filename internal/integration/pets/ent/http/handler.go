// Code generated by entc, DO NOT EDIT.

package http

import (
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/masseelch/elk/internal/integration/pets/ent"
	"go.uber.org/zap"
)

// handler has some convenience methods used on node-handlers.
type handler struct{}

// CategoryHandler handles http crud operations on ent.Category.
type CategoryHandler struct {
	handler

	client    *ent.Client
	log       *zap.Logger
	validator *validator.Validate
}

func NewCategoryHandler(c *ent.Client, l *zap.Logger, v *validator.Validate) *CategoryHandler {
	return &CategoryHandler{
		client:    c,
		log:       l.With(zap.String("handler", "CategoryHandler")),
		validator: v,
	}
}

// RegisterHandlers registers the generated handlers on the given chi router.
func (h *CategoryHandler) RegisterHandlers(r chi.Router) {
	// Do no use r.Route() to avoid wildcard matching.
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Read)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Get("/{id}/pets", h.Pets)
}

// OwnerHandler handles http crud operations on ent.Owner.
type OwnerHandler struct {
	handler

	client    *ent.Client
	log       *zap.Logger
	validator *validator.Validate
}

func NewOwnerHandler(c *ent.Client, l *zap.Logger, v *validator.Validate) *OwnerHandler {
	return &OwnerHandler{
		client:    c,
		log:       l.With(zap.String("handler", "OwnerHandler")),
		validator: v,
	}
}

// RegisterHandlers registers the generated handlers on the given chi router.
func (h *OwnerHandler) RegisterHandlers(r chi.Router) {
	// Do no use r.Route() to avoid wildcard matching.
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Read)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Get("/{id}/pets", h.Pets)
}

// PetHandler handles http crud operations on ent.Pet.
type PetHandler struct {
	handler

	client    *ent.Client
	log       *zap.Logger
	validator *validator.Validate
}

func NewPetHandler(c *ent.Client, l *zap.Logger, v *validator.Validate) *PetHandler {
	return &PetHandler{
		client:    c,
		log:       l.With(zap.String("handler", "PetHandler")),
		validator: v,
	}
}

// RegisterHandlers registers the generated handlers on the given chi router.
func (h *PetHandler) RegisterHandlers(r chi.Router) {
	// Do no use r.Route() to avoid wildcard matching.
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Read)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Get("/{id}/category", h.Category)
	r.Get("/{id}/owner", h.Owner)
	r.Get("/{id}/friends", h.Friends)
}

func (h handler) stripEntError(err error) string {
	return strings.TrimPrefix(err.Error(), "ent: ")
}