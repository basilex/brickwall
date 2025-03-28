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
RUN make app-build
RUN strip bsp

FROM debian:bullseye

RUN apt-get update
RUN apt-get install -y bash
RUN apt-get install -y curl
RUN apt-get install -y wget

WORKDIR /app

COPY --from=builder /build/bsp .
COPY --from=builder /build/resource .
COPY --from=builder /build/.env .

ENV PATH="/app:$PATH"

EXPOSE 8081
