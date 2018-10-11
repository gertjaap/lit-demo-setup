
FROM golang:alpine as build

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc musl-dev
ENV GOROOT=/usr/local/go
COPY vendor/vendor.json /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api/vendor/vendor.json
RUN go get -v github.com/kardianos/govendor
WORKDIR /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api
RUN govendor fetch +external +missing
COPY . /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api
WORKDIR /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api
RUN govendor fetch +missing

## Temporary: Fetch memberweekdemo branch from gertjaap's GIT

WORKDIR /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api/vendor/github.com/mit-dci
RUN rm -rf lit
RUN git clone --single-branch -b memberweekdemo https://github.com/gertjaap/lit
WORKDIR /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api

RUN go get -v ./...
RUN go build

FROM alpine
RUN apk add --no-cache ca-certificates graphviz font-bitstream-type1 ghostscript-fonts
COPY --from=build /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api/admin-api /app/bin/admin-api
EXPOSE 8000
RUN mkdir /app/bin/data
WORKDIR /app/bin
COPY static /app/bin/static 

ENV DOCKER_API_VERSION=1.35

CMD ["./admin-api"]
