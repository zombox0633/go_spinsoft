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
	}, nil
}
