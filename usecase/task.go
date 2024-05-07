package usecase

import (
	"ToDo/domain"
	"ToDo/packages/context"
	"ToDo/resource/request"
)

type TaskInputPort interface {
	CreateTask(ctx context.Context, req *request.CreateTask) error
	GetAllTask(ctx context.Context) error
	GetTaskByID(ctx context.Context, id uint) error
	UpdateTask(ctx context.Context, req *request.UpdateTask) error
	DeleteTask(ctx context.Context, id uint) error
	GetMyTask(ctx context.Context) error
}

type TaskOutputPort interface {
	CreateTask(id uint) error
	GetAllTask(res []domain.Task) error
	GetTaskByID(res domain.Task) error
	UpdateTask(res domain.Task) error
	DeleteTask() error
	GetMyTask(res *domain.Task) error
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *domain.Task) (uint, error)
	GetAllTask(ctx context.Context) ([]domain.Task, error)
	GetTaskByID(ctx context.Context, id uint) (*domain.Task, error)
	UpdateTask(ctx context.Context, task *domain.Task) error
	DeleteTask(ctx context.Context, id uint) error
	GetMyTask(ctx context.Context, id uint) error
}

type task struct {
	outputPort TaskOutputPort
	taskRepo   TaskRepository
}

type TaskInputFactory func(outputPort TaskOutputPort) TaskInputPort

func NewTaskInputFactory(tr TaskRepository) TaskInputFactory {
	return func(o TaskOutputPort) TaskInputPort {
		return &task{
			outputPort: o,
			taskRepo:   tr,
		}
	}
}

func (t task) CreateTask(ctx context.Context, req *request.CreateTask) error {
	newTask := domain.NewTask(req)

	if ctx.IsInValid() {
		return ctx.ValidationError()
	}

	id, err := t.taskRepo.CreateTask(ctx, newTask)
	if err != nil {
		return err
	}

	return t.outputPort.CreateTask(id)
}

func (t task) GetAllTask(ctx context.Context) error {
	tasks, err := t.taskRepo.GetAllTask(ctx)
	if err != nil {
		return err
	}

	return t.outputPort.GetAllTask(tasks)
}

func (t task) GetTaskByID(ctx context.Context, id uint) error {
	task, err := t.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}

	return t.outputPort.GetTaskByID(*task)
}

func (t task) UpdateTask(ctx context.Context, req *request.UpdateTask) error {
	task, err := t.taskRepo.GetTaskByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if req.Title != "" {
		task.Title = req.Title
	}

	if req.Description != "" {
		task.Description = req.Description
	}

	err = t.taskRepo.UpdateTask(ctx, task)
	if err != nil {
		return err
	}

	return t.outputPort.UpdateTask(*task)
}

func (t task) DeleteTask(ctx context.Context, id uint) error {
	err := t.taskRepo.DeleteTask(ctx, id)
	if err != nil {
		return err
	}

	return t.outputPort.DeleteTask()
}

func (t task) GetMyTask(ctx context.Context) error {
	currentUserID := ctx.UID()

	task, err := t.taskRepo.GetTaskByID(ctx, currentUserID)
	if err != nil {
		return err
	}

	return t.outputPort.GetMyTask(task)
}
