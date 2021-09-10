// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"sync"

	"github.com/masseelch/elk/internal/fridge/ent/compartment"
	"github.com/masseelch/elk/internal/fridge/ent/content"
	"github.com/masseelch/elk/internal/fridge/ent/fridge"
	"github.com/masseelch/elk/internal/fridge/ent/predicate"

	"entgo.io/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeCompartment = "Compartment"
	TypeContent     = "Content"
	TypeFridge      = "Fridge"
)

// CompartmentMutation represents an operation that mutates the Compartment nodes in the graph.
type CompartmentMutation struct {
	config
	op              Op
	typ             string
	id              *int
	name            *string
	clearedFields   map[string]struct{}
	fridge          *int
	clearedfridge   bool
	contents        map[int]struct{}
	removedcontents map[int]struct{}
	clearedcontents bool
	done            bool
	oldValue        func(context.Context) (*Compartment, error)
	predicates      []predicate.Compartment
}

var _ ent.Mutation = (*CompartmentMutation)(nil)

// compartmentOption allows management of the mutation configuration using functional options.
type compartmentOption func(*CompartmentMutation)

// newCompartmentMutation creates new mutation for the Compartment entity.
func newCompartmentMutation(c config, op Op, opts ...compartmentOption) *CompartmentMutation {
	m := &CompartmentMutation{
		config:        c,
		op:            op,
		typ:           TypeCompartment,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withCompartmentID sets the ID field of the mutation.
func withCompartmentID(id int) compartmentOption {
	return func(m *CompartmentMutation) {
		var (
			err   error
			once  sync.Once
			value *Compartment
		)
		m.oldValue = func(ctx context.Context) (*Compartment, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Compartment.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withCompartment sets the old Compartment of the mutation.
func withCompartment(node *Compartment) compartmentOption {
	return func(m *CompartmentMutation) {
		m.oldValue = func(context.Context) (*Compartment, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m CompartmentMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m CompartmentMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *CompartmentMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetName sets the "name" field.
func (m *CompartmentMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *CompartmentMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Compartment entity.
// If the Compartment object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *CompartmentMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *CompartmentMutation) ResetName() {
	m.name = nil
}

// SetFridgeID sets the "fridge" edge to the Fridge entity by id.
func (m *CompartmentMutation) SetFridgeID(id int) {
	m.fridge = &id
}

// ClearFridge clears the "fridge" edge to the Fridge entity.
func (m *CompartmentMutation) ClearFridge() {
	m.clearedfridge = true
}

// FridgeCleared reports if the "fridge" edge to the Fridge entity was cleared.
func (m *CompartmentMutation) FridgeCleared() bool {
	return m.clearedfridge
}

// FridgeID returns the "fridge" edge ID in the mutation.
func (m *CompartmentMutation) FridgeID() (id int, exists bool) {
	if m.fridge != nil {
		return *m.fridge, true
	}
	return
}

// FridgeIDs returns the "fridge" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// FridgeID instead. It exists only for internal usage by the builders.
func (m *CompartmentMutation) FridgeIDs() (ids []int) {
	if id := m.fridge; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetFridge resets all changes to the "fridge" edge.
func (m *CompartmentMutation) ResetFridge() {
	m.fridge = nil
	m.clearedfridge = false
}

// AddContentIDs adds the "contents" edge to the Content entity by ids.
func (m *CompartmentMutation) AddContentIDs(ids ...int) {
	if m.contents == nil {
		m.contents = make(map[int]struct{})
	}
	for i := range ids {
		m.contents[ids[i]] = struct{}{}
	}
}

// ClearContents clears the "contents" edge to the Content entity.
func (m *CompartmentMutation) ClearContents() {
	m.clearedcontents = true
}

// ContentsCleared reports if the "contents" edge to the Content entity was cleared.
func (m *CompartmentMutation) ContentsCleared() bool {
	return m.clearedcontents
}

// RemoveContentIDs removes the "contents" edge to the Content entity by IDs.
func (m *CompartmentMutation) RemoveContentIDs(ids ...int) {
	if m.removedcontents == nil {
		m.removedcontents = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.contents, ids[i])
		m.removedcontents[ids[i]] = struct{}{}
	}
}

// RemovedContents returns the removed IDs of the "contents" edge to the Content entity.
func (m *CompartmentMutation) RemovedContentsIDs() (ids []int) {
	for id := range m.removedcontents {
		ids = append(ids, id)
	}
	return
}

// ContentsIDs returns the "contents" edge IDs in the mutation.
func (m *CompartmentMutation) ContentsIDs() (ids []int) {
	for id := range m.contents {
		ids = append(ids, id)
	}
	return
}

// ResetContents resets all changes to the "contents" edge.
func (m *CompartmentMutation) ResetContents() {
	m.contents = nil
	m.clearedcontents = false
	m.removedcontents = nil
}

// Where appends a list predicates to the CompartmentMutation builder.
func (m *CompartmentMutation) Where(ps ...predicate.Compartment) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *CompartmentMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Compartment).
func (m *CompartmentMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *CompartmentMutation) Fields() []string {
	fields := make([]string, 0, 1)
	if m.name != nil {
		fields = append(fields, compartment.FieldName)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *CompartmentMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case compartment.FieldName:
		return m.Name()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *CompartmentMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case compartment.FieldName:
		return m.OldName(ctx)
	}
	return nil, fmt.Errorf("unknown Compartment field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *CompartmentMutation) SetField(name string, value ent.Value) error {
	switch name {
	case compartment.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	}
	return fmt.Errorf("unknown Compartment field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *CompartmentMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *CompartmentMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *CompartmentMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Compartment numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *CompartmentMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *CompartmentMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *CompartmentMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Compartment nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *CompartmentMutation) ResetField(name string) error {
	switch name {
	case compartment.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown Compartment field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *CompartmentMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.fridge != nil {
		edges = append(edges, compartment.EdgeFridge)
	}
	if m.contents != nil {
		edges = append(edges, compartment.EdgeContents)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *CompartmentMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case compartment.EdgeFridge:
		if id := m.fridge; id != nil {
			return []ent.Value{*id}
		}
	case compartment.EdgeContents:
		ids := make([]ent.Value, 0, len(m.contents))
		for id := range m.contents {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *CompartmentMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	if m.removedcontents != nil {
		edges = append(edges, compartment.EdgeContents)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *CompartmentMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case compartment.EdgeContents:
		ids := make([]ent.Value, 0, len(m.removedcontents))
		for id := range m.removedcontents {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *CompartmentMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedfridge {
		edges = append(edges, compartment.EdgeFridge)
	}
	if m.clearedcontents {
		edges = append(edges, compartment.EdgeContents)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *CompartmentMutation) EdgeCleared(name string) bool {
	switch name {
	case compartment.EdgeFridge:
		return m.clearedfridge
	case compartment.EdgeContents:
		return m.clearedcontents
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *CompartmentMutation) ClearEdge(name string) error {
	switch name {
	case compartment.EdgeFridge:
		m.ClearFridge()
		return nil
	}
	return fmt.Errorf("unknown Compartment unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *CompartmentMutation) ResetEdge(name string) error {
	switch name {
	case compartment.EdgeFridge:
		m.ResetFridge()
		return nil
	case compartment.EdgeContents:
		m.ResetContents()
		return nil
	}
	return fmt.Errorf("unknown Compartment edge %s", name)
}

// ContentMutation represents an operation that mutates the Content nodes in the graph.
type ContentMutation struct {
	config
	op                 Op
	typ                string
	id                 *int
	name               *string
	clearedFields      map[string]struct{}
	compartment        *int
	clearedcompartment bool
	done               bool
	oldValue           func(context.Context) (*Content, error)
	predicates         []predicate.Content
}

var _ ent.Mutation = (*ContentMutation)(nil)

// contentOption allows management of the mutation configuration using functional options.
type contentOption func(*ContentMutation)

// newContentMutation creates new mutation for the Content entity.
func newContentMutation(c config, op Op, opts ...contentOption) *ContentMutation {
	m := &ContentMutation{
		config:        c,
		op:            op,
		typ:           TypeContent,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withContentID sets the ID field of the mutation.
func withContentID(id int) contentOption {
	return func(m *ContentMutation) {
		var (
			err   error
			once  sync.Once
			value *Content
		)
		m.oldValue = func(ctx context.Context) (*Content, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Content.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withContent sets the old Content of the mutation.
func withContent(node *Content) contentOption {
	return func(m *ContentMutation) {
		m.oldValue = func(context.Context) (*Content, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ContentMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ContentMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ContentMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetName sets the "name" field.
func (m *ContentMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *ContentMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Content entity.
// If the Content object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ContentMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *ContentMutation) ResetName() {
	m.name = nil
}

// SetCompartmentID sets the "compartment" edge to the Compartment entity by id.
func (m *ContentMutation) SetCompartmentID(id int) {
	m.compartment = &id
}

// ClearCompartment clears the "compartment" edge to the Compartment entity.
func (m *ContentMutation) ClearCompartment() {
	m.clearedcompartment = true
}

// CompartmentCleared reports if the "compartment" edge to the Compartment entity was cleared.
func (m *ContentMutation) CompartmentCleared() bool {
	return m.clearedcompartment
}

// CompartmentID returns the "compartment" edge ID in the mutation.
func (m *ContentMutation) CompartmentID() (id int, exists bool) {
	if m.compartment != nil {
		return *m.compartment, true
	}
	return
}

// CompartmentIDs returns the "compartment" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// CompartmentID instead. It exists only for internal usage by the builders.
func (m *ContentMutation) CompartmentIDs() (ids []int) {
	if id := m.compartment; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetCompartment resets all changes to the "compartment" edge.
func (m *ContentMutation) ResetCompartment() {
	m.compartment = nil
	m.clearedcompartment = false
}

// Where appends a list predicates to the ContentMutation builder.
func (m *ContentMutation) Where(ps ...predicate.Content) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *ContentMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Content).
func (m *ContentMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ContentMutation) Fields() []string {
	fields := make([]string, 0, 1)
	if m.name != nil {
		fields = append(fields, content.FieldName)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ContentMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case content.FieldName:
		return m.Name()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ContentMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case content.FieldName:
		return m.OldName(ctx)
	}
	return nil, fmt.Errorf("unknown Content field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ContentMutation) SetField(name string, value ent.Value) error {
	switch name {
	case content.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	}
	return fmt.Errorf("unknown Content field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ContentMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ContentMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ContentMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Content numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ContentMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ContentMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ContentMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Content nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ContentMutation) ResetField(name string) error {
	switch name {
	case content.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown Content field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ContentMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.compartment != nil {
		edges = append(edges, content.EdgeCompartment)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ContentMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case content.EdgeCompartment:
		if id := m.compartment; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ContentMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ContentMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ContentMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedcompartment {
		edges = append(edges, content.EdgeCompartment)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ContentMutation) EdgeCleared(name string) bool {
	switch name {
	case content.EdgeCompartment:
		return m.clearedcompartment
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ContentMutation) ClearEdge(name string) error {
	switch name {
	case content.EdgeCompartment:
		m.ClearCompartment()
		return nil
	}
	return fmt.Errorf("unknown Content unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ContentMutation) ResetEdge(name string) error {
	switch name {
	case content.EdgeCompartment:
		m.ResetCompartment()
		return nil
	}
	return fmt.Errorf("unknown Content edge %s", name)
}

// FridgeMutation represents an operation that mutates the Fridge nodes in the graph.
type FridgeMutation struct {
	config
	op                  Op
	typ                 string
	id                  *int
	title               *string
	clearedFields       map[string]struct{}
	compartments        map[int]struct{}
	removedcompartments map[int]struct{}
	clearedcompartments bool
	done                bool
	oldValue            func(context.Context) (*Fridge, error)
	predicates          []predicate.Fridge
}

var _ ent.Mutation = (*FridgeMutation)(nil)

// fridgeOption allows management of the mutation configuration using functional options.
type fridgeOption func(*FridgeMutation)

// newFridgeMutation creates new mutation for the Fridge entity.
func newFridgeMutation(c config, op Op, opts ...fridgeOption) *FridgeMutation {
	m := &FridgeMutation{
		config:        c,
		op:            op,
		typ:           TypeFridge,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withFridgeID sets the ID field of the mutation.
func withFridgeID(id int) fridgeOption {
	return func(m *FridgeMutation) {
		var (
			err   error
			once  sync.Once
			value *Fridge
		)
		m.oldValue = func(ctx context.Context) (*Fridge, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Fridge.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withFridge sets the old Fridge of the mutation.
func withFridge(node *Fridge) fridgeOption {
	return func(m *FridgeMutation) {
		m.oldValue = func(context.Context) (*Fridge, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m FridgeMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m FridgeMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *FridgeMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetTitle sets the "title" field.
func (m *FridgeMutation) SetTitle(s string) {
	m.title = &s
}

// Title returns the value of the "title" field in the mutation.
func (m *FridgeMutation) Title() (r string, exists bool) {
	v := m.title
	if v == nil {
		return
	}
	return *v, true
}

// OldTitle returns the old "title" field's value of the Fridge entity.
// If the Fridge object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *FridgeMutation) OldTitle(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldTitle is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldTitle requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTitle: %w", err)
	}
	return oldValue.Title, nil
}

// ResetTitle resets all changes to the "title" field.
func (m *FridgeMutation) ResetTitle() {
	m.title = nil
}

// AddCompartmentIDs adds the "compartments" edge to the Compartment entity by ids.
func (m *FridgeMutation) AddCompartmentIDs(ids ...int) {
	if m.compartments == nil {
		m.compartments = make(map[int]struct{})
	}
	for i := range ids {
		m.compartments[ids[i]] = struct{}{}
	}
}

// ClearCompartments clears the "compartments" edge to the Compartment entity.
func (m *FridgeMutation) ClearCompartments() {
	m.clearedcompartments = true
}

// CompartmentsCleared reports if the "compartments" edge to the Compartment entity was cleared.
func (m *FridgeMutation) CompartmentsCleared() bool {
	return m.clearedcompartments
}

// RemoveCompartmentIDs removes the "compartments" edge to the Compartment entity by IDs.
func (m *FridgeMutation) RemoveCompartmentIDs(ids ...int) {
	if m.removedcompartments == nil {
		m.removedcompartments = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.compartments, ids[i])
		m.removedcompartments[ids[i]] = struct{}{}
	}
}

// RemovedCompartments returns the removed IDs of the "compartments" edge to the Compartment entity.
func (m *FridgeMutation) RemovedCompartmentsIDs() (ids []int) {
	for id := range m.removedcompartments {
		ids = append(ids, id)
	}
	return
}

// CompartmentsIDs returns the "compartments" edge IDs in the mutation.
func (m *FridgeMutation) CompartmentsIDs() (ids []int) {
	for id := range m.compartments {
		ids = append(ids, id)
	}
	return
}

// ResetCompartments resets all changes to the "compartments" edge.
func (m *FridgeMutation) ResetCompartments() {
	m.compartments = nil
	m.clearedcompartments = false
	m.removedcompartments = nil
}

// Where appends a list predicates to the FridgeMutation builder.
func (m *FridgeMutation) Where(ps ...predicate.Fridge) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *FridgeMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Fridge).
func (m *FridgeMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *FridgeMutation) Fields() []string {
	fields := make([]string, 0, 1)
	if m.title != nil {
		fields = append(fields, fridge.FieldTitle)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *FridgeMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case fridge.FieldTitle:
		return m.Title()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *FridgeMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case fridge.FieldTitle:
		return m.OldTitle(ctx)
	}
	return nil, fmt.Errorf("unknown Fridge field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *FridgeMutation) SetField(name string, value ent.Value) error {
	switch name {
	case fridge.FieldTitle:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTitle(v)
		return nil
	}
	return fmt.Errorf("unknown Fridge field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *FridgeMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *FridgeMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *FridgeMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Fridge numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *FridgeMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *FridgeMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *FridgeMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Fridge nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *FridgeMutation) ResetField(name string) error {
	switch name {
	case fridge.FieldTitle:
		m.ResetTitle()
		return nil
	}
	return fmt.Errorf("unknown Fridge field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *FridgeMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.compartments != nil {
		edges = append(edges, fridge.EdgeCompartments)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *FridgeMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case fridge.EdgeCompartments:
		ids := make([]ent.Value, 0, len(m.compartments))
		for id := range m.compartments {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *FridgeMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedcompartments != nil {
		edges = append(edges, fridge.EdgeCompartments)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *FridgeMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case fridge.EdgeCompartments:
		ids := make([]ent.Value, 0, len(m.removedcompartments))
		for id := range m.removedcompartments {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *FridgeMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedcompartments {
		edges = append(edges, fridge.EdgeCompartments)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *FridgeMutation) EdgeCleared(name string) bool {
	switch name {
	case fridge.EdgeCompartments:
		return m.clearedcompartments
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *FridgeMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Fridge unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *FridgeMutation) ResetEdge(name string) error {
	switch name {
	case fridge.EdgeCompartments:
		m.ResetCompartments()
		return nil
	}
	return fmt.Errorf("unknown Fridge edge %s", name)
}