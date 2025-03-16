# Comparison of Relational Databases


## Preface
For the choice of database to use, we have compared some of the most popular RDBMS choices based on the Stack Overflow 2024 Developer Survey[^1]. For the scope of this exercise, we have not considered cloud databases or noSQL databases - since we're seeking to replace the current SQLite database.

A central requirement, is that the RDBM's license is open-source.

## Comparative Analysis

| **Database** | *SQLite* | *PostgreSQL* | *MySQL* | *Oracle* | *Microsoft SQL Server* | *MariaDB*
| --- | ---  | --- | --- | --- | --- | --- |
| **Popularity**    |33.1%[^1]|48.7%[^1]|40.3%[^1]|10.1%[^1]|25.3%[^1]|17.2%[^1]|
| **License**       |Public-Domain[^2]|Open-Source[^3]|Open-Source & Proprietary (Dual)[^4]|Proprietary[^5]|Proprietary[^6]|Open-Source[^7]|
| **Standards Compliance** | Low[^13] | Compliant[^13] | Limited[^13]  |  Unknown | Unknown | Assumed limited - due to being fork of MySQL. | 
| **Max Connections** | 1 | 500,000+[^11] | 100,000+[^11] | Unknown | Unknown | 200,000+[^11] | 
| **Horizontal Scaling** | No | Yes (Citus, Postgres-XL)[^11] | Yes (Out-of-the-box Replication, Sharding)[^11] | Unknown | Unknown | Yes (Galera Cluster)[^11] | 
| **Concurrency Handling** | None | Excellent[^11] | Moderate[^11] | Unknown | Unknown | Strong[^11] |

---

* **Performance** - Generally, we will not take quantitative performance benchmarks into account. Specifically because, most proprietary RDBMS include clauses in their licenses that prohibit users from publishing benchmarks e.g. in Oracle's standard license[^9].

* **Standards Compliance** - Refers to the degree to which the database complies with the international standard of NIST's ISO/IEC 9075[^10].

## Qualitative 


### MySQL
* Since MySQL was bought by Oracle in 2009, large adopters such as Fedora developers and OpenBSD has dropped MySQL[^14], due to concerns of it's acquisition, in favor of MariaDB. 
* MySQL follows a *dual-license* - where some featuers and plugins are available only under the proprietary editions[^13].
* There has been complaints that the development process for MySQL has slowed down significantly[^13].

Due to these concerns, MySQL was has not been chosen as our RDBMS.

### MariaDB vs PostgreSQL

Many resources compares MariaDB and PostgreSQL head-to-head as two robust choices of databases.
* MariaDB is often commended for its MySQL compatability and its advanced storage engines[^12] - often being mentioned as an easy-to-use highly stable drop-in replacement for MySQL[^17].
* PostgreSQL is praised for its SQL-compliance[^13], advanced security setup[^12], and its NoSQL features such as JSON, XML, and HStore[^15].

PostgreSQL was chosen in favor of MariaDB - due to its strong community adoption[^1], and the fact that the contributors of the project are not previously experienced in MySQL - and will therefore not appreciate the interoperability between MariaDB and MySQL.


## Sources

[^1]:[Stack Overflow 2024 Developer Survey - Databases](https://survey.stackoverflow.co/2024/technology#1-databases)

[^2]:[SQLite Licensing](https://www.sqlite.org/copyright.html)

[^3]:[PostgreSQL Licensing](https://www.postgresql.org/about/licence/)

[^4]:[MySQL Licensing](https://www.mysql.com/about/legal/licensing/oem/)

[^5]:[Microsoft SQL Server Licensing](https://www.microsoft.com/en-us/licensing/product-licensing/sql-server)

[^6]:[Oracle Licensing](https://www.oracle.com/a/ocom/docs/database-licensing-070584.pdf)

[^7]:[MariaDB License](https://mariadb.com/kb/en/mariadb-licenses/)

[^8]:[BenchANT Database Rankings](https://benchant.com/ranking/database-ranking)

[^9]:[Oracle Standard License](https://www.oracle.com/downloads/licenses/standard-license.html)

[^10]:[NIST ISO/IEC 9075:2023](https://blog.ansi.org/sql-standard-iso-iec-9075-2023-ansi-x3-135/)

[^11]:[AWS - MariaDB vs PostgreSQL](https://aws.amazon.com/compare/the-difference-between-mariadb-and-postgresql/#:~:text=MariaDB%20is%20a%20modified%20version,indexing%20for%20faster%20read%20performance.)

[^12]:[Medium - MySQL vs PostgreSQL vs MariaDB](https://medium.com/@peymaan.abedinpour/mariadb-vs-mysql-vs-postgresql-vs-sqlite-a-comprehensive-comparison-for-web-applications-0523cc3bc9d8#:~:text=PostgreSQL%20tends%20to%20perform%20better,applications%20with%20stringent%20data%20requirements.)

[^13]:[DigitalOcean - Comparison Of RDBMSs](https://www.digitalocean.com/community/tutorials/sqlite-vs-mysql-vs-postgresql-a-comparison-of-relational-database-management-systems)

[^14]:[Fedoraproject - Replace MySQL with MariaDB](https://fedoraproject.org/wiki/Features/ReplaceMySQLwithMariaDB)

[^15]:[Tooljet - MariaDB vs PostgreSQL](https://blog.tooljet.ai/mariadb-vs-postgresql-a-detailed-comparison-for-developers/)

[^17]:[Kinsta - MariaDB vs PostgreSQL14](https://kinsta.com/blog/mariadb-vs-postgresql/)