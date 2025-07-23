package main

import (
	"auth-svc/config"
	"auth-svc/internal/port"
	"auth-svc/internal/services"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(cfg.DB.Mongo.DSN),
	)
	db := client.Database(cfg.DB.Mongo.Database)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.PWD,
		DB:       cfg.Redis.Database,
	})

	svc := services.NewService(cfg, db, redisClient)

	server := port.NewServer(cfg, svc)
	server.Run()
}
