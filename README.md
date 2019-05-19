# FinChat Server 

This is the backend for FinChat, a bot based chat.

It has following setup:
- Postgres Database: stores users, rooms and messages
- RabbitMQ: used to communicate between the server and bots
- GoServer: Rest API server. It provides authentication with JWT tokens and the API is documented with Swagger
- GoBot: a simple bot who listens to commands sent through rabbit and gets information from stooq.com

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

### Rabbit Management

- URL: http://localhost:15672
- User: rabbitmq
- Password: rabbitmq
- Queues: req_queue and resp_queue

### Go Server API

- URL: http://localhost:8001
- Swagger: http://localhost:8001/swagger/index.html