# Gitea Actions Prometheus Exporter

A Prometheus exporter for Gitea Actions that collects and exposes metrics about action runs. This exporter connects to Gitea's Postgres database and provides metrics that can be scraped by Prometheus.

The primary goal of the project is to provide a real time signal on CI failures when using Gitea Actions.

## Features

- Exposes Gitea Actions metrics in Prometheus format
- Configurable update interval
- Docker support
- PostgreSQL database integration
- REST API endpoint for action runs data

## Prerequisites

- Gitea instance >= v1.23
- Gitea must be using a PostgreSQL database, MySQL and Sqlite are not supported
- Docker

## Getting Started

A dedicated demo and testing project repository is available: [gitea_actions_prometheus_exporter-project](https://github.com/james9001/gitea_actions_prometheus_exporter-project). This repository contains a Docker Compose stack and configuration which will get you up and running quickly.

## Metrics

The exporter exposes the following metrics at the `/metrics` endpoint:

- `action_runs_failure_total`: Total number of all action runs with status "failure"
- `action_runs_not_success_total`: Total number of stopped action runs with status other than "success"
- `action_runs_failure_or_cancelled_total`: Total number of all action runs with status "failure" or "cancelled"

Additionally, for each metric, the following labels are available:
- `repository_name`: The name of the repository that the action run was executed for (e.g. gitea_actions_prometheus_exporter)
- `workflow_id`: The name of the workflow YAML file that the action run was executed for (e.g. build.yaml)

## Configuration

The application is configured using environment variables, which are configured as part of the container. For development reasons, using a `.env` file is also supported. The available configuration options are:

```env
# PostgreSQL Configuration
DB_HOST=localhost          # Database host
DB_PORT=5432               # Database port
DB_USER=postgres           # Database user
DB_PASSWORD=your_password  # Database password
DB_NAME=your_database      # Database name

# Application Configuration
SERVER_PORT=9100           # Port on which the exporter will listen
UPDATE_INTERVAL=1          # Interval in seconds between metric updates
```

## Building and Developing

The easiest way to build the container image locally is to run the `./build.sh` script.

### Local Development

- `pre-commit install` to install hooks
- Install the Pip3 package globally if you have not already `xmlformatter` i.e. `pip3 install xmlformatter --break-system-packages`
- Install Go v1.24.1. https://go.dev/doc/install
- Ensure that your `~/go/bin` is on `$PATH`
- Install goimports: `go install golang.org/x/tools/cmd/goimports@latest`
- Install golangci-lint: `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6`
- Install gofumpt: `go install mvdan.cc/gofumpt@latest`
- Install dlv for local debugging: `go install github.com/go-delve/delve/cmd/dlv@latest`, although personally I wouldn't bother, as the current state of VSCode+Delve is not great (as of 2025-04-05)
- Run this to exec all the pre-commit hooks on the entire project: `pre-commit run --all-files`
- Copy `.env.example` to `.env` and configure it with your settings. Having a `.env` file is required, but environment variables will override entries in `.env`.
- In order for linting to work in VSCode (per the workspace config file), for now, you must copy `~/go/bin/golangci-lint` to `~/go/bin/golangci-lint-v2`

### API Endpoints

- `/metrics`: Prometheus metrics endpoint
- `/action-runs`: JSON endpoint for action runs data

## License

This project is licensed under the MIT License - see the LICENSE file for details.
