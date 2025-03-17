FROM golang:1.24 AS builder

ENV TERM=linux
ENV CGO_ENABLED=0 GOOS=linux GOARCH=arm64

RUN apt-get update

RUN apt-get install git -y
RUN apt-get install build-essential -y

WORKDIR /build

COPY . .
COPY .git .git

RUN go mod tidy
RUN make build
RUN strip bsp

FROM debian:12.9-slim

WORKDIR /app

COPY --from=builder /build/bsp .

ENV PATH="/app:$PATH"

EXPOSE 8081
