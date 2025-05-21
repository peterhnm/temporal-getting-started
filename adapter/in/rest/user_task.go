package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
	"temporal-getting-started/adapter/in/shared"

	"temporal-getting-started/domain"

	"temporal-getting-started/application/port/in"
)

type UserTaskHandler struct {
	useCase in.UserTaskUseCase
}

func NewUserTaskHandler(useCase in.UserTaskUseCase) *UserTaskHandler {
	return &UserTaskHandler{useCase}
}

type TaskCompleteDto struct {
	WorkflowId string `json:"workflowId"`
	RunId      string `json:"runId"`
	Decision   bool   `json:"decision"`
}

type UserTaskDto struct {
	UserId      string  `copier:"must"`
	RecipientId string  `copier:"must"`
	Amount      float64 `copier:"must"`
}

func (handler *UserTaskHandler) GetData(c *gin.Context) {
	var id domain.UserTaskId

	if err := c.BindJSON(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskData, err := handler.useCase.GetData(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		userTaskDto := UserTaskDto{}
		if err := copier.Copy(&userTaskDto, taskData.Data); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		c.JSON(http.StatusOK, userTaskDto)
	}
}

func (handler *UserTaskHandler) Complete(c *gin.Context) {
	var taskCompleteDto TaskCompleteDto

	if err := c.BindJSON(&taskCompleteDto); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := domain.UserTaskId{
		WorkflowId: taskCompleteDto.WorkflowId,
		RunId:      taskCompleteDto.RunId,
	}

	err := handler.useCase.Complete(id, taskCompleteDto.Decision)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		// Send signal to workflow
		temporalClient, err := client.Dial(client.Options{})
		if err != nil {
			log.Fatalln("Unable to create client", err)
		}
		defer temporalClient.Close()

		if err := temporalClient.SignalWorkflow(
			c.Request.Context(),
			taskCompleteDto.WorkflowId,
			taskCompleteDto.RunId,
			shared.ApproveSignal,
			shared.ApproveInput{Approval: taskCompleteDto.Decision},
		); err != nil {
			log.Println("Unable to signal workflow", err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Task completed successfully"})
	}
}
