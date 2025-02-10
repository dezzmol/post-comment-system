FROM golang:1.23.1 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./server.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]