# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/kcapp/api

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/kcapp/api
RUN go install github.com/kcapp/api

# Add wait-for-it
COPY wait-for-it.sh wait-for-it.sh
RUN chmod +x wait-for-it.sh

COPY config/config.yaml config/config.yaml

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/api

# Document that the service listens on port 8001.
EXPOSE 8001
