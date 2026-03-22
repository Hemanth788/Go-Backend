# LOG

## 20-03-2026

- init
- DB Design, SQL queries
- install sqlc, configure and use it to generate CRUD boilerplate for all queries
- spin up PostgreSQL in a docker container, set it up
- Makefile to make my life easy
- test the generated code

## 21-03-2026

### Morning

- create a DB store for DB transactions
- use it to write a money transfer transaction
- test it

### Afternoon

- avoiding deadlock in 2-way money flow concurrent DB transactions
- test it
- I in ACID

### Evening

- GitHub Actions for automated test runs
- REST API with gin
- use viper to load config from .env

### Night

- gomock to mock DB
- table driven testing to increase coverage

## 22-03-2026

### Evening

- add users table with foreign keys to existing account table, 1 user can have 1 account per currency
- generate boilerplate SQL-GO code for utility, mocks for testing
- better error handling for SQL violations

### Night

- add createUser route and handler
- add password hashing
- test hashing
- test the createUser handler
- improve createUser test by writing a custom matcher
