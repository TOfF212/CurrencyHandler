package main

import (
	"myproject/internal/config"
	"myproject/internal/database"
	"myproject/internal/migrations"
	"myproject/internal/services"
	"myproject/internal/redis"

	"log"
	"time"
)

func main() {

	cfg := config.LoadConfig()

	db := database.DataBasePostgres{URL: cfg.DatabaseURL}
	var redisDB redis.RedisDataBase
	redisDB.NewClient()
	migrations.RunMigrations(db)
	db.GetCurrencies()
	go func() {
		for {
			newCurrencies, err := services.GetExchangeRate()
			if err != nil {
				log.Println("failed to get Currencies")
			}
			db.UpdateCurrencies(newCurrencies)

			time.Sleep(24 * time.Hour)

		}
	}()
	time.Sleep(24 * time.Hour) // Спим 24 часа

	log.Println("Database migrated successfully!")
}
