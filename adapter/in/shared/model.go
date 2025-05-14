package shared

const MoneyTransferTaskQueueName = "TRANSFER_MONEY_TASK_QUEUE"

type PaymentDetails struct {
	UserId      string  `json:"userId"`
	RecipientId string  `json:"recipientId"`
	Amount      float64 `json:"amount"`
}
