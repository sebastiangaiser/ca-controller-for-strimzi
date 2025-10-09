FROM golang:1.25.2-alpine as builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o ca-controller-for-strimzi


# final image
FROM alpine:3.22.2

RUN apk add --no-cache ca-certificates \
    && update-ca-certificates

COPY --from=builder /build/ca-controller-for-strimzi /usr/local/bin/

USER 65534:65534

ENTRYPOINT [ "ca-controller-for-strimzi" ]
