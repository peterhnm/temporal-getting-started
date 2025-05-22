# The Concept of Temporal

This article explains the basics of Temporal, a microservice orchestration tool, and how
it works in the context of distributed systems. Temporal helps developers build 
fault-tolerant and stateful applications by managing long-running workflows and 
distributed transactions in a reliable and scalable manner. This article introduces the 
core concepts behind Temporal, such as workflows, activities, task queues, and the 
Temporal service itself, and explains how they interact to provide durable execution and 
high availability in microservices architectures.

Temporal is an open-source platform that enables developers to write and manage complex, 
long-running business logic as workflows. These workflows are written in code (e.g., Go 
or Java), but Temporal ensures they are resilient to process restarts, network issues, 
and service outages by persisting their state. Rather than relying on brittle external 
schedulers or ad-hoc retry logic, Temporal handles retries, state recovery, and execution 
tracking automatically.
* **Workflow** is a function written by the developer that defines the sequence of business 
operations. It is durable and resumable, ensuring that even if a process crashes, the 
workflow can resume from its last known state.
* **Activity** is a unit of work executed by a worker. Activities represent external 
operations, such as API calls or database writes, and can be retried or time out as needed.
* **Worker** is a service process that polls task queues and executes workflow or activity 
code. Workers are stateless and can be scaled independently, allowing Temporal to 
distribute tasks dynamically across instances. When a worker picks up a workflow task from 
the task queue, it replays the workflow’s Event History to reconstruct its current state.
* **Temporal Service** is the backend component that manages workflow state, task queues,
  and execution history.
* **Task Queue** is a logical channel used by the Temporal Service to dispatch tasks to 
available workers. It enables horizontal scaling by allowing multiple workers to poll the
same queue and supports fault isolation by decoupling task producers (Temporal Service) 
from task consumers (workers).
* **Event History** is a complete, append-only log of all events related to a workflow 
execution. It records every state transition, activity call, and decision made, allowing 
for full traceability and deterministic replay of workflows.

Temporal addresses the common pain points of building distributed systems, such as 
ensuring reliability, managing retries and failures, and tracking execution flow across 
services. By implementing Temporal, users can reduce operational complexity, improve 
fault tolerance, and gain clear visibility into their system’s behavior. Developers 
create workflows and activities using the Temporal SDK, which abstracts away the 
complexity of state management and execution tracking.

![temporal_architecture](architecture.png){ style="block" }
*Diagram showing how Temporal orchestrates workflows and activities across services using 
a central Temporal Service, task queues, and distributed workers.*

## Background

The idea of Temporal originated from the growing demand for managing complex, long-running
workflows reliably across distributed microservices. As systems evolved toward event-driven
and service-oriented architectures, developers faced increasing challenges with ensuring 
fault tolerance, consistency, and observability across asynchronous operations. Traditional
solutions often required intricate custom retry logic, state persistence mechanisms, and 
coordination layers, which were error-prone and hard to maintain.

With the rise of microservices and event-driven architectures, the need for things such 
as durable execution, reliable retries, and workflow state management has become paramount. 
Patterns like the Outbox Pattern were introduced to address consistency between services 
and databases during distributed transactions, but they often added complexity and required 
meticulous integration.

Temporal addresses these challenges by providing a unified programming model for orchestrating
asynchronous tasks with built-in reliability, retries, and auditability. It eliminates 
the need for separate coordination systems or transactional messaging hacks by treating 
workflow logic as code and persisting every step of its execution in an Event History. 
This shift dramatically simplifies the development of resilient systems without sacrificing
correctness or visibility.

## Comparison of Temporal and BPMN-based Process Engines

While **Temporal** and BPMN-based process engines like **Camunda 7** and **Camunda 8** 
both serve to orchestrate business processes, they differ significantly in approach, 
developer experience, and operational philosophy. BPMN engines rely on graphical models 
and declarative workflows defined in XML, which are interpreted by the engine at runtime. 
In contrast, Temporal defines workflows in code (e.g., Java or Go), giving developers full 
control over logic, tooling, and versioning, while still providing robust state management 
and fault tolerance.

BPMN is often favored in organizations where business analysts or non-developers are involved
in process modeling, whereas Temporal is optimized for engineering teams who prefer working
directly in code. Camunda 8, with its cloud-native Zeebe engine, narrows the gap by offering
improved scalability, but still relies on BPMN for defining workflows, which can make complex
logic harder to test and maintain compared to Temporal's code-first approach.

**Table: Comparison of Temporal vs BPMN-based Process Engines (Camunda 7/8)**

| Feature                       | Temporal                         | Camunda 7                          | Camunda 8 (Zeebe)                      |
|-------------------------------|----------------------------------|------------------------------------|----------------------------------------|
| Workflow Definition           | Code (Go, Java)                  | BPMN XML                           | BPMN XML                               |
| Target Audience               | Developers                       | Business + Developers              | Business + Developers                  |
| Fault Tolerance               | Built-in with Event History      | Requires manual retry logic        | Built-in, but limited logic            |
| Long-running Workflow Support | Yes (durable, resumable)         | Yes (via database persistence)     | Yes (via distributed log)              |
| Developer Experience          | Native code, testable, IDE-ready | Graphical modeling, external logic | Graphical modeling, limited code hooks |
| Versioning                    | Code-level versioning            | Requires BPMN versioning           | Requires BPMN versioning               |
| State Replay                  | Deterministic replay via history | Not supported                      | Partial replay via event sourcing      |
| Community and Ecosystem       | Growing developer-centric        | Mature, enterprise-focused         | Modern, evolving, cloud-native         |

Temporal provides a code-first, developer-friendly alternative to BPMN engines, particularly
suited for systems where workflow logic is tightly integrated with application code and 
requires strong reliability guarantees.

