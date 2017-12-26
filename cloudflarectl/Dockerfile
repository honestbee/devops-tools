# -- Go Builder Image --
FROM golang:1.9-alpine AS builder

RUN apk add --no-cache git
COPY . /go/src/cloudflarectl
WORKDIR /go/src/cloudflarectl

RUN set -ex \
  && go get -v \
  && go build -v -o "/cloudflarectl"

# -- cloudflarectl Image --
FROM alpine:3.6
RUN set -ex \
  && apk add --no-cache bash ca-certificates git

COPY --from=builder /cloudflarectl /bin/cloudflarectl
ENTRYPOINT [ "/bin/cloudflarectl" ]
