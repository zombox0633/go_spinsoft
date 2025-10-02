package station

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/zombox0633/go_spinsoft/src/utils"
)

type StationService interface {
	ImportFromURL(ctx context.Context, url string) (*StationImportResponse, error)
	FindNearestStation(ctx context.Context, data NearestStationRequest) (*NearestStationResponse, error)
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

	invalidCoordinateCount := 0
	for _, station := range stations {
		if station.WasInvalidated {
			invalidCoordinateCount++
		}
	}

	if err := s.repo.InsertMany(ctx, stations); err != nil {
		return &StationImportResponse{
			Success:            false,
			ImportedCount:      0,
			InvalidCoordinates: invalidCoordinateCount,
			Message:            err.Error(),
		}, nil
	}

	return &StationImportResponse{
		Success:            true,
		ImportedCount:      len(stations),
		InvalidCoordinates: invalidCoordinateCount,
		Message:            "Import completed successfully",
	}, nil
}

// ---------------------------------- Find Nearest Station -------------------------
func (s *stationServiceType) FindNearestStation(ctx context.Context, data NearestStationRequest) (*NearestStationResponse, error) {
	if err := utils.ValidateCoordinates(data.Lat, data.Long); err != nil {
		return nil, err
	}

	if data.Limit < 1 || data.Limit > 100 {
		return nil, fmt.Errorf("invalid limit: must be between 1 and 100")
	}

	stationData, err := s.repo.FindNearestStation(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearest station: %w", err)
	}

	response := &NearestStationResponse{
		Success: true,
		Data:    stationData,
	}

	return response, nil
}

// ---------------------------------- Find Nearest Station Pagination -------------------------
func (s *stationServiceType) FindNearestStationPagination(ctx context.Context, data NearestStationPaginationRequest) (*NearestStationPaginationResponse, error) {
	if err := utils.ValidateCoordinates(data.Lat, data.Long); err != nil {
		return nil, err
	}

	page := data.Page
	pageSize := data.PageSize
	if page < 0 {
		return nil, fmt.Errorf("invalid page: must be greater than 0")
	}

	if pageSize < 1 || pageSize > 100 {
		return nil, fmt.Errorf("invalid page_size: must be between 1 and 100")
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
