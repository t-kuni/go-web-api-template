package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("age"),
		field.String("name"),
		field.String("gender").
			Optional(),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("companies", Company.Type).
			StorageKey(
				edge.Table("employments"),
				edge.Columns("user_id", "company_id"),
			),
		edge.To("emails", Email.Type),
		//edge.To("employments", Employment.Type),
	}
}
