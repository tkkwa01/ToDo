package domain

import (
	"ToDo/domain/vobj"
	"ToDo/packages/context"
	"ToDo/resource/request"
)

type User struct {
	ID       uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserName string        `json:"username" validate:"required" gorm:"type:varchar(255);not null"`
	Password vobj.Password `json:"password" gorm:"type:varchar(255);not null"`
}

func NewUser(ctx context.Context, dto *request.CreateUser) (*User, error) {
	var user = User{
		UserName: dto.UserName,
	}

	ctx.Validate(user)

	password, err := vobj.NewPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user.Password = *password

	return &user, nil
}
