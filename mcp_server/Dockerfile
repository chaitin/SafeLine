FROM golang:1.25.5-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/mcp-server .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && install -d -m 700 /var/lib/safeline-mcp \
    && install -d -m 755 /etc/safeline-mcp

COPY --from=builder /out/mcp-server /app/mcp-server
COPY config.yaml /etc/safeline-mcp/config.yaml

ENV TZ=Asia/Shanghai

EXPOSE 5678

ENTRYPOINT ["/app/mcp-server"]
CMD ["--config", "/etc/safeline-mcp/config.yaml"]
