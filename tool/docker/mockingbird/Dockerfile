FROM ubuntu:18.04
MAINTAINER unders <anders@functionbox.io>

##
## Install Git, Go, wget
##
RUN apt-get update -y \
   && apt-get install git -y \
   && apt-get install wget -y \
   && wget https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz \
   && tar -xvf go1.11.4.linux-amd64.tar.gz \
   && rm go1.11.4.linux-amd64.tar.gz \
   && mv go /usr/local


ENV GOPATH=/go
ENV PATH=/usr/local/go/bin:$GOPATH/bin:$PATH


##
## https://magefile.org/blog/2018/09/mage-in-ci/
##
RUN go get -d github.com/magefile/mage \
     && cd $GOPATH/src/github.com/magefile/mage/ \
     && go run bootstrap.go \
     && go version \
     && mage --version

##
## set mage entrypoint
##
ENTRYPOINT ["mage"]

