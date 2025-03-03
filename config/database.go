package config

import (
    "log"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "sticker-go/models"
)

var DB *gorm.DB

func InitDB() {
    var err error
    DB, err = gorm.Open(sqlite.Open("storage/sticker-go.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    log.Println("Database connected!")
    DB.AutoMigrate(&models.User{})
}
