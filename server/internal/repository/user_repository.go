package repository

import (
	"context"
	"eduhub/server/internal/models"

	"github.com/uptrace/bun"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindByRollNo(ctx context.Context, RollNo string) *models.User
	UpdateUser(ctx context.Context, user *models.User) *models.User
	FreezeUser(ctx context.Context, RollNo string) error
	DeleteUser(ctx context.Context, RollNo string) error
}

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db bun.DB) UserRepository {
	return &userRepository{
		db: &db,
	}
}
