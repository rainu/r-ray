FROM golang:alpine AS builder
RUN apk add --no-cache git gcc musl-dev
ADD . /src
WORKDIR /src
RUN go test ./... && CGO_ENABLED=0 GOOS=linux go install -installsuffix cgo ./cmd/app

FROM alpine:3.10
USER 999:999
COPY --from=builder /go/bin/* /

ENTRYPOINT ["/app"]