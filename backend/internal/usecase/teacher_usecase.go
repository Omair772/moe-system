package usecase

import (
    "context"
    "time"

    "moe-system/ent"
    "moe-system/internal/repository"
)

type TeacherUsecase interface {
    CreateTeacher(ctx context.Context, name, email, subject, phone, specialization, employeeID, userID string, createdAt time.Time) (*ent.Teacher, error)
    ListTeachers(ctx context.Context) ([]*ent.Teacher, error)
    RegisterTeacher(ctx context.Context, name, email, subject, phone, specialization, employeeID, userID string, createdAt time.Time) (*ent.Teacher, error)
    GetTeacher(ctx context.Context, id int) (*ent.Teacher, error)
    GetAllTeachers(ctx context.Context) ([]*ent.Teacher, error)
}

type teacherUsecase struct {
    teacherRepo repository.TeacherRepository
    userRepo    repository.UserRepository
}

func NewTeacherUsecase(tr repository.TeacherRepository, ur repository.UserRepository) TeacherUsecase {
    return &teacherUsecase{
        teacherRepo: tr,
        userRepo:    ur,
    }
}

func (u *teacherUsecase) CreateTeacher(ctx context.Context, name, email, subject, phone, specialization, employeeID, userID string, createdAt time.Time) (*ent.Teacher, error) {
    return u.teacherRepo.CreateTeacher(ctx, name, email, subject, phone, specialization, employeeID, userID, createdAt)
}

func (u *teacherUsecase) RegisterTeacher(ctx context.Context, name, email, subject, phone, specialization, employeeID, userID string, createdAt time.Time) (*ent.Teacher, error) {
    _, err := u.userRepo.CreateUser(ctx, email, "12345678", name, "teacher")
    if err != nil {
        // If user exists, proceed
    }
    return u.teacherRepo.CreateTeacher(ctx, name, email, subject, phone, specialization, employeeID, userID, createdAt)
}

func (u *teacherUsecase) ListTeachers(ctx context.Context) ([]*ent.Teacher, error) {
    return u.teacherRepo.ListTeachers(ctx)
}

func (u *teacherUsecase) GetAllTeachers(ctx context.Context) ([]*ent.Teacher, error) {
    return u.ListTeachers(ctx)
}

func (u *teacherUsecase) GetTeacher(ctx context.Context, id int) (*ent.Teacher, error) {
    return nil, nil
}