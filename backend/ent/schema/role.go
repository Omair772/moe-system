package schema

import (
	"time" // 1. أضفنا استيراد حزمة الوقت هنا

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Role represents the roles table.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Comment("Auto-increment ID"),
		field.String("name").NotEmpty().Comment("Role name (admin, teacher, student)"),
		field.String("description").Comment("Role description"),
		field.Bool("is_active").Default(true).Comment("Role active status"),
		field.Time("created_at").Default(time.Now).Comment("Creation timestamp"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Last update"), // 2. تم تصحيحها إلى UpdateDefault
	}
}