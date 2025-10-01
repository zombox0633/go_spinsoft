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
	Lat  float64 `json:"lat" validate:"required"`
	Long float64 `json:"long" validate:"required"`
}

type NearestStationResponse struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	EnName   string  `json:"en_name"`
	Lat      float64 `json:"lat"`
	Long     float64 `json:"long"`
	Distance float64 `json:"distance"`
}
