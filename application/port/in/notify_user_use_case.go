package in

type NotifyUserUseCase interface {
	Add(cmd NotifyUserCommand) error
}

type NotifyUserCommand struct {
	WorkflowId  string
	RunId       string
	UserId      string
	RecipientId string
	Amount      float64
}
