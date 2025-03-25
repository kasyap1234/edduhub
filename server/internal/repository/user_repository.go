package repository

import (
	"context"
	"eduhub/server/internal/models"

	"github.com/uptrace/bun"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindByRollNo(ctx context.Context, RollNo string) (*models.User,error)
	UpdateUser(ctx context.Context, user *models.User)( *models.User,error)
	FreezeUser(ctx context.Context, RollNo string) error
	DeleteUser(ctx context.Context, RollNo string) error
}

type userRepository struct {
	db  DatabaseRepository[models.User]
}

func NewUserRepository(db DatabaseRepository[models.User]) UserRepository {
	return &userRepository{
		db: db ,
	}
}

func(u*userRepository)CreateUser(ctx context.Context,user*models.User)error {
	return u.db.Create(ctx,user)
}

func(u*userRepository)FindByRollNo(ctx context.Context,RollNo string)(*models.User,error){
	return u.db.FindOne(ctx,"roll_no=?",RollNo)

}

func(u*userRepository)DeleteUser(ctx context.Context,RollNo string)error{

}

func(u*userRepository)FreezeUser(ctx context.Context,RollNo string)error {

}

func(u*userRepository)UpdateUser(ctx context.Context,model*models.User)(*models.User,error){

}