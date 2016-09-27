# S3 Dockerfile
#
# VERSION 0.1

FROM daocloud.io/library/golang:1.7

MAINTAINER Neo "neo1989@outlook.com"

RUN mkdir /go/src/s3/

ADD . /go/src/s3/

RUN mkdir -p /go/src/s3/uploads

RUN curl https://glide.sh/get | sh

WORKDIR /go/src/s3/

RUN glide install

EXPOSE 8888 
EXPOSE 9999 

VOLUME /go/src/s3/uploads

CMD go run s3.go







