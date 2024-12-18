# beacon

Monitoring and analysing multiple device metrics from a cloud dashboard.

| Deliverable | Stack | Function | Status |
| --- | --- | --- | --- |
| [daemon](daemon/) | Go 1.23 | Serve metrics to Aggregator and run commands from Aggregator | [![Daemon CI](https://github.com/bxrne/beacon/actions/workflows/daemon-ci.yaml/badge.svg)](https://github.com/bxrne/beacon/actions/workflows/daemon-ci.yaml) |
| [web](web/) | Go 1.23 | Display metrics, command center and host API for Aggregator | [![CI](https://github.com/bxrne/beacon/actions/workflows/web-ci.yaml/badge.svg)](https://github.com/bxrne/beacon/actions/workflows/web-ci.yaml) | 
| [aggregator](aggregator/) | Go 1.23 | Poll devices for metrics, Poll API for commands | [![Aggregator CI](https://github.com/bxrne/beacon/actions/workflows/aggregator-ci.yaml/badge.svg)](https://github.com/bxrne/beacon/actions/workflows/aggregator-ci.yaml) |
| [diorama](diorama/) | C, FreeRTOS | Interactive diorama of pedestrian crossing | [![PlatformIO CI](https://github.com/bxrne/beacon/actions/workflows/diorama-ci.yaml/badge.svg)](https://github.com/bxrne/beacon/actions/workflows/diorama-ci.yaml) |


## Usage

### Daemon

```sh
cd daemon
go mod download

go run ./cmd

go test ./...
```

### Aggregator

```sh
cd aggregator
go mod download

go run ./cmd

go test ./...
```


### web

```sh
cd web
swag init -g ./cmd/main.go # Generate swagger docs 
go mod download

go run ./cmd

go test ./...
```

## Deployment

### web

web is deployed via `Dockerfile` using [fly](https://fly.io/) which is configured [here](web/fly.toml).

```sh
cd web
fly auth login
fly deploy
```

## Diorama

```sh
cd diorama
pio run -t upload
pio device monitor
```
