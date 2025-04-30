package in

import "temporal-getting-started/domain"

type VerifyUserUseCase interface {
	Verify(user domain.User) error
}
