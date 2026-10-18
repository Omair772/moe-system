package schema

import (
	"time" // 1. أضفنا استيراد حزمة الوقت

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User represents the users table.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Comment("Auto-increment ID"),
		field.String("email").NotEmpty().Comment("User email"),
		field.String("password").NotEmpty().Comment("Hashed password"),
		field.String("full_name").Comment("Full name"),
		field.String("role").Comment("User role (admin, teacher, student)"),
		field.Bool("is_active").Default(true).Comment("Account active status"),
		field.Time("created_at").Default(time.Now).Comment("Creation timestamp"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Last update"), // 2. تم تصحيحها إلى UpdateDefault
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// Add relations as needed
	}
}