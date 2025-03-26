# Create our build image
FROM golang:1.23-alpine AS build_image

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
ARG VERSION=dev
ARG GIT_COMMIT=unknown
RUN CGO_ENABLED=0 go build -o $GOPATH/bin/api \
  -ldflags="-X 'github.com/kcapp/api/models.Version=${VERSION}' -X 'github.com/kcapp/api/models.GitCommit=${GIT_COMMIT}' -extldflags '-static'"

# Separate stage for cloning migrations (non-cacheable)
FROM alpine AS migrations
RUN apk add --no-cache git

RUN git clone https://github.com/kcapp/database /usr/local/kcapp/database
RUN mkdir -p /usr/local/scripts

# Copy and set permission
RUN cp /usr/local/kcapp/database/run_migrations.sh /usr/local/scripts/run_migrations.sh
RUN chmod +x /usr/local/scripts/run_migrations.sh

# Create our actual image
FROM alpine
RUN apk add --no-cache bash git

# Add configuration file
COPY config/config.docker.yaml config/config.yaml

# Add binaries and scripts
COPY --from=build_image /usr/local/scripts/* ./
COPY --from=build_image /go/bin/goose /go/bin/goose
COPY --from=build_image /go/bin/api /go/bin/api

# Force a fresh migration copy by using `MIGRATIONS` stage, avoiding cache
ARG force_migrations_update
RUN --mount=type=cache,target=/var/cache git clone https://github.com/kcapp/database /usr/local/kcapp/database
COPY --from=migrations /usr/local/kcapp/database/migrations /usr/local/kcapp/database/migrations
COPY --from=migrations /usr/local/scripts/run_migrations.sh ./run_migrations.sh
COPY --from=migrations /usr/local/scripts/run_migrations.sh /usr/local/scripts/run_migrations.sh

# Add go binaries to path
ENV PATH="/go/bin:${PATH}"

CMD [ "/go/bin/api" ]
EXPOSE 8001