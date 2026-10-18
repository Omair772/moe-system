package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type HighSchoolResult struct {
	ent.Schema
}

// Fields of the HighSchoolResult.
func (HighSchoolResult) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Comment("Auto-increment ID"),
		field.String("seat_number").Comment("Seat number - MUST BE UNIQUE"),
		field.String("student_name").Comment("Student full name"),
		field.String("section").Comment("Section: SCIENTIFIC, LITERARY, etc."),
		field.String("governorate").Comment("Governorate/province"),
		field.Float("total_score").Comment("Total score (0-100 or 0-1000)"),
		field.Float("percentage").Comment("Percentage (0-100)"),
		field.JSON("grades_json", map[string]interface{}{}).Comment("JSON grades array"), // تم تصحيح الحقل بإضافة نوع البيانات الافتراضي
		field.String("academic_year").Comment("Academic year like '2024'"),
		field.Time("created_at").Default(time.Now).Comment("Creation timestamp"),
	}
}

// Indexes of the HighSchoolResult.
func (HighSchoolResult) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("seat_number").Unique(), // تم تصحيح طريقة كتابة الفهرس الفريد
		index.Fields("governorate"),          // تم تصحيح الفهرس العادي
		index.Fields("academic_year"),        // تم تصحيح الفهرس العادي
	}
}

// Edges of the HighSchoolResult.
func (HighSchoolResult) Edges() []ent.Edge {
	return []ent.Edge{
		// Relations to other schemas if needed
	}
}