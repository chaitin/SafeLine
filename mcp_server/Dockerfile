FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o mcp-server .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/mcp-server .
COPY --from=builder /app/config.yaml .

ENV TZ=Asia/Shanghai

EXPOSE 5678

CMD ["./mcp-server"] 