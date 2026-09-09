package model

type Vehicle struct {
	Trip                TripDescriptor    `json:"trip"`
	Vehicle             VehicleDescriptor `json:"vehicle"`
	Position            Position          `json:"position"`
	CurrentStopSequence uint32            `json:"currentStopSequence"`
	StopID              string            `json:"stopId"`
	Timestamp           uint64            `json:"timestamp"`
}

type VehicleDescriptor struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Position struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
