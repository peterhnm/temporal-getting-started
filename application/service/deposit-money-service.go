package service

import (
	"temporal-getting-started/application/port/in"
	"temporal-getting-started/application/port/out"
	"temporal-getting-started/domain"
)

type DepositMoneyService struct {
	bankRepo out.BankRepository
}

func NewDepositMoneyService(bankRepo out.BankRepository) in.DepositMoneyUseCase {
	return &DepositMoneyService{bankRepo}
}

func (s *DepositMoneyService) Deposit(id string, amount float64) (domain.BankAccount, error) {
	account, err := s.bankRepo.Account(id)
	if err != nil {
		return domain.BankAccount{}, err
	}

	account.Balance += amount

	account, err = s.bankRepo.Update(id, account)
	if err != nil {
		return domain.BankAccount{}, err
	}

	return account, nil
}
