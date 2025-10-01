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

	req := NearestStationRequest{
		Lat:  lat,
		Long: long,
	}

	result, err := c.service.NearestStation(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to find nearest station",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(result)
}
