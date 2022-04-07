package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Book holds the schema definition for the Book entity.
type Book struct {
	ent.Schema
}

// Fields of the Book.
func (Book) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("title"),
		field.Time("published_date"),
		field.String("image_url").Optional(),
		field.String("description").Sensitive(),
		field.Time("created_at").Default(time.Now).Immutable().StructTag(`json:"-"`),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Immutable().StructTag(`json:"-"`),
		field.Time("deleted_at").Optional().Nillable().StructTag(`json:"-"`),
	}
}

// Edges of the Book.
func (Book) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("authors", Author.Type),
	}
}
