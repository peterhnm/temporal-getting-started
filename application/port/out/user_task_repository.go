package out

import "temporal-getting-started/domain"

type UserTaskRepository interface {
	FindById(id domain.UserTaskId) (domain.UserTask, error)
	Add(task domain.UserTask) error
	Update(id domain.UserTaskId, task domain.UserTask) error
}
