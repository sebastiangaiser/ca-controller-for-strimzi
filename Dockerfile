FROM golang:1.22.6-alpine as builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o ca-controller-for-strimzi


# final image
FROM alpine:3.20.2
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/ca-controller-for-strimzi .

ENTRYPOINT [ "./ca-controller-for-strimzi" ]