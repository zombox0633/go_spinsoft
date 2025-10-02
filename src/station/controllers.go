package station

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type StationControllerType struct {
	service StationService
}

func NewStationController(service StationService) *StationControllerType {
	return &StationControllerType{
		service: service,
	}
}

// ---------------------------------- PostImportStationsURL -------------------------
func (c *StationControllerType) PostImportStationsURL(ctx *fiber.Ctx) error {
	var req StationImportRequest

	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := c.service.ImportFromURL(ctx.Context(), req.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(result)
}

// ---------------------------------- Get Nearest Station -------------------------
func (c *StationControllerType) GetNearestStation(ctx *fiber.Ctx) error {
	latStr := ctx.Query("lat")
	longStr := ctx.Query("long")
	limitStr := ctx.Query("limit", "1")

	if latStr == "" || longStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required parameters: lat and long")
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid latitude")
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid longitude")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid limit")
	}

	req := NearestStationRequest{
		Lat:   lat,
		Long:  long,
		Limit: limit,
	}

	result, err := c.service.FindNearestStation(ctx.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(result)
}

// ---------------------------------- Get Nearest Station Pagination -------------------------
func (c *StationControllerType) GetNearestStationPagination(ctx *fiber.Ctx) error {
	latStr := ctx.Query("lat")
	longStr := ctx.Query("long")
	pageStr := ctx.Query("page", "1")
	pageSizeStr := ctx.Query("page_size", "10")

	if latStr == "" || longStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required parameters: lat and long")
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid latitude")
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid longitude")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page format: must be an integer")
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page_size format: must be an integer")
	}

	req := NearestStationPaginationRequest{
		Lat:      lat,
		Long:     long,
		Page:     page,
		PageSize: pageSize,
	}

	result, err := c.service.FindNearestStationPagination(ctx.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}
