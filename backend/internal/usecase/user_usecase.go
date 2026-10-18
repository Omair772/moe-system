package usecase

import (
	"context"

	"moe-system/ent"
	"moe-system/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, fullName, email, password, role string) (*ent.User, error)
	GetUserByID(ctx context.Context, id int) (*ent.User, error)
	GetAllUsers(ctx context.Context) ([]*ent.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (uc *userUseCase) CreateUser(ctx context.Context, fullName, email, password, role string) (*ent.User, error) {
	// تشفير كلمة المرور باستخدام Bcrypt قبل الحفظ
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return uc.userRepo.CreateUser(ctx, fullName, email, string(hashedPassword), role)
}

func (uc *userUseCase) GetUserByID(ctx context.Context, id int) (*ent.User, error) {
	return uc.userRepo.GetUserByID(ctx, id)
}

func (uc *userUseCase) GetAllUsers(ctx context.Context) ([]*ent.User, error) {
	return uc.userRepo.ListUsers(ctx)
}