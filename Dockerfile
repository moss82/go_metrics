FROM golang:1.23-alpine

WORKDIR /app

COPY go_metrics/ .

RUN go mod init app
RUN go mod tidy
RUN go build -o app

EXPOSE 8080

CMD ["./app"]
