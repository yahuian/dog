package model

import (
	"dog/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func Init() error {
	dsn := config.Get().DB.DSN

	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect db err: %w", err)
	}

	db = client

	return nil
}
