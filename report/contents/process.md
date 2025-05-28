# Process perspective

## CI/CD
- Github actions
- Python test were modified to work withour CI chains and now serves as a quiality gate. <!-- This is from a task description:  https://github.com/itu-devops/lecture_notes/blob/master/sessions/session_07/README_TASKS.md -->

## Monitoring 
<!-- Monitoring choice arguments is not a requirement (I checked)  -->
- Prometheus as an Echo middleware, with additional custom made metrics.
    - TODO: make list of custom metrics.
- Grafana
    - As of writing this the dashboards does not work due to swarm scaling. All pictures are from the day of the simulator stopping. 
    - Users:
        - admin user with password shared with the group
        - Helge and Mircea specific login as described on Teams.

Request and response dashboard:

Timeframe: last 30 minutes:
![alt text](image.png)

Timeframe: Last 2 days:
![alt text](image-1.png)

User action dashboards:
Timeframe: Last 7 days:
![alt text](image-2.png)

Virtual memory dashboard:

TODO: insert picture

- User side error logging was given by the Helge and Mircea in form of the Status and Simulator API errors graf. We were encouraged to just use this as our client side error monitoring. <!-- Helge said this in a lecture  -->

## Logging
- The ELK method was implemented but ultimatly scraped in favor of using loki/alloy that intergrate with Grafana which gather out logging and monitoring the same place. 
- Practical Principles:
    - TODO: A process should not worry about storage
    - TODO: A process should log only what is necessary
    - TODO: Logging should be done at the proper level 
    - Logs should be centralised: All logs can be found via Grafana->Drilldown->Logs

## Strategy for scaling and upgrades

## AI use

