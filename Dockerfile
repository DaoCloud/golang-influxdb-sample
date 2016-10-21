FROM google/golang
MAINTAINER Haipeng Wu "haipeng.wu@daocloud.io"

WORKDIR /gopath/app
ENV GOPATH /gopath/app

RUN go get github.com/influxdata/influxdb/client/v2

ADD . /gopath/app/src/golang-influxdb-sample
RUN go install golang-influxdb-sample

EXPOSE 80
CMD ["/gopath/app/bin/golang-influxdb-sample"]
