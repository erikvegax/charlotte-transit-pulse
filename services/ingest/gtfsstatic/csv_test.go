package gtfsstatic

import (
	"os"
	"testing"
)

func TestLoadStaticRouteDataFromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/routes.txt")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	routes, err := LoadStaticRouteDataFromFile(data)
	if err != nil {
		t.Fatalf("LoadStaticRouteDataFromFile returned an error: %v", err)
	}

	if len(routes) != 63 {
		t.Errorf("expected 63 routes, got %d", len(routes))
	}

	first := routes[0]
	if first.RouteID != "1" || first.RouteShortName != "1" || first.RouteLongName != "Mt. Holly Road" || first.RouteType != 3 {
		t.Errorf("unexpected first route: %+v", first)
	}

	// route_id "63x" and similar lettered express routes are why RouteID
	// must stay a string, not an int — regression check for that.
	found63x := false
	for _, r := range routes {
		if r.RouteID == "63x" {
			found63x = true
		}
		if r.RouteID == "" {
			t.Errorf("route has empty route id")
		}
	}
	if !found63x {
		t.Errorf(`expected to find lettered route "63x" in the fixture`)
	}
}

func TestLoadStaticStopDataFromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/stops.txt")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	stops, err := LoadStaticStopDataFromFile(data)
	if err != nil {
		t.Fatalf("LoadStaticStopDataFromFile returned an error: %v", err)
	}

	if len(stops) != 12 {
		t.Errorf("expected 12 stops, got %d", len(stops))
	}

	first := stops[0]
	if first.StopID != "00001" || first.StopName != "7th St Station" {
		t.Errorf("unexpected first stop: %+v", first)
	}
	if first.StopLat != 35.227375 || first.StopLon != -80.838088 {
		t.Errorf("unexpected first stop coordinates: %+v", first)
	}

	// The fixture includes two location_type=1 "station" rows — verify
	// that hierarchy is preserved, not silently dropped.
	stations := 0
	for _, s := range stops {
		if s.LocationType == 1 {
			stations++
			if s.StopID == "" {
				t.Errorf("station stop has empty stop id")
			}
		}
	}
	if stations != 2 {
		t.Errorf("expected 2 station-type stops, got %d", stations)
	}
}

func TestLoadStaticTripDataFromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/trips.txt")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	trips, err := LoadStaticTripDataFromFile(data)
	if err != nil {
		t.Fatalf("LoadStaticTripDataFromFile returned an error: %v", err)
	}

	if len(trips) != 10 {
		t.Errorf("expected 10 trips, got %d", len(trips))
	}

	first := trips[0]
	if first.TripID != "5307109" || first.RouteID != "501" || first.ServiceID != "BLE-Saturday-0-RAIL-26-01-0000010-" {
		t.Errorf("unexpected first trip: %+v", first)
	}
	if first.TripHeadsign != "South to I-485/S. Blvd" || first.DirectionID != 1 {
		t.Errorf("unexpected first trip headsign/direction: %+v", first)
	}
	if first.BlockID != "983906" || first.ShapeID != "5010107" {
		t.Errorf("unexpected first trip block/shape: %+v", first)
	}
}

func TestLoadStaticStopTimeDataFromFile(t *testing.T) {
	data, err := os.ReadFile("testdata/stop_times.txt")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	stopTimes, err := LoadStaticStopTimeDataFromFile(data)
	if err != nil {
		t.Fatalf("LoadStaticStopTimeDataFromFile returned an error: %v", err)
	}

	if len(stopTimes) != 12 {
		t.Errorf("expected 12 stop times, got %d", len(stopTimes))
	}

	first := stopTimes[0]
	if first.TripID != "5307109" || first.ArrivalTime != "23:18:00" || first.StopID != "00090" || first.StopSequence != 1 {
		t.Errorf("unexpected first stop time: %+v", first)
	}

	// GTFS allows arrival/departure times past 24:00:00 for trips that
	// continue after midnight of the same service day — this is why these
	// fields must stay strings, not a parsed time type. Regression check
	// that such rows still parse rather than getting rejected.
	foundPastMidnight := false
	for _, st := range stopTimes {
		if st.ArrivalTime == "24:02:00" {
			foundPastMidnight = true
		}
	}
	if !foundPastMidnight {
		t.Errorf("expected to find a past-midnight arrival time (24:02:00) in the fixture")
	}
}
