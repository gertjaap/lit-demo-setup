
FROM golang:alpine as build

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc musl-dev
ENV GOROOT=/usr/local/go
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN go get github.com/docker/docker/api/types
RUN go get github.com/docker/docker/api/types/container
RUN go get github.com/docker/docker/client
RUN rm -rf /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api
COPY . /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api
WORKDIR /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api
RUN go get -v ./...
RUN go build

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /usr/local/go/src/github.com/gertjaap/lit-demo-setup/admin-api/admin-api /app/bin/admin-api
EXPOSE 8000
RUN mkdir /app/bin/data
WORKDIR /app/bin
COPY static /app/bin/static 

ENV DOCKER_API_VERSION=1.35

CMD ["./admin-api"]