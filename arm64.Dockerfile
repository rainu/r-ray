FROM golang:alpine AS builder
RUN apk add --no-cache git gcc musl-dev
ADD . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go install -installsuffix cgo ./cmd/app

FROM alpine:3.20
USER 999:999
COPY --from=builder /go/bin/* /

ENTRYPOINT ["/app"]