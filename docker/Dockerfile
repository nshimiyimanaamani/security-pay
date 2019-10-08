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
RUN go mod download
COPY . .
RUN make build


#package stage
FROM scratch

EXPOSE 8080
CMD ["/bin/paypack"]
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /bin/paypack /bin/paypack






