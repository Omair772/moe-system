package schema

import (
	"time" // 1. أضفنا استيراد حزمة الوقت

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Teacher represents the teachers table.
type Teacher struct {
	ent.Schema
}

// Fields of the Teacher.
func (Teacher) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Comment("Auto-increment ID"),
		field.String("teacher_id").NotEmpty().Comment("Teacher ID"),
		field.String("first_name").NotEmpty().Comment("First name"),
		field.String("last_name").NotEmpty().Comment("Last name"),
		field.String("email").NotEmpty().Comment("Teacher email"),
		field.String("subject").Comment("Subject taught"),
		field.String("department").Comment("Department"),
		field.String("phone").Comment("Phone number"),
		field.Bool("is_active").Default(true).Comment("Account active status"),
		field.Time("hire_date").Comment("Hire date"),
		field.Time("created_at").Default(time.Now).Comment("Creation timestamp"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Last update"), // 2. تصحيح دالة التحديث
	}
}

// Edges of the Teacher.
func (Teacher) Edges() []ent.Edge {
	return []ent.Edge{
		// Relations can be added here
	} // 3. إضافة الـ return وإغلاق المصفوفة بشكل صحيح
}