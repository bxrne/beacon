# beacon

Monitoring and analysing multiple device metrics from a cloud dashboard.

| Deliverable | Stack | Function | Status |
| --- | --- | --- | --- |
| [daemon](daemon/) | Go 1.23 | Collect device metrics and send to api | [![Daemon CI](https://github.com/bxrne/beacon/actions/workflows/daemon-ci.yaml/badge.svg)](https://github.com/bxrne/beacon/actions/workflows/daemon-ci.yaml) |
| [api](api/) | Go 1.23 | Receive metrics from daemons and store | [![CI](https://github.com/bxrne/beacon/actions/workflows/api-ci.yaml/badge.svg)](https://github.com/bxrne/beacon/actions/workflows/api-ci.yaml) | 
| [web](web/) | HTMX | (TODO) Display metrics from api | TODO |

## Usage

### Daemon

```sh
cd daemon
go run ./cmd

go test ./...
```

### API

```sh
cd api
go run ./cmd

go test ./...
```

## Deployment

### Daemon

```sh
git clone https://github.com/bxrne/beacon.git
cd beacon/daemon

go build -o beacon-daemon ./cmd
./beacon-daemon
```

### API

API is deployed via `Dockerfile` using [fly](https://fly.io/) which is configured [here](api/fly.toml).

```sh
fly auth login
fly deploy
```

