package model

type Stop struct {
	StopID        string  `json:"stopId" bson:"stopId"`
	StopCode      string  `json:"stopCode" bson:"stopCode" `
	StopName      string  `json:"stopName" bson:"stopName"`
	LocationType  int64   `json:"locationType" bson:"locationType"`
	ParentStation string  `json:"parentStation" bson:"parentStation"`
	StopLat       float64 `json:"stopLat" bson:"stopLat"`
	StopLon       float64 `json:"stopLon" bson:"stopLon"`
}
