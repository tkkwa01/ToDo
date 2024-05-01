package presenter

import (
	"ToDo/domain"
	"ToDo/resource/response"
	"ToDo/usecase"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type user struct {
	c *gin.Context
}

type UserOutputFactory func(c *gin.Context) usecase.UserOutputPort

func NewUserOutputFactory() UserOutputFactory {
	return func(c *gin.Context) usecase.UserOutputPort {
		return &user{c: c}
	}
}

func (u user) Create(id uint) error {
	u.c.JSON(http.StatusCreated, gin.H{"id": id})
	return nil
}

func (u user) GetByID(res *domain.User) error {
	u.c.JSON(http.StatusOK, res)
	return nil
}

func (u user) Login(isSession bool, res *response.UserLogin) error {
	if res == nil {
		u.c.Status(http.StatusUnauthorized)
		return nil
	}

	if isSession {
		session := sessions.Default(u.c)
		session.Set("user", res.Token)
		if err := session.Save(); err != nil {
			return err
		}
		u.c.Status(http.StatusOK)
	} else {
		u.c.JSON(http.StatusOK, res)
	}
	return nil
}
