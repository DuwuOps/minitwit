# Process Perspective

## CI/CD
- Github actions
- Python test were modified to work withour CI chains and now serves as a quiality gate. <!-- This is from a task description:  https://github.com/itu-devops/lecture_notes/blob/master/sessions/session_07/README_TASKS.md -->

## Monitoring 
<!-- Monitoring choice arguments is not a requirement (I checked), but added anyway since we had it.  -->
- Prometheus as an Echo middleware, with additional custom made metrics.
    - TODO: make list of custom metrics.
    - Was chosen on the background of:
        - Demonstrated in Class
        - Prior experience for members of the g
- Grafana
    - As of writing this the dashboards does not work due to swarm scaling. All pictures are from the day of the simulator stopping. 
    - Users:
        - Admin user with password shared with the group.
        - Helge and Mircea specific login as described on Teams.

Whitebox Request and response monitoring dashboard:

Timeframe: last 30 minutes:
![Request and response dashboard last 30 minutes](/report/images/monitoring-response-request-t2d.png)

Timeframe: Last 2 days:
![Request and response dashboardLast 2 days](/report/images/monitoring-response-request-t30.png)

Whitebox User action dashboards monitoring:
Timeframe: Last 7 days:
![User action dashboards Last 7 days](/report/images/monitoring-user-actions-t7d.png)

Whitebox Virtual memory dashboard monitoring:
Timframe: last 5 minutes:
![Virtual Memory dashbord Last 5 minutes](/report/images/monitoring-VM-usage-t5.png)

- Black box user side error monitoring was given by the Helge and Mircea in form of the Status and Simulator API errors graf. We were encouraged to just use this as our client side error monitoring. <!-- Helge said this in a lecture  -->

## Logging
- The ELK method was implemented but ultimatly scraped in favor of using loki/alloy that intergrate with Grafana which gather out logging and monitoring the same place. 
- Practical Principles:
    - TODO: A process should not worry about storage
    - TODO: A process should log only what is necessary
    - TODO: Logging should be done at the proper level 
    - Logs should be centralised: All logs can be found via Grafana->Drilldown->Logs

## Strategy for scaling and upgrades
- Used docker swarm using docker stack so that we can leverage the docker compose setup that was already made. Some changes were made to accomidate the swarm set up. These are a network overlay, setting how many replicas per service, where nesecary setting where the service should be placed, and update configs.
- The update config for the minitwit application is set so that it updates one at a time. This is set as we only have two instances of minitwit and if an update fails we don't want more than one instance to be down. On failure we do a rollback. 

## AI use
Throughout the development process, all team members leveraged artificial intelligence tools to varying degrees and for diverse applications. The primary AI systems employed included ChatGPT, Claude, DeepSeek, and GitHub Copilot. Team members provided contextual information regarding code issues or implementation challenges, utilizing AI-generated responses as foundational guidance for problem-solving methodologies rather than direct solution implementation. This methodology facilitated the identification of potential problem domains and remediation strategies while preserving critical assessment of AI-derived recommendations. In accordance with transparency requirements, AI tools have been formally acknowledged as co-authors in relevant version control commits where their contributions influenced the development process.  (This paragraf was written using AI lol)