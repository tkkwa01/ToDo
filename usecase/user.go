package usecase

import (
	"ToDo/config"
	"ToDo/domain"
	"ToDo/packages/auth"
	"ToDo/packages/context"
	"ToDo/resource/request"
	"ToDo/resource/response"
	"strconv"
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

func (u user) Login(ctx context.Context, req *request.UserLogin) error {
	user, err := u.userRepo.GetByUserName(ctx, req.UserName)
	if err != nil {
		return err
	}

	if !user.Password.IsValid(req.Password) {
		ctx.FieldError("Password", "パスワードが違います")
	}

	tokenService := auth.NewJwtTokenService(config.Env.App.Secret)
	token, err := tokenService.IssueJWT(strconv.Itoa(int(user.ID)), user.UserName, "user")
	if err != nil {
		return err
	}
	var res response.UserLogin
	res.Token = token

	return u.outputPort.Login(true, &res)
}
