ARG GOLANG_VERSION=1.12
FROM golang:${GOLANG_VERSION} AS builder

LABEL MAINTAINER="rugwirobaker@gmail.com"

WORKDIR $GOPATH/src/github.com/rugwirobaker/paypack-backend

COPY go.mod go.sum ./

RUN GO111MODULE=on go mod download

COPY . .

ARG VERSION="unset"

RUN GO111MODULE=on CGO_ENABLED=0 go build -o /bin/auditor ./cmd/auditor

#packaging stage
FROM alpine

COPY --from=builder /bin/auditor /bin/auditor

RUN apk add --update ca-certificates tini

ENV GO_ENV=production

RUN adduser -D auditor
USER auditor

ENTRYPOINT [ "/sbin/tini", "-s", "--" ]

CMD ["auditor"]