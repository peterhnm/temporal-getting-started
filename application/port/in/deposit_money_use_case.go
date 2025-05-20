package in

import "temporal-getting-started/domain"

type DepositMoneyUseCase interface {
	Deposit(cmd DepositMoneyCommand) (*domain.BankAccount, error)
}

type DepositMoneyCommand struct {
	AccountId string
	Amount    float64
}
