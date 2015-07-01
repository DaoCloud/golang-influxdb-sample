FROM google/golang
MAINTAINER Haipeng Wu "haipeng.wu@daocloud.io"

WORKDIR /gopath/app
ENV GOPATH /gopath/app
ADD . /gopath/app/src/golang-influxdb-sample

RUN go get -t golang-influxdb-sample
RUN go install golang-mongo-sample

EXPOSE 80
CMD ["/gopath/app/bin/golang-influxdb-sample"]
