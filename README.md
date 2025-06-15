# Aircraft

Combines local ADS-B data with other data sources and notifies about interesting aircraft

## Overview

### Key features

- Triggers Push Notifications when interesting aircraft are spotted by tar1090
- Keeps track of all seen aircraft

### Diagram

```mermaid
flowchart TD
    subgraph NATS["NATS Message Bus"]
        DMessage["'aircraft' subject"]
    end

    Readsb["readsb/tar1090"]
    Discovery["`**Discovery Service**`"]
    Notifier["`**Notifier Service**
    Identifies new and interesting aircraft`"]
    Historian["`**Historian Service**
    Keeps track of all aircraft ever seen`"]
    Pushover["Pushover"]
    MongoDB["MongoDB"]
    Stats["`**Stats Service**`"]
    User["`**User**`"]

    Discovery -->|Polls for currently visible aircraft| Readsb
    Discovery -->|Publishes visible aircraft to| DMessage
    DMessage -->|Forwarded on to| Notifier
    DMessage -->|Forwarded on to| Historian

    Notifier -->|Uses HTTP API to send Push Notifications| Pushover
    Historian -->|Stores all seen aircraft in| MongoDB

    User -->|Uses HTTP API| Stats
    Stats -->|Queries| MongoDB

```

## Getting started

### Using devspace

Install NATS

```bash
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo update
helm install my-nats nats/nats
```

```bash
cd discoverer
devspace dev
```

## TODO list

- [ ] Re-implement`discoverer` service using Go
  - [ ] Figure out best approach for enriching data
  - [ ] Include data about whether aircraft is interesting
- [ ] Re-add NATS
- [ ] Re-implement `notifier` service using Go
