package model

type TripUpdate struct {
	Trip           TripDescriptor    `json:"trip"`
	Vehicle        VehicleDescriptor `json:"vehicle,omitempty"`
	StopTimeUpdate []StopUpdate      `json:"stopTimeUpdate"`
	Timestamp      uint64            `json:"timestamp"`
}

type TripDescriptor struct {
	TripID    string `json:"tripId"`
	RouteID   string `json:"routeId"`
	StartTime string `json:"startTime"`
	StartDate string `json:"startDate"`
}

type StopUpdate struct {
	StopSequence uint32 `json:"stopSequence"`
	StopID       string `json:"stopId"`
	Arrival      int64  `json:"arrival"`
	Departure    int64  `json:"departure"`
}
