package model

type Alert struct {
	ID              string           `json:"id"`
	ActivePeriod    []ActivePeriod   `json:"activePeriod"`
	InformedEntity  []InformedEntity `json:"informedEntity"`
	Cause           string           `json:"cause"`
	Effect          string           `json:"effect"`
	HeaderText      string           `json:"headerText"`
	DescriptionText string           `json:"descriptionText"`
}

type ActivePeriod struct {
	StartTime uint64 `json:"startTime"`
	EndTime   uint64 `json:"endTime"`
}

type InformedEntity struct {
	RouteID string `json:"routeId"`
	StopID  string `json:"stopId"`
}
