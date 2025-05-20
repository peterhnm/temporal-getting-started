package service

import (
	"temporal-getting-started/application/port/in"
	"temporal-getting-started/application/port/out"
	"temporal-getting-started/domain"
)

type WithdrawMoneyService struct {
	bankRepo out.BankRepository
}

func NewWithdrawMoneyService(bankRepo out.BankRepository) in.WithdrawMoneyUseCase {
	return &WithdrawMoneyService{bankRepo}
}

type InsufficientFundsError struct{}

func (e *InsufficientFundsError) Error() string {
	return "Insufficient funds"
}

func (s *WithdrawMoneyService) Withdraw(cmd in.WithdrawMoneyCommand) (*domain.BankAccount, error) {
	account, err := s.bankRepo.Account(cmd.AccountId)
	if err != nil {
		return &domain.BankAccount{}, err
	}

	balance := account.Balance - cmd.Amount

	if balance < 0 {
		return &domain.BankAccount{}, &InsufficientFundsError{}
	}

	account.Balance = balance

	account, err = s.bankRepo.Update(cmd.AccountId, account)
	if err != nil {
		return &domain.BankAccount{}, err
	}

	return &account, nil
}
