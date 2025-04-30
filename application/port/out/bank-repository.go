package out

import "temporal-getting-started/domain"

type BankRepository interface {
	Account(id string) (domain.BankAccount, error)

	Update(id string, account domain.BankAccount) (domain.BankAccount, error)
}
