package db

import (
	"fmt"
	"temporal-getting-started/application/port/out"
	"temporal-getting-started/domain"
)

type UserTaskRepository struct {
	tasks map[domain.UserTaskId]domain.UserTask
}

func NewUserTaskRepository() out.UserTaskRepository {
	return &UserTaskRepository{
		tasks: make(map[domain.UserTaskId]domain.UserTask),
	}
}

func (repo *UserTaskRepository) FindById(id domain.UserTaskId) (domain.UserTask, error) {
	if task, exists := repo.tasks[id]; exists {
		return task, nil
	}
	return domain.UserTask{}, fmt.Errorf("user task with id %s not found", id)
}

func (repo *UserTaskRepository) Add(task domain.UserTask) error {
	if _, exists := repo.tasks[task.Id]; exists {
		return fmt.Errorf("user task with id %s already exists", task.Id)
	}
	repo.tasks[task.Id] = task
	return nil
}

func (repo *UserTaskRepository) Update(id domain.UserTaskId, task domain.UserTask) error {
	if _, exists := repo.tasks[id]; !exists {
		return fmt.Errorf("user task with id %s not found", id)
	}
	repo.tasks[id] = task
	return nil
}
