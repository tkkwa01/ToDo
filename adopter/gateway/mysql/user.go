package mysql

import (
	"ToDo/domain"
	"ToDo/packages/context"
	"ToDo/packages/errors"
	"ToDo/usecase"
	"gorm.io/gorm"
)

type User struct{}

func NewUserRepository() usecase.UserRepository {
	return &User{}
}

func dbError(err error) error {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return errors.NotFound()
	default:
		return errors.NewUnexpected(err)
	}
}

func (u User) Create(ctx context.Context, user *domain.User) (uint, error) {
	db := ctx.DB()

	if err := db.Create(user).Error; err != nil {
		return 0, dbError(err)
	}
	return user.ID, nil
}

func (u User) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	db := ctx.DB()

	var user domain.User
	err := db.Where(&domain.User{ID: id}).First(&user).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &user, nil
}

func (u User) GetByUserName(ctx context.Context, username string) (*domain.User, error) {
	db := ctx.DB()

	var user domain.User
	err := db.Where(&domain.User{UserName: username}).First(&user).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &user, nil
}
