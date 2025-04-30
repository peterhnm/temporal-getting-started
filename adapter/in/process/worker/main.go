package main

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"temporal-getting-started/adapter/in/process"
	"temporal-getting-started/adapter/in/shared"
	"temporal-getting-started/internal/di"
)

func Worker(workflow *process.MoneyTransferWorkflow) {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, shared.MoneyTransferTaskQueueName, worker.Options{})

	w.RegisterWorkflow(workflow.MoneyTransfer)
	w.RegisterActivity(workflow.VerifyUserUseCase.Verify)
	w.RegisterActivity(workflow.WithdrawMoneyUseCase.Withdraw)
	w.RegisterActivity(workflow.DepositMoneyUseCase.Deposit)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func main() {
	workflow := di.Initialize()
	Worker(workflow)
}
