package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Immutable().Default(uuid.New),
		field.String("title").NotEmpty().MaxLen(50),
		field.String("description").NotEmpty().MaxLen(500),
		field.String("status").NotEmpty().GoType(TodoStatus("")),
		field.Bool("completed").Default(false),
		field.Time("completed_at").Optional().Nillable(),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return nil
}

// TodoStatus is used to enforce strong typing on status field.
type TodoStatus string

const (
	StatusNone      TodoStatus = "none"
	StatusProgress  TodoStatus = "progress"
	StatusPending   TodoStatus = "pending"
	StatusCompleted TodoStatus = "completed"
)
