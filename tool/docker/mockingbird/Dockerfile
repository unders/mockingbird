FROM ubuntu:18.04
LABEL maintainer="anders@functionbox.io"


ENV GOPATH=/go
ENV PATH=/usr/local/go/bin:$GOPATH/bin:$PATH


##
## Install Git, Go, wget, unzip
##
RUN apt-get update && apt-get install -y \
   && apt-get install git -y \
   && apt-get install wget -y \
   && apt-get install unzip -y \
   && wget https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz \
   && tar -xvf go1.11.4.linux-amd64.tar.gz \
   && rm go1.11.4.linux-amd64.tar.gz \
   && mv go /usr/local \
   && go version \
   && apt-get clean \
   && rm -rf /tmp/* /var/tmp/* \
   && rm -rf /var/lib/apt/lists/* \
   && rm -f /etc/ssh/ssh_host_*


##
## Install Protobuf 3.6.1
##
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip \
    && unzip protoc-3.6.1-linux-x86_64.zip -d protoc3 \
    && rm protoc-3.6.1-linux-x86_64.zip \
    && mv protoc3/bin/* /usr/local/bin/ \
    && mv protoc3/include/* /usr/local/include/ \
    && protoc --version
    # Optional: change owner
    # sudo chwon [user] /usr/local/bin/protoc
    # sudo chwon -R [user] /usr/local/include/google


##
## https://magefile.org/blog/2018/09/mage-in-ci/
##
RUN go get -d github.com/magefile/mage \
     && cd $GOPATH/src/github.com/magefile/mage/ \
     && go run bootstrap.go \
     && mage --version

##
## set mage entrypoint
##
ENTRYPOINT ["mage"]

