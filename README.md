# License Server
A simple application for license key verification written in Golang.

## Deployment
- **Requirement**
    - Mysql / Sqlite3 / PostgreSQL

## Development

- **Requirement**
    - Go >= 1.13
    - Mysql / Sqlite3 / PostgreSQL
    - dep [Dependency management for Go](https://golang.github.io/dep/)

- Clone the repo
```
git clone https://github.com/hrshadhin/license-server.git
cd license-server
```
- copy `.env.example` file to `.env` and change configuration according to your need.

- **Install dependency** `dep ensure`
- **Run** `go run main.go`
- **Build**
    ```
    go build -o bin/license-server main.go
    cp .env.example bin/.env
    cd bin
    ./license-server
    ```