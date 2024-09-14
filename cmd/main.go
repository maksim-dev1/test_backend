package main

import (
	"log"
	"yummy_mobile_app_backend/configs"
	"yummy_mobile_app_backend/internal/handlers"
	"yummy_mobile_app_backend/internal/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Подключение к базе данных
	db := configs.ConnectDB()

	// Миграция модели в базу данных
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Инициализация Gin
	r := gin.Default()

	// Передаем базу данных в хендлеры
	handlers.SetupRoutes(r, db)

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
