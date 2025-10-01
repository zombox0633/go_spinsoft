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
	Client *mongo.Client
	DBName *mongo.Database
}

var DB *DatabaseType

func InitDatabase(ctx context.Context, cfg *ConfigType) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if cfg.MonGoURL == "" {
		return fmt.Errorf("MongoDB URL is empty: check your .env file")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MonGoURL))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// ping
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(cfg.DataBaseName)

	DB = &DatabaseType{
		Client: client,
		DBName: database,
	}

	log.Println("Successfully connected to MongoDB")
	return nil
}

func (d *DatabaseType) Close(ctx context.Context) error {
	if d.Client != nil {
		log.Println("Closing MongoDB connection...")
		return d.Client.Disconnect(ctx)
	}

	log.Println("MongoDB connection closed successfully")
	return nil
}
