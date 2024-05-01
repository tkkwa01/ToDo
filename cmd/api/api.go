package api

import (
	httpController "ToDo/adopter/controller/http"
	mysqlRepository "ToDo/adopter/gateway/mysql"
	userPresenter "ToDo/adopter/presenter"
	"ToDo/config"
	"ToDo/driver"
	"ToDo/packages/http/middleware"
	"ToDo/packages/http/router"
	"ToDo/packages/log"
	userUsecase "ToDo/usecase"
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Execute() {
	logger := log.Logger()
	defer logger.Sync()

	engine := gin.New()

	engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// cors
	engine.Use(middleware.Cors(nil))

	store := cookie.NewStore([]byte("secret"))
	engine.Use(sessions.Sessions("user", store))

	r := router.New(engine, driver.GetRDB)

	//mysql
	userRepository := mysqlRepository.NewUserRepository()

	//usecase
	userInputFactory := userUsecase.NewUserInputFactory(userRepository)
	userOutputFactory := userPresenter.NewUserOutputFactory()

	//controller
	httpController.NewUser(r, userInputFactory, userOutputFactory)

	//serve
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Env.Port),
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	logger.Info("Succeeded in listen and serve.")
	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Server forced to shutdown: %+v", err))
	}

	logger.Info("Server existing")
}
