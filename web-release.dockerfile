FROM scratch
COPY dist/fin-chat-service /app/
WORKDIR /app
EXPOSE 5000
CMD ["./fin-chat-service"]