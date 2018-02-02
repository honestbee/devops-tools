# -- Go Builder Image --
FROM golang:1.9-alpine AS builder

ENV DEP_VERSION=0.4.1

RUN apk add --no-cache git curl
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 && chmod +x /usr/local/bin/dep
COPY . /go/src/aws-snapshot-cleanup
WORKDIR /go/src/aws-snapshot-cleanup

# https://github.com/golang/dep/blob/master/docs/FAQ.md#how-do-i-use-dep-with-docker
RUN set -ex \
  && dep ensure \
  && go build -v -o "/aws-snapshot-cleanup"

# -- aws-snapshot-cleanup Image --
FROM alpine:3.6
RUN set -ex \
  && apk add --no-cache bash ca-certificates git

COPY --from=builder /aws-snapshot-cleanup /bin/aws-snapshot-cleanup
ENTRYPOINT [ "/bin/aws-snapshot-cleanup", "clear" ]
