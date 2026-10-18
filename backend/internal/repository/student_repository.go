package repository

import (
	"context"
	"time"

	"moe-system/ent"
)

type StudentRepository interface {
	CreateStudent(ctx context.Context, studentID, firstName, lastName, email, phone string, gradeLevel int, enrollmentDate time.Time) (*ent.Student, error)
	GetStudentByID(ctx context.Context, id int) (*ent.Student, error)
	ListStudents(ctx context.Context) ([]*ent.Student, error)
}

type studentRepository struct {
	client *ent.Client
}

func NewStudentRepository(client *ent.Client) StudentRepository {
	return &studentRepository{client: client}
}

func (r *studentRepository) CreateStudent(ctx context.Context, studentID, firstName, lastName, email, phone string, gradeLevel int, enrollmentDate time.Time) (*ent.Student, error) {
	return r.client.Student.
		Create().
		SetStudentID(studentID).
		SetFirstName(firstName).
		SetLastName(lastName).
		SetEmail(email).
		SetPhone(phone).
		SetGradeLevel(gradeLevel).
		SetEnrollmentDate(enrollmentDate).
		Save(ctx)
}

func (r *studentRepository) GetStudentByID(ctx context.Context, id int) (*ent.Student, error) {
	return r.client.Student.Get(ctx, id)
}

func (r *studentRepository) ListStudents(ctx context.Context) ([]*ent.Student, error) {
	return r.client.Student.Query().All(ctx)
}