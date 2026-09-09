package model

type StopTime struct {
	TripID        string `json:"tripId" bson:"tripId"`
	ArrivalTime   string `json:"arrivalTime" bson:"arrivalTime"`
	DepartureTime string `json:"departureTime" bson:"departureTime"`
	StopID        string `json:"stopId" bson:"stopId"`
	StopSequence  int64  `json:"stopSequence" bson:"stopSequence"`
	PickupType    int64  `json:"pickupType" bson:"pickupType"`
	DropOffType   int64  `json:"dropOffType" bson:"dropOffType"`
	Timepoint     int64  `json:"timepoint" bson:"timepoint"`
}
