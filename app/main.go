package main

import (
	"myproject/internal/config"
	"myproject/internal/database"
	"myproject/internal/handlers"
	"myproject/internal/migrations"
	"myproject/internal/redis"
	"myproject/internal/services"

	"log"
	"net/http"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	db := database.DataBasePostgres{URL: cfg.DatabaseURL}
	var rdb redis.RedisDataBase
	rdb.Init()
	migrations.RunMigrations(db)

	go func() {
		for {
			newCurrencies, err := services.GetExchangeRate()
			if err != nil {
				log.Println("failed to get Currencies")
			}
			err = db.UpdateCurrencies(newCurrencies)
			if err != nil {
				log.Fatalf("failed to update database: %v", err)
			}
			currencies := db.GetCurrencies()
			rdb.SetCurrencies(currencies)
			time.Sleep(24 * time.Hour)

		}
	}()
	http.HandleFunc("/convert", handlers.CurrencyTransferHandle)

	log.Println("Database migrated successfully!")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Ошибка при запуске сервера:", err)
	}
}
