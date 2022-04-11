# Bank

## Purpose
Execercies in go

Howto run
- create database bank
- create tables using [sql/initialize.sql](sql/initialize.sql)
- export DATABASE_URL=postgres://testuser:12345@localhost:5432/bank
- go run main.go
  
## BOM

![Bank model](Bank.png)

- **Target**, the target of a transaction. For instance: Gift, Petrol, Electricity, Morgage
- **Transaction**, the transfer of funds from an account to an account for a specific Target
- **Account**, a known account (by id), of an unknown account by account.numer

## Database
### Access database from go
See https://go.dev/doc/tutorial/database-access

### Create database
See https://kb.objectrocket.com/postgresql/how-to-show-databases-in-postgresql-848 

To create a postgres database
- login as postgresql
- create the database

```sql
create database bank

---
--- 
---
USE bank
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM PUBLIC;
GRANT ALL ON ALL TABLES IN SCHEMA public TO testuser;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO testuser;


```
create tables using [sql/initialize.sql](sql/initialize.sql) script



## FAQ
### Howto install a module
To install logrus

```go
go mod init
go mod tidy
```
[Since go 1.17 go get no longer supported](https://go.dev/doc/go-get-install-deprecation)


## References
- [go](https://go.dev/)
- [logrus, example](https://golangdocs.com/logging-in-go-logrus-package)
- [logrus, logging framework](https://github.com/Sirupsen/logrus)
- [postgres driver](https://github.com/lib/pq), this one is old
- [postgres driver] (https://github.com/jackc/pgx)
- [postgres example](https://golangdocs.com/golang-postgresql-example)
- [setup environment](https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-ubuntu-18-04)
