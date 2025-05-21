package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
	"temporal-getting-started/adapter/in/process"
	"temporal-getting-started/adapter/in/shared"
)

type ProcessStartHandler struct {
	workflow *process.MoneyTransferWorkflow
}

func NewProcessStartHandler(workflow *process.MoneyTransferWorkflow) *ProcessStartHandler {
	return &ProcessStartHandler{workflow}
}

type InputDto struct {
	UserId      string `json:"userId"`
	RecipientId string `json:"recipientId"`
	Amount      int    `json:"amount"`
}

type RunIdDto struct {
	RunId string `json:"runId"`
}

func (h *ProcessStartHandler) Start(c *gin.Context) {
	paramId := c.Param("id")

	var inputDto InputDto

	if err := c.BindJSON(&inputDto); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a client to talk to temporal server
	cl, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to connect to temporal server", err)
	}
	defer cl.Close()

	// Set workflowId and which task queue to use
	options := client.StartWorkflowOptions{
		ID:        paramId,
		TaskQueue: shared.MoneyTransferTaskQueueName,
	}

	log.Println("Starting payment workflow")

	// Start the workflow
	input := shared.PaymentDetails{
		UserId:      inputDto.UserId,
		RecipientId: inputDto.RecipientId,
		Amount:      float64(inputDto.Amount),
	}
	we, err := cl.ExecuteWorkflow(
		context.Background(),
		options,
		h.workflow.MoneyTransfer,
		input,
	)
	if err != nil {
		msg := fmt.Sprintf("Unable to start payment workflow %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}

	// Send a message back to the ui
	runIdDto := RunIdDto{
		RunId: we.GetRunID(),
	}
	c.JSON(http.StatusOK, runIdDto)
}
