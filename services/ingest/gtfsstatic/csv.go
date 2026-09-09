package gtfsstatic

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/erikvegax/charlotte-transit-pulse/pkg/model"
)

func LoadStaticStopDataFromFile(file []byte) ([]model.Stop, error) {
	var stops []model.Stop

	reader := csv.NewReader(bytes.NewReader(file))

	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error reading row:", err)
			continue
		}
		if len(row) < 10 {
			log.Println("row has fewer columns than expected")
			continue
		}

		stopId := row[0]
		stopCode := row[1]
		stopName := row[2]
		stopLat := row[4]
		stopLon := row[5]
		locationType := row[8]
		parentStation := row[9]

		stopId = strings.TrimSpace(stopId)
		stopCode = strings.TrimSpace(stopCode)
		stopName = strings.TrimSpace(stopName)
		locationType = strings.TrimSpace(locationType)
		parentStation = strings.TrimSpace(parentStation)
		stopLat = strings.TrimSpace(stopLat)
		stopLon = strings.TrimSpace(stopLon)

		locationTypeInt, err := strconv.ParseInt(locationType, 10, 64)
		if err != nil {
			log.Println("error parsing location type:", err)
			continue
		}

		stopLatFloat, err := strconv.ParseFloat(stopLat, 64)
		if err != nil {
			log.Println("error parsing stop latitude:", err)
			continue
		}

		stopLonFloat, err := strconv.ParseFloat(stopLon, 64)
		if err != nil {
			log.Println("error parsing stop longitude:", err)
			continue
		}

		stop := model.Stop{
			StopID:        stopId,
			StopCode:      stopCode,
			StopName:      stopName,
			LocationType:  locationTypeInt,
			ParentStation: parentStation,
			StopLat:       stopLatFloat,
			StopLon:       stopLonFloat,
		}

		stops = append(stops, stop)
	}

	return stops, nil
}

func LoadStaticRouteDataFromFile(file []byte) ([]model.Route, error) {
	var routes []model.Route

	reader := csv.NewReader(bytes.NewReader(file))

	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error reading row:", err)
			continue
		}
		if len(row) < 5 {
			log.Println("row has fewer columns than expected")
			continue
		}

		routeId := row[0]
		routeShortName := row[1]
		routeLongName := row[2]
		routeType := row[4]

		routeId = strings.TrimSpace(routeId)
		routeShortName = strings.TrimSpace(routeShortName)
		routeLongName = strings.TrimSpace(routeLongName)
		routeType = strings.TrimSpace(routeType)

		routeTypeInt, err := strconv.ParseInt(routeType, 10, 64)
		if err != nil {
			log.Println("error parsing route type:", err)
			continue
		}

		route := model.Route{
			RouteID:        routeId,
			RouteShortName: routeShortName,
			RouteLongName:  routeLongName,
			RouteType:      routeTypeInt,
		}

		routes = append(routes, route)
	}

	return routes, nil
}

func LoadStaticStopTimeDataFromFile(file []byte) ([]model.StopTime, error) {
	var stopTimes []model.StopTime

	reader := csv.NewReader(bytes.NewReader(file))

	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error reading row:", err)
			continue
		}
		if len(row) < 7 {
			log.Println("row has fewer columns than expected")
			continue
		}

		tripId := row[0]
		arrivalTime := row[1]
		departureTime := row[2]
		stopId := row[3]
		stopSequence := row[4]
		pickupType := row[5]
		dropOffType := row[6]
		timepoint := row[7]

		tripId = strings.TrimSpace(tripId)
		arrivalTime = strings.TrimSpace(arrivalTime)
		departureTime = strings.TrimSpace(departureTime)
		stopId = strings.TrimSpace(stopId)
		stopSequence = strings.TrimSpace(stopSequence)
		pickupType = strings.TrimSpace(pickupType)
		dropOffType = strings.TrimSpace(dropOffType)
		timepoint = strings.TrimSpace(timepoint)

		stopSequenceInt, err := strconv.ParseInt(stopSequence, 10, 64)
		if err != nil {
			log.Println("error parsing stop sequence:", err)
			continue
		}

		pickupTypeInt, err := strconv.ParseInt(pickupType, 10, 64)
		if err != nil {
			log.Println("error parsing pickup type:", err)
			continue
		}

		dropOffTypeInt, err := strconv.ParseInt(dropOffType, 10, 64)
		if err != nil {
			log.Println("error parsing drop off type:", err)
			continue
		}

		timepointInt, err := strconv.ParseInt(timepoint, 10, 64)
		if err != nil {
			log.Println("error parsing timepoint:", err)
			continue
		}

		stopTime := model.StopTime{
			TripID:        tripId,
			ArrivalTime:   arrivalTime,
			DepartureTime: departureTime,
			StopID:        stopId,
			StopSequence:  stopSequenceInt,
			PickupType:    pickupTypeInt,
			DropOffType:   dropOffTypeInt,
			Timepoint:     timepointInt,
		}

		stopTimes = append(stopTimes, stopTime)
	}

	return stopTimes, nil
}

func LoadStaticTripDataFromFile(file []byte) ([]model.Trip, error) {
	var trips []model.Trip

	reader := csv.NewReader(bytes.NewReader(file))

	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error reading row:", err)
			continue
		}
		if len(row) < 7 {
			log.Println("row has fewer columns than expected")
			continue
		}

		routeId := row[0]
		serviceId := row[1]
		tripId := row[2]
		tripHeadsign := row[3]
		directionId := row[4]
		blockId := row[5]
		shapeId := row[6]

		routeId = strings.TrimSpace(routeId)
		serviceId = strings.TrimSpace(serviceId)
		tripId = strings.TrimSpace(tripId)
		tripHeadsign = strings.TrimSpace(tripHeadsign)
		directionId = strings.TrimSpace(directionId)
		blockId = strings.TrimSpace(blockId)
		shapeId = strings.TrimSpace(shapeId)

		directionIdInt, err := strconv.ParseInt(directionId, 10, 64)
		if err != nil {
			log.Println("error parsing direction id:", err)
			continue
		}

		trip := model.Trip{
			RouteID:      routeId,
			ServiceID:    serviceId,
			TripID:       tripId,
			TripHeadsign: tripHeadsign,
			DirectionID:  directionIdInt,
			BlockID:      blockId,
			ShapeID:      shapeId,
		}

		trips = append(trips, trip)
	}

	return trips, nil
}
