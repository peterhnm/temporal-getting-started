<div id="top"></div>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
<!-- END OF PROJECT SHIELDS -->

<!-- PROJECT LOGO -->
<br />
<div align="center">
    <a href="https://www.temporal.io">
        <img src="images/temporal-logo.png" alt="Image" height="180">
    </a>
    <p>
        <a href="https://github.com/peterhnm/temporal-getting-started/issues">Report Bug</a>
        ·
        <a href="https://github.com/peterhnm/temporal-getting-started/pulls">Request Feature</a>
    </p>
</div>

## About The Project

Mit diesem Projekt möchte ich nach und nach Temporal besser 
kennenlernen.
Gleichzeitig sind es meine ersten Schritte mit der 
Programmiersprache `Go`.
Die Idee ist es bekannte Konzepte aus BPMN mithilfe von 
Temporal umzusetzen.
Folgende Konzepte sollen intensiver beleuchtet werden:  
- [x] Service Task  
- [x] User Task  
- [ ] Message Correlation  
- [ ] ...  

## Getting started
1. Start a temporal server
   ```bash
   temporal server start-dev --ui-port 8080
   ```

2. Start of the application
   ```bash
   go run .
   ```
   
3. Start a process
   ```bash
   curl -X POST http://127.0.0.1:3000/api/start/pay-invoice-701 \
   -H "Content-Type: application/json" \
   -d '{"userId": "1234", "recipientId": "5678", "amount": 10000}'
   ```
   
4. Complete "user" task (The **runId** is returned by the former command)
   ```bash
   curl -X POST http://127.0.0.1:3000/api/complete \
   -H "Content-Type: application/json" \
   -d '{"workflowId": "pay-invoice-701", "runId": "YOUR_RUN_ID", "decision": false}'
   ```

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/peterhnm/temporal-getting-started.svg?style=for-the-badge
[contributors-url]: https://github.com/peterhnm/temporal-getting-started/graphs/contributors

[stars-shield]: https://img.shields.io/github/stars/peterhnm/temporal-getting-started.svg?style=for-the-badge
[stars-url]: https://github.com/peterhnm/temporal-getting-started/stargazers

[issues-shield]: https://img.shields.io/github/issues/peterhnm/temporal-getting-started.svg?style=for-the-badge
[issues-url]: https://github.com/peterhnm/temporal-getting-started/issues
