package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Command holds the schema definition for the Command entity.
type Command struct {
	ent.Schema
}

// Fields of the Command.
func (Command) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.String("name").MaxLen(128).NotEmpty(),
		field.String("version").MaxLen(128).NotEmpty(),
		field.String("location").MaxLen(512).NotEmpty(),
		field.Bool("activated").Default(false),
	}
}

func (Command) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "version").Unique(),
		index.Fields("activated"),
	}
}

// Edges of the Command.
func (Command) Edges() []ent.Edge {
	return nil
}
