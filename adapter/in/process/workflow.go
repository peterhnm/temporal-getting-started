package process

import (
	"fmt"
	"temporal-getting-started/adapter/in/shared"
	"temporal-getting-started/domain"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"temporal-getting-started/application/port/in"
)

type MoneyTransferWorkflow struct {
	VerifyUserUseCase    in.VerifyUserUseCase
	WithdrawMoneyUseCase in.WithdrawMoneyUseCase
	DepositMoneyUseCase  in.DepositMoneyUseCase
}

func NewMoneyTransferWorkflow(
	verifyUserUseCase in.VerifyUserUseCase,
	withdrawMoneyUseCase in.WithdrawMoneyUseCase,
	depositMoneyUseCase in.DepositMoneyUseCase,
) *MoneyTransferWorkflow {
	return &MoneyTransferWorkflow{
		VerifyUserUseCase:    verifyUserUseCase,
		WithdrawMoneyUseCase: withdrawMoneyUseCase,
		DepositMoneyUseCase:  depositMoneyUseCase,
	}
}

func (w *MoneyTransferWorkflow) MoneyTransfer(ctx workflow.Context, input shared.PaymentDetails) (string, error) {

	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        500,
		NonRetryableErrorTypes: []string{"InvalidUserError", "InvalidAccountError", "InsufficientFundsError"},
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var user domain.User
	var recipient domain.User

	userFuture := workflow.
		ExecuteActivity(ctx, w.VerifyUserUseCase.Verify, input.UserId)
	recipientFuture := workflow.
		ExecuteActivity(ctx, w.VerifyUserUseCase.Verify, input.RecipientId)

	if err := userFuture.Get(ctx, &user); err != nil {
		return "", err
	}
	if err := recipientFuture.Get(ctx, &recipient); err != nil {
		return "", err
	}

	withdrawMoneyCommand := in.WithdrawMoneyCommand{
		AccountId: user.BankAccountId,
		Amount:    input.Amount,
	}
	if err := workflow.
		ExecuteActivity(ctx, w.WithdrawMoneyUseCase.Withdraw, withdrawMoneyCommand).
		Get(ctx, nil); err != nil {
		return "", err
	}

	depositMoneyCommand := in.DepositMoneyCommand{
		AccountId: recipient.BankAccountId,
		Amount:    input.Amount,
	}
	if err := workflow.
		ExecuteActivity(ctx, w.DepositMoneyUseCase.Deposit, depositMoneyCommand).
		Get(ctx, nil); err != nil {
		return "", err
	}

	result := fmt.Sprintf("Money transferred successfully")
	return result, nil
}
