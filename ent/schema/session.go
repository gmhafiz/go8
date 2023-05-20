package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Author.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("token"),
		field.Uint64("user_id").Nillable().Optional(),
		field.Bytes("data"),
		field.Time("expiry"),
	}
}
