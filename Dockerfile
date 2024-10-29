FROM golang:1.23.1-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./main.go

CMD ["./main"]

expose 8080