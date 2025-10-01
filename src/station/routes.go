package station

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func StationRoutes(api fiber.Router, DB *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := DB.Collection("stations")
	stationRepo := NewStationRepository(collection)

	if err := stationRepo.CreateGeoIndex(ctx); err != nil {
		log.Printf("Warning: Failed to create geo index: %v", err)
	}
	log.Print("Station index ready")

	stationService := NewStationService(stationRepo)
	stationController := NewStationController(stationService)

	stationGroup := api.Group("/station")

	stationGroup.Post("/import", stationController.PostImportStationsURL)
	stationGroup.Get("/nearest", stationController.GetNearestStation)
}
