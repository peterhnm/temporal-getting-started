package process

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"temporal-getting-started/adapter/in/shared"
)

type Worker struct {
	workflow *MoneyTransferWorkflow
}

func NewWorker(workflow *MoneyTransferWorkflow) *Worker {
	return &Worker{workflow}
}

func (w *Worker) Run() error {
	c, err := client.Dial(client.Options{})
	if err != nil {
		return err
	}
	defer c.Close()

	wf := worker.New(c, shared.MoneyTransferTaskQueueName, worker.Options{})

	wf.RegisterWorkflow(w.workflow.MoneyTransfer)

	wf.RegisterActivity(w.workflow.VerifyUserUseCase.Verify)
	wf.RegisterActivity(w.workflow.NotifyUserUseCase.Add)
	wf.RegisterActivity(w.workflow.WithdrawMoneyUseCase.Withdraw)
	wf.RegisterActivity(w.workflow.DepositMoneyUseCase.Deposit)

	log.Println("Starting worker")
	err = wf.Run(worker.InterruptCh())
	if err != nil {
		return err
	}

	return nil
}
