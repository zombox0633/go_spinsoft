package station

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func StationRoutes(api fiber.Router, DB *mongo.Database) {
	collection := DB.Collection("stations")

	stationRepo := NewStationRepository(collection)
	stationService := NewStationService(stationRepo)
	stationController := NewStationController(stationService)

	stationGroup := api.Group("/station")

	stationGroup.Post("/import", stationController.ImportStationsURL)
}
