FROM centos:7

# Add Maintainer Info
LABEL maintainer="Luca Ruggieri <luc4.ruggieri@gmail.com>"

#useful for 'install.sh' colored output ^^
ENV TERM=xterm


# Set the Current Working Directory inside the container
WORKDIR /app

RUN yum install sudo -y

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy the source from the current directory to the Working Directory inside the container
COPY scripts ./scripts
RUN chmod +x scripts/*

RUN ./scripts/installEnvironment.sh

ENV GOROOT "/usr/local/go"
ENV PATH "$GOPATH/bin:$GOROOT/bin:$PATH"

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

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