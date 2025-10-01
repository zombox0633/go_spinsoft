package station

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type StationRepository interface {
	InsertMany(ctx context.Context, stations []StationModel) error
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
