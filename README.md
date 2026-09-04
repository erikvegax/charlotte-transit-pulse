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
├── services/
│   ├── ingest/          # Polls GTFS-RT feed, publishes to Kafka
│   ├── consumer/        # Kafka consumers, computes delays, writes to MongoDB
│   └── api/             # REST API serving processed data
├── web/                 # React dashboard
├── docs/
│   └── ARCHITECTURE.md
└── README.md
```

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
docker compose up -d

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

- Data ingestion service (Go): poll GTFS-RT, publish to Kafka; seed static schedule into MongoDB
- Kafka consumers + MongoDB storage: join real-time positions against schedule, persist computed delays
- Go REST API layer: current positions, on-time performance, historical trends
- React dashboard: live map + stats view
- Equity overlay: cross-reference with 311 / CMPD data by neighborhood
- Tests, CI/CD, deploy, docs

## License

MIT
