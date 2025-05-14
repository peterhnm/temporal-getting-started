package in

import "temporal-getting-started/domain"

type WithdrawMoneyUseCase interface {
	Withdraw(cmd WithdrawMoneyCommand) (*domain.BankAccount, error)
}

type WithdrawMoneyCommand struct {
	AccountId string
	Amount    float64
}
