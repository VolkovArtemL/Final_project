FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/todo-app .
COPY --from=builder /app/web ./web

RUN mkdir -p /app/data

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db

CMD ["./todo-app"]