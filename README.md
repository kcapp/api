![kcapp logo](https://raw.githubusercontent.com/kcapp/frontend/master/public/images/logo.png)

[![Go](https://github.com/kcapp/api/actions/workflows/go.yml/badge.svg)](https://github.com/kcapp/api/actions/workflows/go.yml)

# api
Backend API for [kcapp-frontend](https://github.com/kcapp/frontend)

## Install
* Execute `go get github.com/kcapp/api`
* Run `go build` inside `$GOPATH/src/github.com/kcapp/api`
* Now you have a built `executable` which can be run to start the `API`

## Configuration
Configuration is done through `yaml`. All options can be set in `config/config.yaml` or a custom config file by specifying it as a argument to
```bash
./api custom_config.yaml
```

### Database
Information about the database, and its configuration can be found in [kcapp-database](https://github.com/kcapp/database)
