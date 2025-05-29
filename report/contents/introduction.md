# Introduction
## System Depiction

An overview of the system is provided in Figure 1.

![Informal system depiction diagram and a legend](../images/Informal%20System%20Depiction.png){#fig:informal-system-depiction width=80% placement=H}

## Tech-Stack Overview

| **Topic** | **Tech Choice** | **Section** |
| --- | --- | ----- |
| Refactoring | [GoLang](https://go.dev/), [Echo](https://echo.labstack.com/) | **2.1 Programming Language** |
| Orchestrazation | [Docker](https://www.docker.com/) | **2.2 Orchestrazation** |
| Deployment | [DigitalOcean](https://www.digitalocean.com/) | **2.3 VPS** |
| CI/CD | [GitHub Actions](https://github.com/features/actions) | **3.1 CI/CD** |
| Database | [PostgreSQL](https://www.postgresql.org/) | **2.4 Database** |
| Monitoring | [Prometheus](https://prometheus.io/), [Grafana](https://grafana.com/) | **3.2 Monitoring** |
| CI Static Analysis | [golangci-lint](https://github.com/golangci/golangci-lint), [Dependabot](https://github.com/dependabot) | **3.1 CI/CD** |
| Maintainability | [SonarQube](https://www.sonarsource.com/products/sonarqube/), [CodeClimate](https://codeclimate.com/) | **2.5 Maintainability** |
| Logging | [Loki](https://grafana.com/docs/loki/latest/), [Alloy](https://grafana.com/docs/alloy/latest/), [Grafana](https://grafana.com/) | **3.3 Logging** |
| Scalability | [Docker Swarm](https://docs.docker.com/engine/swarm/) | **3.4 Scaling** |
| Security | [CodeQL](https://codeql.github.com/), [Dependabot](https://github.com/dependabot) | **3.5 Security** |  

Table: Overview of tech-stack.