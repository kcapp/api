# Create our build image
FROM golang:1.23-alpine AS BUILD_IMAGE

# Add git, required to install dependencies
RUN apk update && apk add --no-cache git gcc

# Install goose to run database migrations
WORKDIR $GOPATH/src/github.com/pressly/goose
RUN git clone https://github.com/pressly/goose .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -tags='no_clickhouse no_libsql no_mssql no_vertica no_postgres no_sqlite3 no_ydb no_duckdb' -o $GOPATH/bin/goose -a -ldflags '-extldflags "-static"' ./cmd/goose

# Add wait-for-it
COPY wait-for-it.sh /usr/local/scripts/wait-for-it.sh
RUN chmod +x /usr/local/scripts/wait-for-it.sh

WORKDIR $GOPATH/src/github.com/kcapp/api

# Bundle app source
COPY . .

# Install dependencies and build executable
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o $GOPATH/bin/api -a -ldflags '-extldflags "-static"' .

# Separate stage for cloning migrations (non-cacheable)
FROM alpine AS MIGRATIONS
RUN apk add --no-cache git
RUN git clone https://github.com/kcapp/database /usr/local/kcapp/database
RUN cp /usr/local/kcapp/database/run_migrations.sh /usr/local/scripts/run_migrations.sh
RUN chmod +x /usr/local/scripts/run_migrations.sh

# Create our actual image
FROM alpine
RUN apk add --no-cache bash

# Add configuration file
COPY config/config.docker.yaml config/config.yaml

# Add binaries and scripts
COPY --from=BUILD_IMAGE /usr/local/scripts/* ./
COPY --from=BUILD_IMAGE /go/bin/goose /go/bin/goose
COPY --from=BUILD_IMAGE /go/bin/api /go/bin/api

# Force a fresh migration copy by using `MIGRATIONS` stage, avoiding cache
ARG FORCE_MIGRATIONS_UPDATE
RUN --mount=type=cache,target=/var/cache git clone https://github.com/kcapp/database /usr/local/kcapp/database
COPY --from=MIGRATIONS /usr/local/kcapp/database/migrations /usr/local/kcapp/database/migrations
COPY --from=MIGRATIONS /usr/local/scripts/run_migrations.sh /usr/local/scripts/run_migrations.sh

# Add go binaries to path
ENV PATH="/go/bin:${PATH}"

CMD [ "/go/bin/api" ]
EXPOSE 8001