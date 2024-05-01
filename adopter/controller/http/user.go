package http

import (
	"ToDo/adopter/presenter"
	"ToDo/packages/context"
	"ToDo/packages/http/middleware"
	"ToDo/packages/http/router"
	"ToDo/resource/request"
	"ToDo/usecase"
	"github.com/gin-gonic/gin"
)

type user struct {
	inputFactory  usecase.UserInputFactory
	outputFactory func(c *gin.Context) usecase.UserOutputPort
	UserRepo      usecase.UserRepository
}

func NewUser(r *router.Router, inputFactory usecase.UserInputFactory, outputFactory presenter.UserOutputFactory) {
	handler := user{
		inputFactory:  inputFactory,
		outputFactory: outputFactory,
	}

	r.Group("user", nil, func(r *router.Router) {
		r.Post("login", handler.Login)
	})

	r.Group("", []gin.HandlerFunc{middleware.Auth(true, true)}, func(r *router.Router) {
		r.Group("user", nil, func(r *router.Router) {
			r.Get("me", handler.GetMe)
		})
	})
}

func (u user) Login(ctx context.Context, c *gin.Context) error {
	var req request.UserLogin

	if !bind(c, &req) {
		return nil
	}

	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.Login(ctx, &req)
}

func (u user) GetMe(ctx context.Context, c *gin.Context) error {
	outputPort := u.outputFactory(c)
	inputPort := u.inputFactory(outputPort)

	return inputPort.GetByID(ctx, ctx.UID())
}
