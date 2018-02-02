# -- Go Builder Image --
FROM golang:1.9-alpine AS builder

ENV DEP_VERSION=0.4.1
ENV SERVICE_NAME=rds-snapper

RUN apk add --no-cache git curl
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 && chmod +x /usr/local/bin/dep
COPY . /go/src/${SERVICE_NAME}
WORKDIR /go/src/${SERVICE_NAME}

# https://github.com/golang/dep/blob/master/docs/FAQ.md#how-do-i-use-dep-with-docker
RUN set -ex \
  && dep ensure \
  && go build -v -o "/${SERVICE_NAME}"

# -- ${SERVICE_NAME} Image --
FROM alpine:3.6
RUN set -ex \
  && apk add --no-cache bash ca-certificates git

COPY --from=builder /${SERVICE_NAME} /bin/${SERVICE_NAME}
ENTRYPOINT [ "/bin/${SERVICE_NAME}", "clear" ]
