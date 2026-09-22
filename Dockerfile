FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY http-server.go .
RUN go build -o http-server http-server.go


FROM ubuntu:24.04

RUN apt-get update && apt-get install -y --no-install-recommends procps && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/log/script

COPY script.sh /usr/local/bin/script.sh
RUN chmod +x /usr/local/bin/script.sh

COPY --from=builder /app/http-server /usr/local/bin/http-server
RUN chmod +x /usr/local/bin/http-server

EXPOSE 8080

CMD ["/usr/local/bin/http-server"]