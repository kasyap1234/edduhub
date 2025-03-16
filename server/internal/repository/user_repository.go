package repository
import (
	"gorm.io/gorm"
)
type UserRepository interface{
FindByID(ID int)(*User)
FindBy
}


type userRepository struct {
	db *gorm.DB 
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db : db }
}
