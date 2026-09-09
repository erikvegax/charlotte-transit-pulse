# consumer

Subscribes to the Kafka topics published by `ingest`, joins real-time trip updates against the static schedule to compute delay, and writes processed results to MongoDB.

See [`docs/ARCHITECTURE.md`](../../docs/ARCHITECTURE.md) for how this fits into the overall pipeline.

## Responsibilities

- Consume `vehicle-positions`, `trip-updates`, and `service-alerts` from Kafka
- Compute `delay_seconds` per stop per trip (actual vs. scheduled arrival time)
- Upsert current vehicle positions into the `vehicles` collection (used for the live map)
- Append historical delay records into the `stop_events` collection
- Aggregate daily on-time performance per route into `daily_route_stats`

## Running locally

Requires Kafka and MongoDB running (see [`deploy/docker-compose.yml`](../../deploy/docker-compose.yml)), and the `ingest` service publishing events.

```bash
go run .
```

## Environment variables

| Variable | Description | Default |
|---|---|---|
| `KAFKA_BROKERS` | Comma-separated Kafka broker address(es) | `localhost:9092` |
| `KAFKA_CONSUMER_GROUP` | Consumer group ID | `transit-pulse-consumer` |
| `MONGO_URI` | MongoDB connection string | `mongodb://localhost:27017/transit_pulse` |