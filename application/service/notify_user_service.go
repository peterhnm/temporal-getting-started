package service

import (
	"temporal-getting-started/application/port/in"
	"temporal-getting-started/application/port/out"
	"temporal-getting-started/domain"
)

type NotifyUserService struct {
	userTaskRepo out.UserTaskRepository
}

func NewNotifyUserService(userTaskRepo out.UserTaskRepository) in.NotifyUserUseCase {
	return &NotifyUserService{userTaskRepo}
}

func (n *NotifyUserService) Add(cmd in.NotifyUserCommand) error {
	userTask := domain.UserTask{
		Id: domain.UserTaskId{
			WorkflowId: cmd.WorkflowId,
			RunId:      cmd.RunId,
		},
		Data: domain.UserTaskData{
			UserId:      cmd.UserId,
			RecipientId: cmd.RecipientId,
			Amount:      cmd.Amount,
			Decision:    nil,
		},
	}

	return n.userTaskRepo.Add(userTask)
}
