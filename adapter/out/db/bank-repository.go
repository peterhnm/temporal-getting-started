package db

import (
	"fmt"
	"log"
	"temporal-getting-started/application/port/out"
	"temporal-getting-started/domain"
)

type BankRepository struct {
	accounts map[string]domain.BankAccount
}

func NewBankRepository() out.BankRepository {
	return &BankRepository{
		accounts: map[string]domain.BankAccount{
			"12-3456": {"12-3456", 100},
			"56-7890": {"56-7890", 200},
		},
	}
}

type InvalidAccountError struct {
	id string
}

func (err InvalidAccountError) Error() string {
	return fmt.Sprintf("Account with id %s does not exist", err.id)
}

func (repo *BankRepository) Account(id string) (domain.BankAccount, error) {
	if account, exists := repo.accounts[id]; exists {
		log.Printf("Found account with id %s", id)
		return account, nil
	}
	return domain.BankAccount{}, &InvalidAccountError{id}
}

func (repo *BankRepository) Update(id string, account domain.BankAccount) (domain.BankAccount, error) {
	repo.accounts[id] = account
	log.Printf("Updated account with id %s", id)
	return account, nil
}
