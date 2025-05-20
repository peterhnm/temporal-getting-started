package service

import (
	"fmt"
	"log"
	"temporal-getting-started/application/port/in"
	"temporal-getting-started/application/port/out"
	"temporal-getting-started/domain"
)

// VerificationService The activity object is the struct
// that maintains shared state across Activities.
// If the Worker crashes, this Activity object loses its state.
type VerificationService struct {
	userRepo out.UserRepository
}

// NewVerificationService The factory method to wire the implementation to the interface.
func NewVerificationService(userRepo out.UserRepository) in.VerifyUserUseCase {
	return &VerificationService{userRepo}
}

func (svc *VerificationService) Verify(userId string) (*domain.User, error) {
	log.Println(fmt.Sprintf("[USER] Verifying user %s", userId))
	user, err := svc.userRepo.FindById(userId)
	if err != nil {
		return &domain.User{}, err
	}
	return &user, nil
}
