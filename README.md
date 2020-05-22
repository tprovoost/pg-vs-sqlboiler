# pg-vs-sqlboiler

This project is a comparison to evaluate the complexity and the performances of two "ORM" modules called pg and sql boiler for a PostgreSQL database.

# Installation

First and foremost, I would advise to clone this project into the `GOPATH/src` for a maximum of simplicity. If you know your ways with Go, be free to do it the way you want. ;-)

## Database

First you need a database on your machine. The database is called pgguide and will require you to install the dump: http://postgresguide.com/setup/example.html.

## SQL Boiler
SQLBoiler is a tool to generate a Go ORM tailored to your database schema. Please check the following installation instructions: https://github.com/volatiletech/sqlboiler#getting-started

In my case, I decided to install SQL Boiler directly by running the commands within my `GOPATH`.

This is my file:

```toml
output = "./src/orm_compare/database_models"
wipe = true
no_tests = true

[psql]
dbname = "pgguide"
host   = "localhost"
port   = 5432
user = "thomasprovoost"
sslmode = "disable"
```

If you have any problem with imports, check the [FAQ of SQL Boiler](https://github.com/volatiletech/sqlboiler#missing-imports-for-generated-package).

# PG

PG is a PostgreSQL client and an ORM for Golang.

https://github.com/go-pg/pg

## Installation

Everything should already be in the `go.mod` file. Inc ase of missing imports, go to the main folder of the project and run the command `go get -u -t` in a terminal.