package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"moe-system/ent"
	"moe-system/ent/user"
	"moe-system/pkg/auth"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Login(ctx context.Context, email, password string) (string, string, error)
}

type authUseCase struct {
	client      *ent.Client
	redisClient *redis.Client
}

func NewAuthUseCase(client *ent.Client, redisClient *redis.Client) AuthUseCase {
	return &authUseCase{
		client:      client,
		redisClient: redisClient,
	}
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (string, string, error) {
	// 1. البحث عن المستخدم بواسطة البريد الإلكتروني في قاعدة البيانات
	u, err := uc.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// 2. التحقق من تطابق كلمة المرور المشفرة (Bcrypt)
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// 3. توليد الـ Tokens (Access Token صالح 15 دقيقة & Refresh Token صالح 7 أيام)
	accessToken, refreshToken, err := auth.GenerateTokens(u.ID, u.Role)
	if err != nil {
		return "", "", err
	}

	// 4. حفظ الـ Refresh Token في Redis لربط الجلسة
	redisKey := fmt.Sprintf("refresh_token:%d", u.ID)
	err = uc.redisClient.Set(ctx, redisKey, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		return "", "", errors.New("failed to store session in redis")
	}

	return accessToken, refreshToken, nil
}