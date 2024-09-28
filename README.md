# Kiddy Service
## Intro
Golang home task from this source:

https://drive.google.com/file/d/1rSni7pBvn80k36sHhN80FgNZTMmADZvS/view
## How to run
```shell
make up
```
Script will run application, storage and provider container.

See Makefile for details.

## How to use
1. via HTTP

Use CLI HTTP client(like curl or wget)
```shell
curl http://localhost:48001/ready
```
Port 48001 is used by default config.

See `api/openapi.yml` for description.

2. via GRPC

Use client, instructions inside.
```
./tools/grpc_client
```

## About
### Description
Application (app for short) is working in Docker containers:
1. processor - the golang application itself
2. storage - Postgres DB
3. lines-provider - container from requirements


### Configuration
App configuration options can be found in default `config/config.json` file.

App can be configured with:
- config file
- environment variables

#### Defaults

App is configured with default values.

You can find default values in `config/config.json` file.

#### How to use different config file
Create new `config_x.yml` and provide env `APP_ENV=x` to application.

E.g, we want to use `config_prod.json`:
```shell
APP_ENV=prod go run github.com/supressionstop/xenking_test_1/cmd/processor

# or

export APP_ENV=prod
go run github.com/supressionstop/xenking_test_1/cmd/processor
```

#### How to use environment variables
Use screaming snake case as env names. E.g.:
```text
APP_NAME=foo
DB_URL=sqlite://...
```
