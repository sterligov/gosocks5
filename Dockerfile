FROM golang:1.14 as builder

ARG WORKPATH=/app
WORKDIR $WORKPATH

COPY go.mod $WORKPATH
RUN go mod download

COPY . $WORKPATH

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o socks5 cmd/main.go

FROM scratch
COPY --from=builder /app/logs /logs
COPY --from=builder /app/socks5 /socks5
ENTRYPOINT ["./socks5"]