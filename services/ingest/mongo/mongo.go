package mongo

import (
	"context"
	"log"

	"github.com/erikvegax/charlotte-transit-pulse/pkg/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	StaticDataDatabaseName  = "gtfsstatic"
	StopsCollectionName     = "stops"
	StopTimesCollectionName = "stop_times"
	TripsCollectionName     = "trips"
	RoutesCollectionName    = "routes"
)

type MongoClient struct {
	Client *mongo.Client
}

func NewMongoClient(ctx context.Context, uri string) *MongoClient {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &MongoClient{
		Client: client,
	}
}

func (m *MongoClient) Disconnect(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

func (m *MongoClient) SaveStops(ctx context.Context, stops []model.Stop) error {
	collection := m.Client.Database(StaticDataDatabaseName).Collection(StopsCollectionName)

	models := make([]mongo.WriteModel, 0, len(stops))
	for _, stop := range stops {
		models = append(models, mongo.NewReplaceOneModel().
			SetFilter(bson.M{"stopId": stop.StopID}).
			SetReplacement(stop).
			SetUpsert(true))
	}

	_, err := collection.BulkWrite(ctx, models)

	return err
}

func (m *MongoClient) SaveRoutes(ctx context.Context, routes []model.Route) error {
	collection := m.Client.Database(StaticDataDatabaseName).Collection(RoutesCollectionName)

	models := make([]mongo.WriteModel, 0, len(routes))
	for _, route := range routes {
		models = append(models, mongo.NewReplaceOneModel().
			SetFilter(bson.M{"routeId": route.RouteID}).
			SetReplacement(route).
			SetUpsert(true))
	}

	_, err := collection.BulkWrite(ctx, models)

	return err
}

func (m *MongoClient) SaveStopTimes(ctx context.Context, stopTimes []model.StopTime) error {
	collection := m.Client.Database(StaticDataDatabaseName).Collection(StopTimesCollectionName)

	models := make([]mongo.WriteModel, 0, len(stopTimes))
	for _, st := range stopTimes {
		models = append(models, mongo.NewReplaceOneModel().
			SetFilter(bson.M{"tripId": st.TripID, "stopSequence": st.StopSequence}).
			SetReplacement(st).
			SetUpsert(true))
	}

	_, err := collection.BulkWrite(ctx, models)
	return err
}

func (m *MongoClient) SaveTrips(ctx context.Context, trips []model.Trip) error {
	collection := m.Client.Database(StaticDataDatabaseName).Collection(TripsCollectionName)

	models := make([]mongo.WriteModel, 0, len(trips))
	for _, trip := range trips {
		models = append(models, mongo.NewReplaceOneModel().
			SetFilter(bson.M{"tripId": trip.TripID}).
			SetReplacement(trip).
			SetUpsert(true))
	}

	_, err := collection.BulkWrite(ctx, models)

	return err
}

// EnsureIndexes creates the unique indexes each collection's upsert filter
// relies on. Without these, every upsert in a bulk write has to scan the
// collection to check for an existing match, which is fine at small scale
// but becomes extremely slow as a collection grows into the hundreds of
// thousands of documents (as stop_times does). Creating an index that
// already exists is a safe no-op, so this can run on every startup.
func (m *MongoClient) EnsureIndexes(ctx context.Context) error {
	db := m.Client.Database(StaticDataDatabaseName)

	_, err := db.Collection(RoutesCollectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "routeId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = db.Collection(StopsCollectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "stopId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = db.Collection(TripsCollectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tripId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = db.Collection(StopTimesCollectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tripId", Value: 1}, {Key: "stopSequence", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	return nil
}
