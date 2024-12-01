package database

import (
	"log"
	"myproject/internal/config"

	"myproject/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDataBasePostrges(conf config.Config) *DataBasePostgres {
	return &DataBasePostgres{URL: conf.DatabaseURL}
}

type DataBasePostgres struct {
	URL      string
	DataBase *gorm.DB
}

func (d *DataBasePostgres) Open() {
	db, err := gorm.Open(postgres.Open(d.URL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	d.DataBase = db
}

func (d *DataBasePostgres) Close() {
	sqlDB, err := d.DataBase.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("failed to close database connection: %v", err)
	}
}

func (d *DataBasePostgres) GetCurrencies() []models.Currency {
	d.Open()
	defer d.Close()
	var currencies []models.Currency
	if err := d.DataBase.Find(&currencies).Error; err != nil {
		log.Fatalf("failed to fetch currencies: %v", err)
	}

	return currencies
}

func (d *DataBasePostgres) UpdateCurrencies(newCurrencies map[string]float64) error {
	d.Open()
	defer d.Close()

	for currencyCode, rate := range newCurrencies {
		newCurrency := models.Currency{Currency: currencyCode, Rate: rate}
		var currency models.Currency
		err := d.DataBase.Where("currency = ?", currencyCode).First(&currency).Error
		if err != nil {
			log.Fatalf("failed to get currencies from database, %v", err)
		}
		if currency.ID == 0 {
			err = d.DataBase.Save(&newCurrency).Error
			if err != nil {
				log.Printf("failed to save new currency %s: %v", currencyCode, err)
				return err
			}
			continue
		}

		if currency.Rate != rate {
			currency.Rate = rate
			if err := d.DataBase.Save(&currency).Error; err != nil {
				log.Fatalf("error when executing the request: %v", err)
				return err
			}
		}
	}
	return nil
}
