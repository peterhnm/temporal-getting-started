# Workflow

Workflows are the core abstraction Temporal provides for modeling long-running, stateful 
business logic in a way that is resilient to failures, restarts, and timeouts. This section 
covers what workflows are, how they behave, and how developers use them to model business 
processes, especially in distributed systems.

A **Workflow** in Temporal is a special type of function written in code (e.g., Go or Java)
that defines a series of steps or decisions in a business process. Unlike regular functions,
workflows are durable: they are persisted and replayed from an append-only **Event History** 
to ensure deterministic execution. This means workflows can run for minutes, days, or even 
months, surviving failures, redeployments, or system crashes without losing their state.

- **Workflow** represents the orchestrated logic of a process — it defines the order in 
which **Activities** are executed, as well as branching, retries, timers, signals, and 
other control flow logic.
- **Workflow Execution** is the running instance of a workflow definition, uniquely 
identified by a Workflow ID and Run ID. Each execution is tracked and fully auditable via 
its event history.
- **Determinism** is essential: workflows must produce the same results when replayed from 
history. For this reason, only a restricted subset of operations (e.g., no random numbers 
or non-deterministic APIs) can be used inside workflow code.
- **Workflow State** is not stored in memory or a database; instead, it is reconstructed 
by replaying past events, which ensures high reliability and consistency.

Temporal workflows address the common pain points of coordinating distributed systems, 
especially those requiring retry logic, manual compensation, long timeouts, and reliable 
state tracking. By implementing workflows in Temporal, users can replace brittle glue code 
and ad hoc orchestration logic with a consistent, scalable programming model.

By using workflows, developers gain a fault-tolerant and observable execution environment 
where failures are first-class citizens, and where business logic can evolve safely through 
code versioning and deterministic replays.

To use a **Workflow** in Temporal, you create a workflow function, register it with a
**Worker**, and start it via the Temporal service using a **WorkflowStub** or API call. 
The service then ensures the workflow progresses step by step — reliably, durably, and at 
scale.

## Comparison of Temporal Workflow and BPMN Process

While both **Temporal Workflows** and **BPMN Processes** aim to model business processes 
and orchestrate tasks, they take fundamentally different approaches. BPMN uses a graphical 
language defined in XML, relying on drag-and-drop tools to model control flow elements like 
gateways and tasks. In contrast, Temporal workflows are implemented directly in code, which 
shifts control flow from visual diagrams to standard programming constructs like `if`, 
`switch`, and `for` statements.

This difference has practical implications. For example, a **BPMN XOR Gateway** that models 
a branching decision becomes a simple `if-else` block in a Temporal workflow. 
A **Parallel Gateway** in BPMN, which initiates concurrent paths, is represented in Temporal 
by executing activities asynchronously using the `async` or `Future` APIs and then joining 
them with `await` or `get`. This gives developers precise control over concurrency and logic 
flow without needing a separate modeling tool.

While BPMN is highly expressive and accessible to non-developers, it can become difficult 
to manage and test as complexity grows. Temporal, by contrast, integrates naturally into 
developer workflows, offering the full power of a programming language — including testability, 
version control, and IDE support — at the cost of a steeper learning curve for business 
analysts.

**Table: Comparison of Temporal Workflow vs BPMN Process Flow Constructs**

| Flow Construct       | BPMN                              | Temporal Workflow (Code-based)                                          |
|----------------------|-----------------------------------|-------------------------------------------------------------------------|
| Sequence Flow        | Visual connectors between tasks   | Sequential function calls                                               |
| XOR Gateway          | Diamond with branching conditions | `if` / `else if` / `else` statements                                    |
| Parallel Gateway     | Forks multiple paths in parallel  | Start activities with `async` / `Future` and join with `await` or `get` |
| Looping (e.g. while) | Loop marker or gateways           | `while`, `for`, or recursive calls                                      |

Temporal workflows provide a powerful and flexible code-first alternative to BPMN processes. 
While they lack the immediate visual clarity of BPMN diagrams, they offer significantly 
more control, debuggability, and maintainability for developer-centric teams.

## Implementation

For our specific example, this means the following:

The process you already saw on the starting page can look like this code in `go`.

![bpmn process](process.svg){ width="200" thumbnail="true" }

```go
func (w *MoneyTransferWorkflow) MoneyTransfer(
    ctx workflow.Context,
    input shared.PaymentDetails
) (string, error) {

    retryPolicy := &temporal.RetryPolicy{
        InitialInterval:        time.Second,
        BackoffCoefficient:     2.0,
        MaximumInterval:        100 * time.Second,
        MaximumAttempts:        500,
        NonRetryableErrorTypes: []string{
            "InvalidUserError",
            "InvalidAccountError",
            "InsufficientFundsError",
        },
    }

    options := workflow.ActivityOptions{
        StartToCloseTimeout: time.Minute,
        RetryPolicy:         retryPolicy,
    }

    ctx = workflow.WithActivityOptions(ctx, options)

    // Verify user information
    var user domain.User
    var recipient domain.User

    // non-blocking execution of two activities in parallel
    userFuture := workflow.ExecuteActivity(
        ctx,
        w.VerifyUserUseCase.Verify,
        input.UserId)
    recipientFuture := workflow.ExecuteActivity(
        ctx,
        w.VerifyUserUseCase.Verify,
        input.RecipientId)

    // Wait for the result
    if err := userFuture.Get(ctx, &user); err != nil {
        return "", err
    }
    if err := recipientFuture.Get(ctx, &recipient); err != nil {
        return "", err
    }

    // Check transfer if amount > 10.000 €
    if input.Amount >= 10000 {
        var approveInput shared.ApproveInput
        
        // Create a signal channel
        sig := workflow.GetSignalChannel(ctx, shared.ApproveSignal)

        // inform human worker about the new task
        workflowInfo := workflow.GetInfo(ctx)
        notifyUserCommand := in.NotifyUserCommand{
            WorkflowId:  workflowInfo.WorkflowExecution.ID,
            RunId:       workflowInfo.WorkflowExecution.RunID,
            UserId:      input.UserId,
            RecipientId: input.RecipientId,
            Amount:      input.Amount,
        }
        if err := workflow.ExecuteActivity(
                ctx,
                w.NotifyUserUseCase.Add,
                notifyUserCommand,
            ).
            Get(ctx, nil); err != nil {
                
            return "", err
        }

        // Wait for incoming signal (human interaction)
        sig.Receive(ctx, &approveInput)
        approval := approveInput.Approval
        if !approval {
            return "Money transfer could not be approved!", nil
        }
    }

    // Withdraw money
    withdrawMoneyCommand := in.WithdrawMoneyCommand{
        AccountId: user.BankAccountId,
        Amount:    input.Amount,
    }
    if err := workflow.ExecuteActivity(
            ctx,
            w.WithdrawMoneyUseCase.Withdraw,
            withdrawMoneyCommand,
        ).
        Get(ctx, nil); err != nil {
            
        return "", err
    }

    // Deposit money
    depositMoneyCommand := in.DepositMoneyCommand{
        AccountId: recipient.BankAccountId,
        Amount:    input.Amount,
    }
    if err := workflow.ExecuteActivity(
                ctx,
                w.DepositMoneyUseCase.Deposit,
                depositMoneyCommand
            ).
            Get(ctx, nil); err != nil {
                
            return "", err
    }

    result := fmt.Sprintf("Money transferred successfully")
    return result, nil
}
```

For a full example, please check out the repository on [GitHub](https://github.com/peterhnm/temporal-getting-started).

<seealso>
    <category ref="temp">
        <a href="https://docs.temporal.io/workflows">Workflow</a>
        <a href="https://docs.temporal.io/develop/go/core-application#develop-workflows">Go SDK</a>
        <a href="https://docs.temporal.io/develop/java/core-application#develop-workflows">Java SDK</a>
    </category>
</seealso>
