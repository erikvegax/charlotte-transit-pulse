# ingest

Polls CATS' GTFS-Realtime feed on a fixed interval, decodes the protobuf payload, and publishes normalized events to Kafka. Also performs a one-time (and periodic re-sync) load of the static GTFS schedule into MongoDB as reference data.

See [`docs/ARCHITECTURE.md`](../../docs/ARCHITECTURE.md) for how this fits into the overall pipeline.

## Responsibilities

- Poll `VehiclePositions`, `TripUpdates`, and `ServiceAlerts` feeds from CATS' GTFS-RT source
- Decode protobuf into internal event structs
- Publish to the `vehicle-positions`, `trip-updates`, and `service-alerts` Kafka topics
- Load/refresh the static GTFS schedule (routes, stops, trips, stop_times) into MongoDB

## Running locally

Requires Kafka and MongoDB running (see [`deploy/docker-compose.yml`](../../deploy/docker-compose.yml)).

```bash
go run .
```

## Environment variables

CATS publishes GTFS-Realtime as three separate feed endpoints (not one combined feed), so each gets its own URL:

| Variable | Description | Default |
|---|---|---|
| `GTFS_RT_VEHICLE_POSITIONS_URL` | CATS VehiclePositions feed URL | — (required) |
| `GTFS_RT_TRIP_UPDATES_URL` | CATS TripUpdates feed URL | — (required) |
| `GTFS_RT_ALERTS_URL` | CATS ServiceAlerts feed URL | — (required) |
| `GTFS_STATIC_URL` | CATS static GTFS schedule (zip) URL | — (required) |
| `POLL_INTERVAL_SECONDS` | How often to poll the real-time feeds | `40` |
| `STATIC_RESYNC_INTERVAL_HOURS` | How often to re-sync the static GTFS schedule into MongoDB | `48` |
| `KAFKA_BROKERS` | Comma-separated Kafka broker address(es) | `localhost:9092` |
| `MONGO_URI` | MongoDB connection string | `mongodb://localhost:27017/transit_pulse` |

Required variables with no default (marked above) must be set — the service fails fast at startup rather than starting in a broken state.