package database

import (
	"log"

	"gorm.io/gorm"

	"github.com/undercode99/article_service/config"
	"github.com/undercode99/article_service/internal/app/article"
	"gorm.io/driver/postgres"
)

func NewDatabase(cfg *config.Config) *gorm.DB {
	cfgDsn := cfg.Database.Dsn
	log.Println("Creating a new database client")
	// Open a connection to the database using the provided database configuration.
	log.Println("Opening a connection to the database")
	db, err := gorm.Open(postgres.Open(cfgDsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Connection to the database opened successfully")
	return db
}

func MigrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(&article.Article{})
}
