# Aircraft

Combines local ADS-B data with other data sources and notifies about interesting aircraft.

## Overview

### Key features

- Triggers Push Notifications when interesting aircraft are spotted by tar1090
- Keeps track of all seen aircraft

### Diagram

![C4 Model-style "Container" diagram](docs/Aircraft-Excalidraw-2025-07-03-1721.svg)

## Getting started

### Using `devspace`

#### Install NATS

```bash
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo update
helm install my-nats nats/nats
```

#### To develop a service:

The target service will start in `dev` mode and dependent services will be deployed alongside it in `deploy` mode:

```bash
cd notifier
devspace dev
go run main.go
```

## TODO list

- [x] Re-implement`discoverer` service using Go
  - [x] Publish aircraft when found
- [x] Re-add NATS
- [x] Implement `enricher` service
  - [x] Enrich with HexDB data
  - [ ] Use Context appropriately with enrichers (e.g. set deadline, cancel gracefully)
  - [x] Enrich with PlaneAlertDb data (i.e. whether it's an interesting aircraft and why)
  - [ ] What should happen if an enricher fails? Should it continue? Later enricher may be able to fill in gaps
  - [ ] Investigate any other potential data sources
- [ ] Implement `evaluator` service
  - [ ] Implement logic to identify interesting aircraft
- [x] Implement `notifier` service
  - [x] Publish notifications using Pushover
- [ ] Try out https://github.com/caarlos0/env
- [ ] Add `historian` service
- [ ] Add `stats` service
- [ ] Add a `monitoring` service to keep track of enrichment failures (could use this to compare sources and find backup's for when one source doesn't have any details)
