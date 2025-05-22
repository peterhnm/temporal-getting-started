# Activity and Service Tasks

Activities in Temporal are the units of work that perform external or side effect 
operations — such as calling a service, writing to a database, or sending an email. 
Understanding how activities are defined, invoked, and managed is essential to building 
reliable and maintainable Temporal workflows.

An **Activity** in Temporal is a regular function or method, typically defined in a 
separate interface or module from the workflow itself, that performs an operation which 
is not deterministic or has side effects. 
Unlike workflows, activities are not replayed;
instead, they are executed by workers, and their results are persisted to Temporal's 
history for durability and auditability.

- **Activity** represents the actual execution of business logic, similar to a 
**Service Task** in BPMN, which delegates work to an external system or application.
- **Activities** are connected to workflows via remote procedure calls (RPC) that 
Temporal orchestrates, ensuring retries, timeouts, and fault handling automatically.
- **Retry Policies** and **Timeouts** can be configured per activity, offering robust 
error-handling mechanisms out of the box.
- **Heartbeats** can be used for long-running activities to report progress and detect 
failures early.

Temporal activities address the common pain points of managing unreliable external 
systems, transient errors, and complex retry logic. 
By implementing activities within the Temporal model, developers can isolate side 
effects from workflow logic, maintain clean separation of concerns, and gain automatic 
durability and observability.

By using activities, developers gain fine-grained control over error handling, resilience,
and operational visibility. Temporal ensures that activities are retried when they fail, 
timed out when they hang, and tracked end-to-end via the event history.

To use an **Activity** in Temporal, you create an activity interface, provide an 
implementation, register it with a worker, and invoke it from a workflow using an 
activity stub. 
Temporal handles the invocation, retry logic, and result persistence transparently, 
allowing you to focus on business logic instead of infrastructure plumbing.

<include from="Workflow.md" element-id="generic-github-link" />

1. Implement the activity
   ```go
   func (svc *VerificationService) Verify(userId string) (*domain.User, error) {
       log.Println(fmt.Sprintf("[USER] Verifying user %s", userId))
       user, err := svc.userRepo.FindById(userId)
       if err != nil {
           return &domain.User{}, err
       }
       return &user, nil
   }
   ```

2. Register the activity with a worker
   ```go
   func (w *Worker) Run() error {
       c, err := client.Dial(client.Options{})
       if err != nil {
           return err
       }
       defer c.Close()

       wf := worker.New(c, shared.MoneyTransferTaskQueueName, worker.Options{})

       // Register the Workflow
       wf.RegisterWorkflow(w.workflow.MoneyTransfer)

       // Register one or more activities
       wf.RegisterActivity(w.workflow.VerifyUserUseCase.Verify)
       // ...
   }
   ```

3. Call the activity from the workflow
   ```go
   func (w *MoneyTransferWorkflow) MoneyTransfer(
       ctx workflow.Context,
       input shared.PaymentDetails
   ) (string, error) {
       // ...
       var result domain.User
       userFuture := workflow.
           ExecuteActivity(ctx, w.VerifyUserUseCase.Verify, input.UserId).
           Get(ctx, &user)
       // ...
   }
   ```

<seealso>
    <category ref="temp">
        <a href="https://docs.temporal.io/activities">Activity</a>
        <a href="https://docs.temporal.io/develop/go/core-application#develop-activities">Go SDK</a>
        <a href="https://docs.temporal.io/develop/java/core-application#develop-activities">Java SDK</a>
    </category>
</seealso>
