# cat-api

Solution to [challenge](challenge.pdf).

## Requirements

- Go 1.15

## Setup

In development, you can use docker-compose to build database infrastructure:

```
docker-compose up -d
```

Copy the environment variables file and adjust it accordingly:

```
cp .env.example .env
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
curl 127.0.0.1:8080/breeds?name=sib -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-IDcSemACt8x4iTMCda8Yhe3iZaWbvV5XKSTbuAn0M'
```

Or you can test the entire app using the builtin test tool:

```
go test
```

## Author

| [<img src="https://avatars2.githubusercontent.com/u/1824706?s=120&v=4" width=120><br><sub>@willystadnick</sub>](https://github.com/willystadnick) |
| :---: |
