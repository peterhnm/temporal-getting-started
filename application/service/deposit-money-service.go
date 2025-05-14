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

func (s *DepositMoneyService) Deposit(cmd in.DepositMoneyCommand) (*domain.BankAccount, error) {
	account, err := s.bankRepo.Account(cmd.AccountId)
	if err != nil {
		return &domain.BankAccount{}, err
	}

	account.Balance += cmd.Amount

	account, err = s.bankRepo.Update(cmd.AccountId, account)
	if err != nil {
		return &domain.BankAccount{}, err
	}

	return &account, nil
}
