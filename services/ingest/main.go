package main

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/erikvegax/charlotte-transit-pulse/pkg/config"
	"github.com/erikvegax/charlotte-transit-pulse/pkg/model"
	"github.com/erikvegax/charlotte-transit-pulse/services/ingest/gtfsrt"
	"github.com/erikvegax/charlotte-transit-pulse/services/ingest/gtfsstatic"
	"github.com/erikvegax/charlotte-transit-pulse/services/ingest/kafka"
	"github.com/erikvegax/charlotte-transit-pulse/services/ingest/mongo"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
)

const stopsFileName = "stops.txt"
const routesFileName = "routes.txt"
const stopTimesFileName = "stop_times.txt"
const tripsFileName = "trips.txt"

type StaticData interface {
	SaveStops(context.Context, []model.Stop) error
	SaveRoutes(context.Context, []model.Route) error
	SaveStopTimes(context.Context, []model.StopTime) error
	SaveTrips(context.Context, []model.Trip) error
}

type RealTimeData interface {
	PublishVehicles(ctx context.Context, vehicles []model.Vehicle) error
	PublishTripUpdates(ctx context.Context, tripUpdates []model.TripUpdate) error
	PublishAlerts(ctx context.Context, alerts []model.Alert) error
}

func main() {
	log.Println("starting ingest service")
	godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	client := mongo.NewMongoClient(context.Background(), config.MongoURI)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	if err := client.EnsureIndexes(context.Background()); err != nil {
		log.Fatal(err)
	}

	writer := kafka.NewKafkaWriter(config.KafkaBrokers)

	scheduler, err := gocron.NewScheduler(gocron.WithStopTimeout(60 * time.Second))
	if err != nil {
		log.Fatal(err)
	}

	// static data job
	scheduler.NewJob(
		gocron.DurationJob(time.Duration(config.StaticResyncIntervalInHours)*time.Hour),
		gocron.NewTask(func() {
			loadStaticData(config.CatsStaticURL, client)
		}),
		gocron.WithStartAt(gocron.WithStartImmediately()),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)

	// real time feed job
	scheduler.NewJob(
		gocron.DurationJob(time.Duration(config.PollIntervalInSeconds)*time.Second),
		gocron.NewTask(func() {
			loadRealTimeFeedData(config.CatsVehiclePositionURL, config.CatsAlertsURL, config.CatsTripUpdateURL, writer)
		}),
		gocron.WithStartAt(gocron.WithStartImmediately()),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)

	scheduler.Start()
	defer scheduler.Shutdown()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
}

func loadStaticData(url string, client StaticData) {
	log.Println("loading static data into mongo")
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error in http call to %v: %v", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading body: %v", err)
		return
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Printf("error creating reader: %v", err)
		return
	}

	var unzippedStopsFile, unzippedRoutesFile, unzippedStopTimesFile, unzippedTripsFile []byte

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			log.Printf("error reading zip file %s: %v", zipFile.Name, err)
			continue
		}

		switch zipFile.Name {
		case stopsFileName:
			unzippedStopsFile = unzippedFileBytes
		case routesFileName:
			unzippedRoutesFile = unzippedFileBytes
		case stopTimesFileName:
			unzippedStopTimesFile = unzippedFileBytes
		case tripsFileName:
			unzippedTripsFile = unzippedFileBytes
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		if unzippedRoutesFile == nil {
			log.Println("routes file not found in zip archive")
			return
		}

		routes, err := gtfsstatic.LoadStaticRouteDataFromFile(unzippedRoutesFile)
		if err != nil {
			log.Printf("error loading static route data: %v", err)
			return
		}

		err = client.SaveRoutes(context.Background(), routes)
		if err != nil {
			log.Printf("error saving routes to mongo: %v", err)
			return
		}
		log.Printf("successfully saved %d routes to mongo", len(routes))
	}()

	go func() {
		defer wg.Done()
		if unzippedStopsFile == nil {
			log.Println("stops file not found in zip archive")
			return
		}

		stops, err := gtfsstatic.LoadStaticStopDataFromFile(unzippedStopsFile)
		if err != nil {
			log.Printf("error loading static stop data: %v", err)
			return
		}

		err = client.SaveStops(context.Background(), stops)
		if err != nil {
			log.Printf("error saving stops to mongo: %v", err)
			return
		}
		log.Printf("successfully saved %d stops to mongo", len(stops))
	}()

	go func() {
		defer wg.Done()
		if unzippedStopTimesFile == nil {
			log.Println("stop times file not found in zip archive")
			return
		}

		stopTimes, err := gtfsstatic.LoadStaticStopTimeDataFromFile(unzippedStopTimesFile)
		if err != nil {
			log.Printf("error loading static stop time data: %v", err)
			return
		}

		err = client.SaveStopTimes(context.Background(), stopTimes)
		if err != nil {
			log.Printf("error saving stop times to mongo: %v", err)
			return
		}
		log.Printf("successfully saved %d stop times to mongo", len(stopTimes))
	}()

	go func() {
		defer wg.Done()
		if unzippedTripsFile == nil {
			log.Println("trips file not found in zip archive")
			return
		}

		trips, err := gtfsstatic.LoadStaticTripDataFromFile(unzippedTripsFile)
		if err != nil {
			log.Printf("error loading static trip data: %v", err)
			return
		}

		err = client.SaveTrips(context.Background(), trips)
		if err != nil {
			log.Printf("error saving trips to mongo: %v", err)
			return
		}
		log.Printf("successfully saved %d trips to mongo", len(trips))
	}()

	wg.Wait()

	log.Println("finished loading static data into mongo")
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func loadRealTimeFeedData(vehiclesUrl, alertsUrl, tripUpdatesUrl string, writer RealTimeData) {
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		resp, err := http.Get(vehiclesUrl)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		vehicles, err := gtfsrt.LoadVehiclesFromFeed(body)
		if err != nil {
			log.Println(err)
			return
		}

		err = writer.PublishVehicles(context.Background(), vehicles)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("published %d vehicles to %v", len(vehicles), "vehicle-positions")
	}()

	go func() {
		defer wg.Done()

		resp, err := http.Get(alertsUrl)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		alerts, err := gtfsrt.LoadAlertsFromFeed(body)
		if err != nil {
			log.Println(err)
			return
		}

		err = writer.PublishAlerts(context.Background(), alerts)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("published %d alerts to %v", len(alerts), "service-alerts")
	}()

	go func() {
		defer wg.Done()

		resp, err := http.Get(tripUpdatesUrl)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		tripUpdates, err := gtfsrt.LoadTripUpdatesFromFeed(body)
		if err != nil {
			log.Println(err)
			return
		}

		err = writer.PublishTripUpdates(context.Background(), tripUpdates)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("published %d trip updates to %v", len(tripUpdates), "trip-updates")
	}()

	wg.Wait()
}
