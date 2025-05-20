package in

import "temporal-getting-started/domain"

type VerifyUserUseCase interface {
	Verify(userId string) (*domain.User, error)
}
