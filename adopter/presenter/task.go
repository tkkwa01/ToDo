package presenter

import (
	"ToDo/domain"
	"ToDo/usecase"
	"github.com/gin-gonic/gin"
)

type task struct {
	c *gin.Context
}

type TaskOutputFactory func(c *gin.Context) usecase.TaskOutputPort

func NewTaskOutputFactory() TaskOutputFactory {
	return func(c *gin.Context) usecase.TaskOutputPort {
		return &task{c: c}
	}
}

func (t task) CreateTask(id uint) error {
	t.c.JSON(200, gin.H{"id": id})
	return nil
}

func (t task) GetAllTask(res []domain.Task) error {
	t.c.JSON(200, res)
	return nil
}

func (t task) GetTaskByID(res domain.Task) error {
	t.c.JSON(200, res)
	return nil
}

func (t task) UpdateTask(res domain.Task) error {
	t.c.JSON(200, res)
	return nil
}

func (t task) DeleteTask() error {
	t.c.JSON(200, gin.H{"message": "task delete successfully"})
	return nil
}

func (t task) GetMyTask(res *domain.Task) error {
	t.c.JSON(200, res)
	return nil
}
