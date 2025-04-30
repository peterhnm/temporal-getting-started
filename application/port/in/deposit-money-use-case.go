package in

import "temporal-getting-started/domain"

type DepositMoneyUseCase interface {
	Deposit(id string, amount float64) (domain.BankAccount, error)
}
