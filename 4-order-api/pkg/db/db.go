package db

import (
	"demo/order-api/configs"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Db *gorm.DB
}

func NewDB(config *configs.Config) *DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Db.Host,
		config.Db.Port,
		config.Db.User,
		config.Db.Password,
		config.Db.DbName,
		config.Db.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return &DB{Db: db}
}

func (d *DB) Migrate(model interface{}) {
	d.Db.AutoMigrate(model)
}
