//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"temporal-getting-started/internal"

	"temporal-getting-started/adapter/in/process"
	"temporal-getting-started/adapter/in/rest"
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
	TemporalSet = wire.NewSet(
		process.NewMoneyTransferWorkflow,
		process.NewWorker,
	)
	WebSet = wire.NewSet(
		rest.NewServerHTTP,
		rest.NewUserTaskHandler,
		rest.NewProcessStartHandler,
	)
)

func InitializeApplication() (*internal.Application, error) {
	wire.Build(
		RepositorySet,
		ServiceSet,
		TemporalSet,
		WebSet,
		wire.Struct(new(internal.Application), "*"),
	)
	return &internal.Application{}, nil
}

// func InitializeTemporal() (*process.Worker, error) {
// 	wire.Build(
// 		RepositorySet,
// 		ServiceSet,
// 		WorkflowSet,
// 		WorkerSet,
// 	)
// 	return &process.Worker{}, nil
// }
//
// func InitializeApi() (*rest.ServerHTTP, error) {
// 	wire.Build(
// 		RepositorySet,
// 		ServiceSet,
// 		WorkflowSet,
// 		WebSet,
// 	)
//
// 	return &rest.ServerHTTP{}, nil
// }
