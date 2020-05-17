# FROM ubuntu:18.04

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

ENV dir /home/dev
ENV app cmd

WORKDIR ${dir}

# Copy the local package files to the container's workspace.
ADD . .

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install ./${app}

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/${app}

# Document that the service listens on port 8080.
EXPOSE 1323
