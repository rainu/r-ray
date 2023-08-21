FROM golang:alpine AS builder
RUN apk add --no-cache git gcc musl-dev
ADD . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go install -installsuffix cgo ./cmd/lambda

FROM public.ecr.aws/lambda/provided:al2

COPY --from=builder /go/bin/* /

ENTRYPOINT ["/lambda"]

