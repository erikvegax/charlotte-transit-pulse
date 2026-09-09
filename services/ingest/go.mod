module github.com/erikvegax/charlotte-transit-pulse/services/ingest

go 1.25.0

require (
	github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs v1.0.0
	github.com/erikvegax/charlotte-transit-pulse/pkg v0.0.0
	github.com/go-co-op/gocron/v2 v2.22.0
	github.com/joho/godotenv v1.5.1
	github.com/segmentio/kafka-go v0.4.51
	go.mongodb.org/mongo-driver/v2 v2.9.0
	google.golang.org/protobuf v1.26.0
)

require (
	github.com/caarlos0/env/v11 v11.4.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/klauspost/compress v1.19.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/text v0.39.0 // indirect
)

replace github.com/erikvegax/charlotte-transit-pulse/pkg => ../../pkg
