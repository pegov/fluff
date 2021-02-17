FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go get ./...
RUN GOOS=linux go build -o ./bin/web ./cmd/server/server.go

FROM alpine:3.13
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
ENTRYPOINT ["/go/bin/web"]