package database

import (
	"fmt"

	"github.com/atharv-bhadange/producer_consumer/configs"
	"github.com/atharv-bhadange/producer_consumer/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func InitDatabase() error {

	pg_user, pg_password, pg_db, host, err := configs.LoadDatabaseConfig()

	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", pg_user, pg_password, host, "5432", pg_db)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	db.AutoMigrate(&models.User{}, &models.Product{})

	DB = Dbinstance{
		Db: db,
	}

	return nil
}
