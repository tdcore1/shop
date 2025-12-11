package schema

import "entgo.io/ent"
import 	"entgo.io/ent/schema/field"
import "time"

// Cart holds the schema definition for the Cart entity.
type Cart struct {
	ent.Schema
}

// Fields of the Cart.
func (Cart) Fields() []ent.Field {
	return  []ent.Field{
		field.Int("product_id"),
		field.Int("count"),
		field.Int("user_id"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Cart.
func (Cart) Edges() []ent.Edge {
	return nil
}
