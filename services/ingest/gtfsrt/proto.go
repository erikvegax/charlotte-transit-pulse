package gtfsrt

import (
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/erikvegax/charlotte-transit-pulse/pkg/model"
	"google.golang.org/protobuf/proto"
)

func LoadVehiclesFromFeed(data []byte) ([]model.Vehicle, error) {
	var feed gtfs.FeedMessage
	err := proto.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	var vehicles []model.Vehicle
	for _, entity := range feed.GetEntity() {
		if vehicle := entity.GetVehicle(); vehicle != nil {
			vehicles = append(vehicles, model.Vehicle{
				Trip: model.TripDescriptor{
					TripID:    vehicle.GetTrip().GetTripId(),
					RouteID:   vehicle.GetTrip().GetRouteId(),
					StartTime: vehicle.GetTrip().GetStartTime(),
					StartDate: vehicle.GetTrip().GetStartDate(),
				},
				Vehicle: model.VehicleDescriptor{
					ID:    vehicle.GetVehicle().GetId(),
					Label: vehicle.GetVehicle().GetLabel(),
				},
				Position: model.Position{
					Latitude:  vehicle.GetPosition().GetLatitude(),
					Longitude: vehicle.GetPosition().GetLongitude(),
				},
				CurrentStopSequence: vehicle.GetCurrentStopSequence(),
				StopID:              vehicle.GetStopId(),
				Timestamp:           vehicle.GetTimestamp(),
			})
		}
	}
	return vehicles, nil
}

func LoadAlertsFromFeed(data []byte) ([]model.Alert, error) {
	var feed gtfs.FeedMessage
	err := proto.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	var alerts []model.Alert
	for _, entity := range feed.GetEntity() {
		if alert := entity.GetAlert(); alert != nil {
			alerts = append(alerts, model.Alert{
				ID: entity.GetId(),
				ActivePeriod: func() []model.ActivePeriod {
					var periods []model.ActivePeriod
					for _, ap := range alert.GetActivePeriod() {
						periods = append(periods, model.ActivePeriod{
							StartTime: ap.GetStart(),
							EndTime:   ap.GetEnd(),
						})
					}
					return periods
				}(),
				InformedEntity: func() []model.InformedEntity {
					var entities []model.InformedEntity
					for _, ie := range alert.GetInformedEntity() {
						entities = append(entities, model.InformedEntity{
							RouteID: ie.GetRouteId(),
							StopID:  ie.GetStopId(),
						})
					}
					return entities
				}(),
				Cause:  alert.GetCause().String(),
				Effect: alert.GetEffect().String(),
				HeaderText: func() string {
					for _, t := range alert.GetHeaderText().GetTranslation() {
						if t.GetLanguage() == "en" {
							return t.GetText()
						}
					}
					return ""
				}(),
				DescriptionText: func() string {
					for _, t := range alert.GetDescriptionText().GetTranslation() {
						if t.GetLanguage() == "en" {
							return t.GetText()
						}
					}
					return ""
				}(),
			})
		}
	}
	return alerts, nil
}

func LoadTripUpdatesFromFeed(data []byte) ([]model.TripUpdate, error) {
	var feed gtfs.FeedMessage
	err := proto.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	var trips []model.TripUpdate

	for _, entity := range feed.GetEntity() {
		if tu := entity.GetTripUpdate(); tu != nil {
			trips = append(trips, model.TripUpdate{
				Trip: model.TripDescriptor{
					TripID:    tu.GetTrip().GetTripId(),
					RouteID:   tu.GetTrip().GetRouteId(),
					StartTime: tu.GetTrip().GetStartTime(),
					StartDate: tu.GetTrip().GetStartDate(),
				},
				Vehicle: model.VehicleDescriptor{
					ID:    tu.GetVehicle().GetId(),
					Label: tu.GetVehicle().GetLabel(),
				},
				StopTimeUpdate: func() []model.StopUpdate {
					var updates []model.StopUpdate
					for _, stu := range tu.GetStopTimeUpdate() {
						updates = append(updates, model.StopUpdate{
							StopSequence: stu.GetStopSequence(),
							StopID:       stu.GetStopId(),
							Arrival:      stu.GetArrival().GetTime(),
							Departure:    stu.GetDeparture().GetTime(),
						})
					}
					return updates
				}(),
				Timestamp: tu.GetTimestamp(),
			})
		}
	}

	return trips, nil
}
