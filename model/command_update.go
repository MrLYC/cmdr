// Code generated by entc, DO NOT EDIT.

package model

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/mrlyc/cmdr/model/command"
	"github.com/mrlyc/cmdr/model/predicate"
)

// CommandUpdate is the builder for updating Command entities.
type CommandUpdate struct {
	config
	hooks    []Hook
	mutation *CommandMutation
}

// Where appends a list predicates to the CommandUpdate builder.
func (cu *CommandUpdate) Where(ps ...predicate.Command) *CommandUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetName sets the "name" field.
func (cu *CommandUpdate) SetName(s string) *CommandUpdate {
	cu.mutation.SetName(s)
	return cu
}

// SetVersion sets the "version" field.
func (cu *CommandUpdate) SetVersion(s string) *CommandUpdate {
	cu.mutation.SetVersion(s)
	return cu
}

// SetLocation sets the "location" field.
func (cu *CommandUpdate) SetLocation(s string) *CommandUpdate {
	cu.mutation.SetLocation(s)
	return cu
}

// SetActivated sets the "activated" field.
func (cu *CommandUpdate) SetActivated(b bool) *CommandUpdate {
	cu.mutation.SetActivated(b)
	return cu
}

// SetNillableActivated sets the "activated" field if the given value is not nil.
func (cu *CommandUpdate) SetNillableActivated(b *bool) *CommandUpdate {
	if b != nil {
		cu.SetActivated(*b)
	}
	return cu
}

// Mutation returns the CommandMutation object of the builder.
func (cu *CommandUpdate) Mutation() *CommandMutation {
	return cu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *CommandUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(cu.hooks) == 0 {
		if err = cu.check(); err != nil {
			return 0, err
		}
		affected, err = cu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CommandMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = cu.check(); err != nil {
				return 0, err
			}
			cu.mutation = mutation
			affected, err = cu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(cu.hooks) - 1; i >= 0; i-- {
			if cu.hooks[i] == nil {
				return 0, fmt.Errorf("model: uninitialized hook (forgotten import model/runtime?)")
			}
			mut = cu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CommandUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CommandUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CommandUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cu *CommandUpdate) check() error {
	if v, ok := cu.mutation.Name(); ok {
		if err := command.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("model: validator failed for field \"name\": %w", err)}
		}
	}
	if v, ok := cu.mutation.Version(); ok {
		if err := command.VersionValidator(v); err != nil {
			return &ValidationError{Name: "version", err: fmt.Errorf("model: validator failed for field \"version\": %w", err)}
		}
	}
	if v, ok := cu.mutation.Location(); ok {
		if err := command.LocationValidator(v); err != nil {
			return &ValidationError{Name: "location", err: fmt.Errorf("model: validator failed for field \"location\": %w", err)}
		}
	}
	return nil
}

func (cu *CommandUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   command.Table,
			Columns: command.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: command.FieldID,
			},
		},
	}
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: command.FieldName,
		})
	}
	if value, ok := cu.mutation.Version(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: command.FieldVersion,
		})
	}
	if value, ok := cu.mutation.Location(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: command.FieldLocation,
		})
	}
	if value, ok := cu.mutation.Activated(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: command.FieldActivated,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{command.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// CommandUpdateOne is the builder for updating a single Command entity.
type CommandUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *CommandMutation
}

// SetName sets the "name" field.
func (cuo *CommandUpdateOne) SetName(s string) *CommandUpdateOne {
	cuo.mutation.SetName(s)
	return cuo
}

// SetVersion sets the "version" field.
func (cuo *CommandUpdateOne) SetVersion(s string) *CommandUpdateOne {
	cuo.mutation.SetVersion(s)
	return cuo
}

// SetLocation sets the "location" field.
func (cuo *CommandUpdateOne) SetLocation(s string) *CommandUpdateOne {
	cuo.mutation.SetLocation(s)
	return cuo
}

// SetActivated sets the "activated" field.
func (cuo *CommandUpdateOne) SetActivated(b bool) *CommandUpdateOne {
	cuo.mutation.SetActivated(b)
	return cuo
}

// SetNillableActivated sets the "activated" field if the given value is not nil.
func (cuo *CommandUpdateOne) SetNillableActivated(b *bool) *CommandUpdateOne {
	if b != nil {
		cuo.SetActivated(*b)
	}
	return cuo
}

// Mutation returns the CommandMutation object of the builder.
func (cuo *CommandUpdateOne) Mutation() *CommandMutation {
	return cuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *CommandUpdateOne) Select(field string, fields ...string) *CommandUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Command entity.
func (cuo *CommandUpdateOne) Save(ctx context.Context) (*Command, error) {
	var (
		err  error
		node *Command
	)
	if len(cuo.hooks) == 0 {
		if err = cuo.check(); err != nil {
			return nil, err
		}
		node, err = cuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CommandMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = cuo.check(); err != nil {
				return nil, err
			}
			cuo.mutation = mutation
			node, err = cuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(cuo.hooks) - 1; i >= 0; i-- {
			if cuo.hooks[i] == nil {
				return nil, fmt.Errorf("model: uninitialized hook (forgotten import model/runtime?)")
			}
			mut = cuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CommandUpdateOne) SaveX(ctx context.Context) *Command {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *CommandUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CommandUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cuo *CommandUpdateOne) check() error {
	if v, ok := cuo.mutation.Name(); ok {
		if err := command.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("model: validator failed for field \"name\": %w", err)}
		}
	}
	if v, ok := cuo.mutation.Version(); ok {
		if err := command.VersionValidator(v); err != nil {
			return &ValidationError{Name: "version", err: fmt.Errorf("model: validator failed for field \"version\": %w", err)}
		}
	}
	if v, ok := cuo.mutation.Location(); ok {
		if err := command.LocationValidator(v); err != nil {
			return &ValidationError{Name: "location", err: fmt.Errorf("model: validator failed for field \"location\": %w", err)}
		}
	}
	return nil
}

func (cuo *CommandUpdateOne) sqlSave(ctx context.Context) (_node *Command, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   command.Table,
			Columns: command.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: command.FieldID,
			},
		},
	}
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Command.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, command.FieldID)
		for _, f := range fields {
			if !command.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != command.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: command.FieldName,
		})
	}
	if value, ok := cuo.mutation.Version(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: command.FieldVersion,
		})
	}
	if value, ok := cuo.mutation.Location(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: command.FieldLocation,
		})
	}
	if value, ok := cuo.mutation.Activated(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: command.FieldActivated,
		})
	}
	_node = &Command{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{command.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
