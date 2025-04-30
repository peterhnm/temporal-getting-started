package db

import (
	"fmt"
	"log"
	"temporal-getting-started/application/port/out"
	"temporal-getting-started/domain"
)

type UserRepository struct {
	users map[string]domain.User
}

func NewUserRepository() out.UserRepository {
	return &UserRepository{
		users: map[string]domain.User{
			"1234": {"1234", "12-3456"},
			"5678": {"5678", "56-7890"},
		},
	}
}

type InvalidUserError struct {
	id string
}

func (e *InvalidUserError) Error() string {
	return fmt.Sprintf("Invalid user id: %s", e.id)
}

func (repo *UserRepository) FindById(id string) (domain.User, error) {
	if user, exists := repo.users[id]; exists {
		log.Printf("Found user with id %s", id)
		return user, nil
	}
	return domain.User{}, &InvalidUserError{id}
}
