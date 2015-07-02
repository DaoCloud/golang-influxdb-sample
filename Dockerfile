FROM google/golang
MAINTAINER Haipeng Wu "haipeng.wu@daocloud.io"

WORKDIR /gopath/app
ENV GOPATH /gopath/app

RUN go get github.com/influxdb/influxdb/client

ADD . /gopath/app/src/golang-influxdb-sample
RUN go install golang-influxdb-sample

EXPOSE 80
CMD ["/gopath/app/bin/golang-influxdb-sample"]
