// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/google/wire"
	"temporal-getting-started/adapter/in/process"
	"temporal-getting-started/adapter/out/db"
	"temporal-getting-started/application/service"
)

// Injectors from wire.go:

func Initialize() *process.MoneyTransferWorkflow {
	userRepository := db.NewUserRepository()
	verifyUserUseCase := service.NewVerificationService(userRepository)
	bankRepository := db.NewBankRepository()
	withdrawMoneyUseCase := service.NewWithdrawMoneyService(bankRepository)
	depositMoneyUseCase := service.NewDepositMoneyService(bankRepository)
	moneyTransferWorkflow := process.NewMoneyTransferWorkflow(verifyUserUseCase, withdrawMoneyUseCase, depositMoneyUseCase)
	return moneyTransferWorkflow
}

// wire.go:

var (
	RepositorySet = wire.NewSet(db.NewUserRepository, db.NewBankRepository)
	ServiceSet    = wire.NewSet(service.NewVerificationService, service.NewWithdrawMoneyService, service.NewDepositMoneyService)
	WorkflowSet   = wire.NewSet(process.NewMoneyTransferWorkflow)
)
