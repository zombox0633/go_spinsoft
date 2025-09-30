package station

type StationImportRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type StationImportResponse struct {
	Success       bool   `json:"success"`
	ImportedCount int    `json:"imported_count"`
	Message       string `json:"message"`
}
