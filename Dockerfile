#build stage
FROM golang:1.12.9 AS build-stage 

LABEL MAINTAINER="rugwirobaker@gmail.com"
LABEL service="paypack-backend"
LABEL daemon="/bin/backend"
ENV GO111MODULE=on

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
CMD ["/bin/backend"]
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /user/group /user/passwd /etc/
COPY --from=build-stage /src/bin/app /bin/backend






