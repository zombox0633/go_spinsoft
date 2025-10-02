package station

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

type StationService interface {
	ImportFromURL(ctx context.Context, url string) (*StationImportResponse, error)
	FindNearestStation(ctx context.Context, data NearestStationRequest) ([]NearestStationResponse, error)
	FindNearestStationPagination(ctx context.Context, data NearestStationPaginationRequest) (*NearestStationPaginationResponse, error)
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

// ---------------------------------- ImportFromURL -------------------------
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

// ---------------------------------- Find NearestStation -------------------------
func (s *stationServiceType) FindNearestStation(ctx context.Context, data NearestStationRequest) ([]NearestStationResponse, error) {
	if data.Lat < -90 || data.Lat > 90 {
		return nil, fmt.Errorf("invalid latitude: must be between -90 and 90")
	}

	if data.Long < -180 || data.Long > 180 {
		return nil, fmt.Errorf("invalid longitude: must be between -90 and 90")
	}

	if data.Limit < 0 || data.Limit > 100 {
		return nil, fmt.Errorf("invalid limit: must be between 0 and 100")
	}

	responses, err := s.repo.FindNearestStation(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearest station: %w", err)
	}

	return responses, nil
}

func (s *stationServiceType) FindNearestStationPagination(ctx context.Context, data NearestStationPaginationRequest) (*NearestStationPaginationResponse, error) {
	if data.Lat < -90 || data.Lat > 90 {
		return nil, fmt.Errorf("invalid latitude: must be between -90 and 90")
	}

	if data.Long < -180 || data.Long > 180 {
		return nil, fmt.Errorf("invalid longitude: must be between -90 and 90")
	}

	page := data.Page
	pageSize := data.PageSize
	if page < 0 {
		return nil, fmt.Errorf("invalid page: must be greater than 0")
	}
	if pageSize < 0 {
		return nil, fmt.Errorf("invalid page_size: must be greater than 0")
	}
	if pageSize > 100 {
		return nil, fmt.Errorf("invalid page_size: must be less than or equal to 100")
	}

	station, totalItems, err := s.repo.FindNearestStationPagination(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearest stations: %w", err)
	}

	itemStart := (page-1)*pageSize + 1
	itemEnd := itemStart + len(station) - 1
	TotalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	response := &NearestStationPaginationResponse{
		Success:    true,
		Page:       page,
		PageSize:   pageSize,
		PagesItems: len(station),
		ItemStart:  itemStart,
		ItemEnd:    itemEnd,
		TotalItems: totalItems,
		TotalPages: TotalPages,
		Data:       station,
	}
	return response, nil
}
