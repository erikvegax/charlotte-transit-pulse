package gtfsrt

import (
	"os"
	"testing"
)

func TestLoadVehiclesFromFeed(t *testing.T) {
	data, err := os.ReadFile("testdata/vehiclepositions.pb")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	vehicles, err := LoadVehiclesFromFeed(data)
	if err != nil {
		t.Fatalf("LoadVehiclesFromFeed returned an error: %v", err)
	}

	if len(vehicles) != 84 {
		t.Errorf("expected 84 vehicles, got %d", len(vehicles))
	}

	first := vehicles[0]
	if first.Trip.TripID != "5441986" {
		t.Errorf("expected first vehicle trip id 5441986, got %s", first.Trip.TripID)
	}
	if first.Vehicle.ID != "4133" || first.Vehicle.Label != "2301" {
		t.Errorf("expected vehicle descriptor {4133 2301}, got %+v", first.Vehicle)
	}
	if first.Position.Latitude != 35.30795 || first.Position.Longitude != -80.72122 {
		t.Errorf("expected position {35.30795 -80.72122}, got %+v", first.Position)
	}

	// Sanity check every decoded vehicle, not just the first one — catches
	// e.g. a getter chain that silently degrades to a zero value.
	for _, v := range vehicles {
		if v.Position.Latitude < 34.5 || v.Position.Latitude > 36.0 {
			t.Errorf("vehicle %s latitude %v is outside the Charlotte area", v.Vehicle.ID, v.Position.Latitude)
		}
		if v.Position.Longitude < -81.5 || v.Position.Longitude > -80.0 {
			t.Errorf("vehicle %s longitude %v is outside the Charlotte area", v.Vehicle.ID, v.Position.Longitude)
		}
	}
}

func TestLoadTripUpdatesFromFeed(t *testing.T) {
	data, err := os.ReadFile("testdata/tripupdates.pb")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	trips, err := LoadTripUpdatesFromFeed(data)
	if err != nil {
		t.Fatalf("LoadTripUpdatesFromFeed returned an error: %v", err)
	}

	if len(trips) != 131 {
		t.Errorf("expected 131 trip updates, got %d", len(trips))
	}

	missingVehicle := 0
	for _, tu := range trips {
		if tu.Vehicle.ID == "" {
			missingVehicle++
		}
		for _, stu := range tu.StopTimeUpdate {
			if stu.StopID == "" {
				t.Errorf("trip %s has a stop time update with an empty stop id", tu.Trip.TripID)
			}
		}
	}

	// This fixture is a frozen real-world snapshot known to contain trip
	// updates with no vehicle assigned yet. Regression test for the bug
	// where mapping Vehicle via *tu.GetVehicle().Id panicked on these.
	if missingVehicle != 9 {
		t.Errorf("expected 9 trip updates with no vehicle assigned, got %d", missingVehicle)
	}
}

func TestLoadAlertsFromFeed(t *testing.T) {
	data, err := os.ReadFile("testdata/alerts.pb")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	alerts, err := LoadAlertsFromFeed(data)
	if err != nil {
		t.Fatalf("LoadAlertsFromFeed returned an error: %v", err)
	}

	if len(alerts) != 8 {
		t.Errorf("expected 8 alerts, got %d", len(alerts))
	}

	first := alerts[0]
	if first.ID != "23744" {
		t.Errorf("expected first alert id 23744, got %s", first.ID)
	}
	if first.Cause != "CONSTRUCTION" || first.Effect != "DETOUR" {
		t.Errorf("expected cause/effect CONSTRUCTION/DETOUR, got %s/%s", first.Cause, first.Effect)
	}
	if first.HeaderText != "Construction will create Detour" {
		t.Errorf("unexpected header text: %s", first.HeaderText)
	}
	if len(first.InformedEntity) != 1 || first.InformedEntity[0].RouteID != "9" || first.InformedEntity[0].StopID != "45104" {
		t.Errorf("unexpected informed entity: %+v", first.InformedEntity)
	}

	for _, a := range alerts {
		if a.ID == "" {
			t.Errorf("alert has empty id")
		}
		if a.HeaderText == "" {
			t.Errorf("alert %s has empty header text", a.ID)
		}
	}
}
