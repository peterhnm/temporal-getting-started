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

func (s *WithdrawMoneyService) Withdraw(id string, amount float64) (domain.BankAccount, error) {
	account, err := s.bankRepo.Account(id)
	if err != nil {
		return domain.BankAccount{}, err
	}

	balance := account.Balance - amount

	if balance < 0 {
		return domain.BankAccount{}, &InsufficientFundsError{}
	}

	account.Balance = balance

	account, err = s.bankRepo.Update(id, account)
	if err != nil {
		return domain.BankAccount{}, err
	}

	return account, nil
}
