#build stage
ARG GOLANG_VERSION=1.12
FROM golang:${GOLANG_VERSION} AS builder

LABEL MAINTAINER="rugwirobaker@gmail.com"

WORKDIR $GOPATH/src/github.com/rugwirobaker/paypack-backend

COPY . .

ARG VERSION="unset"

RUN DATE="$(date -u +%Y-%m-%d-%H:%M:%S-%Z)" && GO111MODULE=on CGO_ENABLED=0 go build -mod vendor -ldflags "-X github.com/rugwirobaker/paypack-backend/pkg/build.version=$VERSION -X github.com/rugwirobaker/pkg/paypack-backend/build.buildDate=$DATE" -o /bin/paypack ./cmd/paypack


#package stage
FROM alpine

COPY --from=builder /bin/paypack /bin/paypack

RUN apk add --update ca-certificates tini

ENV GO_ENV=production

RUN adduser -D paypack
USER paypack

ENTRYPOINT [ "/sbin/tini", "-s", "--" ]

CMD ["paypack"]








