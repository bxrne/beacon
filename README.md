# beacon

Monitoring and analysing multiple device metrics from a cloud dashboard.

| Deliverable | Stack | Function | Status |
| --- | --- | --- | --- |
| [daemon](daemon/) | Go 1.23 | Collect device metrics and send to api | [![Daemon CI](https://github.com/bxrne/beacon/actions/workflows/daemon-ci.yaml/badge.svg)](https://github.com/bxrne/beacon/actions/workflows/daemon-ci.yaml) |
| [api](api/) | Go 1.23 | Receive metrics from daemons and store | [![CI](https://github.com/bxrne/beacon/actions/workflows/api-ci.yaml/badge.svg)](https://github.com/bxrne/beacon/actions/workflows/api-ci.yaml) | 
| [diorama](diorama/) | C, FreeRTOS | Interactive diorama of pedestrian crossing | [![PlatformIO CI](https://github.com/bxrne/beacon/actions/workflows/diorama-ci.yaml/badge.svg)](https://github.com/bxrne/beacon/actions/workflows/diorama-ci.yaml) |
| [dashboard](dashboard/) | React | (TODO) Display metrics from api | TODO |

## Usage

### Daemon

```sh
cd daemon
go mod download

go run ./cmd

go test ./...
```

### API

```sh
cd api
swag init -g ./cmd/main.go # Generate swagger docs 
go mod download

go run ./cmd

go test ./...
```

## Deployment

### Daemon

```sh
cd daemon

go build -o beacon-daemon ./cmd
./beacon-daemon

TODO: Add systemd service file
```

### API

API is deployed via `Dockerfile` using [fly](https://fly.io/) which is configured [here](api/fly.toml).

```sh
cd api
fly auth login
fly deploy
```

## Diorama

```sh
cd diorama
pio run -t upload
pio device monitor
```
