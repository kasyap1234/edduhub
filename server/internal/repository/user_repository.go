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

func(u*userRepository)CreateUser(ctx context.Context,user*models.User)error {
	userModel :=new(models.User)
	_,err :=u.db.NewInsert().Model(userModel).Exec(ctx)
	return err 
}

func(u*userRepository)FindByRollNo(ctx context.Context,RollNo string)*models.User{
	user :=new(models.User)
	err :=u.db.NewSelect().Model(user).Where("roll_no = ?",RollNo).Scan(ctx)
	if err !=nil {
		return nil 
	}
	return user 
}

