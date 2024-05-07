package http

import (
	"ToDo/adopter/presenter"
	"ToDo/packages/context"
	"ToDo/packages/http/middleware"
	"ToDo/packages/http/router"
	"ToDo/resource/request"
	"ToDo/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type task struct {
	inputFactory  usecase.TaskInputFactory
	outputFactory func(c *gin.Context) usecase.TaskOutputPort
	taskRepo      usecase.TaskRepository
}

func NewTask(r *router.Router, inputFactory usecase.TaskInputFactory, outputFactory presenter.TaskOutputFactory) {
	handler := task{
		inputFactory:  inputFactory,
		outputFactory: outputFactory,
	}

	r.Get("tasks", handler.GetAll)

	r.Group("", []gin.HandlerFunc{middleware.Auth(true, true)}, func(r *router.Router) {
		r.Group("tasks", nil, func(r *router.Router) {
			r.Post("", handler.Create)
			r.Put("/:taskId", handler.Update)
			r.Delete("/:taskId", handler.Delete)
			r.Get("/myTasks", handler.GetMy)
		})
	})
}

func (t task) Create(ctx context.Context, c *gin.Context) error {
	var req request.CreateTask

	if !bind(c, &req) {
		return nil
	}

	outputPort := t.outputFactory(c)
	inputPort := t.inputFactory(outputPort)

	return inputPort.CreateTask(ctx, &req)
}

func (t task) GetAll(ctx context.Context, c *gin.Context) error {
	outputPort := t.outputFactory(c)
	inputPort := t.inputFactory(outputPort)

	return inputPort.GetAllTask(ctx)
}

func (t task) Update(ctx context.Context, c *gin.Context) error {
	var req request.UpdateTask

	id, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return err
	}
	req.ID = uint(id)

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return err
	}

	outputPort := t.outputFactory(c)
	inputPort := t.inputFactory(outputPort)

	return inputPort.UpdateTask(ctx, &req)
}

func (t task) Delete(ctx context.Context, c *gin.Context) error {
	id, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return err
	}

	outputPort := t.outputFactory(c)
	inputPort := t.inputFactory(outputPort)

	return inputPort.DeleteTask(ctx, uint(id))
}

func (t task) GetMy(ctx context.Context, c *gin.Context) error {
	outputPort := t.outputFactory(c)
	inputPort := t.inputFactory(outputPort)

	return inputPort.GetMyTask(ctx)
}
