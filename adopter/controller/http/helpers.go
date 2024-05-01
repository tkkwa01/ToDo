package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func bind(c *gin.Context, request interface{}) (ok bool) {
	if err := c.BindJSON(request); err != nil {
		c.Status(http.StatusBadRequest)
		return false
	} else {
		return true
	}
}
