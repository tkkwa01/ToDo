package api

import (
	"ToDo/packages/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Execute() {
	logger := log.Logger()
	defer logger.Sync()

	engine := gin.New()

	engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })
}
