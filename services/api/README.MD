# api

REST API serving processed transit data from MongoDB to the React dashboard.

See [`docs/ARCHITECTURE.md`](../../docs/ARCHITECTURE.md) for how this fits into the overall pipeline.

## Endpoints (planned)

| Method | Path | Description |
|---|---|---|
| `GET` | `/vehicles` | Current vehicle positions, for the live map |
| `GET` | `/routes/:id/performance` | On-time performance for a route over a time range |
| `GET` | `/routes/performance` | All routes, sortable by reliability |
| `GET` | `/equity/overlay` | Reliability data joined with 311/CMPD data by neighborhood (Phase 2) |

## Running locally

Requires MongoDB running and populated by the `consumer` service (see [`deploy/docker-compose.yml`](../../deploy/docker-compose.yml)).

```bash
go run .
```

## Environment variables

| Variable | Description | Default |
|---|---|---|
| `MONGO_URI` | MongoDB connection string | `mongodb://localhost:27017/transit_pulse` |
| `PORT` | Port the API listens on | `8080` |


## Observability

Instrumented with Prometheus metrics (request latency, Kafka consumer lag surfaced from `consumer`, ingest poll success/failure rate). Metrics endpoint TBD.