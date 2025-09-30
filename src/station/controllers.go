package station

import "github.com/gofiber/fiber/v2"

type StationControllerType struct {
	service StationService
}

func NewStationController(service StationService) *StationControllerType {
	return &StationControllerType{
		service: service,
	}
}

func (c *StationControllerType) ImportStationsURL(ctx *fiber.Ctx) error {
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
