//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"temporal-getting-started/adapter/in/process"
	"temporal-getting-started/adapter/out/db"
	"temporal-getting-started/application/service"
)

var (
	RepositorySet = wire.NewSet(
		db.NewUserRepository,
		db.NewBankRepository,
		db.NewUserTaskRepository,
	)
	ServiceSet = wire.NewSet(
		service.NewVerificationService,
		service.NewNotifyUserService,
		service.NewUserTaskService,
		service.NewWithdrawMoneyService,
		service.NewDepositMoneyService,
	)
	WorkflowSet = wire.NewSet(
		process.NewMoneyTransferWorkflow,
	)
)

func Initialize() *process.MoneyTransferWorkflow {
	wire.Build(
		RepositorySet,
		ServiceSet,
		WorkflowSet,
	)
	return nil
}
