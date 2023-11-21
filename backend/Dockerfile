FROM golang:1.21 as build

WORKDIR /app

COPY .  .

RUN go mod tidy && \
    CGO_ENABLED=0 go build -o server .

FROM alpine:latest

WORKDIR /app

COPY --from=0 /app/server .

CMD ["./server"]
