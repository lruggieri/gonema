# Start from the latest golang base image
FROM centos:7

# Add Maintainer Info
LABEL maintainer="Luca Ruggieri <luc4.ruggieri@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

RUN yum install sudo -y
#RUN useradd -m docker && echo "docker:docker" | chpasswd && usermod -aG wheel docker
#RUN mkdir -p /etc/sudoers.d
#RUN echo "docker ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/docker && \
#    chmod 0440 /etc/sudoers.d/docker
# USER docker

# Install golang
RUN yum install wget -y
RUN yum install git -y
RUN wget https://dl.google.com/go/go1.12.7.linux-amd64.tar.gz
RUN tar -xzf go1.12.7.linux-amd64.tar.gz
RUN mv go /usr/local
#useful for 'install.sh' colored output ^^
ENV TERM=xterm
ENV GOROOT=/usr/local/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY scripts ./scripts
RUN chmod +x scripts/*

RUN ./scripts/install_environment.sh

# Copy everything else in the Working Directory
COPY . .

RUN chmod +x scripts/install_dependencies.sh
RUN ./scripts/install_dependencies.sh

# Build the Go app
RUN go build -o gonema cmd/visual_resource_server/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./gonema"]