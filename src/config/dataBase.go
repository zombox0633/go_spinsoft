package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseType struct {
	Client     *mongo.Client
	DBName     *mongo.Database
	Collection *mongo.Collection
}

var DB *DatabaseType

func InitDatabase(ctx context.Context, cfg *ConfigType) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if cfg.MonGoURL == "" {
		return fmt.Errorf("MongoDB URL is empty - check your .env file")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MonGoURL))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// ping
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println(cfg.CollectionName)

	database := client.Database(cfg.DataBaseName)
	collection := database.Collection(cfg.CollectionName)

	DB = &DatabaseType{
		Client:     client,
		DBName:     database,
		Collection: collection,
	}

	log.Println("Successfully connected to MongoDB")
	return nil
}

func (d *DatabaseType) Close(ctx context.Context) error {
	if d.Client != nil {
		log.Panicln("Closing MongoDB connection...")
		return d.Client.Disconnect(ctx)
	}
	return nil
}
