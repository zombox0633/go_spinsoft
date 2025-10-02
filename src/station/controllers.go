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

func (c *StationControllerType) PostImportStationsURL(ctx *fiber.Ctx) error {
	var req StationImportRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	result, err := c.service.ImportFromURL(ctx.Context(), req.URL)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (c *StationControllerType) GetNearestStation(ctx *fiber.Ctx) error {
	latStr := ctx.Query("lat")
	longStr := ctx.Query("long")
	limitStr := ctx.Query("limit", "1")

	if latStr == "" || longStr == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"error":   "Missing required parameters",
		})
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid latitude",
		})
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid longitude",
		})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid limit",
		})
	}

	req := NearestStationRequest{
		Lat:   lat,
		Long:  long,
		Limit: limit,
	}

	result, err := c.service.FindNearestStation(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to find nearest station",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (c *StationControllerType) GetNearestStationPagination(ctx *fiber.Ctx) error {
	latStr := ctx.Query("lat")
	longStr := ctx.Query("long")
	pageStr := ctx.Query("page", "1")
	pageSizeStr := ctx.Query("page_size", "10")

	if latStr == "" || longStr == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"error":   "Missing required parameters",
		})
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid latitude",
		})
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid longitude",
		})
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid page format: must be an integer",
		})
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid page_size format: must be an integer",
		})
	}

	req := NearestStationPaginationRequest{
		Lat:      lat,
		Long:     long,
		Page:     page,
		PageSize: pageSize,
	}

	result, err := c.service.FindNearestStationPagination(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}
