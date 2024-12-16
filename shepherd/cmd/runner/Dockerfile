# syntax=docker/dockerfile:1

FROM golang:1.17-alpine AS builder

# Please note that this Dockerfile references files from outside directories,
# therefore, it need to be built with context set to the project's root.

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /build ./cmd/runner

FROM alpine:latest

WORKDIR /app

# This is required for any bash scripts to run - alpine does not 
# come with bash by default, and all scripts set interpreter to bash in shebang.
RUN apk add --no-cache bash \
  # Python3, dateutil and requests are required for vast CLI to work, and Linux alpine does not come with it.
  && apk add --no-cache python3 py3-pip \
  # We need to use "--break-system-packages" flag to allow pip to install the lib globally. \
  # @TODO: Consider using virtual env.
  && pip install python-dateutil requests --break-system-packages \
  # We also need to install vast CLI executable. Note that we are adding it directly to /usr/bin, so
  # it will be available in PATH.
  && wget https://raw.githubusercontent.com/vast-ai/vast-python/master/vast.py -O /usr/bin/vast; chmod +x /usr/bin/vast;

COPY --from=builder /build/.env /build/runner ./

CMD ["./runner"]
