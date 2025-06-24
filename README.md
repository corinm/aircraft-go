# Aircraft

Combines local ADS-B data with other data sources and notifies about interesting aircraft

## Overview

### Key features

- In progress: Triggers Push Notifications when interesting aircraft are spotted by tar1090
- Keeps track of all seen aircraft

### Diagram

![C4 Model-style "Container" diagram](docs/Aircraft-Excalidraw-2025-06-24-0942.svg)

## Getting started

### Using devspace

Install NATS

```bash
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo update
helm install my-nats nats/nats
```

Start desired service using dev:

```bash
cd discoverer
devspace dev
```

For other services:

```bash
devspace enter
# Select the service in the menu
# Wait for a shell
go run main.go
```

## TODO list

- [x] Re-implement`discoverer` service using Go
  - [x] Publish aircraft when found
- [x] Re-add NATS
- [ ] Implement `enricher` service
  - [x] Enrich with HexDB data
  - [ ] Enrich with PlaneAlertDb data (i.e. whether it's an interesting aircraft and why)
  - [ ] Investigate any other potential data sources
- [ ] Implement `evaluator` service
  - [ ] Implement logic to identify interesting aircraft
- [ ] Implement `notifier` service
  - [ ] Publish notifications using Pushover
- [ ] Add `historian` service
- [ ] Add `stats` service
