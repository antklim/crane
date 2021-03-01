# 1. Build
FROM golang:alpine as builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o . ./...

# 2. Run
FROM gcr.io/distroless/base

COPY --from=builder /build/static /

# Command to run
ENTRYPOINT ["/static"]