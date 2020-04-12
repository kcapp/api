# Create our build image
FROM golang:alpine AS BUILD_IMAGE

# Add git, required to install dependencies
RUN apk update && apk add --no-cache git gcc

# Install goose to run database migrations
WORKDIR $GOPATH/src/github.com/pressly/goose
RUN git clone https://github.com/pressly/goose .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -tags='no_postgres no_mysql no_sqlite3' -i -o $GOPATH/bin/goose -a -ldflags '-extldflags "-static"' ./cmd/goose

# Add script to run migrations
RUN mkdir -p /usr/local/scripts
RUN git clone https://github.com/kcapp/database /usr/local/kcapp/database
RUN cp /usr/local/kcapp/database/run_migrations.sh /usr/local/scripts/run_migrations.sh
RUN chmod +x /usr/local/scripts//run_migrations.sh

# Add wait-for-it
COPY wait-for-it.sh /usr/local/scripts/wait-for-it.sh
RUN chmod +x /usr/local/scripts/wait-for-it.sh

WORKDIR $GOPATH/src/github.com/kcapp/api

# Bundle app source
COPY . .

# Install dependencies and build executable
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o $GOPATH/bin/api -a -ldflags '-extldflags "-static"' .

# Create our actual image
FROM alpine

RUN apk add --no-cache bash

# Add configuration file
COPY config/config.docker.yaml config/config.yaml

# Add binaries and scripts
COPY --from=BUILD_IMAGE /usr/local/scripts/* ./
COPY --from=BUILD_IMAGE /go/bin/goose /go/bin/goose
COPY --from=BUILD_IMAGE /usr/local/kcapp/database/migrations /usr/local/kcapp/database/migrations
COPY --from=BUILD_IMAGE /go/bin/api /go/bin/api

# Add go binaries to path
ENV PATH="/go/bin:${PATH}"

CMD [ "/go/bin/api" ]
EXPOSE 8001