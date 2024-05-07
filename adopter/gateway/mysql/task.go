package mysql

import (
	"ToDo/domain"
	"ToDo/packages/context"
	"ToDo/usecase"
)

type task struct{}

func NewTaskRepository() usecase.TaskRepository {
	return &task{}
}

func (t task) CreateTask(ctx context.Context, task *domain.Task) (uint, error) {
	db := ctx.DB()

	if err := db.Create(task).Error; err != nil {
		return 0, dbError(err)
	}
	return task.ID, nil
}

func (t task) GetAllTask(ctx context.Context) ([]domain.Task, error) {
	db := ctx.DB()

	var tasks []domain.Task
	if err := db.Find(&tasks).Error; err != nil {
		return nil, dbError(err)
	}
	return tasks, nil
}

func (t task) GetTaskByID(ctx context.Context, id uint) (*domain.Task, error) {
	db := ctx.DB()

	var task domain.Task
	err := db.Where(&domain.Task{ID: id}).First(&task).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &task, nil
}

func (t task) UpdateTask(ctx context.Context, task *domain.Task) error {
	db := ctx.DB()

	if err := db.Save(task).Error; err != nil {
		return dbError(err)
	}
	return nil
}

func (t task) DeleteTask(ctx context.Context, id uint) error {
	db := ctx.DB()

	var task domain.Task
	res := db.Where("id = ?", id).Delete(&task)
	if res.Error != nil {
		return dbError(res.Error)
	}
	return nil
}

func (t task) GetMyTask(ctx context.Context, id uint) error {
	db := ctx.DB()

	var task domain.Task
	err := db.Where(&domain.Task{ID: id}).First(&task).Error
	if err != nil {
		return dbError(err)
	}
	return nil
}
