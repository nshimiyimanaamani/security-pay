#build stage
FROM golang:1.12 AS builder

LABEL MAINTAINER="rugwirobaker@gmail.com"

WORKDIR $GOPATH/src/github.com/rugwirobaker/paypack-backend

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

COPY . .

ARG VERSION="unset"

RUN DATE="$(date -u +%Y-%m-%d-%H:%M:%S-%Z)" && GO111MODULE=on CGO_ENABLED=0 go build -mod vendor -ldflags "-X github.com/rugwirobaker/paypack-backend/pkg/build.version=$VERSION -X github.com/rugwirobaker/pkg/paypack-backend/build.buildDate=$DATE" -o /bin/paypack ./cmd/paypack


#package stage
FROM alpine

COPY --from=builder /bin/paypack /bin/paypack

RUN apk add --update tini

ENV GO_ENV=production

ENTRYPOINT [ "/sbin/tini", "--" ]

CMD ["paypack"]








