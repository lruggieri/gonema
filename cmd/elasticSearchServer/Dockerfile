FROM centos:7

LABEL maintainer="Luca Ruggieri <luc4.ruggieri@gmail.com>"

WORKDIR /elasticSearchAPI

RUN yum install -y wget
RUN yum install -y curl
RUN yum install -y git

# install go
RUN wget https://dl.google.com/go/go1.12.7.linux-amd64.tar.gz && \
    tar -xzf go1.12.7.linux-amd64.tar.gz && \
    mv go /usr/local

ENV GOROOT='/usr/local/go'
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

COPY go.mod go.sum ./

RUN go mod download
RUN go mod vendor

COPY cmd/elasticSearchServer cmd/elasticSearchServer
COPY pkg pkg

RUN go build -o elasticSearchServer cmd/elasticSearchServer/main.go

EXPOSE 8080

# Command to run visualResourceServer
CMD ["./elasticSearchServer"]