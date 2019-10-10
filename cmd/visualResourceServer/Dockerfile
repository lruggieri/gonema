FROM centos:7

# Add Maintainer Info
LABEL maintainer="Luca Ruggieri <luc4.ruggieri@gmail.com>"

WORKDIR /visualResourceServer

# gcc and g++ are necessary for gosseract and underling libraries
RUN yum install -y sudo wget
RUN yum install -y curl
RUN yum install -y gcc
RUN yum install -y gcc-c++
RUN yum install -y epel-release
RUN yum install -y tesseract-devel
RUN yum install -y leptonica-devel
RUN yum install -y git

# Google chrome is needed as the main tool used by chromedp to scrape the web
RUN curl https://intoli.com/install-google-chrome.sh | bash

# install go
RUN wget https://dl.google.com/go/go1.12.7.linux-amd64.tar.gz && \
    tar -xzf go1.12.7.linux-amd64.tar.gz && \
    mv go /usr/local

# put here every other package to install
# RUN yum install something_else

ENV GOROOT='/usr/local/go'
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH


COPY cmd/visualResourceServer cmd/visualResourceServer
COPY configs configs
COPY pkg pkg
COPY go.mod go.sum ./
COPY .env ./

RUN go mod download
RUN go mod vendor

# Build visualResourceServer app
RUN go build -o visualResourceServer cmd/visualResourceServer/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run visualResourceServer
CMD ["./visualResourceServer"]