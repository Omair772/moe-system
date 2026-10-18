package usecase

import (
    "context"
    "time"

    "moe-system/ent"
    "moe-system/internal/repository"
)

type StudentUsecase interface {
    CreateStudent(ctx context.Context, name, email, grade, section, studentID string, userID int, createdAt time.Time) (*ent.Student, error)
    ListStudents(ctx context.Context) ([]*ent.Student, error)
    RegisterStudent(ctx context.Context, name, email, grade, section, studentID string, userID int, createdAt time.Time) (*ent.Student, error)
    GetStudent(ctx context.Context, id int) (*ent.Student, error)
    GetAllStudents(ctx context.Context) ([]*ent.Student, error)
}

type studentUsecase struct {
    studentRepo repository.StudentRepository
    userRepo    repository.UserRepository
}

func NewStudentUsecase(sr repository.StudentRepository, ur repository.UserRepository) StudentUsecase {
    return &studentUsecase{
        studentRepo: sr,
        userRepo:    ur,
    }
}

func (u *studentUsecase) CreateStudent(ctx context.Context, name, email, grade, section, studentID string, userID int, createdAt time.Time) (*ent.Student, error) {
    return u.studentRepo.CreateStudent(ctx, name, email, grade, section, studentID, userID, createdAt)
}

func (u *studentUsecase) RegisterStudent(ctx context.Context, name, email, grade, section, studentID string, userID int, createdAt time.Time) (*ent.Student, error) {
    _, _ = u.userRepo.CreateUser(ctx, email, "12345678", name, "student")
    return u.studentRepo.CreateStudent(ctx, name, email, grade, section, studentID, userID, createdAt)
}

func (u *studentUsecase) ListStudents(ctx context.Context) ([]*ent.Student, error) {
    return u.studentRepo.ListStudents(ctx)
}

func (u *studentUsecase) GetAllStudents(ctx context.Context) ([]*ent.Student, error) {
    return u.ListStudents(ctx)
}

func (u *studentUsecase) GetStudent(ctx context.Context, id int) (*ent.Student, error) {
    return nil, nil
}