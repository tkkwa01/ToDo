package usecase

import (
	"ToDo/config"
	"ToDo/domain"
	"ToDo/packages/context"
	"ToDo/packages/errors"
	"ToDo/resource/request"
	"ToDo/resource/response"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type UserInputPort interface {
	Create(ctx context.Context, req *request.CreateUser) error
	GetByID(ctx context.Context, id uint) error
	Login(ctx context.Context, req *request.UserLogin) error
}

type UserOutputPort interface {
	Create(id uint) error
	GetByID(res *domain.User) error
	Login(isSession bool, res *response.UserLogin) error
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (uint, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByUserName(ctx context.Context, username string) (*domain.User, error)
}

type user struct {
	outputPort UserOutputPort
	userRepo   UserRepository
}

type UserInputFactory func(outputPort UserOutputPort) UserInputPort

func NewUserInputFactory(ur UserRepository) UserInputFactory {
	return func(o UserOutputPort) UserInputPort {
		return &user{
			outputPort: o,
			userRepo:   ur,
		}
	}
}

func (u user) Create(ctx context.Context, req *request.CreateUser) error {
	newUser, err := domain.NewUser(ctx, req)
	if err != nil {
		return err
	}

	if ctx.IsInValid() {
		return ctx.ValidationError()
	}

	id, err := u.userRepo.Create(ctx, newUser)
	if err != nil {
		return err
	}

	return u.outputPort.Create(id)
}

func (u user) GetByID(ctx context.Context, id uint) error {
	res, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return u.outputPort.GetByID(res)
}

func issueJwtToken(userID, userName, realm, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":       userID,
		"user_name": userName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
		"realm":     realm,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u user) Login(ctx context.Context, req *request.UserLogin) error {
	user, err := u.userRepo.GetByUserName(ctx, req.UserName)
	if err != nil {
		return err
	}

	if user.Password.IsValid(req.Password) {
		token, err := issueJwtToken(strconv.Itoa(int(user.ID)), user.UserName, "user", config.Env.App.Secret)
		if err != nil {
			return errors.NewUnexpected(err)
		}

		var res response.UserLogin
		res.Token = token

		return u.outputPort.Login(req.Session, &res)
	} else {
		fmt.Println("Invalid password")
	}
	return u.outputPort.Login(req.Session, nil)
}
