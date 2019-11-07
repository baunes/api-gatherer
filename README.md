# API Gatherer
[![Go Report Card](https://goreportcard.com/badge/github.com/baunes/api-gatherer)](https://goreportcard.com/report/github.com/baunes/api-gatherer)

**Disclaimer**: The repository is just a side project to play with golang.\
Don't use for production.

### Run all test in the project
    go test -cover -v github.com/baunes/api-gatherer/...

### Run
    docker-compose -f docker/docker-compose.yml up -d

    go run main.go -url <url> -db.database <db_name> -db.collection <col_name>

Example:

    go run main.go -url https://quote-garden.herokuapp.com/quotes/random -db.database quotes -db.collection random
