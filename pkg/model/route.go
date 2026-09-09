package model

type Route struct {
	RouteID        string `json:"routeId" bson:"routeId"`
	RouteShortName string `json:"routeShortName" bson:"routeShortName"`
	RouteLongName  string `json:"routeLongName" bson:"routeLongName"`
	RouteType      int64  `json:"routeType" bson:"routeType"`
}
