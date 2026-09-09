# Architecture

## Overview

Charlotte Transit Pulse is an event-driven pipeline: a Go service polls CATS' GTFS-Realtime feed, publishes normalized events to Kafka, a consumer service joins those events against the static schedule and writes computed reliability data to MongoDB, and a Go API serves that data to a React dashboard.

```
CATS GTFS-RT feed
       │  (poll every N seconds)
       ▼
┌──────────────┐      publish       ┌───────┐      consume       ┌───────────────┐
│  ingest svc  │ ─────────────────► │ Kafka │ ─────────────────► │  consumer svc  │
│    (Go)      │                    └───────┘                    │     (Go)       │
└──────────────┘                                                 └───────┬────────┘
                                                                          │ write
                                                                          ▼
                                                                    ┌──────────┐
                                                                    │ MongoDB  │
                                                                    └────┬─────┘
                                                                         │ read
                                                                         ▼
                                                                   ┌──────────┐
                                                                   │  API svc │
                                                                   │   (Go)   │
                                                                   └────┬─────┘
                                                                        │ REST
                                                                        ▼
                                                                  ┌───────────┐
                                                                  │  React    │
                                                                  │ dashboard │
                                                                  └───────────┘
```

## Components

### 1. Ingest service (Go)

Polls the CATS GTFS-Realtime feed on a fixed interval (e.g. every 15–30 seconds), decodes the protobuf payload, and publishes normalized JSON events to Kafka.

Responsibilities:
- Poll `VehiclePositions`, `TripUpdates`, and `ServiceAlerts` feeds
- Decode protobuf into internal event structs
- Publish to Kafka topics (see below)
- One-time (and periodic re-sync) load of the static GTFS schedule — routes, stops, trips, stop_times — used as reference data by the consumer service

Failure handling: if a poll fails or the feed is unreachable, log and retry with backoff rather than crash the loop — CATS' feed availability isn't guaranteed.

### 2. Kafka topics

| Topic | Key | Payload | Notes |
|---|---|---|---|
| `vehicle-positions` | `vehicle_id` | lat/lon, route_id, trip_id, timestamp | High volume, short retention |
| `trip-updates` | `trip_id` | stop-level arrival/departure predictions | Used to compute delay |
| `service-alerts` | `alert_id` | affected routes/stops, description | Low volume |

Partition count and retention are tuned for a portfolio-scale deployment (not production transit-agency scale) — a few partitions per topic and a short retention window (e.g. 24–48h) is enough, since MongoDB is the durable store.

### 3. Consumer service (Go)

Subscribes to the Kafka topics, joins real-time `trip-updates` against the static schedule (expected arrival time) to compute delay, and writes results to MongoDB.

Responsibilities:
- Compute `delay_seconds` per stop per trip (actual vs. scheduled)
- Aggregate into rolling on-time performance metrics per route
- Upsert current vehicle positions (for the live map) and append historical delay records (for trend analysis)

### 4. MongoDB schema (initial draft)

- `vehicles` — current position per vehicle, upserted on each update (used for the live map)
- `stop_events` — one document per (trip, stop) with scheduled vs. actual time and computed delay (used for historical/aggregate queries)
- `routes` / `stops` — static reference data from the GTFS schedule feed
- `daily_route_stats` — precomputed daily aggregates (on-time %, avg delay) per route, to keep dashboard queries fast

### 5. API service (Go)

REST API over the MongoDB data. Candidate endpoints:
- `GET /vehicles` — current positions for the live map
- `GET /routes/:id/performance?range=` — on-time performance over a time range
- `GET /routes/performance` — all routes, sortable by reliability
- `GET /equity/overlay` — reliability data joined with 311/CMPD data by neighborhood (Phase 2)

Instrumented with Prometheus metrics (request latency, Kafka consumer lag, ingest poll success/failure rate) since that observability stack is already familiar territory.

### 6. React dashboard

- Live map view (vehicle positions, refreshed on a short poll or via websocket in a later iteration)
- Route performance table/chart (on-time % by route, sortable)
- Equity overlay view (Phase 2)

## Equity overlay (Phase 2)

Charlotte's 311 and CMPD open datasets are keyed by location/neighborhood. The overlay work joins `daily_route_stats` (by route/stop geography) against these datasets at the neighborhood level to explore whether transit reliability correlates with other neighborhood-level indicators — the same kind of "equity scorecard" angle discussed as a civic tech direction. This is exploratory and additive; it doesn't block the core pipeline from being useful and demoable on its own.

## Infrastructure

No self-managed infrastructure is required for this project. Kafka and MongoDB are both used as managed cloud services, and each Go service is deployed independently as a small app that connects out to them over the public internet using a connection string + credentials — the same pattern as connecting to any hosted database or API. There's no private networking, service discovery, or container orchestration to set up.

### Local development

Kafka and MongoDB run locally via Docker Compose (`deploy/docker-compose.yml`) while building. See [`deploy/README.md`](../deploy/README.md) for the full command reference (start, verify, logs, reset) — quick version:

```bash
cd deploy
docker compose up -d
```

- Kafka broker (single-node, KRaft mode — no Zookeeper): `localhost:9092`
- MongoDB: `mongodb://localhost:27017/transit_pulse`

Each Go service reads these as environment variables (`KAFKA_BROKERS`, `MONGO_URI`), so switching from local Docker services to the managed cloud versions below is just a matter of changing env vars — no code changes.

### Managed services (deployment)

| Piece | Provider | Why |
|---|---|---|
| Kafka | [Upstash Kafka](https://upstash.com/) | Serverless, pay-per-request, generous free tier — no broker to operate |
| MongoDB | [MongoDB Atlas](https://www.mongodb.com/atlas) | Free-forever M0 tier (512MB) is enough for this project's scale |
| `ingest`, `consumer`, `api` services | [Railway](https://railway.app/) | Connects to the GitHub repo, builds each service from its Dockerfile, assigns each a URL, handles env vars in a simple dashboard |
| React dashboard | [Vercel](https://vercel.com/) | Same deploy pattern already used for the NC Native Plant Bloom Planner |

Each Go service gets its own environment variables set in Railway's dashboard (Kafka broker address + auth, Mongo connection string). The `api` service gets a public URL from Railway, which the React app calls directly — no additional routing layer needed at this scale.

### Cost

| Piece | Expected cost |
|---|---|
| MongoDB Atlas (M0) | $0 |
| Upstash Kafka | $0 at this project's request volume |
| Railway (3 always-on Go services) | ~$5–15/month once past free-tier usage |
| Vercel (React) | $0 |

Total: free while developing locally; roughly $5–15/month once deployed and running continuously. If cost becomes a concern, the `ingest` service is the one most compatible with scale-to-zero / scheduled-run patterns, since it doesn't need to be always-on the way the API does — though that trades off how current the live data stays.

## Open questions / future work

- Websocket push for live vehicle updates instead of polling from the frontend
- Longer-term historical storage/rollups if daily aggregates outgrow MongoDB's comfort zone for the free tier
- Whether to expose the equity overlay as a separate public-facing view, given the sensitivity of crime data at fine geographic granularity