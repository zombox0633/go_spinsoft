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

	if _, err := r.collection.InsertMany(ctx, data); err != nil {
		return fmt.Errorf("failed to insert stations: %w", err)
	}

	return nil
}
