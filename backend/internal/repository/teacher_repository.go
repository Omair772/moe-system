package repository

import (
	"context"
	"time"

	"moe-system/ent"
)

type TeacherRepository interface {
	CreateTeacher(ctx context.Context, teacherID, firstName, lastName, email, subject, department, phone string, hireDate time.Time) (*ent.Teacher, error)
	GetTeacherByID(ctx context.Context, id int) (*ent.Teacher, error)
	ListTeachers(ctx context.Context) ([]*ent.Teacher, error)
}

type teacherRepository struct {
	client *ent.Client
}

func NewTeacherRepository(client *ent.Client) TeacherRepository {
	return &teacherRepository{client: client}
}

func (r *teacherRepository) CreateTeacher(ctx context.Context, teacherID, firstName, lastName, email, subject, department, phone string, hireDate time.Time) (*ent.Teacher, error) {
	return r.client.Teacher.
		Create().
		SetTeacherID(teacherID).
		SetFirstName(firstName).
		SetLastName(lastName).
		SetEmail(email).
		SetSubject(subject).
		SetDepartment(department).
		SetPhone(phone).
		SetHireDate(hireDate).
		Save(ctx)
}

func (r *teacherRepository) GetTeacherByID(ctx context.Context, id int) (*ent.Teacher, error) {
	return r.client.Teacher.Get(ctx, id)
}

func (r *teacherRepository) ListTeachers(ctx context.Context) ([]*ent.Teacher, error) {
	return r.client.Teacher.Query().All(ctx)
}