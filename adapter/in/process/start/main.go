package main

import (
	"context"
	"go.temporal.io/sdk/client"
	"log"
	"temporal-getting-started/adapter/in/process"
	"temporal-getting-started/adapter/in/shared"
	"temporal-getting-started/internal/di"
)

func Start(moneyTransferWorkflow *process.MoneyTransferWorkflow) {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to connect to temporal server", err)
	}

	defer c.Close()

	input := shared.PaymentDetails{
		UserId:      "1234",
		RecipientId: "5678",
		Amount:      10000.00,
	}

	options := client.StartWorkflowOptions{
		ID:        "pay-invoice-701",
		TaskQueue: shared.MoneyTransferTaskQueueName,
	}

	log.Println("Starting payment workflow")

	we, err := c.ExecuteWorkflow(context.Background(), options, moneyTransferWorkflow.MoneyTransfer, input)
	if err != nil {
		log.Fatalln("Unable to start payment workflow", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result string

	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get result from workflow", err)
	}

	log.Printf("Workflow Result: %s", result)
}

func main() {
	workflow := di.Initialize()
	Start(workflow)
}
