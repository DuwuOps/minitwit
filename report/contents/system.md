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
Two sequence diagrams have been created to show the flow of information through the system, from a "Follow" request by a user, to the system's returned response. 

The first version shows the processes involved when the request is sent via. the UI, whereas the second version shows the processes involved when sent via. the API. 

![Sequence diagram - Follow request via UI](../images/sequence_diagram_follow_UI.png)

![Sequence diagram - Follow request via API](../images/sequence_diagram_follow_API.png)

Note that the two versions use different endpoints to interact with the same API.

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
To streamline deployment, Docker, docker-compose, Docker Swarm, and Terraform are used. 

The Dockerfile copies all source code from the `src` package to a binary image of the program.

There are two docker-compose files, `docker-compose.yml` and `docker-compose.deploy.yml`. Both define the six central services of the system: app, prometheus, alloy, loki, grafana, and database. 

`docker-compose.yml` is needed for local deployment. It uses localhost IP-adresses and has default usernames and passwords. 

`docker-compose.deploy.yml` is used for remote deployment. It builds on `docker-compose.yml`, but replaces information where relevant. 
It specifies the configuration of a Docker Swarm with one manager and two workers: The app runs on two worker replicas, while logging and monitoring services are constrained to only run on the manager node (though alloy collects logs from everywhere). This enables horizontal scaling. 

Infrastructure-as-Code is used simplify the setup of the Docker Swarm remotely. Terraform files can be found in `.infrastructure/infrastructure-as-code`. Automatic deployment via. Terraform works as illustrated in the sequence diagram below. 

![Sequence diagram of IaC](../images/sequence_diagram_IaC.png)

## Deployment

### VPS

To host the system on a remote server, [DigitalOcean](https://www.digitalocean.com/products/droplets) was chosen as the VPS provider. This choice was based on pricing (see @tbl:vps-comparison), its apparent ease-of-use[@Quinn_2022] [@aliamin7] [@Finder_2023], its familiarity to the group.

| **VPS**                   | **DigitalOcean**                      |  **Microsoft Azure**          | **Oracle**                        | **AWS (Lightsail)**                   |
|---------------------------|---------------------------------------|-------------------------------|-----------------------------------|---------------------------------------|
| **Virtual Machine Price** | ca. $12/mo [@digitalocean_price]      | ca. $15/mo [@azure_price]     | $13/mo [@oracle_price]            | ca. $12/mo [@aws_lightsail_price]     |
| **Storage Price**         | 50GB included [@digitalocean_price]   | ca. $5 (64GB) [@azure_price]  | ca. $2.5 (50GB) [@oracle_price]   | ca. $12/mo [@aws_lightsail_price]  |
| **Total Price**           | ca. $12/mo                            | ca. $20/mo                    | ca. $15.5/mo                      | ca. $12/mo                            |

: Price comparison of VPS providers. {#tbl:vps-comparison}


### Allocation viewpoint

![Deployment diagram](../images/deployment_diagram.png)
