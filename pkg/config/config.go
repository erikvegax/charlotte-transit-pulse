package config

import "github.com/caarlos0/env/v11"

type Config struct {
	CatsVehiclePositionURL      string   `env:"GTFS_RT_VEHICLE_POSITIONS_URL,required"`
	CatsTripUpdateURL           string   `env:"GTFS_RT_TRIP_UPDATES_URL,required"`
	CatsAlertsURL               string   `env:"GTFS_RT_ALERTS_URL,required"`
	CatsStaticURL               string   `env:"GTFS_STATIC_URL,required"`
	PollIntervalInSeconds       int64    `env:"POLL_INTERVAL_SECONDS" envDefault:"40"`
	StaticResyncIntervalInHours int64    `env:"STATIC_RESYNC_INTERVAL_HOURS" envDefault:"48"`
	KafkaBrokers                []string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	MongoURI                    string   `env:"MONGO_URI" envDefault:"mongodb://localhost:27017/transit_pulse"`
}

func Load() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
