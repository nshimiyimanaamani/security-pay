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
FROM scratch

ENV GO_ENV=production

EXPOSE 8080

CMD /paypack

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /bin/paypack /paypack






