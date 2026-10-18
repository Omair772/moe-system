package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Student represents the students table.
type Student struct {
	ent.Schema
}

// Fields of the Student.
func (Student) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Comment("Auto-increment ID"),
		field.String("student_id").NotEmpty().Comment("Student ID"),
		field.String("first_name").NotEmpty().Comment("First name"),
		field.String("last_name").NotEmpty().Comment("Last name"),
		field.String("email").NotEmpty().Comment("Student email"),
		field.String("phone").Comment("Phone number"),
		field.Int("grade_level").Comment("Grade level"),
		field.Bool("is_active").Default(true).Comment("Student active status"),
		field.Time("enrollment_date").Comment("Enrollment date"),
		field.Time("created_at").Default(time.Now).Comment("Creation timestamp"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Last update"),
	}
}

// Edges of the Student.
func (Student) Edges() []ent.Edge {
	return []ent.Edge{
		// Relation to User (one-to-one)
		// edge.To("user", schema.User.Type),
	}
}