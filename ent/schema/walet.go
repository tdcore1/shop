package schema

import "entgo.io/ent"
import 	"entgo.io/ent/schema/field"
import "time"


// Walet holds the schema definition for the Walet entity.
type Walet struct {
	ent.Schema
}

// Fields of the Walet.
func (Walet) Fields() []ent.Field {
	return []ent.Field{
		field.Int("amount"),
		field.Int("user_id"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Walet.
func (Walet) Edges() []ent.Edge {
	return nil
}
