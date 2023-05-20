package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the Author.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("first_name").Optional(),
		field.String("middle_name").Optional(),
		field.String("last_name").Optional(),
		field.String("email"),
		field.String("password"),
		field.Time("verified_at").Optional().Nillable().StructTag(`json:"-"`),
	}
}
