# FinChat Server 

This is the backend for FinChat, a bot based chat.

It's deployed in AWS using the following services:
- RDS: relational database to store users, rooms and messages
- SNS/SQS: used to communicate between the server and bot
- ECS: **GoServer** (Rest API server) runs on a container in an ECS cluster
It provides authentication with JWT tokens and the API is documented with Swagger
- Lambda: **GoBot** is a lambda function triggered by SNS messages. It gets and parses information from stooq.com and then is published in a SQS queue.

## CircleCI

[![CircleCI](https://circleci.com/gh/hernanrocha/fin-chat.svg?style=svg)](https://circleci.com/gh/hernanrocha/fin-chat)

## Run with Docker

```
docker-compose up
```

## How to Use it

After you run all containers, you could interact with the API through Postman or Swagger. The file _FinChat.postman_collection.json_ is a Postman collection with all supported endpoints

Basically, you have to create an account with _/register_ and then generate a JWT token with _/login_.
The rest of the endpoints should be called with the value 'Bearer <TOKEN>' in _Authorization_ header.

## Links

### Adminer

- URL: http://localhost:18080
- DB Type: PostgreSQL
- DB Host: db
- DB Port: 5432
- DB User: postgres
- DB Password: postgres
- DB Name: finchat


### Go Server API

- URL: http://finchat-loadbalancer-1974477651.us-east-2.elb.amazonaws.com/
- Swagger: http://finchat-loadbalancer-1974477651.us-east-2.elb.amazonaws.com/swagger/index.html


### Commands

Run server: `go run service/main.go`
Run bot: `go run service/main.go`

### Generate documentation 

```sh
go get -u github.com/swaggo/swag/cmd/swag
cd service
swag init
cd ..
```

### TODOs

- Add health endpoint
- Add monitoring tools
- Add clean shutdown (wait for AWS signal)
- Add HTTPS support
- Run GinGonic on release mode
- Read JWT secret from environment variable
- Create load tests