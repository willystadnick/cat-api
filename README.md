# cat-api

Cat API made with Go 1.15

## Requirements

- Go 1.15

## Setup

```
docker-compose up -d
```

## Run

```
go run main.go
```

## Test

With the app running, you can test the endpoints with curl:

```
curl 127.0.0.1:8080/ping
curl 127.0.0.1:8080/login -d '{"password":"@#$RF@!718","username":"admin"}'
```

Or you can test the entire app using the builtin test tool:

```
go test
```
