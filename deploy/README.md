# deploy

Local development infrastructure: Kafka + MongoDB, run via Docker Compose. See [`docker-compose.yml`](./docker-compose.yml) for the service definitions, and [`docs/ARCHITECTURE.md`](../docs/ARCHITECTURE.md) for how they fit into the pipeline.

## Prerequisites

Docker Desktop installed and running (the daemon must be up — check the whale icon in the menu bar, or run `docker info`).

## Start

```bash
cd deploy
docker compose up -d
```

First run pulls the Kafka and Mongo images, so it takes longer than subsequent runs. `-d` runs containers in the background (detached) instead of holding your terminal.

- Kafka broker: `localhost:9092`
- MongoDB: `mongodb://localhost:27017/transit_pulse`

## Verify it's actually working

`docker compose ps` showing `Up` only means the process started — worth confirming each service responds before pointing app code at it.

```bash
# containers running + ports mapped?
docker compose ps

# Mongo responding?
docker exec transit-pulse-mongo mongosh --quiet --eval "db.runCommand({ping:1})"
# → { ok: 1 }

# Kafka responding? (lists existing topics — empty is normal until
# a producer publishes, since KAFKA_AUTO_CREATE_TOPICS_ENABLE=true)
docker exec transit-pulse-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
```

## Logs

```bash
docker compose logs -f kafka
docker compose logs -f mongo
```

First place to look if a service won't connect.

## Stop / reset

```bash
# stop containers, keep data (kafka-data / mongo-data volumes persist)
docker compose down

# stop and wipe all data — nuclear reset if state gets weird
docker compose down -v
```
