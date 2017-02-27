FROM ubuntu:16.04

MAINTAINER Julien CORIOLAND

ENTRYPOINT ["/go/src/govorlonjs/main"]

EXPOSE 82

RUN apt-get update -y 
RUN apt-get install -y git && apt-get install -y docker.io && apt-get install -y wget
RUN wget https://storage.googleapis.com/golang/go1.7.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.7.5.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin 
ENV GOPATH=/go

COPY *.go /go/src/govorlonjs/

WORKDIR /go/src/govorlonjs
RUN go get
RUN go build main.go utils.go dockercli.go