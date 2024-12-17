# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

# Please note that this Dockerfile references files from outside directories,
# therefore, it need to be built with context set to the project's root.

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /build ./cmd/api

FROM alpine:latest

WORKDIR /app

# This is required for any bash scripts to run - alpine does not 
# come with bash by default, and all scripts set interpreter to bash in shebang.
RUN apk add --no-cache bash

COPY --from=builder /build/.env /build/api ./
COPY --from=builder /build/tools ./tools

CMD ["./api"]
