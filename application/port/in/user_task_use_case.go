package in

import "temporal-getting-started/domain"

type UserTaskUseCase interface {
	GetData(id domain.UserTaskId) (domain.UserTask, error)
	Complete(id domain.UserTaskId, decision bool) error
}
