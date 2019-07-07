runservice:
	go run service/main.go

runbot:
	go run bot/main.go

genswagger:
	go get -u github.com/swaggo/swag/cmd/swag
	swag init

build:
	docker build -t fin-chat-web -f web-release.dockerfile .
	docker build -t fin-chat-bot -f bot-release.dockerfile .

dockerservice:
	docker run --env RABBIT_CONNECTION=amqp://rabbitmq:rabbitmq@192.168.1.36:5672/ \
	--env DB_CONNECTION="host=192.168.1.36 port=15432 user=postgres password=postgres dbname=finchat sslmode=disable" \
	fin-chat-web

dockerbot:
	docker run --env RABBIT_CONNECTION=amqp://rabbitmq:rabbitmq@192.168.1.36:5672/ fin-chat-bot