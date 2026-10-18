package repository

import (
    "context"

    "moe-system/ent"
)

type UserRepository interface {
    CreateUser(ctx context.Context, fullName, email, password, role string) (*ent.User, error)
    GetUserByID(ctx context.Context, id int) (*ent.User, error)
    ListUsers(ctx context.Context) ([]*ent.User, error)
}

type userRepository struct {
    client *ent.Client
}

func NewUserRepository(client *ent.Client) UserRepository {
    return &userRepository{client: client}
}

func (r *userRepository) CreateUser(ctx context.Context, fullName, email, password, role string) (*ent.User, error) {
    return r.client.User.
        Create().
        SetFullName(fullName).
        SetEmail(email).
        SetPassword(password).
        SetRole(role).
        Save(ctx)
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (*ent.User, error) {
    return r.client.User.Get(ctx, id)
}

func (r *userRepository) ListUsers(ctx context.Context) ([]*ent.User, error) {
    return r.client.User.Query().All(ctx)
}
