FROM alpine as setup

RUN apk add -U --no-cache ca-certificates

FROM alpine

ADD bin/paypack /bin/

COPY --from=setup /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["paypack"]


