FROM centos:7

# Add Maintainer Info
LABEL maintainer="Luca Ruggieri <luc4.ruggieri@gmail.com>"

WORKDIR /web

# gcc and g++ are necessary for gosseract and underling libraries
RUN yum install -y wget
RUN yum install -y git

# install go
RUN wget https://dl.google.com/go/go1.12.7.linux-amd64.tar.gz && \
    tar -xzf go1.12.7.linux-amd64.tar.gz && \
    mv go /usr/local

# put here every other package to install
# RUN yum install something_else

ENV GOROOT='/usr/local/go'
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH


COPY pkg/util pkg/util
COPY go.mod go.sum ./
RUN go mod download
RUN go mod vendor
COPY website website


# Build visualResourceServer app
RUN go build -o gonemaweb website/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

#ENV TEMPLATES_DIR=website/templates
#ENV STATIC_ASSET_DIR=website/static
#ENV GONEMAPI_HOST=gonemapi.ruggieri.tech
#ENV GONEMAPI_PORT=8080

# Command to run visualResourceServer
CMD ["./gonemaweb"]