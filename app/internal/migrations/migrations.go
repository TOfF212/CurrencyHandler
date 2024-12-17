package migrations

import (
	"api/internal/database"
	"log"

	"gorm.io/gorm"
)

var migrationsToApply = []string{
	"MigrateCurrency",
}

func MigrateTrackingTable(db *gorm.DB) error {
	query := `CREATE TABLE IF NOT EXISTS migrations (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	return db.Exec(query).Error
}

func MigrateCurrency(db *gorm.DB) error {
	query := `CREATE TABLE IF NOT EXISTS currencies (
        id SERIAL PRIMARY KEY,
        currency VARCHAR(3),
        rate DECIMAL(10,4)
    );`
	return db.Exec(query).Error
}

func recordMigration(db *gorm.DB, name string) {
	query := `INSERT INTO migrations (name) VALUES (?)`
	if err := db.Exec(query, name).Error; err != nil {
		log.Fatalf("Error recording migration %s: %v", name, err)
	}
	log.Printf("Migration %s applied successfully.", name)
}

func RunMigrations(db database.DataBasePostgres) {
	db.Open()
	defer db.Close()
	if err := MigrateTrackingTable(db.DataBase); err != nil {
		log.Fatalf("Error creating migrations table: %v", err)
	}

	appliedMigrations := make(map[string]bool)

	rows, err := db.DataBase.Table("migrations").Select("name").Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			rows.Scan(&name)
			appliedMigrations[name] = true
		}
	} else {
		log.Fatalf("Error fetching applied migrations: %v", err)
	}

	for _, migration := range migrationsToApply {
		if !appliedMigrations[migration] {
			switch migration {
			case "MigrateCurrency":
				if err := MigrateCurrency(db.DataBase); err == nil {
					recordMigration(db.DataBase, migration)
				} else {
					log.Fatalf("Error executing migration %s: %v", migration, err)
				}
			}
		} else {
			log.Printf("Migration %s has already been applied.", migration)
		}
	}

	log.Println("Migrations processed successfully.")
}
