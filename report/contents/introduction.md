# Introduction
## System Depiction

An informal depiction of the system as an initial overview is provided in @fig:informal-system-depiction.

![Informal system depiction and legend](../images/Informal%20System%20Depiction.png){#fig:informal-system-depiction width=100% placement=H}

## Tech-Stack Overview

| **Topic** | **Tech Choice** | **Section** |
| --- | --- | ----- |
| Refactoring | [GoLang](https://go.dev/), [Echo](https://echo.labstack.com/) | **2.1 Minitwit** |
| Orchestrazation | [Docker](https://www.docker.com/) | **2.2 Orchestration** |
| Deployment | [DigitalOcean](https://www.digitalocean.com/) | **2.3 Deployment** |
| CI/CD | [GitHub Actions](https://github.com/features/actions) | **3.1 CI/CD** |
| Database | [PostgreSQL](https://www.postgresql.org/) | **2.4 Database** |
| Monitoring | [Prometheus](https://prometheus.io/), [Grafana](https://grafana.com/) | **3.2 Monitoring** |
| CI Static Analysis | [golangci-lint](https://github.com/golangci/golangci-lint), [Dependabot](https://github.com/dependabot) | **3.1 CI/CD** |
| Maintainability | [SonarQube](https://www.sonarsource.com/products/sonarqube/), [CodeClimate](https://codeclimate.com/) | **2.1 Minitwit** |
| Logging | [Loki](https://grafana.com/docs/loki/latest/), [Alloy](https://grafana.com/docs/alloy/latest/), [Grafana](https://grafana.com/) | **3.3 Logging** |
| Scalability | [Docker Swarm](https://docs.docker.com/engine/swarm/) | **3.4 Strategy for Scaling and Upgrades** |
| Security | [CodeQL](https://codeql.github.com/), [Dependabot](https://github.com/dependabot) | **4 Security Assessment** |  

Table: Overview of tech-stack.