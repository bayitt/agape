FROM golang:1.22-alpine

RUN mkdir -p /app

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ./agape

EXPOSE 8080

CMD ["./agape"]