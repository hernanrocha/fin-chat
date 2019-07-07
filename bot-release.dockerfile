FROM golang:1.12-alpine as builder

# To fix go get and build with cgo
RUN apk add --no-cache --virtual .build-deps \
    bash \
    gcc \
    git \
    musl-dev

RUN mkdir build
WORKDIR /build

COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o botserver ./bot/main.go
RUN adduser -S -D -H -h /build botserver
USER botserver

FROM scratch
COPY --from=builder /build/botserver /app/
WORKDIR /app
EXPOSE 5000
CMD ["./botserver"]