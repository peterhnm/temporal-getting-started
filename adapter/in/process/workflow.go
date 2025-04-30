package process

import (
	"fmt"
	"temporal-getting-started/adapter/in/shared"
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

	userErr := workflow.ExecuteActivity(ctx, w.VerifyUserUseCase.Verify, input).Get(ctx, nil)
	if userErr != nil {
		return "", userErr
	}

	withdrawErr := workflow.ExecuteActivity(ctx, w.WithdrawMoneyUseCase.Withdraw, input).Get(ctx, nil)
	if withdrawErr != nil {
		return "", withdrawErr
	}

	depositErr := workflow.ExecuteActivity(ctx, w.DepositMoneyUseCase.Deposit, input).Get(ctx, nil)
	if depositErr != nil {
		return "", depositErr
	}

	result := fmt.Sprintf("Money transferred successfully")
	return result, nil
}
