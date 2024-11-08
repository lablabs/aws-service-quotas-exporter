ARG ALPINE_VERSION=3.19
FROM python:3.10-alpine${ALPINE_VERSION} as builder_aws_cli

ARG AWS_CLI_VERSION=2.15.19
RUN apk add --no-cache git unzip groff build-base libffi-dev cmake
RUN git clone --single-branch --depth 1 -b ${AWS_CLI_VERSION} https://github.com/aws/aws-cli.git

WORKDIR aws-cli
RUN ./configure --with-install-type=portable-exe --with-download-deps
RUN make
RUN make install

# reduce image size: remove autocomplete and examples
RUN rm -rf \
    /usr/local/lib/aws-cli/aws_completer \
    /usr/local/lib/aws-cli/awscli/data/ac.index \
    /usr/local/lib/aws-cli/awscli/examples
RUN find /usr/local/lib/aws-cli/awscli/data -name completions-1*.json -delete
RUN find /usr/local/lib/aws-cli/awscli/botocore/data -name examples-1.json -delete
RUN (cd /usr/local/lib/aws-cli; for a in *.so*; do test -f /lib/$a && rm $a; done)

FROM golang:1.22-alpine${ALPINE_VERSION} as builder_golang

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

FROM alpine:${ALPINE_VERSION}

RUN apk --no-cache add jq bash

COPY --from=builder_aws_cli /usr/local/lib/aws-cli/ /usr/local/lib/aws-cli/
RUN ln -s /usr/local/lib/aws-cli/aws /usr/local/bin/aws

COPY --from=builder_golang /aws-service-quotas-exporter .

RUN addgroup -S -g 1001 exporter && \
    adduser -S -u 1001 -G exporter exporter

USER 1001:1001

ENTRYPOINT ["/aws-service-quotas-exporter"]
