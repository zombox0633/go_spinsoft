package station

import (
	"context"
	"fmt"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StationRepository interface {
	InsertMany(ctx context.Context, stations []StationModel) error
	FindNearestStation(ctx context.Context, data NearestStationRequest) (*NearestStationResponse, error)
	CreateGeoIndex(ctx context.Context) error
}

type stationRepositoryType struct {
	collection *mongo.Collection
}

func NewStationRepository(collection *mongo.Collection) StationRepository {
	return &stationRepositoryType{
		collection: collection,
	}
}

func (r *stationRepositoryType) InsertMany(ctx context.Context, stations []StationModel) error {
	if len(stations) == 0 {
		return nil
	}

	data := make([]interface{}, len(stations))
	for i, station := range stations {
		data[i] = station
	}

	result, err := r.collection.InsertMany(ctx, data)

	if err != nil {
		return fmt.Errorf("failed to insert stations: %w", err)
	}

	fmt.Printf("Successfully inserted %d stations\n", len(result.InsertedIDs))
	return nil
}

func (r *stationRepositoryType) FindNearestStation(ctx context.Context, data NearestStationRequest) (*NearestStationResponse, error) {
	searchPoint := bson.M{
		"type":        "Point",
		"coordinates": []float64{data.Long, data.Lat},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$geoNear", Value: bson.M{
			"near":          searchPoint,
			"distanceField": "distance",
			"maxDistance":   100000, //100km
			"spherical":     true,
			"query": bson.M{
				"active":   1,
				"location": bson.M{"$exists": true},
			},
		}}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute geoNear: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		StationModel `bson:",inline"`
		Distance     float64 `bson:"distance"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no stations found within 100km")
	}

	station := results[0]

	distanceKm := math.Round((station.Distance/1000)*100) / 100

	response := &NearestStationResponse{
		ID:       station.StationID,
		Name:     station.Name,
		EnName:   station.EnName,
		Lat:      station.Lat,
		Long:     station.Long,
		Distance: distanceKm,
	}

	return response, nil
}

func (r *stationRepositoryType) CreateGeoIndex(ctx context.Context) error {
	indexStation := mongo.IndexModel{
		Keys: bson.D{
			{Key: "location", Value: "2dsphere"},
		},
		Options: options.Index().SetName("location_2dsphere"),
	}

	if _, err := r.collection.Indexes().CreateOne(ctx, indexStation); err != nil {
		return fmt.Errorf("failed to create geo index: %w", err)
	}

	return nil
}
