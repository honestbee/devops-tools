# -- Go Builder Image --
FROM golang:1.9-alpine AS builder

ENV DEP_VERSION=0.4.1

RUN apk add --no-cache git curl
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 && chmod +x /usr/local/bin/dep
COPY . /go/src/rds-snapper
WORKDIR /go/src/rds-snapper

# https://github.com/golang/dep/blob/master/docs/FAQ.md#how-do-i-use-dep-with-docker
RUN set -ex \
  && dep ensure \
  && go build -v -o "/rds-snapper"

# -- rds-snapper Image --
FROM alpine:3.6
RUN set -ex \
  && apk add --no-cache bash ca-certificates git

COPY --from=builder /rds-snapper /bin/rds-snapper
ENTRYPOINT [ "/bin/rds-snapper" ]
