FROM golang:alpine as builder
MAINTAINER Zichao Li <zichao@haxii.com>

# build go-swagger
WORKDIR /go/src/app
COPY . .

RUN mkdir -p /go/src/app/vendor/github.com/haxii/go-swagger-ui/static &&\
    mv static/static.go /go/src/app/vendor/github.com/haxii/go-swagger-ui/static &&\
    go build -ldflags "-X main.Build=9bcba46 -X main.Version=v3.23.11" swagger.go &&\
    mkdir -p /swagger &&\
    mv swagger /go-swagger &&\
    cd .. &&\
    rm -rf /go/src/app/*


FROM alpine
COPY --from=builder /go-swagger /go-swagger

EXPOSE 8080
VOLUME /swagger
CMD [ "/go-swagger" ]
