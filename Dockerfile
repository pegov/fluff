FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

EXPOSE 8080

ENV GIN_MODE="release"

CMD ["/main"]
