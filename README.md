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
- **Build & Run**
    ```
    make 
    ./bin/license-server
    ```
## Cross Platform Build
- Cross compilation is hard, and docker is help us in that way! Install docker and pull
    docker image `docker pull karalabe/xgo-latest` and install a go package `go get github.com/karalabe/xgo`
- For build most of the platforms binary use 
    ```
    xgo github.com/hrshadhin/license-server
    ```
- Or specific platform
    ```
    xgo --targets=linux/amd64 github.com/hrshadhin/license-server 
    ```
- After build is finished you should have all platforms binary in your
current directory.
- More build details find [here](https://github.com/karalabe/xgo)
    
## API Docs
[HERE](https://documenter.getpostman.com/view/8901126/SVzz2eF1)