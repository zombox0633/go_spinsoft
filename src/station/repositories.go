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
	FindNearestStation(ctx context.Context, data NearestStationRequest) ([]NearestStationResponse, error)
	FindNearestStationPagination(ctx context.Context, data NearestStationPaginationRequest) ([]NearestStationResponse, int, error)
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

// ---------------------------------- ImportFromURL -------------------------
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

// ---------------------------------- Find NearestStation -------------------------
func (r *stationRepositoryType) FindNearestStation(ctx context.Context, data NearestStationRequest) ([]NearestStationResponse, error) {
	searchPoint := bson.M{
		"type":        "Point",
		"coordinates": []float64{data.Long, data.Lat},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$geoNear", Value: bson.M{
			"near":          searchPoint,
			"distanceField": "distance",
			"maxDistance":   10000, //10km
			"spherical":     true,
			"query": bson.M{
				"active":   1,
				"location": bson.M{"$exists": true},
			},
		}}},
		{{Key: "$limit", Value: data.Limit}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute geoNear: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		StationID int     `bson:"id"`
		Name      string  `bson:"name"`
		EnName    string  `bson:"en_name"`
		Lat       float64 `bson:"lat"`
		Long      float64 `bson:"long"`
		Distance  float64 `bson:"distance"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no stations found within 100km")
	}

	responses := make([]NearestStationResponse, len(results))
	for i, station := range results {
		distanceKm := math.Round((station.Distance/1000)*1000) / 1000

		responses[i] = NearestStationResponse{
			ID:       station.StationID,
			Name:     station.Name,
			EnName:   station.EnName,
			Lat:      station.Lat,
			Long:     station.Long,
			Distance: distanceKm,
		}
	}

	return responses, nil
}

func (r *stationRepositoryType) FindNearestStationPagination(ctx context.Context, data NearestStationPaginationRequest) ([]NearestStationResponse, int, error) {

	pageSize := data.PageSize

	start := (data.Page - 1) * pageSize

	searchPoint := bson.M{
		"type":        "Point",
		"coordinates": []float64{data.Long, data.Lat},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$geoNear", Value: bson.M{
			"near":          searchPoint,
			"distanceField": "distance",
			"spherical":     true,
			"query": bson.M{
				"active":   1,
				"location": bson.M{"$exists": true},
			},
		}}},
		{{Key: "$project", Value: bson.M{
			"id":       1,
			"name":     1,
			"en_name":  1,
			"lat":      1,
			"long":     1,
			"distance": 1,
		}}},
		{{Key: "$facet", Value: bson.M{
			"metadata": []bson.M{
				{"$count": "total"},
			},
			"data": []bson.M{
				{"$skip": start},
				{"$limit": pageSize},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute geoNear: %w", err)
	}
	defer cursor.Close(ctx)

	var pipelineResult []struct {
		Metadata []struct {
			Total int `bson:"total"`
		} `bson:"metadata"`
		Data []struct {
			StationID int     `bson:"id"`
			Name      string  `bson:"name"`
			EnName    string  `bson:"en_name"`
			Lat       float64 `bson:"lat"`
			Long      float64 `bson:"long"`
			Distance  float64 `bson:"distance"`
		} `bson:"data"`
	}

	if err := cursor.All(ctx, &pipelineResult); err != nil {
		return nil, 0, fmt.Errorf("failed to decode results: %w", err)
	}

	if len(pipelineResult) == 0 {
		return nil, 0, fmt.Errorf("no results returned")
	}

	result := pipelineResult[0]
	totalItems := result.Metadata[0].Total

	if totalItems == 0 {
		return nil, 0, fmt.Errorf("no results returned")
	}

	responses := make([]NearestStationResponse, len(result.Data))
	for i, station := range result.Data {
		distanceKm := math.Round((station.Distance/1000)*1000) / 1000

		responses[i] = NearestStationResponse{
			ID:       station.StationID,
			Name:     station.Name,
			EnName:   station.EnName,
			Lat:      station.Lat,
			Long:     station.Long,
			Distance: distanceKm,
		}
	}

	return responses, totalItems, nil
}

// ---------------------------------- CreateGeoIndex -------------------------
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
