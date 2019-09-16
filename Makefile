build:
	go build -o fin-chat-service service/main.go
	go build -o fin-chat-bot bot/main.go

lambdaupload:
	zip function.zip fin-chat-bot 
	aws lambda update-function-code --function-name StooqParse --zip-file fileb://function.zip

runservice:
	go run service/main.go

runbot:
	go run bot/main.go

genswagger:
	go get -u github.com/swaggo/swag/cmd/swag
	swag init

test:
	go test github.com/hernanrocha/fin-chat... --cover -count=1

dockerbuild:
	docker build -t 089576757282.dkr.ecr.us-east-2.amazonaws.com/finchat -f web-release.dockerfile .
	docker push 089576757282.dkr.ecr.us-east-2.amazonaws.com/finchat
#	docker build -t fin-chat-bot -f bot-release.dockerfile .

dockerservice:
	docker run --env RABBIT_CONNECTION=amqp://rabbitmq:rabbitmq@192.168.1.36:5672/ \
	--env DB_CONNECTION="host=192.168.1.36 port=15432 user=postgres password=postgres dbname=finchat sslmode=disable" \
	fin-chat-web

dockerbot:
	docker run --env RABBIT_CONNECTION=amqp://rabbitmq:rabbitmq@192.168.1.36:5672/ fin-chat-bot