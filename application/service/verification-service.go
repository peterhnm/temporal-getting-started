package service

import (
	"temporal-getting-started/application/port/in"
	"temporal-getting-started/application/port/out"
	"temporal-getting-started/domain"
)

type VerificationService struct {
	userRepo out.UserRepository
}

func NewVerificationService(userRepo out.UserRepository) in.VerifyUserUseCase {
	return &VerificationService{userRepo}
}

func (svc *VerificationService) Verify(user domain.User) error {
	_, err := svc.userRepo.FindById(user.Id)
	if err != nil {
		return err
	}
	return nil
}
