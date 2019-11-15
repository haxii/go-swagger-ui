FROM golang:1.12 as builder
MAINTAINER Zichao Li <zichao@haxii.com>

# Set Environment Variables
ENV CGO_ENABLED 0
ENV GOOS linux

# build go-swagger
WORKDIR /go/src/app
COPY . .

RUN GO111MODULE=on go mod download &&\
    GO111MODULE=on go build -ldflags "-X main.Build=450ce4a -X main.Version=v3.24.2" swagger.go &&\
    mkdir -p /swagger &&\
    mv swagger /go-swagger &&\
    cd .. &&\
    rm -rf /go/src/app/*


FROM alpine
COPY --from=builder /go-swagger /go-swagger

EXPOSE 8080
VOLUME /swagger
CMD [ "/go-swagger" ]
