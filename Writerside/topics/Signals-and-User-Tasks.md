# Signals and User Tasks

While BPMN explicitly defines User Tasks as points where human interaction is required, 
Temporal handles such asynchronous external inputs through the use of Signals. 
This section explores what signals are, how they are used in workflows, and how they 
enable dynamic interaction with long-running processes.

A **Signal** in Temporal is a mechanism for sending asynchronous, external input to a 
running workflow execution. 
Signals are used to communicate information or trigger logic inside a workflow after it 
has already started — such as updating internal state, making a decision, or resuming 
from a wait condition. 
They are a key building block for implementing workflows that react to human actions or 
events from external systems.

Signals ...
- represents an asynchronous, external input to a workflow execution.
- are connected to workflows through signal handlers, which can update workflow state or 
influence control flow.
- are similar to BPMN **User Tasks**, which pause the process until a user completes a 
task or provides input.
- solve the challenge of integrating human input or third-party system events into a 
running, durable workflow — without sacrificing consistency or fault tolerance.

By implementing signals, developers can build workflows that pause and wait for input 
without busy-waiting or polling, and then react immediately once the input arrives.
This pattern enables building human-in-the-loop systems, approval flows, and reactive 
business processes with minimal complexity.

By using signals, developers gain a clean, event-driven way to interact with workflows 
after they've started. 
Temporal ensures that signal delivery is reliable and that signals are persisted in the 
workflow's event history, so the workflow can always resume correctly — even after 
failures or restarts.

<include from="Workflow.md" element-id="generic-github-link" />

1. Create a signal channel inside your workflow
   ```go
   func (w *MoneyTransferWorkflow) MoneyTransfer(
       ctx workflow.Context,
       input shared.PaymentDetails
   ) (string, error) {
       // ...
       sig := workflow.GetSignalChannel(ctx, shared.UserSignal)
       // ...
   }
   ```
   
2. Inform the user about a new task (e.g. Entry to a database, Email, etc.)
   ```go
   func (w *MoneyTransferWorkflow) MoneyTransfer(
       ctx workflow.Context,
       input shared.PaymentDetails
   ) (string, error) {
       // ...
       sig := workflow.GetSignalChannel(ctx, shared.UserSignal)
       // save a user task to the database
       workflowInfo := workflow.GetInfo(ctx)
       userTask := UserTask{
            WorkflowId:  workflowInfo.WorkflowExecution.ID,
            RunId:       workflowInfo.WorkflowExecution.RunID,
            Data:        // Add data that the user needs
       }
       if err := workflow.
           ExecuteActivity(ctx, w.SaveUserTask.Save, userTask).
           Get(ctx, nil); err != nil {
           return "", err
       }
       // ...
   }
   ```
   
3. Wait for the user input
   ```go
   func (w *MoneyTransferWorkflow) MoneyTransfer(
       ctx workflow.Context,
       input shared.PaymentDetails
   ) (string, error) {
       // ...
       sig := workflow.GetSignalChannel(ctx, shared.UserSignal)
       // save a user task to the database
       workflowInfo := workflow.GetInfo(ctx)
       userTask := UserTask{
            WorkflowId:  workflowInfo.WorkflowExecution.ID,
            RunId:       workflowInfo.WorkflowExecution.RunID,
            Data:        // Add data that the user needs
       }
       if err := workflow.
           ExecuteActivity(ctx, w.SaveUserTask.Save, userTask).
           Get(ctx, nil); err != nil {
           return "", err
       }
   
       // Wait for user input
       var userInput shared.UserInput
       sig.Receive(ctx, &userInput)
       // ...
   }
   ```
   
4. Implement sending a signal
   ```go
   func CompleteUserTask(
       ctx context.Context,
       taskCompletionDto TaskCompletionDto
   ) {
       // Send signal to workflow
       temporalClient, err := client.Dial(client.Options{})
       if err != nil {
           log.Fatalln("Unable to create client", err)
       }
       defer temporalClient.Close()

       if err := temporalClient.SignalWorkflow(
               ctx,
               taskCompleteDto.WorkflowId,
               taskCompleteDto.RunId,
               shared.UserSignal,
               shared.UserInput{Approval: taskCompleteDto.Decision},
           ); err != nil {
   
           log.Println("Unable to signal workflow", err)
       }
   }
   ```

<seealso>
    <category ref="temp">
        <a href="https://docs.temporal.io/develop/go/message-passing#send-signal-from-workflow">Go SDK</a>
        <a href="https://docs.temporal.io/develop/java/message-passing#send-signal-from-workflow">Java SDK</a>
    </category>
</seealso>
