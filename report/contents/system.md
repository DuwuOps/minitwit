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

## Database (PostgreSQL)

Our setup includes two PostgreSQL databases: one for production and one for testing. Each runs on a separate, containerized droplet, with access restricted via a firewall to ensure security and isolation between environments (see Figure 1).

[PostgreSQL](https://www.postgresql.org/) was to replace the SQLite setup, due to strong SQL standards compliance [@do_dbcomparison], high community adoption [@stackoverflow_survey_2024], advanced features (e.g., JSON, HStore, Security) [@tooljet_mariavspostgres], [@Medium_Peymaan_DB_Comparison].

### Choice of Technology - Database

To replace our current SQLite setup, we compared leading relational databases based on the Stack Overflow 2024 Developer Survey [@stackoverflow_survey_2024]. Only open-source, self-hosted RDBMSs were considered—excluding NoSQL and cloud services.

| **Database** | **SQLite** | **PostgreSQL** | **MySQL** | **Oracle** | **SQL Server** | **MariaDB** |
| --- | --- | --- | --- | --- | --- | --- |
| **Popularity** | 33.1% [@stackoverflow_survey_2024] | 49.7% [@stackoverflow_survey_2024] | 40.3% [@stackoverflow_survey_2024] | 10.1% [@stackoverflow_survey_2024] | 25.3% [@stackoverflow_survey_2024] | 17.2% [@stackoverflow_survey_2024] |
| **License** | Public-Domain [@sqlite_license] | Open-Source [@postgresql_license] | Open-Source & Proprietary [@MySQL_license] | Proprietary | Proprietary [@microsoftsqlserver_license] | Open-Source [@mariadb_license] |
| **Standards Compliance** [@SQL_Standard_ISO] | Low [@do_dbcomparison] | Compliant [@do_dbcomparison] | Limited [@do_dbcomparison] | *Unknown* | *Unknown* | Fork of MySQL; Assumed limited |
| **Max Connections** | 1 | 500,000+ [@Medium_Peymaan_DB_Comparison] | 100,000+ [@Medium_Peymaan_DB_Comparison] | *Unknown* | *Unknown* | 200,000+ [@Medium_Peymaan_DB_Comparison] |
| **Horizontal Scaling** | No | Yes [@Medium_Peymaan_DB_Comparison] | Yes [@Medium_Peymaan_DB_Comparison] | *Unknown* | *Unknown* | Yes [@Medium_Peymaan_DB_Comparison] |
| **Concurrency Handling** | None | Excellent [@Medium_Peymaan_DB_Comparison] | Moderate [@Medium_Peymaan_DB_Comparison] | *Unknown* | *Unknown* | Strong [@Medium_Peymaan_DB_Comparison] |

Table: Comparison of RDBMSs.
**Note**: Performance benchmarks are excluded due to license restrictions placed on benchmarking by licensing of proprietary DBMSs [@Oracle_Network_License].
 
MySQL was ruled out due to licensing issues and development concerns post-Oracle acquisition [@Fedora_MariaDB], [@do_dbcomparison].