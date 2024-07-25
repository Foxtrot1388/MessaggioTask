FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /usr/local/go/src/
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
RUN go build -mod=mod -o app.exe main.go

#lightweight docker container with binary
FROM alpine:latest

ARG POSTGRES_USER=$POSTGRES_USER
ARG POSTGRES_PASSWORD=$POSTGRES_PASSWORD
ARG POSTGRES_HOST=$POSTGRES_HOST
ARG POSTGRES_PORT=$POSTGRES_PORT
ARG POSTGRES_DB=$POSTGRES_DB
ARG LOG_LEVEL=$LOG_LEVEL
ARG KAFKA_HOST=$KAFKA_HOST
ARG KAFKA_TOPIC=&KAFKA_TOPIC

RUN apk --no-cache add ca-certificates

COPY --from=builder /usr/local/go/src/app.exe /
COPY /migrations /migrations

EXPOSE 8080

CMD [ "/app.exe"]