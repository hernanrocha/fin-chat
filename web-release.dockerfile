FROM ubuntu

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY dist/fin-chat-service /app/

WORKDIR /app

EXPOSE 5000

CMD ["./fin-chat-service"]