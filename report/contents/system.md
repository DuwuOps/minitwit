# System perspective

## Database (PostgreSQL)

Our setup includes two PostgreSQL databases: one for production and one for testing. Each runs on a separate, containerized droplet, with access restricted via a firewall to ensure security and isolation between environments (see Figure 1).

[PostgreSQL](https://www.postgresql.org/) was to replace the SQLite setup, due to strong SQL standards compliance[^13], high community adoption[^1], advanced features (e.g., JSON, HStore, Security)[^12][^15].

### Choice of Technology - Database

To replace our current SQLite setup, we compared leading relational databases based on the Stack Overflow 2024 Developer Survey[^1]. Only open-source, self-hosted RDBMSs were considered—excluding NoSQL and cloud services.

| **Database**     | **License**         | **Popularity**[^1] | **SQL Compliance**[^10] | **Max Connections**[^11] | **Scaling**[^11]             | **Concurrency**[^11] | **Notes** |
|------------------|---------------------|---------------------|--------------------------|---------------------------|-------------------------------|----------------------|-----------|
| **SQLite**       | Public Domain[^2]   | 33.1%               | Low                      | 1                         | No                            | None                 | File-based, lightweight |
| **PostgreSQL**   | Open-Source[^3]     | 48.7%               | High                     | 500,000+                  | Yes (Citus, Postgres-XL)      | Excellent            | Strong standards, JSON/XML support |
| **MariaDB**      | Open-Source[^7]     | 17.2%               | Moderate (MySQL fork)    | 200,000+                  | Yes (Galera Cluster)          | Strong               | MySQL-compatible, stable |
| **MySQL**        | Dual-License[^4]    | 40.3%               | Limited                  | 100,000+                  | Yes (native sharding/replica) | Moderate             | Dropped by Fedora[^14], dev slowdown[^13] |
| **Oracle DB**    | Proprietary[^5]     | 10.1%               | *Unknown*                  | *Unknown*                   | *Unknown*                       | *Unknown*              | -|
| **SQL Server**   | Proprietary[^6]     | 25.3%               | *Unknown*                  | *Unknown*                   | *Unknown*                       | *Unknown*              | -         |

> **Note**: Performance benchmarks are excluded due to license restrictions placed on benchmarking by licensing of proprietary DBMSs[^9].

While MariaDB is stable and compatible with MySQL, its advantages rely on prior MySQL knowledge, which our team lacks. MySQL was ruled out due to licensing issues and development concerns post-Oracle acquisition[^13][^14].


## References
[^1]: [Stack Overflow 2024 Developer Survey](https://survey.stackoverflow.co/2024/technology#1-databases)  
[^2]: [SQLite Licensing](https://www.sqlite.org/copyright.html)  
[^3]: [PostgreSQL License](https://www.postgresql.org/about/licence/)  
[^4]: [MySQL License](https://www.mysql.com/about/legal/licensing/oem/)  
[^5]: [Oracle Licensing](https://www.oracle.com/a/ocom/docs/database-licensing-070584.pdf)  
[^6]: [SQL Server Licensing](https://www.microsoft.com/en-us/licensing/product-licensing/sql-server)  
[^7]: [MariaDB License](https://mariadb.com/kb/en/mariadb-licenses/)  
[^9]: [Oracle Standard License](https://www.oracle.com/downloads/licenses/standard-license.html)  
[^10]: [ISO/IEC 9075 SQL Standard](https://blog.ansi.org/sql-standard-iso-iec-9075-2023-ansi-x3-135/)  
[^11]: [AWS - MariaDB vs PostgreSQL](https://aws.amazon.com/compare/the-difference-between-mariadb-and-postgresql/#:~:text=MariaDB%20is%20a%20modified%20version,indexing%20for%20faster%20read%20performance.)  
[^12]: [Medium - RDBMS Comparison](https://medium.com/@peymaan.abedinpour/mariadb-vs-mysql-vs-postgresql-vs-sqlite-a-comprehensive-comparison-for-web-applications-0523cc3bc9d8)  
[^13]: [DigitalOcean - RDBMS Comparison](https://www.digitalocean.com/community/tutorials/sqlite-vs-mysql-vs-postgresql-a-comparison-of-relational-database-management-systems)  
[^14]: [Fedora Drops MySQL](https://fedoraproject.org/wiki/Features/ReplaceMySQLwithMariaDB)  
[^15]: [Tooljet - MariaDB vs PostgreSQL](https://blog.toolje)