package station

// Station Import
type StationImportRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type StationImportResponse struct {
	Success       bool   `json:"success"`
	ImportedCount int    `json:"imported_count"`
	Message       string `json:"message"`
}

// Find Near Station
type NearestStationRequest struct {
	Lat   float64 `json:"lat" validate:"required"`
	Long  float64 `json:"long" validate:"required"`
	Limit int     `json:"limit,omitempty"`
}

type NearestStationResponse struct {
	Success bool                 `json:"success"`
	Data    []NearestStationData `json:"data"`
}

type NearestStationData struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	EnName   string  `json:"en_name"`
	Lat      float64 `json:"lat"`
	Long     float64 `json:"long"`
	Distance float64 `json:"distance_km"`
}

// Pagination
type NearestStationPaginationRequest struct {
	Lat      float64 `json:"lat" validate:"required"`
	Long     float64 `json:"long" validate:"required"`
	Page     int     `json:"page,omitempty"`
	PageSize int     `json:"page_size,omitempty"`
}

type NearestStationPaginationResponse struct {
	Success    bool                 `json:"success"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	PagesItems int                  `json:"pages_items"`
	ItemStart  int                  `json:"item_start"`
	ItemEnd    int                  `json:"item_end"`
	TotalPages int                  `json:"total_pages"`
	TotalItems int                  `json:"total_items"`
	Data       []NearestStationData `json:"data"`
}
