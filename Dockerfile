# build stage
FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz \
    | tar -xvz -C /app

# run stage
FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
# RUN chmod +x /app/migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
