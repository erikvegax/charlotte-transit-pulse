# Charlotte Transit Pulse

Real-time transit reliability tracking and analytics for CATS (Charlotte Area Transit System), built on an event-driven Go/Kafka/MongoDB backend with a React dashboard.

Charlotte Transit Pulse ingests CATS' live GTFS-Realtime feed, computes on-time performance against the published schedule, and surfaces route-level reliability data — including a neighborhood-level view that layers in Charlotte's open 311 and public-safety data to explore whether service reliability tracks with other indicators of neighborhood investment.

## Why this project

Public transit agencies publish real-time vehicle data, but rarely make reliability trends easy to see. This project turns a raw GTFS-RT feed into something people can actually look at: which routes run on time, which don't, and when. It's also a testbed for an event-driven architecture pattern (Go services + Kafka + MongoDB) applied to a public dataset instead of internal company data.

## Features

- **Live vehicle tracking** — near-real-time map of CATS buses and light rail
- **On-time performance metrics** — computed delay per route, aggregated by time of day and day of week
- **Historical trends** — which routes are reliable, which are not, and whether that's changing over time
- **Equity overlay** — reliability data cross-referenced against Charlotte 311 and CMPD open data by neighborhood

See [ARCHITECTURE.md](./ARCHITECTURE.md) for how the pieces fit together.

## Tech stack

| Layer | Technology |
|---|---|
| Ingestion & processing | Go |
| Messaging | Kafka |
| Storage | MongoDB |
| API | Go (REST) |
| Frontend | React |
| Observability | Prometheus + Grafana |
| Deployment | TBD (API/consumers: Fly.io or Railway; frontend: Vercel) |

## Data sources

- [CATS GTFS-Realtime feed](https://www.transit.land/feeds/f-dnq-charlotteareatransitsystem~rt) (via Transitland) — live vehicle positions, trip updates, service alerts
- [CATS static GTFS schedule](https://mobilitydatabase.org/feeds/gtfs/mdb-2265/map) (via Mobility Database) — routes, stops, scheduled trip times
- [Charlotte 311 Service Requests](https://data.charlottenc.gov/datasets/service-requests-311/api) — City of Charlotte Open Data Portal
- [CMPD Violent Crime dataset](https://data.charlottenc.gov/datasets/charlotte::cmpd-violent-crime/about) — City of Charlotte Open Data Portal

## Project structure

```
charlotte-transit-pulse/
├── .github/
│   └── workflows/
│       ├── ci.yml              # lint + test on PR for Go services
│       └── ci-web.yml          # lint + test on PR for React app
├── services/
│   ├── ingest/
│   │   ├── main.go
│   │   ├── gtfsrt/              # protobuf decoding, feed polling logic
│   │   ├── kafka/                # producer setup
│   │   ├── go.mod
│   │   └── README.md
│   ├── consumer/
│   │   ├── main.go
│   │   ├── delay/                 # delay computation logic
│   │   ├── mongo/                 # storage layer
│   │   ├── kafka/                  # consumer setup
│   │   ├── go.mod
│   │   └── README.md
│   └── api/
│       ├── main.go
│       ├── handlers/
│       ├── mongo/
│       ├── go.mod
│       └── README.md
├── pkg/                          # shared Go code (types used by multiple services)
│   ├── models/                    # Vehicle, StopEvent, Route structs
│   └── config/                    # shared env/config loading
├── web/
│   ├── src/
│   ├── package.json
│   └── README.md
├── deploy/
│   ├── docker-compose.yml         # local Kafka + MongoDB for dev
│   ├── ingest.Dockerfile
│   ├── consumer.Dockerfile
│   └── api.Dockerfile
├── docs/
│   └── ARCHITECTURE.md
├── .gitignore
├── LICENSE
└── README.md
```

Each Go service (`services/ingest`, `services/consumer`, `services/api`) has its own `go.mod`, so it can be built and deployed independently. Shared types (e.g. `Vehicle`, `StopEvent`) live in `pkg/` rather than being duplicated across services. Service-specific setup lives in each service's own README; this root README stays high-level.

## Getting started

### Prerequisites

- Go 1.22+
- Node 20+
- Kafka (local via Docker Compose, or a managed instance)
- MongoDB (local via Docker Compose, or Atlas)

### Setup

```bash
# clone
git clone https://github.com/erikvegax/charlotte-transit-pulse.git
cd charlotte-transit-pulse

# start Kafka + MongoDB locally
cd deploy && docker compose up -d && cd ..

# run the ingestion service
cd services/ingest && go run .

# run the consumer service
cd services/consumer && go run .

# run the API
cd services/api && go run .

# run the frontend
cd web && npm install && npm run dev
```

Environment variables and per-service configuration are documented in each service's own README (TBD as services are built).

## Roadmap

- [1] Data ingestion service (Go): poll GTFS-RT, publish to Kafka; seed static schedule into MongoDB
- [2] Kafka consumers + MongoDB storage: join real-time positions against schedule, persist computed delays
- [3] Go REST API layer: current positions, on-time performance, historical trends
- [4] React dashboard: live map + stats view
- [5] Equity overlay: cross-reference with 311 / CMPD data by neighborhood
- [6] Tests, CI/CD, deploy, docs

## License
MIT