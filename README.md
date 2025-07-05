# Aircraft

Combines local ADS-B data with other data sources and notifies about interesting aircraft.

## Overview

### Key features

- Triggers Push Notifications when interesting aircraft are spotted by tar1090
- Keeps track of all seen aircraft

### Diagram

![C4 Model-style "Container" diagram](docs/Aircraft-Excalidraw-2025-07-03-1721.svg)

## Getting started

### Using `devspace` and K8s

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

## Tests

### Testing strategy

This is a personal project I'm using to learn and experiment. I want to be able to make changes as easily as possible, therefore I'm intentionally keeping automated testing very minimal for now. I’ll add tests where they make sense and where they help me understand or validate something specific, but I’m not aiming for the type of test coverage I'd expect from production-grade software.

### Running unit tests

Currently only the `enricher` service has unit tests

```bash
make unit-tests
```

## TODO list

- [x] Re-implement`discoverer` service using Go
  - [x] Publish aircraft when found
- [x] Re-add NATS
- [x] Implement `enricher` service
  - [x] Enrich with HexDB data
  - [x] Enrich with PlaneAlertDb data (i.e. whether it's an interesting aircraft and why)
  - [x] What should happen if an enricher fails? Should it continue? Later enricher may be able to fill in gaps
  - [ ] Use Context appropriately with enrichers (e.g. set deadline, cancel gracefully)
  - [ ] Investigate any other potential data sources
- [x] Implement `evaluator` service
  - [x] Implement logic to identify interesting aircraft
- [x] Implement `notifier` service
  - [x] Publish notifications using Pushover
- [ ] Add Postgres
- [ ] Add `historian` service
- [ ] Add `stats` service
- [ ] Add a `monitoring` service to keep track of enrichment failures (could use this to compare sources and find backup enrichers for when one source doesn't have any details)
- [ ] Other / refactoring / future
  - [ ] Try out https://github.com/caarlos0/env
  - [ ] Back-off approach for enrichers if they fail with certain error codes (not 404)
  - [ ] Add a "system context" C4-style diagram
  - [ ] Add some architecture docs explaining design choices
