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

RUN ./scripts/installEnvironment.sh

# Copy everything else in the Working Directory
COPY . .

RUN chmod +x scripts/installDependencies.sh
RUN ./scripts/installDependencies.sh

# Build visualResourceServer app
RUN go build -o visualResourceServer cmd/visualResourceServer/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run visualResourceServer
CMD ["./visualResourceServer"]