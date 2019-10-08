#build stage
FROM golang:1.12.9 AS builder

LABEL MAINTAINER="rugwirobaker@gmail.com"

WORKDIR $GOPATH/src/github.com/rugwirobaker/paypack-backend

ENV GO111MODULE=on

#ENV GOPROXY="https://proxy.golang.org"

WORKDIR /src
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

COPY go.mod .

COPY go.sum .

RUN go mod tidy

COPY . .

ARG VERSION="unset"

RUN DATE="$(date -u +%Y-%m-%d-%H:%M:%S-%Z)" && CGO_ENABLED=0 go build -ldflags "-X github.com/rugwirobaker/paypack-backend/build.version=$VERSION -X github.com/rugwirobaker/paypack-backend/build.buildDate=$DATE" -o /bin/paypack ./cmd/.


#package stage
FROM scratch

EXPOSE 8080

CMD ["paypack"]

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /bin/paypack /bin/paypack






