package model

type Trip struct {
	TripID       string `json:"tripId" bson:"tripId"`
	RouteID      string `json:"routeId" bson:"routeId"`
	ServiceID    string `json:"serviceId" bson:"serviceId"`
	TripHeadsign string `json:"tripHeadsign" bson:"tripHeadsign"`
	DirectionID  int64  `json:"directionId" bson:"directionId"`
	BlockID      string `json:"blockId" bson:"blockId"`
	ShapeID      string `json:"shapeId" bson:"shapeId"`
}
