package service

import (
	"temporal-getting-started/application/port/in"
	"temporal-getting-started/application/port/out"
	"temporal-getting-started/domain"
)

type UserTaskService struct {
	userTaskRepo out.UserTaskRepository
}

func NewUserTaskService(userTaskRepo out.UserTaskRepository) in.UserTaskUseCase {
	return &UserTaskService{userTaskRepo}
}

func (u *UserTaskService) GetData(id domain.UserTaskId) (domain.UserTask, error) {
	return u.userTaskRepo.FindById(id)
}

func (u *UserTaskService) Complete(id domain.UserTaskId, decision bool) error {
	userTask, err := u.GetData(id)
	if err != nil {
		return err
	}

	userTask.Data.Decision = &decision
	if err := u.userTaskRepo.Update(id, userTask); err != nil {
		return err
	}

	return nil
}
