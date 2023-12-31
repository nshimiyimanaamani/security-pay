FROM node:14.20.0-alpine as static
WORKDIR /app
RUN apk add --no-cache python3 py3-pip
RUN apk add curl make
COPY . .
RUN touch .env && make ui


FROM  golang:1.16.0-alpine as builder
ENV GO111MODULE=on
WORKDIR /paypack
RUN apk add git curl
RUN apk add --no-cache tzdata
COPY go.mod go.sum  ./
RUN GOPROXY="https://proxy.golang.org" go mod download
COPY . .

COPY --from=static  /app/web/dist ./web/dist/

ARG VERSION="unset"
RUN DATE="$(date -u +%Y-%m-%d-%H:%M:%S-%Z)" \ 
    VERSION="$(git rev-parse --short HEAD)" \
    LDFLAGS="-X github.com/nshimiyimanaamani/paypack-backend/pkg/build.version=$VERSION -X github.com/quarksgroup/paypack-engine/pkg/build.buildDate=$DATE" \
    CGO_ENABLED=0 go build -ldflags=$LDFLAGS -o /bin/paypack ./cmd/paypack

ARG VERSION="unset"
RUN DATE="$(date -u +%Y-%m-%d-%H:%M:%S-%Z)" \ 
    VERSION="$(git rev-parse --short HEAD)" \
    LDFLAGS="-X github.com/nshimiyimanaamani/paypack-backend/pkg/build.version=$VERSION -X github.com/quarksgroup/paypack-engine/pkg/build.buildDate=$DATE" \
    CGO_ENABLED=0 go build -ldflags=$LDFLAGS -o /bin/worker ./cmd/worker

FROM scratch
WORKDIR /
EXPOSE 8080
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/* /bin/
ENTRYPOINT ["/bin/paypack"]
