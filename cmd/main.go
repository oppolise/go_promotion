package main

import (
	"goPromotion/cmd/database"
	serverconfig "goPromotion/config/server_config"
	"goPromotion/handler"
	"goPromotion/pkg/repository"
	"goPromotion/pkg/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.InitDatabase()

	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	if err := database.Migration(db); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	app := fiber.New()

	orderPepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderImpService(orderPepo)
	orderHandler := handler.NewHttpOrderHandler(orderService)

	app.Get("/order/:id", orderHandler.GetOrder)

	app.Listen(":" + serverconfig.ServerConfig().PORT)
}
