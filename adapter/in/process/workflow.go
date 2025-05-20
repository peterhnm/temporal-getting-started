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
	NotifyUserUseCase    in.NotifyUserUseCase
	WithdrawMoneyUseCase in.WithdrawMoneyUseCase
	DepositMoneyUseCase  in.DepositMoneyUseCase
}

func NewMoneyTransferWorkflow(
	verifyUserUseCase in.VerifyUserUseCase,
	notifyUserUseCase in.NotifyUserUseCase,
	withdrawMoneyUseCase in.WithdrawMoneyUseCase,
	depositMoneyUseCase in.DepositMoneyUseCase,
) *MoneyTransferWorkflow {
	return &MoneyTransferWorkflow{
		VerifyUserUseCase:    verifyUserUseCase,
		NotifyUserUseCase:    notifyUserUseCase,
		WithdrawMoneyUseCase: withdrawMoneyUseCase,
		DepositMoneyUseCase:  depositMoneyUseCase,
	}
}

const ApproveSignal = "ApproveMoneyTransfer"

type ApproveInput struct {
	approval bool
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

	// Verify user information
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

	// Check transfer if amount > 10.000 €
	if input.Amount > 10000 {
		approval := false
		var approveInput ApproveInput
		sig := workflow.GetSignalChannel(ctx, ApproveSignal)

		// inform human worker about the new task
		workflowInfo := workflow.GetInfo(ctx)
		notifyUserCommand := in.NotifyUserCommand{
			WorkflowId:  workflowInfo.WorkflowExecution.ID,
			RunId:       workflowInfo.WorkflowExecution.RunID,
			UserId:      input.UserId,
			RecipientId: input.RecipientId,
			Amount:      input.Amount,
		}
		if err := workflow.
			ExecuteActivity(ctx, w.NotifyUserUseCase, notifyUserCommand).
			Get(ctx, nil); err != nil {
			return "", err
		}

		sig.Receive(ctx, &approveInput)
		approval = approveInput.approval
		if !approval {
			return "Money transfer could not be approved!", nil
		}
	}

	// Withdraw money
	withdrawMoneyCommand := in.WithdrawMoneyCommand{
		AccountId: user.BankAccountId,
		Amount:    input.Amount,
	}
	if err := workflow.
		ExecuteActivity(ctx, w.WithdrawMoneyUseCase.Withdraw, withdrawMoneyCommand).
		Get(ctx, nil); err != nil {
		return "", err
	}

	// Deposit money
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
