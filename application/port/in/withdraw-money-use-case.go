package in

import "temporal-getting-started/domain"

type WithdrawMoneyUseCase interface {
	Withdraw(id string, amount float64) (domain.BankAccount, error)
}
