# Process Perspective

## CI/CD (GitHub Actions)

[GitHub Actions](https://github.com/features/actions) was chosen based on its simplicity, familiarity, and pricing (free for public repositories). A motivating factor, was the suite of services supported natively in GitHub. Therefore:

* GitHub *Secrets & Variables* was used to store environment variables, and deployment keys..
* GitHub *Tags*, *Releases*, and *Artifact Storage* were utilized, in order to create a clean version history of our application. 
* GitHub integrations, such as *Dependabot*, *SonarCube*, and *Webhooks*. 

**Write about**:
* *Linters*
* *Python testing*

A total of **7** workflows are established, and can be found under `/.github/workflows/`. These are:

| File    | Purpose | Runs when |
| -------- | ------- |------- |
| `continous-development.yml`  | Primary CI/CD flow for continous integration & delivery. Consists of steps `Tests`, `Build & Push`, `Release`, and `Deploy` | Any changes to `main` |
| `codeql.yml` | Analyzes GoLang source code using [CodeQL Analysis tool](github/codeql-action/analyze@v3) | Any push to `main`. Any pull-request to `main`. Once a week in cron-job. | 
| `generate-report.yml`| Generates `report.pdf` from markdown files & images in `/report/`     | Any changes to `/report/*` recursively.   | | 
| `linter-workflow.yml`| Runs [golangci-lint](https://github.com/golangci/golangci-lint) linter on GoLang source code. Configured by `/golangci.yml`. | Push to `main` or any action to pull-requests. | 
| `pull-request-tests.yml` | Runs python tests. | Any actions to pull-requests |
| `test-deployment.yml`    | Secondary CI/CD flow for continous integration & delivery against TEST-environment. Consists of `Tests`, `Build & Push`, `Deploy` | On push with tag `test-env*`. |Identical to `continous-development.yml` but does not include a `Release` step. | 
| `sonarcube_analysis.yml` | Analyses GoLang source code using SonarCloud. | On pull-requests to `main`. |
> **Table**: GitHub Action workflows employed.

![Visualization of continous-development.yml](../images/github_actions-continuous-development.png)
> **Figure**: Visualization of `continous-development.yml`

#### Choice of CI/CD

Since GitHub was chosen as the git repository management site, options such as [GitLab CI/CD](https://docs.gitlab.com/ci/) and [BitBucket Pipelines](https://www.atlassian.com/software/bitbucket/features/pipelines) were discarded as candidates, as they are specific to alternative git repository management sites. As such, the choice was between GitHub's native [GitHub Actions](https://github.com/features/actions) or CI/CD systems agnostic to repository management sites. 

Furthermore, commercial automation tools such as [Azure DevOps](https://azure.microsoft.com/en-us/products/devops) and [TeamCity](https://www.jetbrains.com/teamcity/) were discarded due to the pricing and limitations of their free plans. An overview of the comparison performed can be seen in the table below.

| **CI/CD Tool / Platform**       | **GitHub Actions**                                                                 | **Jenkins**                                                         | **Azure DevOps**                                | **TeamCity (JetBrains)**                       |
|----------------------------------|-------------------------------------------------------------------------------------|----------------------------------------------------------------------|--------------------------------------------------|--------------------------------------------------|
| **Ease-of-use**                 | Simple [¹](#ref1)                                                                  | Medium [¹](#ref1)                                                   | *Undetermined*                                  | *Undetermined*                                  |
| **Version Control**             | Native GitHub Integration [²](#ref2)                                               | Agnostic [²](#ref2)                                                 | Agnostic [²](#ref2)                             | Agnostic [²](#ref2)                             |
| **Hosting**                     | Primarily cloud-based [²](#ref2)                                                   | Self-hosted [²](#ref2)                                              | Cloud-based [²](#ref2)                          | Cloud-based or self-hosted [²](#ref2)          |
| **Pricing Model**               | Free for public repositories, tiered for private [²](#ref2)                        | Open-source (MIT License), only cost is for hosting [²](#ref2)      | Commercial with a limited free tier [²](#ref2)  | Commercial [²](#ref2)                          |

> **Table**: Comparison between CI/CD systems.

It was decided that time-to-production, in the case of establishing working CI/CD pipelines, was the biggest priority. As an alternative, the self-hosted automation system [Jenkins](https://www.jenkins.io/) was considered, but the perceived learning curve along with the self-hosted infrastructure setup [¹](#ref1) dissuaded it as the choice of CI/CD system.

## Monitoring 
<!-- Monitoring choice arguments is not a requirement (I checked), but added anyway since we had it.  -->
- We use Prometheus as an Echo middleware, with additional custom made metrics to scrape our application every 5 seconds.
    - Custom metrics: 
        - User follower (gauge)
        - User followees (gauge)
        - VM CPU usage (gauge)
        - Messages posted (by time) (counter)
        - Messages posted (by user) (gauge)
        - Mesages flagged (by user) (gauge)
        - New user (counter)
        - Total users (gauge)
    <!-- Der er det her i vores kode som jeg ikke helt ved hvad er og om det burde komme med her:
    MemoryUsage.WithLabelValues("UsedPercent").Set(vm.UsedPercent)
                MemoryUsage.WithLabelValues("Used").Set(float64(vm.Used))
                MemoryUsage.WithLabelValues("Available").Set(float64(vm.Available))
                MemoryUsage.WithLabelValues("Total").Set(float64(vm.Total))
    -->
    - Prometheus was chosen on the background of:
        - Demonstrated in Class
        - Easy integration with golang/echo via. middleware
        - Wide spread usage and easy to integrate with e.g. Grafana
        - Free to use
- Grafana
    - As of writing this, the dashboards does not work due to swarm scaling. All pictures are from the day of the simulator stopping. 
    - Users:
        - Admin user with password shared with the group.
        - Helge and Mircea specific login as described on Teams.
    - Was chosen on the background of:
        - Demonstrated in Class
        - Rich Visualization 
        - Free to use 
        <!-- Har fjernet Alerting Functionality (fra overleaf listen) fordi vi alerter gennem botten og ikke gennem grafana -->

### Grafana Dashboards

**Whitebox Request and response monitoring dashboard:**

Timeframe: last 30 minutes:

![Request and response dashboard last 30 minutes](../images/monitoring-response-request-t2d.png)

Timeframe: Last 2 days:

![Request and response dashboardLast 2 days](../images/monitoring-response-request-t30.png)

**Whitebox User action dashboards monitoring:**

Timeframe: Last 7 days:

![User action dashboards Last 7 days](../images/monitoring-user-actions-t7d.png)

**Whitebox Virtual memory dashboard monitoring:**

Timframe: last 5 minutes:

![Virtual Memory dashbord Last 5 minutes](../images/monitoring-VM-usage-t5.png)

### Black box monitoring

Black box user side error monitoring was given by the Helge and Mircea in form of the Status and Simulator API errors graf. We were encouraged to just use this as our client side error monitoring. <!-- Helge said this in a lecture  -->

### DigitalOcean monitoring

DigitalOcean provides some monitoring capabilities (Bandwidth, CPU usage, and Disk I/O). This did help to identify an attack. More on that [Insert refrence here]

### Alert System
An alert system was set up via a Discord bot that on the server via a cronjob that checks every 5 miniutes. If the application is not up it sends a Discord message and tags everyone on our group server. 
![Alert bot example](../images/alert-example.png)
 <!-- Jeg syntes det var sødt med et billede af vores discord, hvis nogen er uenige så bare fjern <3> -->

## Logging
- The ELK method was implemented but ultimatly scraped in favor of using loki/alloy that intergrate with Grafana which gather our logging and monitoring the same place. 
- Practical Principles:
    - TODO: A process should not worry about storage
    - TODO: A process should log only what is necessary
    - TODO: Logging should be done at the proper level: Mention emoji use
    - Logs should be centralised: All logs can be found via Grafana->Drilldown->Logs

## Strategy for scaling and upgrades
- We used docker swarm with docker stack so that we can leverage the docker compose setup that was already made. Some changes were made to accommodate the swarm set up. These are: a network overlay, setting how many replicas per service, where nesecary setting where the service should be placed, and update confirguration.
- The update config for the minitwit application is set so that it updates one at a time. This is set as we only have two instances of minitwit and if an update fails we don't want more than one instance to be down. On failure we do a rollback. 
- We do rolling updates as this is nativly supported on docker swarm.

## AI use
Throughout the development process, all team members leveraged artificial intelligence tools to varying degrees and for diverse applications. The primary AI systems employed included ChatGPT, Claude, DeepSeek, and GitHub Copilot. Team members provided contextual information regarding code issues or implementation challenges, utilizing AI-generated responses as foundational guidance for problem-solving methodologies rather than direct solution implementation. This methodology facilitated the identification of potential problem domains and remediation strategies while preserving critical assessment of AI-derived recommendations. In accordance with transparency requirements, AI tools have been formally acknowledged as co-authors in relevant version control commits where their contributions influenced the development process.  (This paragraf was written using AI lol)
<!-- Det her er bare mig der syntes det er thihi fnis as få AI til at skrive AI afsnittet. Hvis nogle er uenige så bare lav det om <3 > -->


---

### References

<a name="ref1">¹</a>: githubactions_vs_jenkins  
<a name="ref2">²</a>: 20_cicd_comparison