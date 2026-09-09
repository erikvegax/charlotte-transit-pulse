package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/erikvegax/charlotte-transit-pulse/pkg/model"
	se "github.com/segmentio/kafka-go"
)

const (
	vehiclesTopic      = "vehicle-positions"
	tripUpdatesTopic   = "trip-updates"
	serviceAlertsTopic = "service-alerts"
)

type KafkaWriter struct {
	Writer *se.Writer
}

func NewKafkaWriter(brokers []string) *KafkaWriter {
	return &KafkaWriter{
		Writer: se.NewWriter(se.WriterConfig{Brokers: brokers}),
	}
}

func (w *KafkaWriter) PublishVehicles(ctx context.Context, vehicles []model.Vehicle) error {
	messages := []se.Message{}

	for _, v := range vehicles {
		value, err := json.Marshal(v)
		if err != nil {
			log.Printf("error marshaling %v: %v", v.Vehicle.ID, err)
			continue
		}

		messages = append(messages, se.Message{
			Key:   []byte(v.Vehicle.ID),
			Topic: vehiclesTopic,
			Value: value,
		})
	}

	return w.Writer.WriteMessages(ctx, messages...)
}

func (w *KafkaWriter) PublishTripUpdates(ctx context.Context, tripUpdates []model.TripUpdate) error {
	messages := []se.Message{}

	for _, tu := range tripUpdates {
		value, err := json.Marshal(tu)
		if err != nil {
			log.Printf("error marshaling %v: %v", tu.Trip.TripID, err)
			continue
		}

		messages = append(messages, se.Message{
			Key:   []byte(tu.Trip.TripID),
			Topic: tripUpdatesTopic,
			Value: value,
		})
	}

	return w.Writer.WriteMessages(ctx, messages...)
}

func (w *KafkaWriter) PublishAlerts(ctx context.Context, alerts []model.Alert) error {
	messages := []se.Message{}

	for _, a := range alerts {
		value, err := json.Marshal(a)
		if err != nil {
			log.Printf("error marshaling %v: %v", a.ID, err)
			continue
		}

		messages = append(messages, se.Message{
			Key:   []byte(a.ID),
			Topic: serviceAlertsTopic,
			Value: value,
		})
	}

	return w.Writer.WriteMessages(ctx, messages...)
}
