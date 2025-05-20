package domain

type UserTask struct {
	Id   UserTaskId   `json:"id"`
	Data UserTaskData `json:"data"`
}

type UserTaskId struct {
	WorkflowId string
	RunId      string
}

type UserTaskData struct {
	UserId      string
	RecipientId string
	Amount      float64
	Decision    *bool // will be set after the user decides
}
