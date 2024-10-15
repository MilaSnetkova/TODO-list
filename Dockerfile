FROM golang:1.23 AS builder

RUN apt-get update && apt-get install -y build-essential

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 go build -o main ./main.go

FROM alpine:3.18

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

WORKDIR /app

COPY --from=builder /app/main .  
COPY --from=builder /app/web ./web 
COPY --from=builder /app/.env .

CMD ["./main"]