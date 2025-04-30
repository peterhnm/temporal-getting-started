package out

import "temporal-getting-started/domain"

type UserRepository interface {
	FindById(id string) (domain.User, error)
}
