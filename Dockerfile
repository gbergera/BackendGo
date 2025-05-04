FROM golang:latest


RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN /go/bin/swag init

RUN go build -o main .

EXPOSE 9090

CMD ["./main"]
