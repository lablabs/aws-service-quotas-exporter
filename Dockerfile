FROM golang:1.22-alpine3.19 as builder

ARG GOOS=linux
ARG GOARCH=amd64

RUN apk --update add ca-certificates

WORKDIR $GOPATH/src/github.com/lablabs/aws-service-quotas-exporter
COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN go mod vendor
RUN go mod verify

RUN cd cmd/exporter && \
    GOOS=$GOOS GOARCH=$GOARCH \
    CGO_ENABLED=0 \
    go build -o /aws-service-quotas-exporter .

FROM alpine:3.19.1 as security

RUN apk add -U --no-cache \
    tzdata \
    ca-certificates

RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

FROM scratch

COPY --from=security /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=security /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=security /etc/passwd /etc/passwd
COPY --from=security /etc/group /etc/group

COPY --from=builder /aws-service-quotas-exporter .

ENTRYPOINT ["/aws-service-quotas-exporter"]

