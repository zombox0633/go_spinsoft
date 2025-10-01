package station

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type StationService interface {
	ImportFromURL(ctx context.Context, url string) (*StationImportResponse, error)
	NearestStation(ctx context.Context, data NearestStationRequest) ([]NearestStationResponse, error)
}

type stationServiceType struct {
	repo       StationRepository
	httpClient *http.Client
}

func NewStationService(repo StationRepository) StationService {
	return &stationServiceType{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *stationServiceType) ImportFromURL(ctx context.Context, url string) (*StationImportResponse, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	var stations []StationModel
	if err := json.Unmarshal(body, &stations); err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	if err := s.repo.InsertMany(ctx, stations); err != nil {
		return &StationImportResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &StationImportResponse{
		Success:       true,
		ImportedCount: len(stations),
		Message:       "Import completed successfully",
	}, nil
}

func (s *stationServiceType) NearestStation(ctx context.Context, data NearestStationRequest) ([]NearestStationResponse, error) {
	if data.Lat < -90 || data.Lat > 90 {
		return nil, fmt.Errorf("invalid latitude: must be between -90 and 90")
	}

	if data.Long < -180 || data.Long > 180 {
		return nil, fmt.Errorf("invalid longitude: must be between -90 and 90")
	}

	if data.Limit < 0 || data.Limit > 1000 {
		return nil, fmt.Errorf("invalid limit: must be between 0 and 100")
	}

	responses, err := s.repo.FindNearestStation(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearest station: %w", err)
	}

	return responses, nil
}
