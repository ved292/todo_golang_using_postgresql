FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . ./

RUN go build -o main .

CMD ["./main"]
