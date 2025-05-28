# System perspective
This section presents the system.

## Design and architecture 
The system is refactored with the Go programming language. 

### Module diagram
An overview of the modules of the codebase is presented by the following package diagram.
The package diagram models the internal structure of the codebase from the src folder (not infrastructure).
Though it should be noted that within the handlers folder, the classes auth.go, message.go and user.go are presented, and its dependencies. This is to depict the complexity of that modules (since it is the biggest and most central module in the system).

![Module diagram](../images/module_diagram.png)

In the diagram it can be seen, that the main.go file orchestrates the system. It (in this context) has the responsibility for:
1. Render the template (frontend)
2. Initialize a new instance of the database object
3. Setup middleware
4. Setup routes, which have the responsibility of exposing the endpoints that further orchestrates to the handlers module for the logic of the API.

### Sequence diagrams
Below, two sequence diagrams showcase how the different parts of the system interact when procecssing a "Follow" request. The first version shows the processes involved when the request is sent via. the UI, whereas the second version shows the processes involved when sent via. the API. 

![Sequence diagram - Follow request via UI](../images/sequence_diagram_follow_UI.png)

![Sequence diagram - Follow request via API](../images/sequence_diagram_follow_API.png)

Note that while both versions use the same API, they use different endpoints.

## Dependencies

## System interactions

## Current state of the system
### SonarQube analysis summary

The following table summarizes key code quality metrics from SonarQube analysis:

| Metric                | Value                  |
|-----------------------|------------------------|
| Lines of Code (LOC)   | 1,591                  |
| Code Duplication      | 4.1%                   |
| Security Hotspots     | 8                      |
| Overall Rating        | A (Excellent quality)  |
| Cyclomatic Complexity | 216 (handlers: 151)    |
| Technical Debt        | ~1 hour 7 minutes      |

### Code climate

The following table summarizes key code quality metrics from Code Climate analysis:

| Metric                | Value                  |
|-----------------------|------------------------|
| Lines of Code (LOC)   | 1,912                  |
| Code Duplication      | 0 %                   |
| Overall Rating        | A (Excellent quality)  |
| Complexity | 299 (handlers: 196)    |
| Technical Debt        | ~1 day 2 hours      |

### Overall assesment
Both tools show a high complexity in the handlers module

## Orchestration
![Sequence diagram of IaC](../images/sequence_diagram_IaC.png)

## Deployment

### Allocation viewpoint

![Deployment diagram](../images/deployment_diagram.png)
